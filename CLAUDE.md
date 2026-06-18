# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project context

This is the "parte 2" demo repository for a GolangSP (São Paulo Go community) talk on **Golang + AI + MCP**. Part 1 covered Go as a productivity/tooling language; this talk extends that into the AI agent space.

The central narrative is a **before/after comparison** of MCP approaches, live-coded in VS Code with the Cline extension:

- **Scene 1 — Traditional MCP**: raw Go MCP server, showing boilerplate/overhead (~47K tokens, ~$1.43 cloud cost)
- **Scene 2 — Kronk + Cline**: same capability via the Kronk framework with local inference ($0.00, local CPU)

The **WOW moment** is showing Cline's token counter + USD cost side by side between the two scenes.

## Actual structure

```
parte2/
  fake_project/internal/inventory/   # target Go app with intentional race condition bug
  gemini_version/                    # Kronk hub demo (Go module, stub implementation)
  obsidian_vault/tasks/              # business requirements in Markdown (Obsidian)
  GEMINI.md                          # fuller context file (see also)
```

> `traditional-mcp/` and `kronk-demo/` directories (referenced in earlier planning) are not yet created.

## Commands

All commands run from the relevant subdirectory (each demo is a self-contained Go module).

```sh
# Run the Kronk hub demo
cd gemini_version && go run main.go

# Run the fake project tests (detect the race condition)
cd fake_project && go test ./internal/inventory/... -race -v

# Build check
cd gemini_version && go build ./...
```

## Demo scenario

The "fake project" is an e-commerce inventory service (`fake_project/internal/inventory/service.go`) with a deliberate race condition in `DeductStock` — no mutex around the stock check + subtraction. This causes negative stock under concurrent load.

**TASK-001** (`obsidian_vault/tasks/TASK-001.md`) is the Obsidian "ticket" the AI agent reads to understand what to fix. The agent (Cline) reads the task, locates the code, applies a mutex fix, and verifies with `go test -race`.

## Finalized stack

| Component | Choice | Notes |
|---|---|---|
| Editor/stage | VS Code | Cline runs natively here |
| Agent | Cline (VS Code ext.) | Token counter + USD cost visible in UI |
| Local inference | Kronk (Ardan Labs) | llama.cpp in Go |
| Model | Qwen3 14B Q4_K_M | Acceptable speed on CPU |
| Requirements | Obsidian | ACs and TAs in Markdown |
| Target code | `fake_project/` | Go service with race condition |
| Hardware | Linux, i7, 96 GB RAM, no GPU | CPU-only inference via llama.cpp |

## Talk constraints

- All demos must run **offline** — conference Wi-Fi is unreliable. No cloud API calls during scenes.
- Keep startup times fast — demos switch live during the talk.
- Target audience: Go developers, not necessarily familiar with MCP or AI tooling.
