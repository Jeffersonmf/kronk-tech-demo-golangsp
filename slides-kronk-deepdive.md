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

## Slide 3 — Por dentro do Kronk Server: o que sobe quando você dá `kronk server start`

**Tela:** Trecho real do log de startup (`/tmp/kronk.log` de hoje), com 4 linhas em destaque:
```
"status":"api router started","host":"0.0.0.0:11435"
"status":"debug v1 router started","host":"0.0.0.0:11445"
"status":"local mcp server started","host":"127.0.0.1:9000"
"status":"start auth server"
```

**Roteiro:**
"O Slide 2 mostrou o Kronk como uma caixa preta no meio do diagrama. Vamos abrir essa caixa rapidinho, com log real de hoje, não slide teórico. Quando você dá `kronk server start`, sobem pelo menos quatro coisas: o **router da API** (porta 11435, é o `/v1/chat/completions` compatível com OpenAI que o Cline fala), um **router de debug** (11445, métricas e introspecção), um **serviço de auth** (JWT/API keys — hoje desligado porque é uso local, mas existe pra produção), e — detalhe curioso — um **servidor MCP nativo do próprio Kronk**, na porta 9000. Não é nenhum dos nossos dois MCPs de demo; é o Kronk se expondo via MCP por conta própria, pra outros agentes consultarem o estado dele direto.

O ponto de fundo: isso não é um wrapper fininho em volta do llama.cpp. É um processo de servidor com superfície de operação real — roteamento, auth, observabilidade — parecido com qualquer outro serviço backend que vocês já rodam em Go."

---

## Slide 4 — Pool de modelos: por que carregar um modelo não é só "abrir o arquivo"

**Tela:** Linha real do log de hoje:
```
"status":"resman-init","budget-percent":80,"ram-budget":"38.7GB",
"max-models-in-pool":10
```

**Roteiro:**
"O Kronk não trata modelo como 'um arquivo que você abre uma vez'. Ele tem um **pool de modelos** — igual um pool de conexões de banco, só que pra modelos GGUF carregados na memória. Três coisas que esse pool faz por você:

Primeiro, **orçamento de memória**: por padrão, o Kronk só usa 80% da RAM/VRAM disponível pra modelos — aqui, hoje, isso virou um orçamento de 38.7GB. Se um modelo não cabe no orçamento, o Kronk **recusa carregar antes de tentar** em vez de só travar feio na hora de alocar — foi exatamente esse mecanismo de segurança que pegou a configuração original dessa talk, preparada num Linux com GPU dedicada de 12GB, antes da gente migrar pra esse Mac com memória unificada (mais sobre isso no Slide 7, quando chegarmos em MoE).

Segundo, **múltiplos modelos ao mesmo tempo**: o pool aceita até 10 modelos residentes simultaneamente. Não é teórico — agora mesmo, nesta máquina, tem dois modelos carregados ao mesmo tempo: o `gpt-oss-20b` pra chat, e o `embeddinggemma` pra busca semântica do `vault_smart`. Cada um com seu próprio orçamento de memória, dividindo o mesmo processo Kronk.

Terceiro, **TTL de ociosidade**: modelo sem uso por um tempo configurável é descarregado automaticamente, liberando memória pra outra coisa — sem você precisar gerenciar isso manualmente."

---

## Slide 5 — Cache incremental: por que o Cline não fica mais lento a cada turno

**Tela:** Trecho do `model_config.yaml` com a opção comentada:
```yaml
incremental-cache: false   # Incremental message caching for agentic
                           # workflows (Cline, OpenCode)
```
E o valor real que o servidor usa hoje, confirmado no log: `IncrementalCache[true]`.

**Roteiro:**
"Aqui está um detalhe que conecta direto com a demo de hoje. Uma conversa com o Cline não é uma pergunta isolada — é uma sequência de turnos, e cada turno novo reenvia **todo o histórico anterior** mais a mensagem nova. Sem cache, o modelo teria que reprocessar a conversa inteira do zero em todo turno — cada vez mais lento conforme a conversa cresce, porque o prefixo só aumenta.

