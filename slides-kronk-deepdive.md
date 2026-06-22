# Roteiro — Seção Deep Dive: Kronk

Seção técnica aprofundada sobre o Kronk. Entra depois da abertura (`slides-abertura.md`) e antes das demos ao vivo (Cena 1 vs Cena 2, que ficam pro final da talk — ver ordem completa no topo de `slides-abertura.md`). Cada slide tem: título, o que mostrar na tela, e o roteiro (o que falar).

---

## Slide 1 — O que é o Kronk?

**Tela:** Logo/nome do projeto (github.com/ardanlabs/kronk) + uma frase de efeito.

**Roteiro:**
"Kronk é um SDK em Go, da Ardan Labs, que envolve o llama.cpp — ou seja, ele roda modelos de linguagem localmente, sem precisar de Python, sem precisar de internet, sem precisar mandar nada pra nuvem. A grande sacada: ele expõe uma API compatível com a da OpenAI. Isso significa que qualquer ferramenta que já fala com GPT — Cline, Cursor, scripts seus — fala com o Kronk sem mudar uma linha de código além da URL."

**Ponto chave pra reforçar:** Kronk não é "mais um wrapper de chatbot". É infraestrutura — server, pool de modelos, cache, autenticação, métricas. Pensado pra produção, não só demo.

---

## Slide 2 — Arquitetura: onde o Kronk entra

**Tela:** Diagrama simples: [Seu Agente/Cline] → (API OpenAI-compatible) → [Kronk Server] → (llama.cpp backend: CPU / CUDA / Metal / ROCm / Vulkan) → [Modelo GGUF]

**Roteiro:**
"O Kronk fica no meio. De um lado, qualquer cliente que fale o protocolo da OpenAI — chat completions, embeddings, reranking. Do outro, ele escolhe automaticamente o melhor backend de hardware disponível: CPU puro, CUDA se você tem NVIDIA, Metal se for Mac, ROCm pra AMD, Vulkan como fallback genérico. Quem decide isso é uma flag, `--processor`. E o modelo em si é um arquivo `.gguf` — o formato de quantização que saiu do ecossistema llama.cpp."

---

## Slide 3 — Conceito: Quantização (a base de tudo)

**Tela:** Tabela rápida — F16 (sem perda) vs Q8_0 (quase sem perda) vs Q4_K_M (mais compacto, alguma perda).

**Roteiro:**
"Antes de falar de MoE e Hybrid, precisamos alinhar um conceito: quantização. Um modelo de linguagem é, no fundo, bilhões de números de ponto flutuante. Quantização é comprimir esses números pra ocupar menos memória. Q8_0 usa 8 bits por número — praticamente sem perda de qualidade. Q4_K_M usa 4 bits — bem mais compacto, com uma perdinha de precisão que na prática quase não se nota em tarefas de código. Hoje usamos Qwen3-8B em Q8_0: 8 bilhões de parâmetros, ocupando uns 8.7GB em disco."

---

## Slide 4 — Conceito: MoE (Mixture of Experts)

**Tela:** Diagrama: um modelo denso (todas as "fatias" do cérebro acendem pra cada token) vs um MoE (só algumas "fatias/especialistas" acendem por token, um roteador escolhe quais).

**Roteiro:**
"Mixture of Experts é uma arquitetura onde o modelo não é um bloco monolítico — ele é dividido em várias sub-redes, os 'especialistas'. Para cada token gerado, um componente roteador decide quais especialistas ativar, normalmente só uma fração pequena do total. O ganho: você pode ter um modelo gigante em número de parâmetros — gpt-oss-20b, por exemplo, que testamos hoje, tem 20 bilhões de parâmetros — mas o custo computacional por token é parecido com um modelo bem menor, porque só uma fatia dos especialistas é usada a cada vez.

