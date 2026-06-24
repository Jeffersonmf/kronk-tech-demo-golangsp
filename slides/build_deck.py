"""Builds the GolangSP "Golang + AI + MCP" deck as a .pptx, ready to open in
Google Slides (File > Open, or upload to Drive and "Open with Google Slides").

Source of truth for content is the roteiro markdown files at the repo root:
slides-abertura.md, slides-kronk-deepdive.md, slides-demos-finais.md. On-screen
bullets here are a condensed version of each slide's "Tela" + "Roteiro"; the
full spoken "Roteiro" text goes into the PPTX speaker notes, not on the slide.

Run: python3 slides/build_deck.py
"""

from pptx import Presentation
from pptx.util import Inches, Pt
from pptx.enum.text import PP_ALIGN
from pptx.dml.color import RGBColor

OUT_PATH = "slides/GolangSP-Kronk-Talk.pptx"
GOPHER_PATH = "slides/assets/gopher.png"  # official Go gopher (Renee French, CC BY 3.0)

DARK = RGBColor(0x1A, 0x1A, 0x1A)
ACCENT = RGBColor(0x00, 0x6F, 0xE6)
CODE_BG = RGBColor(0x24, 0x24, 0x24)
CODE_FG = RGBColor(0xE0, 0xE0, 0xE0)

