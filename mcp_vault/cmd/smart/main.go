// Command smart is the "Cena 2" MCP server: it loads every Markdown file in
// the Obsidian vault into an in-memory DuckDB vector index (embeddings via
// Kronk, running in-process) and exposes a single tool, search_vault, that
// returns only the top-K most relevant documents for a query.
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/ardanlabs/kronk/sdk/kronk"
	"github.com/ardanlabs/kronk/sdk/kronk/model"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jeff/mcp-vault/internal/store"
	"github.com/jeff/mcp-vault/internal/vault"
)

const defaultVaultPath = "/Users/jeffersonferreira/Developer/kronk-tech-demo-golangsp/obsidian_vault"

func embedModelPath() string {
	if v := os.Getenv("KRONK_EMBED_MODEL"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return home + "/.kronk/models/ggml-org/embeddinggemma-300m-qat-q8_0-GGUF/embeddinggemma-300m-qat-Q8_0.gguf"
}

const addr = ":9002"

const defaultTopK = 3

var callCount atomic.Int64

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	vaultPath := os.Getenv("VAULT_PATH")
	if vaultPath == "" {
		vaultPath = defaultVaultPath
	}

	log.Println("[vault_smart] carregando modelo de embeddings...")
	modelStart := time.Now()

	if err := kronk.Init(); err != nil {
		return fmt.Errorf("kronk init: %w", err)
	}

	krnEmbed, err := kronk.New(model.WithModelFiles([]string{embedModelPath()}))
	if err != nil {
		return fmt.Errorf("kronk new: %w", err)
	}
	defer krnEmbed.Unload(context.Background())

	log.Printf("[vault_smart] modelo carregado em %s", time.Since(modelStart).Round(time.Millisecond))

	log.Println("[vault_smart] lendo vault e construindo índice vetorial...")
	indexStart := time.Now()

	docs, err := vault.LoadAll(vaultPath)
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
	}

	// Compute full vault stats for later comparison in search logs.
	var vaultTotalChars int
	for _, d := range docs {
		vaultTotalChars += len(d.Content)
	}

	dims, err := embeddingDims(krnEmbed)
	if err != nil {
		return fmt.Errorf("measure embedding dims: %w", err)
	}

	ctx := context.Background()

	db, err := store.LoadVault(ctx, krnEmbed, docs, dims)
	if err != nil {
		return fmt.Errorf("load vault into store: %w", err)
	}
	defer db.Close()

	indexDur := time.Since(indexStart)
	log.Printf("[vault_smart] índice pronto: docs=%d  dims=%d  vault_chars=%d  vault_est_tokens=%d  index_time=%s",
		len(docs), dims, vaultTotalChars, estTokens(vaultTotalChars), indexDur.Round(time.Millisecond))

	server := mcp.NewServer(&mcp.Implementation{Name: "vault_smart", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_vault",
		Description: "Searches the Obsidian vault and returns only the top-K Markdown documents most relevant to the query, ranked by semantic similarity.",
	}, searchVaultHandler(krnEmbed, db, vaultTotalChars))

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)

	log.Printf("[vault_smart] listening on %s/mcp  vault=%s", addr, vaultPath)

	return http.ListenAndServe(addr, mux)
}

func embeddingDims(krnEmbed *kronk.Kronk) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := krnEmbed.Embeddings(ctx, model.D{"input": "dimension probe"})
	if err != nil {
		return 0, fmt.Errorf("embeddings: %w", err)
	}

	if len(resp.Data) == 0 {
		return 0, fmt.Errorf("empty embeddings response")
	}

	return len(resp.Data[0].Embedding), nil
}

type searchVaultParams struct {
	Query string `json:"query" jsonschema:"the search query describing what to look for in the vault"`
	TopK  int    `json:"top_k,omitempty" jsonschema:"number of results to return (default 3)"`
}

type searchVaultResult struct {
	Results []store.Result `json:"results"`
}

func searchVaultHandler(krnEmbed *kronk.Kronk, db *sql.DB, vaultTotalChars int) func(context.Context, *mcp.CallToolRequest, searchVaultParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params searchVaultParams) (*mcp.CallToolResult, any, error) {
		n := callCount.Add(1)
		start := time.Now()
		topK := params.TopK
		if topK <= 0 {
			topK = defaultTopK
		}

		log.Printf("[vault_smart] call #%d — search_vault  query=%q  top_k=%d", n, params.Query, topK)

		// Phase 1: embed query.
		embedStart := time.Now()
		queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		resp, err := krnEmbed.Embeddings(queryCtx, model.D{"input": params.Query})
		if err != nil {
			return nil, nil, fmt.Errorf("embed query: %w", err)
		}
		if len(resp.Data) == 0 {
			return nil, nil, fmt.Errorf("empty embedding for query")
		}
		embedDur := time.Since(embedStart)
		log.Printf("[vault_smart] call #%d — embed_query: dims=%d  duration=%s",
			n, len(resp.Data[0].Embedding), embedDur.Round(time.Millisecond))

		// Phase 2: vector search.
		searchStart := time.Now()
		results, err := store.Search(db, resp.Data[0].Embedding, topK)
		if err != nil {
			return nil, nil, fmt.Errorf("search: %w", err)
		}
		searchDur := time.Since(searchStart)
		log.Printf("[vault_smart] call #%d — hnsw_search: results=%d  duration=%s",
			n, len(results), searchDur.Round(time.Millisecond))

		// Log each result with similarity and size.
		var responseChars int
		for i, r := range results {
			chars := len(r.Text)
			responseChars += chars
			log.Printf("[vault_smart] call #%d — result[%d]: path=%q  similarity=%.4f  chars=%d  est_tokens=%d",
				n, i, r.Path, r.Similarity, chars, estTokens(chars))
		}

		// Phase 3: marshal response.
		out := searchVaultResult{Results: results}
		sb, err := json.Marshal(out)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal results: %w", err)
		}
		payloadChars := len(sb)
		totalDur := time.Since(start)

		// Summary with comparison against full vault dump.
		vaultEstTok := estTokens(vaultTotalChars)
		responseEstTok := estTokens(payloadChars)
		reduction := 100.0 * (1.0 - float64(responseEstTok)/float64(vaultEstTok))

		log.Printf("[vault_smart] call #%d — payload: chars=%d  est_tokens=%d  total_duration=%s",
			n, payloadChars, responseEstTok, totalDur.Round(time.Millisecond))
		log.Printf("[vault_smart] call #%d — 📊 vault_tokens=%d  response_tokens=%d  reducao=%.1f%%",
			n, vaultEstTok, responseEstTok, reduction)

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(sb)}},
		}, out, nil
	}
}

// estTokens estimates the token count from a character count (~3.5 chars/token).
func estTokens(chars int) int {
	return int(float64(chars) / 3.5)
}