Isso importa pro Kronk porque ele tem suporte nativo a colocar os especialistas em dispositivos diferentes — uma opção de configuração chamada `moe.mode`, que pode ser `auto`, `experts_cpu` (manda todos os especialistas pra CPU e deixa só a parte 'roteadora' na GPU), `experts_gpu`, ou `keep_top_n` (mantém só os N especialistas mais usados na GPU). Isso é uma ferramenta poderosíssima pra rodar modelos MoE gigantes em GPUs pequenas — você não precisa caber o modelo inteiro na VRAM, só a parte que realmente compute por token."

**Nota de bastidor (se quiser contar a história real):** "A gente tem o gpt-oss-20b baixado aqui, mas pra hoje usamos o Qwen3-8B, que é denso — sem MoE. Por quê? [ver Slide 7 — a pegadinha de VRAM]."

---

## Slide 5 — Conceito: Hybrid CPU/GPU (offload de camadas)

**Tela:** Diagrama de um modelo "empilhado" em camadas (36 camadas, no nosso caso), com uma linha cortando: X camadas verdes (GPU) embaixo, resto cinza (CPU) acima.

**Roteiro:**
"Todo modelo de linguagem é uma pilha de camadas — transformer blocks. O Kronk deixa você decidir, camada por camada, onde cada uma roda: GPU ou CPU. Isso se chama offload de camadas, ou modo híbrido.

A configuração é `ngpu-layers`. E aqui tem uma pegadinha que a gente caiu na cara hoje, ao vivo, então vale contar: a convenção do Kronk é meio invertida do que você imaginaria. `0` ou deixar vazio significa 'todas as camadas na GPU'. `-1` significa 'nenhuma camada na GPU, força CPU'. E qualquer número positivo é o número exato de camadas a colocar na GPU. Documentação é sua amiga aqui — sem ler, você acaba invertendo tudo, igual eu fiz."

**Por que isso importa:** "Modelo inteiro na GPU é mais rápido, mas exige que ele caiba todo na VRAM — pesos do modelo, mais o cache de contexto (que cresce junto com a conversa), mais buffers de computação. Se não cabe, ou você usa CPU (mais devagar, mas a memória RAM costuma ser bem maior), ou divide: parte na GPU, parte na CPU. É um dial de ajuste fino entre velocidade e memória disponível."

---

## Slide 6 — Nossos números reais (testados hoje, nesta máquina)

**Tela:** Tabela com os números reais (peguei do benchmarks-apresentacao.md):

| Configuração | Tokens/s | Tempo até 1º token | VRAM usada |
|---|---|---|---|
| CPU only (prompt pequeno) | 6.3 tok/s | 520 ms | — |
| GPU 100% — 36/36 camadas (contexto pequeno) | **50.8 tok/s** | **54 ms** | 10.3GB / 12.28GB |
| GPU 100% — contexto grande (16k tokens) | ❌ falha (CUDA out of memory) | — | — |
| Híbrido 18/36 camadas (prompt real, ~12.6K tokens) | 5.2 tok/s | 12.4 s | 7.87GB / 12.28GB |
| Híbrido 28/36 camadas (mesmo prompt) | 9.5 tok/s | 7.3 s | 10.15GB / 12.28GB |

**Roteiro:**
"Esses números não são de benchmark de paper, são de hoje, nesta máquina, RTX 4070 Ti de 12GB. Olha a primeira linha contra a segunda: GPU 100% é 8 vezes mais rápido que CPU pura pra gerar texto, e quase 10 vezes mais rápido pra começar a responder. Isso é o tipo de ganho que faz local-first ser viável de verdade.

Mas — terceira linha. Quando o prompt cresce, porque o agente carrega ferramentas MCP, histórico de conversa, etc., GPU 100% simplesmente quebra. Sem VRAM suficiente. E aí entra o híbrido: a gente foi testando quantas camadas colocar na GPU até achar um equilíbrio entre velocidade e margem de segurança de memória — 18 camadas deixou bastante folga, 28 camadas foi mais rápido mas quase sem margem."