O **cache incremental de mensagens (IMC)** do Kronk guarda o estado interno (KV-cache) da conversa depois de cada turno. No turno seguinte, ele reaproveita esse estado e só processa de fato os tokens **novos** — a pergunta atual, não a conversa inteira de novo. Isso já vem **ligado por padrão** pensando exatamente em clientes agênticos como Cline e OpenCode, que reenviam histórico crescente a cada chamada.

É por isso que os números de TTFT que vamos ver no próximo slide — tipo ~2.2s pra um prompt de 2 mil tokens — tendem a ficar ainda melhores turno a turno numa conversa real com o Cline, não pioram."

---

## Slide 6 — Conceito: Quantização (a base de tudo)

**Tela:** Tabela rápida — F16 (sem perda) vs Q8_0 (quase sem perda) vs Q4_K_M (mais compacto, alguma perda).

**Roteiro:**
"Antes de falar de MoE e Hybrid, precisamos alinhar um conceito: quantização. Um modelo de linguagem é, no fundo, bilhões de números de ponto flutuante. Quantização é comprimir esses números pra ocupar menos memória. Q8_0 usa 8 bits por número — praticamente sem perda de qualidade. Q4_K_M usa 4 bits — bem mais compacto, com uma perdinha de precisão que na prática quase não se nota em tarefas de código. Hoje usamos o gpt-oss-20b em Q8_0: 20 bilhões de parâmetros, ocupando uns 12GB em disco."

---

## Slide 7 — Conceito: MoE (Mixture of Experts)

**Tela:** Diagrama: um modelo denso (todas as "fatias" do cérebro acendem pra cada token) vs um MoE (só algumas "fatias/especialistas" acendem por token, um roteador escolhe quais).

**Roteiro:**
"Mixture of Experts é uma arquitetura onde o modelo não é um bloco monolítico — ele é dividido em várias sub-redes, os 'especialistas'. Para cada token gerado, um componente roteador decide quais especialistas ativar, normalmente só uma fração pequena do total. O ganho: você pode ter um modelo gigante em número de parâmetros — o gpt-oss-20b que estamos usando hoje tem 20 bilhões de parâmetros, 32 especialistas — mas o custo computacional por token é parecido com um modelo bem menor, porque só 4 especialistas são ativados a cada vez.

Isso importa pro Kronk porque ele tem suporte nativo a colocar os especialistas em dispositivos diferentes — uma opção de configuração chamada `moe.mode`, que pode ser `auto`, `experts_cpu` (manda todos os especialistas pra CPU e deixa só a parte 'roteadora' na GPU), `experts_gpu`, ou `keep_top_n` (mantém só os N especialistas mais usados na GPU). Isso é uma ferramenta poderosíssima pra rodar modelos MoE gigantes em GPUs pequenas — você não precisa caber o modelo inteiro na VRAM, só a parte que realmente computa por token."

**Nota de bastidor:** "Originalmente preparamos essa talk num Linux com GPU dedicada de 12GB, e lá o gpt-oss-20b não cabia junto com um contexto grande — tivemos que usar um modelo denso menor. Aqui no Mac, com memória unificada, isso simplesmente não é um problema: o MoE de 20B cabe inteiro com folga de sobra. É um ótimo exemplo de como a mesma arquitetura de modelo se comporta diferente dependendo do tipo de memória da máquina."

---

## Slide 8 — Conceito: Hybrid CPU/GPU (offload de camadas)

**Tela:** Diagrama de um modelo "empilhado" em camadas (24 camadas, no nosso caso — gpt-oss-20b), com uma linha cortando: X camadas verdes (GPU) embaixo, resto cinza (CPU) acima.

**Roteiro:**
"Todo modelo de linguagem é uma pilha de camadas — transformer blocks. O Kronk deixa você decidir, camada por camada, onde cada uma roda: GPU ou CPU. Isso se chama offload de camadas, ou modo híbrido.

A configuração é `ngpu-layers`. E aqui tem uma pegadinha que a gente caiu na cara ao preparar essa talk, então vale contar: a convenção do Kronk é meio invertida do que você imaginaria. `0` ou deixar vazio significa 'todas as camadas na GPU'. `-1` significa 'nenhuma camada na GPU, força CPU'. E qualquer número positivo é o número exato de camadas a colocar na GPU. Documentação é sua amiga aqui — sem ler, você acaba invertendo tudo, igual a gente fez."

