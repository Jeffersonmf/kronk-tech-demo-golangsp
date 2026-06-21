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
	"time"

	"github.com/ardanlabs/kronk/sdk/kronk"
	"github.com/ardanlabs/kronk/sdk/kronk/model"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jeff/mcp-vault/internal/store"
	"github.com/jeff/mcp-vault/internal/vault"
)

const defaultVaultPath = "/home/jmarchetti/Developer/kronk-tech-demo-golangsp/obsidian_vault"

// embedModelPath points to the embeddinggemma GGUF already downloaded via
// `kronk model pull --local`. Dimension measured at runtime via Embeddings().
const embedModelPath = "/home/jmarchetti/.kronk/models/ggml-org/embeddinggemma-300m-qat-q8_0-GGUF/embeddinggemma-300m-qat-Q8_0.gguf"

const addr = ":9002"

const defaultTopK = 3

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

	fmt.Println("Loading embedding model...")

	if err := kronk.Init(); err != nil {
		return fmt.Errorf("kronk init: %w", err)
	}

	krnEmbed, err := kronk.New(model.WithModelFiles([]string{embedModelPath}))
	if err != nil {
		return fmt.Errorf("kronk new: %w", err)
	}
	defer krnEmbed.Unload(context.Background())

	fmt.Println("Reading vault and building vector index...")

	docs, err := vault.LoadAll(vaultPath)
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
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

	fmt.Printf("Indexed %d documents (dims=%d)\n", len(docs), dims)

	server := mcp.NewServer(&mcp.Implementation{Name: "vault_smart", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_vault",
		Description: "Searches the Obsidian vault and returns only the top-K Markdown documents most relevant to the query, ranked by semantic similarity.",
	}, searchVaultHandler(krnEmbed, db))

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)

	fmt.Printf("vault_smart MCP server listening on %s/mcp (vault: %s)\n", addr, vaultPath)

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

func searchVaultHandler(krnEmbed *kronk.Kronk, db *sql.DB) func(context.Context, *mcp.CallToolRequest, searchVaultParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params searchVaultParams) (*mcp.CallToolResult, any, error) {
		topK := params.TopK
		if topK <= 0 {
			topK = defaultTopK
		}

		queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		resp, err := krnEmbed.Embeddings(queryCtx, model.D{"input": params.Query})
		if err != nil {
			return nil, nil, fmt.Errorf("embed query: %w", err)
		}
		if len(resp.Data) == 0 {
			return nil, nil, fmt.Errorf("empty embedding for query")
		}

		results, err := store.Search(db, resp.Data[0].Embedding, topK)
		if err != nil {
			return nil, nil, fmt.Errorf("search: %w", err)
		}

		out := searchVaultResult{Results: results}

		sb, err := json.Marshal(out)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal results: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(sb)}},
		}, out, nil
	}
}