**Gancho pra próximo slide:** "E essa decisão — quanto offload, quanto contexto — não é teórica. A gente bateu de cara com isso ao vivo preparando essa talk."

---

## Slide 7 — A pegadinha que vivemos hoje (war story)

**Tela:** Print/trecho do log de erro real: `acquire: reserve: request[Qwen3-8B-Q8_0] needs vram=10.3GB but largest device budget is 9.7GB` e depois `unable to init context: failed to initialize model`.

**Roteiro:**
"Conto rapidinho o que aconteceu preparando essa demo, porque é uma aula prática melhor que qualquer slide teórico. Configuramos GPU 100%, contexto generoso pro agente conseguir conversar — e a coisa simplesmente não inicializava. O Kronk tem um sistema de orçamento de memória, calcula antes de carregar se vai caber. Recusou. Aumentamos o orçamento à força — passou da estimativa, mas na hora de alocar de verdade, o CUDA devolveu erro de memória. A trava de segurança do Kronk, na verdade, nos salvou de travar o processo de forma feia.

A solução não foi 'forçar mais GPU'. Foi entender que offload parcial (híbrido) é a ferramenta certa pra esse problema — sacrificar uma fatia de velocidade por uma fatia de memória livre, de forma controlada."

---

## Slide 8 — Conceito: "Thinking mode" e o custo escondido do raciocínio

**Tela:** Comparação visual: uma resposta direta ("ok") vs um bloco gigante de "pensamento" antes da resposta.

**Roteiro:**
"Último conceito, e esse foi uma surpresa real durante a preparação. Modelos como o Qwen3 têm um modo de 'raciocínio' — antes de agir, eles geram um monólogo interno, tipo um rascunho de pensamento, que pode ter centenas de tokens. Isso melhora a qualidade da decisão — vimos na prática o modelo acertar a lógica certa de correção de bug só depois de 'pensar' bastante.

O problema: cada token de pensamento custa tempo de geração. Em CPU, isso pode ser a diferença entre uma resposta em segundos e uma resposta em minutos. O Kronk deixa você desligar isso via configuração (`enable_thinking`). Testamos os dois lados: desligado, o modelo respondia rápido, mas cometia erros bobos de estrutura de código — esquecia de importar um pacote, por exemplo. Ligado, ele acertava a lógica completa, mas demorava muito mais.

É um tradeoff real que qualquer um rodando IA local vai enfrentar: velocidade contra profundidade de raciocínio."

---

## Slide 9 — Resumo / Takeaways

**Tela:** Bullet points grandes, tipo "cartões".

**Roteiro (ler cada um com uma pausa):**
- "Local-first com Kronk é rápido de verdade — quando cabe na GPU."
- "MoE e offload híbrido são as ferramentas pra rodar modelos maiores que sua VRAM."
- "A convenção de camadas do Kronk é invertida — leia a doc antes de configurar."
- "Modo de raciocínio é um dial, não uma mágica: mais qualidade custa mais tempo."
- "Tudo isso é configurável em um arquivo YAML, sem precisar recompilar nada."

**Fechamento:** "E é exatamente essa flexibilidade — CPU, GPU, híbrido, MoE, raciocínio ligado ou desligado — que torna o Kronk uma peça de infraestrutura séria pra rodar IA local em produção, não só em demo de palco."

---

---

# Roteiro — Seção: Por dentro do código dos 2 MCPs

Essa seção entra bem antes do deep dive do Kronk (ou logo depois do "momento WOW") — é o "abre a caixa" do que realmente faz a diferença de tokens/custo entre as duas cenas. Os três arquivos-fonte são: `mcp_vault/internal/vault/reader.go` (compartilhado pelos dois), `mcp_vault/cmd/dumb/main.go` (Cena 1), `mcp_vault/cmd/smart/main.go` + `mcp_vault/internal/store/duckdb.go` (Cena 2).

---

## Slide 10 — A base compartilhada: lendo o vault

**Tela:** trecho de código de `internal/vault/reader.go` — a função `LoadAll`.