SLIDES = [
    {"kind": "title", "title": "IA sem Nuvem, em Golang",
     "subtitle": "Golang + AI + MCP — GolangSP, parte 2",
     "notes": "Boa noite, pessoal! Essa é a parte 2 da talk sobre Go. Na "
               "parte 1 falamos de Go como ferramenta de produtividade. "
               "Agora: Go encontrando IA, e o protocolo que conecta os "
               "dois — o MCP."},

    # --- Abertura ---------------------------------------------------
    {"title": "IA sem Nuvem, em Golang",
     "bullets": ["Parte 2 da talk sobre Go",
                 "Parte 1: Go como ferramenta de produtividade",
                 "Parte 2: Go encontrando IA — e o protocolo MCP"],
     "notes": "Boa noite, pessoal! Essa é a parte 2 da talk que comecei "
               "sobre Go. Na parte 1 a gente falou de Go como ferramenta "
               "de produtividade. Agora vamos um passo além: Go "
               "encontrando IA, e o protocolo que está conectando os "
               "dois — o MCP."},
    {"title": "Recap rápido da Parte 1",
     "bullets": ["Go não é só linguagem de backend — é ferramenta de "
                 "produtividade no dia a dia (scripts, automação, tooling)",
                 "Essa parte 2 é praticamente independente",
                 "O que importa saber: Go é uma linguagem séria pra "
                 "construir infraestrutura — hoje, infraestrutura de IA"],
     "notes": "Resumo de 30 segundos pra quem chegou agora: na parte 1 "
               "mostramos Go não só como linguagem de backend, mas como "
               "ferramenta de produtividade no dia a dia. O que importa "
               "saber é uma coisa: Go é uma linguagem séria pra construir "
               "infraestrutura. E hoje vamos provar isso construindo "
               "infraestrutura de IA."},
    {"title": "O problema que ninguém fala em voz alta",
     "bullets": ["Toda demo de agente de IA que você viu este ano "
                 "rodou na nuvem.",
                 "E custou dinheiro a cada clique."],
     "notes": "Quero começar com uma provocação. Quantos de vocês já "
               "usaram um agente de IA pra programar? Quantos sabem "
               "exatamente quanto cada interação custa, em tokens, em "
               "dinheiro? É essa lacuna que essa talk ataca. Não é sobre "
               "'IA é mágica'. É sobre entender o motor por dentro."},
    {"title": "CONTEXTO",
     "bullets": ["O que você manda pro modelo é o que você paga",
                 "— e o que decide se ele acerta.",
                 "Como dar pro modelo só o que ele precisa, quando "
                 "precisa? A resposta moderna: MCP."],
     "notes": "Toda interação com um modelo de linguagem tem um custo "
               "proporcional a uma coisa: quanto texto você manda pra "
               "ele, o contexto. Quanto mais contexto irrelevante, mais "
               "caro, mais lento, e às vezes pior a qualidade da "
               "resposta."},
    {"title": "O que é MCP, de verdade (sem hype)",
     "bullets": ["Model Context Protocol — protocolo aberto, da Anthropic",
                 "Padroniza como um agente de IA pede informação e "
                 "executa ações fora de si mesmo",
                 "O \"USB\" da IA: escreve um servidor MCP uma vez, "
                 "qualquer cliente compatível conversa com ele",
                 "Escrever um servidor MCP é, literalmente, escrever "
                 "um servidor — e Go faz isso muito bem"],
     "notes": "MCP é um protocolo aberto que padroniza como um agente de "
               "IA pede informação e executa ações fora de si mesmo. "
               "Antes do MCP, cada ferramenta tinha sua forma proprietária "
               "de se conectar a dados. O pulo do gato pra essa "
               "comunidade Go: escrever um servidor MCP é escrever um "
               "servidor, coisa que Go faz muito bem há muito tempo."},
    {"title": "Existe um jeito certo e um jeito ingênuo de fazer isso",
     "bullets": ["Jeito ingênuo: manda tudo, deixa o modelo se virar",
                 "Jeito certo: manda só o que importa — busca semântica "
                 "de verdade, local, sem custar nuvem",
                 "Primeiro as peças (Kronk), depois os dois lado a "
                 "lado, ao vivo, com custo na tela"],
     "notes": "Vou mostrar dois jeitos de construir um servidor MCP que "
               "resolve o MESMO problema. Um é ingênuo: pega tudo, manda "
               "tudo. O outro usa busca semântica de verdade, rodando "
               "localmente, sem custar nuvem."},
    {"title": "Roteiro de hoje",
     "bullets": ["1. O motor por baixo: Kronk (quantização, MoE, "
                 "híbrido CPU/GPU)",
                 "2. Números reais de hoje, nesta máquina",
                 "3. Por dentro do código dos dois servidores MCP",
                 "4. Ao vivo: as duas cenas, lado a lado, com custo "
                 "na tela"],
     "notes": "Esse é o mapa de hoje. Vamos primeiro entender o motor — "
               "o Kronk. Depois abrir o código dos dois servidores MCP. "
               "E só no final, com tudo entendido, ver os dois rodando "
               "ao vivo, com o contador de tokens e custo do Cline na "
               "tela. Bora?"},

    # --- Deep Dive: Kronk --------------------------------------------
    {"title": "Slide 1 — O que é o Kronk?",
     "bullets": ["SDK em Go (Ardan Labs) que envolve o llama.cpp",
                 "Roda LLMs localmente — sem Python, sem internet, "
                 "sem nuvem",
                 "API compatível com OpenAI — Cline/Cursor funcionam só "
                 "mudando a URL",
                 "Infraestrutura: server, pool de modelos, cache, auth, "
                 "métricas — não só demo"],
     "notes": "Kronk é um SDK em Go, da Ardan Labs, que envolve o "
               "llama.cpp. A grande sacada: ele expõe uma API compatível "
               "com a da OpenAI. Kronk não é 'mais um wrapper de "
               "chatbot'. É infraestrutura, pensada pra produção."},
    {"title": "Slide 2 — Arquitetura: onde o Kronk entra",
     "bullets": ["[Agente/Cline] → API OpenAI-compatible → [Kronk Server]",
                 "→ llama.cpp (CPU / CUDA / Metal / ROCm / Vulkan)",
                 "→ [Modelo GGUF]",
                 "Backend escolhido pela flag --processor"],
     "notes": "O Kronk fica no meio. De um lado, qualquer cliente que "
               "fale o protocolo da OpenAI. Do outro, ele escolhe "
               "automaticamente o melhor backend de hardware disponível. "
               "O modelo é um arquivo .gguf, formato do ecossistema "
               "llama.cpp."},
    {"title": "Slide 3 — Por dentro do Kronk Server",
     "bullets": ["API router — porta 11435 (chat completions, "
                 "compatível OpenAI)",
                 "Debug v1 router — porta 11445 (métricas, introspecção)",
                 "Auth server (JWT/API keys) — desligado hoje, existe "
                 "pra produção",
                 "Servidor MCP nativo do próprio Kronk — porta 9000 "
                 "(não é nenhum dos nossos MCPs de demo!)"],
     "code_lang": "json",
     "code": '"status":"api router started","host":"0.0.0.0:11435"\n'
             '"status":"debug v1 router started","host":"0.0.0.0:11445"\n'
             '"status":"local mcp server started","host":"127.0.0.1:9000"\n'
             '"status":"start auth server"',
     "notes": "Vamos abrir a caixa preta do Slide 2 com log real de hoje. "
               "Quando você dá kronk server start, sobem pelo menos "
               "quatro coisas: router da API, router de debug, serviço "
               "de auth, e um servidor MCP nativo do próprio Kronk. Não "
               "é um wrapper fininho — é um processo de servidor com "
               "superfície de operação real."},
    {"title": "Slide 4 — Pool de modelos",
     "bullets": ["Orçamento de memória: só 80% da RAM/VRAM disponível "
                 "(hoje: 38.7GB)",
                 "Modelo que não cabe é recusado antes de tentar carregar",
                 "Até 10 modelos residentes ao mesmo tempo no pool",
                 "Agora: gpt-oss-20b (chat) + embeddinggemma "
                 "(embedding) juntos",
                 "TTL de ociosidade — modelo sem uso é descarregado "
                 "automaticamente"],
     "code_lang": "json",
     "code": '"status":"resman-init","budget-percent":80,\n'
             '"ram-budget":"38.7GB","max-models-in-pool":10',
     "notes": "O Kronk não trata modelo como 'um arquivo que você abre "
               "uma vez'. Ele tem um pool de modelos — orçamento de "
               "memória, múltiplos modelos simultâneos, TTL de "
               "ociosidade. Foi esse mecanismo de segurança que travou "
               "a configuração original dessa talk no Linux com GPU "
               "de 12GB."},
    {"title": "Slide 5 — Cache incremental (IMC)",
     "bullets": ["Conversa agêntica reenvia todo o histórico a cada turno",
                 "Sem cache: reprocessar tudo do zero, cada vez mais lento",
                 "IMC guarda o KV-cache entre turnos — só processa "
                 "tokens novos",
                 "Ligado por padrão — pensado pra Cline/OpenCode"],
     "code_lang": "yaml",
     "code": "incremental-cache: false   # Incremental message caching\n"
             "                           # for agentic workflows\n"
             "                           # (Cline, OpenCode)\n"
             "# valor real hoje, confirmado no log: IncrementalCache[true]",
     "notes": "Uma conversa com o Cline reenvia todo o histórico "
               "anterior a cada turno. O cache incremental de mensagens "
               "guarda o estado da conversa e só processa os tokens "
               "novos no turno seguinte. Vem ligado por padrão pra "
               "clientes agênticos."},
    {"title": "Slide 6 — Quantização (a base de tudo)",
     "table": [["Formato", "Bits", "Perda de qualidade"],
               ["F16", "16", "Nenhuma"],
               ["Q8_0", "8", "Quase nenhuma"],
               ["Q4_K_M", "4", "Pequena, raramente perceptível em código"]],
     "bullets": ["Hoje: gpt-oss-20b em Q8_0 — 20B parâmetros, ~12GB "
                 "em disco"],
     "notes": "Quantização é comprimir os números de ponto flutuante "
               "do modelo pra ocupar menos memória. Q8_0 usa 8 bits — "
               "praticamente sem perda. Q4_K_M usa 4 bits — mais "
               "compacto, perdinha quase imperceptível em código."},
    {"title": "Slide 7 — MoE (Mixture of Experts)",
     "bullets": ["Modelo dividido em sub-redes \"especialistas\"; um "
                 "roteador escolhe quais ativar por token",
                 "gpt-oss-20b: 20B parâmetros, 32 especialistas, só 4 "
                 "ativados por vez",
                 "Kronk: moe.mode = auto / experts_cpu / experts_gpu / "
                 "keep_top_n",
                 "Permite rodar MoE gigante em GPU pequena"],
     "notes": "MoE: para cada token, um roteador decide quais "
               "especialistas ativar — só uma fração do total. O "
               "gpt-oss-20b tem 20B parâmetros mas custo computacional "
               "de um modelo bem menor. Nota de bastidor: no Linux com "
               "GPU de 12GB, o gpt-oss-20b não cabia com contexto "
               "grande — usamos um modelo denso menor. No Mac, com "
               "memória unificada, o MoE de 20B cabe com folga."},
    {"title": "Slide 8 — Hybrid CPU/GPU (offload de camadas)",
     "bullets": ["ngpu-layers decide, camada por camada, GPU ou CPU",
                 "Convenção invertida: 0 = todas na GPU, -1 = nenhuma "
                 "(força CPU), N = N camadas",
                 "GPU dedicada: modelo precisa caber todo na VRAM, "
                 "senão precisa offload parcial",
                 "Mac (memória unificada): não existe \"VRAM separada\" "
                 "pra estourar — ngpu-layers:0 sem medo"],
     "notes": "Pegadinha real que vivemos: a convenção do Kronk é "
               "invertida do que se imagina — 0 significa todas as "
               "camadas na GPU. Numa GPU dedicada isso é um dial real "
               "entre velocidade e memória; no Mac com memória "
               "unificada, o gargalo simplesmente não existe da mesma "
               "forma."},
    {"title": "Slide 9 — Nossos números reais (hoje, nesta máquina)",
     "table": [["Configuração", "Prompt tokens", "TTFT", "Tokens/s"],
               ["CPU only", "pequeno", "—", "~6-8 tok/s"],
               ["Metal full offload — curto", "81", "~0.2-0.4s",
                "~70 tok/s"],
               ["Metal full offload — médio (Cline real)", "2.034",
                "~2.2s", "~65-68 tok/s"],
               ["Metal full offload — stress (MCPs grandes)", "9.084",
                "~11-12s", "~63 tok/s"]],
     "bullets": ["Metal full offload é ~9-10x mais rápido que CPU pura",
                 "Nem o prompt de stress (~9k tokens) falha por falta "
                 "de memória"],
     "notes": "MacBook M4 Pro, 48GB de memória unificada. Metal full "
               "offload é uns 9-10x mais rápido que CPU pura. E nenhuma "
               "linha falha por falta de memória — vantagem real da "
               "memória unificada. Mas teve uma pegadinha, e ela ensina "
               "uma lição melhor que qualquer slide teórico."},
    {"title": "Slide 10 — A pegadinha que vivemos hoje (war story)",
     "bullets": ["TTFT de ~12s no prompt grande → suspeita: batch "
                 "pequeno demais",
                 "Dobramos, depois quadruplicamos nbatch/nubatch",
                 "Resultado: só 6% de ganho (12.36s → 11.58s)",
                 "Lição real: prefill é limitado por compute "
                 "(roteamento MoE), não por batch",
                 "Meça antes de assumir o gargalo"],
     "notes": "Vendo o TTFT de ~12s, nosso instinto foi 'aumentar o "
               "batch'. Resultado: 6% de ganho, quase nada. A lição "
               "real: pra esse modelo MoE no Metal, o prefill é "
               "limitado por compute (roteamento de especialistas), não "
               "por throughput de batch. Só descobrimos porque medimos "
               "antes de assumir."},
    {"title": 'Slide 11 — "Thinking mode" e o custo do raciocínio',
     "bullets": ["gpt-oss-20b gera um \"rascunho de pensamento\" antes "
                 "de responder",
                 "Mesmo um \"oi\" gastou 40 tokens de raciocínio",
                 "Controlado por requisição: reasoning_effort "
                 "(low/medium/high)",
                 "Tradeoff real: velocidade vs profundidade de raciocínio"],
     "notes": "O gpt-oss-20b tem um modo de raciocínio que melhora a "
               "qualidade da decisão, mas cada token de pensamento "
               "custa tempo. O Kronk deixa controlar isso por "
               "requisição, com reasoning_effort — não é configuração "
               "fixa do modelo."},
    {"title": "Slide 12 — Trocar de modelo é runtime, não startup",
     "code_lang": "bash",
     "code": "kronk model list --local\n"
             "kronk model pull <nome>\n"
             "\n"
             "# depois, só aponta o cliente pro novo nome:\n"
             "# Cline -> Settings -> Model ID\n"
             "# curl  -> campo \"model\" no JSON",
     "bullets": ["kronk server start sobe o pool, não fixa um modelo",
                 "Modelo é escolhido por requisição (campo \"model\")",
                 "Não baixado? kronk model pull <nome>",
                 "Pool carrega/descarrega sob demanda, sem reiniciar"],
     "notes": "Detalhe que confunde quem vem de outras stacks: você "
               "não escolhe o modelo no boot do servidor. O pool "
               "decide por requisição, igual a API da OpenAI. Pra "
               "trocar: confere com model list --local, baixa com "
               "model pull se precisar, e só aponta o cliente pro "
               "novo nome — sem reiniciar o kronk server start."},
    {"title": "Slide 13 — Maior ou menor? O dial de tamanho no model_config.yaml",
     "table": [["Modelo", "Tamanho", "Status"],
               ["Qwopus3.5-4B-Coder", "menor (4B)",
                "configurado, não baixado"],
               ["gpt-oss-20b-Q8_0", "usado hoje (20B, MoE)",
                "baixado e quente"],
               ["gemma-4-26B-A4B-it", "maior (26B-A4B, MoE)",
                "configurado, não baixado"],
               ["Qwen3.6-35B-A3B", "maior (35B-A3B, MoE)",
                "configurado, não baixado"]],
     "notes": "Tamanho de modelo também é um dial, não decisão "
               "travada em pedra. Mesma API, só muda o nome no campo "
               "model (e baixar antes, se for novo). Já tenho tuning "
               "pronto pra dois modelos bem maiores que o de hoje e "
               "um menor (4B) pra hardware mais limitado — a config "
               "já existe no model_config.yaml desta máquina, só "
               "falta o pull."},
    {"title": "Slide 14 — Resumo / Takeaways",
     "bullets": ["Local-first com Kronk é rápido de verdade — uns "
                 "9-10x vs CPU pura",
                 "Memória unificada muda o jogo do offload — mas MoE e "
                 "camadas valem em qualquer GPU dedicada",
                 "A convenção de camadas do Kronk é invertida — leia "
                 "a doc",
                 "Não assuma o gargalo, meça",
                 "Raciocínio é um dial, não mágica — configurável por "
                 "requisição",
                 "Tudo configurável em YAML ou na chamada da API, sem "
                 "recompilar"],
     "notes": "É exatamente essa flexibilidade — CPU, GPU, híbrido, "
               "MoE, raciocínio ligado ou desligado — que torna o "
               "Kronk uma peça de infraestrutura séria pra rodar IA "
               "local em produção, não só em demo de palco."},
    {"title": "Slide 15 — A base compartilhada: lendo o vault",
     "code_lang": "go",
     "code": 'func LoadAll(root string) ([]Document, error) {\n'
             '    var docs []Document\n'
             '    filepath.WalkDir(root, func(path string, d os.DirEntry,'
             ' err error) error {\n'
             '        if d.IsDir() || !strings.HasSuffix(d.Name(), '
             '".md") {\n'
             '            return nil\n'
             '        }\n'
             '        content, _ := os.ReadFile(path)\n'
             '        rel, _ := filepath.Rel(root, path)\n'
             '        docs = append(docs, Document{Path: rel, '
             'Content: string(content)})\n'
             '        return nil\n'
             '    })\n'
             '    return docs, nil\n'
             '}',
     "bullets": ["Função compartilhada pelos dois MCPs",
                 "26 documentos passam pela mesma função nos dois "
                 "servidores",
                 "A diferença está em O QUE cada servidor faz depois"],
     "notes": "LoadAll anda pela pasta do vault, pega todo .md, lê o "
               "conteúdo bruto e devolve uma lista de documentos. "
               "Simples, sem mágica. A diferença entre as duas cenas "
               "não está em COMO os dados são lidos."},
    {"title": 'Slide 16 — Cena 1: o servidor "burro" (vault_dumb)',
     "code_lang": "go",
     "code": 'mcp.AddTool(server, &mcp.Tool{\n'
             '    Name: "read_vault",\n'
             '    Description: "Returns the full contents of every '
             'Markdown file...",\n'
             '}, readVaultHandler(vaultPath))\n\n'
             'func readVaultHandler(vaultPath string) ... {\n'
             '    docs, _ := vault.LoadAll(vaultPath)\n'
             '    // concatena tudo: --- FILE: path ---\\ncontent\\n\\n\n'
             '}',
     "bullets": ["Uma única ferramenta: read_vault (sem parâmetros)",
                 "Lê TODOS os documentos e concatena tudo",
                 "Sem filtro, sem busca, sem relevância",
                 "~47.000 tokens consumidos · <20 linhas de código"],
     "notes": "O agente pediu 'me ajuda com a TASK-001' e essa "
               "ferramenta devolve as 26 tasks completas. É proposital: "
               "representa o jeito ingênuo de integrar LLM com dados. "
               "É exatamente essa simplicidade que custa caro."},
    {"title": "Slide 17 — Cena 2, parte 1: montando o índice vetorial",
     "code_lang": "go",
     "code": 'krnEmbed, _ := kronk.New(model.WithModelFiles(\n'
             '    []string{embedModelPath}))\n'
             'docs, _ := vault.LoadAll(vaultPath)\n'
             'dims, _ := embeddingDims(krnEmbed)\n'
             'db, _ := store.LoadVault(ctx, krnEmbed, docs, dims)',
     "bullets": ["Carrega um 2º modelo: embeddinggemma (especializado "
                 "em embeddings)",
                 "26 documentos → vetores de 768 números",
                 "DuckDB em memória + extensão VSS + índice HNSW",
                 "Indexação acontece uma vez, na inicialização"],
     "notes": "Antes de aceitar qualquer pergunta, a Cena 2 carrega um "
               "modelo de embedding e transforma cada documento num "
               "vetor — 'coordenadas de significado' num espaço "
               "matemático. Isso vai pro DuckDB com índice HNSW."},
    {"title": "Slide 18 — Cena 2, parte 2: a busca de verdade",
     "code_lang": "sql",
     "code": "SELECT path, text,\n"
             "       array_cosine_similarity(embedding, "
             "?::FLOAT[768]) AS similarity\n"
             "FROM items\n"
             "ORDER BY similarity DESC\n"
             "LIMIT 3;",
     "bullets": ["Pergunta → vetor de 768 números (mesmo modelo)",
                 "Top-3 por similaridade de cosseno — busca de "
                 "sentido, não palavra-chave",
                 "26 documentos completos → só os 3 mais relevantes",
                 "~47.000 tokens → ~4.000 tokens"],
     "notes": "search_vault transforma a pergunta no mesmo tipo de "
               "vetor e busca os 3 documentos mais similares por "
               "cosseno. É por isso que perguntar sobre 'dinheiro e "
               "pagamentos' acha a task de cupom mesmo sem a palavra "
               "'dinheiro' aparecer."},
    {"title": "Slide 19 — Lado a lado: mesma pergunta, dois caminhos",
     "table": [["", "vault_dumb (Cena 1)", "vault_smart (Cena 2)"],
               ["Ferramenta MCP", "read_vault (sem parâmetros)",
                "search_vault(query, top_k)"],
               ["O que faz", "Concatena tudo",
                "Embedding + busca por similaridade"],
               ["Preparação", "Nenhuma", "Indexa o vault uma vez"],
               ["Infra extra", "Nenhuma", "Embedding + DuckDB + HNSW"],
               ["Tokens (teste)", "~47.000", "~4.000"],
               ["Linhas de código", "~20", "~80"]],
     "notes": "Cena 1 é 4x mais simples de código, mas custa 10x mais "
               "tokens. Cena 2 exige infraestrutura extra, mas paga "
               "isso de volta multiplicado em economia de contexto — e "
               "tudo local, via Kronk, sem custar um centavo."},

    # --- Demos Finais -------------------------------------------------
    {"title": "Slide 20 — Chegou a hora",
     "bullets": ["Mesma pergunta. Mesmo agente.",
                 "Duas formas de buscar contexto.",
                 "Cena 1  vs  Cena 2"],
     "notes": "Lembram do gancho do começo — 'tem um jeito ingênuo e um "
               "jeito inteligente'? Vou rodar exatamente a mesma tarefa "
               "duas vezes, mudando só o MCP ativo e o modelo por trás."},
    {"title": "Slide 21 — Apresentando a Cena 1",
     "bullets": ["Provider: conta paga na nuvem",
                 "MCP ativo: vault_dumb",
                 "Contador de tokens e custo visíveis — hoje: zero"],
     "notes": "Cena 1 representa o jeito mais comum hoje: modelo na "
               "nuvem, pago por token, conectado a um MCP que não "
               "filtra nada. Olhem os números antes de eu começar: "
               "zero."},
    {"title": "Slide 22 — Cena 1 ao vivo  [DEMO]",
     "bullets": ["Prompt: \"Leia a TASK-001 no vault do Obsidian e "
                 "resolva o problema descrito nela.\"",
                 "Cline chama read_vault (única ferramenta disponível)",
                 "Observar o contador de tokens/custo subir"],
     "notes": "Reparem: o Cline está decidindo chamar read_vault — só "
               "tem essa ferramenta disponível. Contingência: se a "
               "geração de nuvem travar num tool-call depois, o ponto "
               "já foi feito no instante em que read_vault devolve a "
               "resposta — pode pausar aí e seguir pra Cena 2."},
    {"title": "Slide 23 — O número da Cena 1",
     "bullets": ["[Print real do contador do Cline ao final da rodada]",
                 "Tokens e custo em destaque"],
     "notes": "Aqui está o resultado da Cena 1: [ler os números reais "
               "capturados]. Token e dinheiro de verdade, gastos só "
               "porque a ferramenta de busca não soube filtrar nada. "
               "Guardem esse número."},
    {"title": "Slide 24 — Apresentando a Cena 2",
     "bullets": ["Provider: Kronk local (http://localhost:11435/v1)",
                 "MCP ativo: vault_smart",
                 "Mesmo prompt exato da Cena 1"],
     "notes": "Troquei o provider do Cline pra apontar pro Kronk, "
               "rodando aqui, agora, sem internet envolvida. E troquei "
               "o MCP ativo pra vault_smart."},
    {"title": "Slide 25 — Cena 2 ao vivo  [DEMO]",
     "bullets": ["Mesmo prompt, colado de novo",
                 "Cline chama search_vault",
                 "Retorna só os documentos relevantes",
                 "Custo: zero — roda local"],
     "notes": "Reparem que a ferramenta disponível agora é outra: "
               "search_vault. Contingência: o modelo local com "
               "raciocínio pode demorar mais que a nuvem pra terminar "
               "o fix completo — o ponto da demo é o contraste de "
               "tokens/custo na chamada da ferramenta, que aparece "
               "rápido."},
    {"title": "Slide 26 — Lado a lado: o contraste final",
     "table": [["", "Cena 1 (nuvem + dumb)", "Cena 2 (local + smart)"],
               ["Tokens consumidos", "[número real]", "[número real]"],
               ["Custo", "[$ real]", "$0,00"],
               ["Onde rodou", "Servidor de terceiro",
                "Esta máquina, agora"],
               ["Internet necessária", "Sim", "Não"]],
     "notes": "Essa é a imagem que eu queria que vocês levassem pra "
               "casa. Mesma tarefa, mesmo agente. A diferença está em "
               "como a informação chega até o modelo, e onde o modelo "
               "roda."},
    {"title": "Slide 27 — Fechamento",
     "gopher": True,
     "bullets": ["Go não é só a linguagem do seu backend.",
                 "É a linguagem da sua infraestrutura de IA também.",
                 "Repositório: github.com/[seu-usuario]/"
                 "kronk-tech-demo-golangsp"],
     "notes": "Tudo que vocês viram hoje — servidor MCP, motor de "
               "inferência, busca vetorial — é Go, de ponta a ponta. "
               "Não precisamos de Python, não precisamos de nuvem. "
               "Obrigado!"},
]


