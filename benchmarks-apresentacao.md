# Benchmarks Kronk local — gpt-oss-20b-Q8_0 (MacBook M4 Pro, Metal, 48GB unificada)

Números medidos ao vivo durante a preparação da talk (2026-06-22), pra usar nos slides de comparação CPU vs Metal.

> Nota de bastidor: a talk foi originalmente preparada num Linux com RTX 4070 Ti (12GB de VRAM dedicada) usando o `Qwen3-8B-Q8_0` — ver histórico do repo se quiser a história completa de VRAM-out-of-memory daquela máquina. A gravação real acabou migrando pro MacBook M4 Pro, e isso muda o personagem da história: memória unificada (48GB compartilhada entre CPU e GPU) não tem o teto rígido de VRAM dedicada, então o `gpt-oss-20b-Q8_0` (MoE, 20B parâmetros, ~12GB em disco) cabe inteiro e sobra bastante folga — não precisamos mais do offload híbrido como muleta de memória. O gargalo muda de lugar: virou compute do kernel Metal, não capacidade de memória.

## CPU only (`--processor cpu`)

| Métrica | Valor |
|---|---|
| Velocidade de geração | ~6-8 tokens/s (estimado, mesma ordem de grandeza do Linux original) |

## Metal full offload (24/24 camadas, `ngpu-layers: 0` = todas)

Testado com três tamanhos de prompt, simulando desde uma pergunta simples até um prompt real do Cline com MCPs ativos:

| Cenário | Prompt tokens | Tempo até 1º token | Velocidade de geração |
|---|---|---|---|
| Curto | 81 | ~205-430 ms | ~70 tokens/s |
| Médio (estilo Cline real) | 2.034 | ~2.2 s | ~65-68 tokens/s |
| Stress (simulando Cline com MCPs grandes) | 9.084 | ~11.2-12.4 s | ~63 tokens/s |

> Comparado ao CPU puro, full offload no Metal é ~9-10x mais rápido pra gerar texto — sem nenhum drama de "não cabe na memória", porque é memória unificada.

## Configuração final escolhida: full offload, sem híbrido

Diferente da máquina original (Linux + RTX 4070 Ti), aqui não existe a pegadinha de VRAM
que forçava offload parcial. `ngpu-layers: 0` (todas as camadas na GPU) funciona até no
prompt de stress de ~9k tokens, sem falhar por falta de memória.

Config em `~/.kronk/models/model_config.yaml`:

```yaml
gpt-oss-20b-Q8_0:
  ngpu-layers: 0          # 0 = todas as camadas na GPU (convenção invertida do Kronk)
  flash-attention: enabled # forçado, não "auto" — gpt-oss mistura atenção full/sliding-window por camada
  offload-kqv: true
  op-offload: true
  nbatch: 4096
  nubatch: 2048
  context-window: 32768
```

Servidor: `KRONK_PROCESSOR=metal kronk server start`

### O que tentamos e não ajudou tanto: aumentar batch size

Testamos `nbatch`/`nubatch` em 512/1024/2048/4096 esperando um ganho grande na velocidade
de prefill (tempo até o 1º token) do prompt de stress. O ganho real foi de só **~6%**
(12.36s → 11.58s). Conclusão: pra esse modelo MoE, no Metal, o prefill é limitado pelo
*compute* do kernel — o roteamento de especialistas — não pelo tamanho do batch. É um bom
contraponto pra história de "só aumentar batch resolve": medimos antes de assumir.