```go
func LoadAll(root string) ([]Document, error) {
    var docs []Document
    filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
        if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
            return nil
        }
        content, _ := os.ReadFile(path)
        rel, _ := filepath.Rel(root, path)
        docs = append(docs, Document{Path: rel, Content: string(content)})
        return nil
    })
    return docs, nil
}
```

**Roteiro:**
"Antes de comparar os dois MCPs, vale mostrar que os dois partem exatamente do mesmo lugar: essa função aqui, `LoadAll`. Ela só anda pela pasta do vault do Obsidian, pega todo arquivo `.md`, lê o conteúdo bruto e devolve uma lista de documentos — caminho relativo mais texto. Simples, sem mágica, sem dependência externa. Os 26 documentos do nosso vault de teste passam por essa mesma função nos dois servidores. A diferença entre as duas cenas não está em COMO os dados são lidos — está em O QUE cada servidor faz com eles depois."

---

## Slide 11 — Cena 1: o servidor "burro" (`vault_dumb`)

**Tela:** código completo do handler de `cmd/dumb/main.go` (é bem curto, cabe na tela):

```go
mcp.AddTool(server, &mcp.Tool{
    Name:        "read_vault",
    Description: "Returns the full contents of every Markdown file...",
}, readVaultHandler(vaultPath))

func readVaultHandler(vaultPath string) ... {
    return func(...) (*mcp.CallToolResult, any, error) {
        docs, _ := vault.LoadAll(vaultPath)
        var sb strings.Builder
        for _, doc := range docs {
            fmt.Fprintf(&sb, "--- FILE: %s ---\n%s\n\n", doc.Path, doc.Content)
        }
        return &mcp.CallToolResult{
            Content: []mcp.Content{&mcp.TextContent{Text: sb.String()}},
        }, nil, nil
    }
}
```

**Roteiro:**
"Esse é o servidor inteiro da Cena 1 — uma única ferramenta MCP, `read_vault`, que não recebe nenhum parâmetro. Quando o agente chama essa ferramenta, ela faz uma coisa só: lê TODOS os documentos do vault e concatena tudo num texto gigante, com um cabeçalho `--- FILE: caminho ---` antes de cada um, e devolve isso inteiro como resposta.

Não tem filtro. Não tem busca. Não tem relevância. O agente pediu uma informação específica — 'me ajuda com a TASK-001' — e essa ferramenta devolve as 26 tasks completas, junto com o README do vault, tudo. É proposital: representa o jeito que muita gente integra LLM com dados hoje — 'enfia tudo no contexto e deixa o modelo achar o que precisa'. Funciona, mas é isso que gera aquele consumo de ~47 mil tokens que vimos no contador do Cline."

**Linha de efeito:** "Olha o tamanho do código: menos de 20 linhas relevantes. É o MCP mais simples que existe. E é exatamente essa simplicidade ingênua que custa caro."

---

## Slide 12 — Cena 2, parte 1: montando o índice vetorial (`vault_smart`, setup)

**Tela:** trecho de `cmd/smart/main.go` (a função `run`, parte de inicialização) lado a lado com um diagrama: Documentos → [Modelo de Embedding] → Vetores → [DuckDB + índice HNSW].

```go
krnEmbed, _ := kronk.New(model.WithModelFiles([]string{embedModelPath}))
docs, _ := vault.LoadAll(vaultPath)
dims, _ := embeddingDims(krnEmbed)
db, _ := store.LoadVault(ctx, krnEmbed, docs, dims)
```

**Roteiro:**
"A Cena 2 começa de forma bem diferente. Antes mesmo de aceitar qualquer pergunta do agente, ela faz um trabalho de preparação: carrega um segundo modelo — não o modelo de chat, um modelo *de embedding*, o `embeddinggemma`, bem menor e especializado só em transformar texto em vetores numéricos. Cada um dos 26 documentos do vault passa por esse modelo e se transforma num vetor de 768 números — pensem nisso como as 'coordenadas de significado' daquele texto num espaço matemático.