**Por que isso importa:** "Numa GPU dedicada com VRAM própria — uma RTX, por exemplo — modelo inteiro na GPU é mais rápido, mas exige que ele caiba todo na VRAM: pesos do modelo, mais o cache de contexto (que cresce junto com a conversa), mais buffers de computação. Se não cabe, ou você usa CPU, ou divide: parte na GPU, parte na CPU. É um dial de ajuste fino entre velocidade e memória disponível.

Num Mac com Apple Silicon, como o que estamos usando hoje, isso muda de figura: CPU e GPU compartilham a mesma memória unificada — não existe uma 'VRAM separada' pra estourar. Por isso aqui simplesmente deixamos `ngpu-layers: 0` (tudo na GPU) sem medo. O conceito de offload híbrido continua existindo e é igualmente importante em GPU dedicada — só não é o gargalo da nossa máquina hoje."

---

## Slide 9 — Nossos números reais (testados hoje, nesta máquina)

**Tela:** Tabela com os números reais (peguei do benchmarks-apresentacao.md):

| Configuração | Prompt tokens | Tempo até 1º token | Tokens/s |
|---|---|---|---|
| CPU only | pequeno | — | ~6-8 tok/s |
| Metal full offload — curto | 81 | ~0.2-0.4 s | ~70 tok/s |
| Metal full offload — médio (estilo Cline real) | 2.034 | ~2.2 s | ~65-68 tok/s |
| Metal full offload — stress (MCPs grandes) | 9.084 | ~11-12 s | ~63 tok/s |

**Roteiro:**
"Esses números não são de benchmark de paper, são de hoje, nesta máquina, MacBook M4 Pro com 48GB de memória unificada. Olha a primeira linha contra a segunda: Metal full offload é uns 9-10 vezes mais rápido que CPU pura pra gerar texto. Isso é o tipo de ganho que faz local-first ser viável de verdade.

E repara que aqui não tem nenhuma linha de 'falhou por falta de memória' — mesmo o prompt de stress, simulando um Cline carregado de ferramentas MCP com ~9 mil tokens de contexto, só fica mais lento pra começar a responder (TTFT sobe pra ~11-12s), mas não quebra. É a vantagem de memória unificada: o teto que existia numa GPU dedicada simplesmente não existe aqui da mesma forma."

**Gancho pra próximo slide:** "Mas isso não significa que não tem pegadinha nenhuma. Teve uma, e ela ensina uma lição melhor que qualquer slide teórico."

---

## Slide 10 — A pegadinha que vivemos hoje (war story)

**Tela:** A tabela do Slide 9 de novo, com destaque na linha de stress — TTFT de ~12s pra um prompt de 9k tokens.

**Roteiro:**
"Conto rapidinho o que aconteceu preparando essa demo. Vendo aquele TTFT de ~12 segundos no prompt grande, nosso primeiro instinto foi: 'deve ser o batch size pequeno demais, o Kronk não está mandando chunks grandes o suficiente pro Metal de uma vez'. Dobramos, depois quadruplicamos o `nbatch` e o `nubatch` — os parâmetros que controlam quantos tokens do prompt são processados por vez. Resultado: **6% de ganho**. De 12.36s pra 11.58s. Quase nada.

A lição real: pra esse modelo MoE, no Metal, o prefill não está limitado pelo tamanho do batch — está limitado pelo *compute* do kernel, especificamente o roteamento de especialistas (lembra do Slide 7? só 4 dos 32 especialistas são ativados por token, e decidir quais ativar tem um custo). Aumentar batch não ajuda quando o gargalo é compute, não throughput de dados. A gente só descobriu isso porque *medimos* antes de assumir — é fácil mudar um parâmetro, ver o número melhorar um pouquinho, e declarar vitória sem checar se valeu o esforço."

---

## Slide 11 — Conceito: "Thinking mode" e o custo escondido do raciocínio

**Tela:** Comparação visual: uma resposta direta ("ok") vs um bloco gigante de "pensamento" antes da resposta.