def add_bullets(tf, bullets, first_clear=True):
    if first_clear:
        tf.text = bullets[0]
        rest = bullets[1:]
    else:
        rest = bullets
    for b in rest:
        p = tf.add_paragraph()
        p.text = b
    for p in tf.paragraphs:
        p.font.size = Pt(20)
        p.font.color.rgb = DARK
        p.level = 0


def add_code_box(slide, code, left, top, width, height):
    box = slide.shapes.add_textbox(left, top, width, height)
    box.fill.solid()
    box.fill.fore_color.rgb = CODE_BG
    tf = box.text_frame
    tf.word_wrap = True
    lines = code.split("\n")
    tf.text = lines[0]
    for line in lines[1:]:
        p = tf.add_paragraph()
        p.text = line
    for p in tf.paragraphs:
        p.font.size = Pt(13)
        p.font.name = "Menlo"
        p.font.color.rgb = CODE_FG
    return box


def add_table(slide, rows, left, top, width, height):
    n_rows, n_cols = len(rows), len(rows[0])
    gtable = slide.shapes.add_table(n_rows, n_cols, left, top, width,
                                     height).table
    for c, header in enumerate(rows[0]):
        cell = gtable.cell(0, c)
        cell.text = header
        cell.text_frame.paragraphs[0].font.bold = True
        cell.text_frame.paragraphs[0].font.size = Pt(14)
    for r in range(1, n_rows):
        for c in range(n_cols):
            cell = gtable.cell(r, c)
            cell.text = rows[r][c]
            cell.text_frame.paragraphs[0].font.size = Pt(13)
    return gtable