Esses vetores vão pra dentro de um banco DuckDB, criado na memória, na hora, com uma extensão chamada VSS — Vector Similarity Search — e um índice HNSW por cima, que é uma estrutura de dados feita pra busca de vizinho-mais-próximo em alta dimensão, rápida mesmo com muitos vetores. Esse trabalho de indexação acontece **uma vez**, na inicialização do servidor — não a cada pergunta."

---

## Slide 13 — Cena 2, parte 2: a busca de verdade (`search_vault`)

**Tela:** trecho de `internal/store/duckdb.go`, a função `Search`, com destaque na query SQL.

```sql
SELECT path, text, array_cosine_similarity(embedding, ?::FLOAT[768]) AS similarity
FROM items
ORDER BY similarity DESC
LIMIT 3;
```

**Roteiro:**
"Aqui está o coração da Cena 2. Quando o agente chama a ferramenta `search_vault` com uma pergunta — por exemplo, 'qual a race condition no inventário?' — o servidor pega essa pergunta, transforma ela no mesmo tipo de vetor de 768 números usando o mesmo modelo de embedding, e manda essa query SQL pro DuckDB: 'me dá os 3 documentos cujo vetor tem a maior similaridade de cosseno com esse vetor de pergunta'.

Similaridade de cosseno, resumindo, é uma forma de medir o quão 'parecidos em significado' dois vetores são — não é busca de palavra-chave, é busca de *sentido*. É por isso que perguntar sobre 'dinheiro e pagamentos', como testamos antes, consegue achar a task de cupom de desconto mesmo sem a palavra 'dinheiro' aparecer lá.

O resultado: em vez de 26 documentos completos, o agente recebe só os 3 mais relevantes. Isso é o que derruba o consumo de ~47 mil tokens pra perto de 4 mil."

---

## Slide 14 — Lado a lado: mesma pergunta, dois caminhos

**Tela:** tabela comparativa.

| | `vault_dumb` (Cena 1) | `vault_smart` (Cena 2) |
|---|---|---|
| Ferramenta MCP | `read_vault` (sem parâmetros) | `search_vault(query, top_k)` |
| O que faz | Concatena tudo | Embedding + busca por similaridade |
| Trabalho de preparação | Nenhum | Indexa o vault uma vez na subida |
| Infraestrutura extra | Nenhuma | Modelo de embedding + DuckDB + HNSW |
| Tokens pro modelo (nosso teste) | ~47.000 | ~4.000 |
| Linhas de código (aprox.) | ~20 | ~80 (incluindo o índice) |

**Roteiro:**
"Resumindo o trade-off: a Cena 1 é 4x mais simples de código, mas custa 10x mais tokens. A Cena 2 exige montar uma pecinha extra de infraestrutura — um modelo de embedding, um banco vetorial — mas paga isso de volta multiplicado em economia de contexto. E o detalhe mais bonito: como tudo isso roda local, via Kronk, esse trabalho extra de indexação não custa um centavo, só tempo de CPU/GPU nossa, que já vimos nos números do deep dive técnico. Esse é o argumento de fundo da talk: orquestrar bem onde e como você busca contexto vale mais do que confiar que o modelo vai 'dar um jeito' com tudo jogado na cara dele."

---

## Notas de produção (não falar em voz alta)

- Os números do Slide 6 vêm de `benchmarks-apresentacao.md` neste mesmo repo — conferir se não mudaram caso re-teste antes da talk.
- Slide 7 é ótimo gancho de humor/autenticidade — sugiro tom leve, "rir da própria dor", já que é uma talk técnica e isso humaniza.
- Se houver tempo, o Slide 4 (MoE) pode ganhar uma demonstração rápida ao vivo trocando `moe.mode` no `model_config.yaml` do gpt-oss-20b — mas só se already tiver sido testado antes, não arriscar ao vivo sem ensaio.
