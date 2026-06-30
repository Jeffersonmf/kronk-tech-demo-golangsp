// Command dumb is the "Cena 1" MCP server: it has a single tool, read_vault,
// that dumps the full contents of every Markdown file in the Obsidian vault.
// No retrieval, no filtering — the client receives everything.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jeff/mcp-vault/internal/vault"
)

const defaultVaultPath = "/Users/jeffersonferreira/Developer/kronk-tech-demo-golangsp/obsidian_vault"

const addr = ":9001"

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

	server := mcp.NewServer(&mcp.Implementation{Name: "vault_dumb", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "read_vault",
		Description: "Returns the full contents of every Markdown file in the Obsidian vault, concatenated with file path headers.",
	}, readVaultHandler(vaultPath))

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)

	log.Printf("[vault_dumb] listening on %s/mcp  vault=%s", addr, vaultPath)

	return http.ListenAndServe(addr, mux)
}

// readVaultParams has no fields: read_vault takes no arguments.
type readVaultParams struct{}

func readVaultHandler(vaultPath string) func(context.Context, *mcp.CallToolRequest, readVaultParams) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, _ readVaultParams) (*mcp.CallToolResult, any, error) {
		n := callCount.Add(1)
		start := time.Now()
		log.Printf("[vault_dumb] call #%d — read_vault iniciada", n)

		loadStart := time.Now()
		docs, err := vault.LoadAll(vaultPath)
		if err != nil {
			return nil, nil, fmt.Errorf("load vault: %w", err)
		}
		loadDur := time.Since(loadStart)

		buildStart := time.Now()
		var sb strings.Builder
		for _, doc := range docs {
			fmt.Fprintf(&sb, "--- FILE: %s ---\n%s\n\n", doc.Path, doc.Content)
		}
		payload := sb.String()
		buildDur := time.Since(buildStart)

		totalChars := len(payload)
		estTok := estTokens(payload)
		totalDur := time.Since(start)

		log.Printf("[vault_dumb] call #%d — docs=%d  chars=%d  est_tokens=%d",
			n, len(docs), totalChars, estTok)
		log.Printf("[vault_dumb] call #%d — load=%s  build=%s  total=%s",
			n, loadDur.Round(time.Millisecond), buildDur.Round(time.Millisecond), totalDur.Round(time.Millisecond))
		log.Printf("[vault_dumb] call #%d — ⚠️  despejando vault inteiro: %d tokens injetados no contexto do agente",
			n, estTok)

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: payload}},
		}, nil, nil
	}
}

// estTokens estimates the token count for a string (~3.5 chars per token).
func estTokens(s string) int {
	return int(float64(len(s)) / 3.5)
}