def build():
    prs = Presentation()
    prs.slide_width = Inches(13.333)
    prs.slide_height = Inches(7.5)

    title_layout = prs.slide_layouts[0]
    blank_title_layout = prs.slide_layouts[5]  # "Title Only"

    for spec in SLIDES:
        if spec.get("kind") == "title":
            slide = prs.slides.add_slide(title_layout)
            slide.shapes.title.text = spec["title"]
            slide.placeholders[1].text = spec["subtitle"]
            slide.shapes.title.text_frame.paragraphs[0].font.size = Pt(44)
            slide.shapes.add_picture(GOPHER_PATH, Inches(11.1), Inches(0.4),
                                      height=Inches(1.8))
        else:
            slide = prs.slides.add_slide(blank_title_layout)
            slide.shapes.title.text = spec["title"]
            slide.shapes.title.text_frame.paragraphs[0].font.size = Pt(28)
            slide.shapes.title.text_frame.paragraphs[0].font.color.rgb = ACCENT

            has_code = "code" in spec
            has_table = "table" in spec
            bullets = spec.get("bullets")

            content_top = Inches(1.4)

            if bullets:
                bw = Inches(6.0) if (has_code or has_table) else Inches(12.0)
                box = slide.shapes.add_textbox(Inches(0.6), content_top,
                                                bw, Inches(5.6))
                add_bullets(box.text_frame, bullets)

            side_left = Inches(6.9) if bullets else Inches(0.6)
            side_width = Inches(6.0) if bullets else Inches(12.0)

            if has_code:
                add_code_box(slide, spec["code"], side_left, content_top,
                             side_width, Inches(5.2))
            elif has_table:
                rows = spec["table"]
                add_table(slide, rows, side_left, content_top,
                          side_width, Inches(0.5 * len(rows)))

            if spec.get("gopher"):
                slide.shapes.add_picture(GOPHER_PATH, Inches(11.1),
                                          Inches(5.6), height=Inches(1.6))

        notes = spec.get("notes")
        if notes:
            slide.notes_slide.notes_text_frame.text = notes

    prs.save(OUT_PATH)
    print(f"Wrote {len(prs.slides._sldIdLst)} slides to {OUT_PATH}")


if __name__ == "__main__":
    build()