**Roteiro:**
"Último conceito, e esse foi uma surpresa real durante a preparação. O gpt-oss-20b, como outros modelos de raciocínio, tem um modo de 'pensamento' — antes de responder, ele gera um monólogo interno, tipo um rascunho de pensamento, que pode ter dezenas ou centenas de tokens (no nosso teste mais simples, uma pergunta de 'oi' já gastou 40 tokens de raciocínio antes da resposta). Isso melhora a qualidade da decisão — vimos na prática o modelo acertar a lógica certa de correção de bug só depois de 'pensar' bastante.

O problema: cada token de pensamento custa tempo de geração. O Kronk deixa controlar isso direto na chamada da API, com o parâmetro `reasoning_effort` (`low`, `medium`, `high`) — não é uma configuração fixa do modelo, é por requisição. `low` responde rápido mas pensa pouco; `high` demora mais, mas tem mais 'espaço' pra raciocinar antes de decidir. (Se der tempo de validar antes da talk, vale repetir o teste de bug-fix real com os dois extremos e trocar essa frase por um número real, em vez de ficar na teoria.)

É um tradeoff real que qualquer um rodando IA local vai enfrentar: velocidade contra profundidade de raciocínio."

---

## Slide 12 — Resumo / Takeaways

**Tela:** Bullet points grandes, tipo "cartões".

**Roteiro (ler cada um com uma pausa):**
- "Local-first com Kronk é rápido de verdade — uns 9-10x mais rápido que CPU pura nesta máquina."
- "Memória unificada (Apple Silicon) muda o jogo do offload — mas o conceito de MoE e camadas continua valendo em qualquer GPU dedicada."
- "A convenção de camadas do Kronk é invertida — leia a doc antes de configurar."
- "Não assuma o gargalo, meça: batch size maior só ganhou 6% porque o limite real era compute, não throughput."
- "Modo de raciocínio é um dial, não uma mágica: mais qualidade custa mais tempo, e é configurável por requisição."
- "Tudo isso é configurável em um arquivo YAML ou direto na chamada da API, sem precisar recompilar nada."

**Fechamento:** "E é exatamente essa flexibilidade — CPU, GPU, híbrido, MoE, raciocínio ligado ou desligado — que torna o Kronk uma peça de infraestrutura séria pra rodar IA local em produção, não só em demo de palco."

---

---

# Roteiro — Seção: Por dentro do código dos 2 MCPs

Essa seção entra bem antes do deep dive do Kronk (ou logo depois do "momento WOW") — é o "abre a caixa" do que realmente faz a diferença de tokens/custo entre as duas cenas. Os três arquivos-fonte são: `mcp_vault/internal/vault/reader.go` (compartilhado pelos dois), `mcp_vault/cmd/dumb/main.go` (Cena 1), `mcp_vault/cmd/smart/main.go` + `mcp_vault/internal/store/duckdb.go` (Cena 2).

---

## Slide 13 — A base compartilhada: lendo o vault

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

## Slide 14 — Cena 1: o servidor "burro" (`vault_dumb`)

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

## Slide 15 — Cena 2, parte 1: montando o índice vetorial (`vault_smart`, setup)

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

## Slide 16 — Cena 2, parte 2: a busca de verdade (`search_vault`)

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

## Slide 17 — Lado a lado: mesma pergunta, dois caminhos

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

- Os números do Slide 9 vêm de `benchmarks-apresentacao.md` neste mesmo repo — conferir se não mudaram caso re-teste antes da talk.
- Slide 10 mudou de "war story de VRAM" pra "lição de não assumir o gargalo" — tom ainda pode ser leve/autêntico, mas o ponto agora é mais de engenharia (medir antes de otimizar) do que de drama de travamento.
- Se houver tempo, o Slide 7 (MoE) pode ganhar uma demonstração rápida ao vivo trocando `moe.mode` no `model_config.yaml` do gpt-oss-20b — mas só se já tiver sido testado antes, não arriscar ao vivo sem ensaio.
- Slides 3-5 (arquitetura interna, pool de modelos, cache incremental) são novos — ensaiar com o log real (`/tmp/kronk.log`) aberto numa aba, caso queiram mostrar a linha de log de verdade em vez de só o trecho colado no slide.
