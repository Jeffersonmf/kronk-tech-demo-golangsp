# kronk-tech-demo-golangsp

Repositório de demonstração da talk **"Golang + AI + MCP"**, parte 2, apresentada na comunidade **GolangSP**. Mostra duas formas de construir um servidor MCP (Model Context Protocol) que busca contexto num backlog de notas (Obsidian) para um agente de IA — uma ingênua (carrega tudo) e uma com busca semântica de verdade (recupera só o relevante) — rodando localmente, sem nuvem, via [Kronk](https://github.com/ardanlabs/kronk).

## Estrutura

```
fake_project/      # app Go alvo, com um bug de concorrência proposital
mcp_vault/          # os dois servidores MCP da demo (Cena 1 "dumb" vs Cena 2 "smart")
obsidian_vault/     # backlog de tasks em Markdown, consumido pelos MCPs
```

### `fake_project/`

Serviço de inventário de e-commerce (`internal/inventory/service.go`) com uma race condition real em `DeductStock` — checagem de estoque sem lock, causando estoque negativo sob concorrência. É o código que o agente de IA corrige ao vivo na demo, guiado pela `TASK-001` do vault.

```sh
cd fake_project && go test ./internal/inventory/... -race -v
```

### `mcp_vault/`

Os dois servidores MCP que sustentam a comparação central da talk:

- **`cmd/dumb`** (Cena 1, porta `:9001`) — ferramenta `read_vault`, sem parâmetros: devolve o vault inteiro concatenado. Simples, mas caro em tokens.
- **`cmd/smart`** (Cena 2, porta `:9002`) — ferramenta `search_vault(query, top_k)`: embute cada documento com o Kronk e indexa num DuckDB com índice HNSW (similaridade de cosseno), devolvendo só os documentos mais relevantes pra pergunta.

```sh
cd mcp_vault && ./run.sh
```

### `obsidian_vault/`

25 tasks em Markdown (`tasks/TASK-001.md` a `TASK-025.md`) simulando o backlog de um time. `TASK-001` é a única conectada a código real (a race condition); as demais são ruído de backlog proposital — ver `obsidian_vault/tasks/README.md` para o porquê desse tamanho.

## Rodando a demo localmente

Pré-requisitos: [`kronk`](https://github.com/ardanlabs/kronk) instalado e os modelos baixados (`kronk model pull`).

```sh
# 1. Suba o servidor de inferência local
kronk server start --processor cpu   # ou --processor cuda, se tiver GPU NVIDIA

# 2. Suba os dois servidores MCP
cd mcp_vault && ./run.sh

# 3. Configure seu cliente MCP (Cline, etc.) apontando para:
#    - Provider OpenAI-compatible: http://localhost:11435/v1
#    - MCP servers: http://localhost:9001/mcp e http://localhost:9002/mcp
```

Detalhes de configuração do modelo (offload de GPU, tamanho de contexto, cache) ficam em `~/.kronk/models/model_config.yaml`.

## Documentação adicional

- [`benchmarks-apresentacao.md`](benchmarks-apresentacao.md) — números reais de desempenho (CPU vs GPU vs híbrido) medidos nesta máquina.
- [`slides-abertura.md`](slides-abertura.md), [`slides-kronk-deepdive.md`](slides-kronk-deepdive.md), [`slides-demos-finais.md`](slides-demos-finais.md) — roteiro completo da talk, slide a slide.
- [`CLAUDE.md`](CLAUDE.md) — contexto do projeto para assistentes de IA trabalhando neste repositório.
