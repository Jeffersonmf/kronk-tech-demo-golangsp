# GEMINI Project Context: GolangSP Presentation - Part 2

## Project Overview
This project is a technical presentation for **GolangSP** titled "Golang + AI + Kronk". The core objective is to demonstrate the evolution from traditional Model Context Protocol (MCP) implementations to a unified, high-performance Go-centric approach using the **Kronk** framework and the **Cline** AI agent.

### Key Technologies
- **Go (Golang):** Primary language for building the unified AI hub and MCP servers.
- **Kronk:** Ardan Labs framework for local LLM inference and streamlined MCP integration.
- **Cline:** Autonomous AI agent acting as the MCP client (VS Code extension).
- **VS Code:** Editor and presentation panel — Cline's native environment, shows token count and USD cost in real time in the UI.
- **Zed:** Descartado como palco principal. Cline não tem extensão para Zed (suporte oficial: VS Code, JetBrains, Cursor, Windsurf, CLI). Pode ser mencionado na palestra como "também funciona com Zed AI + Kronk", mas não é o palco da demo.
- **Obsidian:** Used as the "Source of Truth" for business rules, Acceptance Criteria (ACs), and Technical Specifications (TAs).

### Architecture & Comparison Strategy
- **Traditional MCP (`traditional-mcp/`):** Vanilla Go implementation using `stdio`. Demonstrates boilerplate complexity and management overhead.
- **Kronk + Cline (`kronk-demo/`):** SDK-based implementation. Demonstrates DX improvements, local-first inference, and token/cost efficiency.
- **The "Fake Project" (`fake-project/`):** [CREATED] Target Go application with a critical race condition bug in the inventory service.
- **The "Obsidian Vault" (`obsidian_vault/`):** [CREATED] Business requirements (TASK-001) for fixing the inventory race condition.

## Building and Running
*Note: This project is in active development for the GolangSP demo.*

### Prerequisites
- Go 1.22+
- Zed Editor & Obsidian
- Local LLM (GGUF) compatible with Kronk.

### Commands
- `go run traditional-mcp/main.go`: Start the raw MCP server.
- `go run kronk-demo/main.go`: Start the Kronk hub and MCP server.

## Stack Final (decisão consolidada)

| Componente | Escolha | Motivo |
|---|---|---|
| Editor/palco | VS Code | Cline roda nativamente |
| Agente | Cline (VS Code ext.) | Token counter + custo USD nativo na UI |
| Inference local | Kronk | llama.cpp em Go |
| Modelo | Qwen3 14B Q4_K_M | Velocidade aceitável no i7 CPU |
| Requisitos | Obsidian | Issues, ACs, TAs em markdown |
| Código alvo | fake-project (Go) | Projeto com bugs/features que o agente trabalha |
| Hardware | Linux, i7, 96GB RAM, sem GPU | Inference CPU via llama.cpp |

## Development Conventions
- **Offline-First:** All demos must be runnable without internet (using local models).
- **Scene-Based:** Use VS Code tasks or Markdown files to trigger presentation stages.
- **Pragmatic AI:** Focus on real-world Go problems (refactoring, testing, profiling).
- **WOW Moment:** Mostrar o contador de tokens e custo USD do Cline lado a lado nas duas cenas — Cena 1 (~47K tokens, ~$1.43) vs Cena 2 (~4K tokens, $0.00).

## Persistent Memory Location
- Private notes and transient state are stored in: `/home/jmarchetti/.gemini/tmp/parte2/memory/MEMORY.md`
