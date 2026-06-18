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

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jeff/mcp-vault/internal/vault"
)

const defaultVaultPath = "/home/jmarchetti/Developer/golangsp/parte2/obsidian_vault"

const addr = ":9001"

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

	fmt.Printf("vault_dumb MCP server listening on %s/mcp (vault: %s)\n", addr, vaultPath)

	return http.ListenAndServe(addr, mux)
}

// readVaultParams has no fields: read_vault takes no arguments.
type readVaultParams struct{}

func readVaultHandler(vaultPath string) func(context.Context, *mcp.CallToolRequest, readVaultParams) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, _ readVaultParams) (*mcp.CallToolResult, any, error) {
		docs, err := vault.LoadAll(vaultPath)
		if err != nil {
			return nil, nil, fmt.Errorf("load vault: %w", err)
		}

		var sb strings.Builder
		for _, doc := range docs {
			fmt.Fprintf(&sb, "--- FILE: %s ---\n%s\n\n", doc.Path, doc.Content)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: sb.String()}},
		}, nil, nil
	}
}
