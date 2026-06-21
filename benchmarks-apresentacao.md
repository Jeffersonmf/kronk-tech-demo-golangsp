# Benchmarks Kronk local — Qwen3-8B-Q8_0 (RTX 4070 Ti, 12GB)

Números medidos ao vivo durante a preparação da talk (2026-06-21), pra usar nos slides de comparação CPU vs GPU.

## CPU only (`--processor cpu`)

| Métrica | Valor |
|---|---|
| Velocidade de geração | **6.3 tokens/s** |
| Tempo até o 1º token (prompt pequeno, ~18 tokens) | 520 ms |

## GPU full offload (36/36 camadas), contexto pequeno

| Métrica | Valor |
|---|---|
| Velocidade de geração | **50.8 tokens/s** (~8x vs CPU) |
| Tempo até o 1º token (prompt pequeno, ~18 tokens) | **54 ms** (~10x vs CPU) |
| VRAM usada | 10.3GB / 12.28GB |

> Esse modo (GPU 100%) só funciona com contexto pequeno (testado com 4096). Com contexto maior (16384, necessário pro Cline real com os MCPs ativos) a inicialização **falha** — a placa não tem VRAM suficiente pra pesos completos (8.7GB) + cache de contexto grande.

## Configuração final escolhida: GPU parcial (18/36 camadas) + contexto 16384

Testado com prompt real grande (~12.6K tokens, simulando o Cline com 2 MCP servers ativos):

| Métrica | Valor |
|---|---|
| Velocidade de geração | 5.2 tokens/s |
| Tempo até o 1º token (prompt grande, ~12.6K tokens) | 12.4 s |
| VRAM usada | 7.87GB / 12.28GB |
| Margem livre | 3.88GB |

### Ponto de comparação testado (não escolhido): 28/36 camadas

| Métrica | Valor |
|---|---|
| Velocidade de geração | 9.5 tokens/s |
| Tempo até o 1º token (mesmo prompt grande) | 7.3 s |
| VRAM usada | 10.15GB / 12.28GB |
| Margem livre | 1.65GB (descartado — margem parecida com a falha anterior) |

## Decisão final

Ficou em **18 de 36 camadas na GPU**, contexto 16384, `cache-type-k/v: q8_0`, `--budget-percent 92`. Prioriza margem de VRAM (3.88GB livres) sobre velocidade máxima, já que a gravação tolera retake mas não tolera ficar reconfigurando no meio da sessão.

Config em `~/.kronk/models/model_config.yaml`:

```yaml
Qwen3-8B-Q8_0:
  context-window: 16384
  cache-type-k: q8_0
  cache-type-v: q8_0
  ngpu-layers: 18
```

Servidor: `kronk server start --processor cuda --budget-percent 92`
