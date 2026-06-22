# Roteiro — Abertura / Intro (Parte 2)

Ordem da talk completa, pra não perder o fio:
1. **Abertura/Intro** (este arquivo)
2. Deep dive Kronk (`slides-kronk-deepdive.md`, slides 1-9)
3. Por dentro do código dos 2 MCPs (`slides-kronk-deepdive.md`, slides 10-14)
4. **Demos ao vivo — Cena 1 e Cena 2** (no final, é o pagamento de tudo que foi explicado antes)

---

## Slide A1 — Título

**Tela:** Nome da talk + seu nome/handle. Algo como: "Golang + AI + MCP: tirando a IA da nuvem" (ajuste pro título real que vocês escolheram pro CFP).

**Roteiro:**
"Boa noite, pessoal! Essa é a parte 2 da talk que comecei [ou: que fulano começou] sobre Go. Na parte 1 a gente falou de Go como ferramenta de produtividade. Agora vamos um passo além: Go encontrando IA, e o protocolo que está conectando os dois — o MCP."

---

## Slide A2 — Recap rápido da Parte 1 (pra quem não viu)

**Tela:** 2-3 bullets do que foi a parte 1 (ajustar pro conteúdo real).

**Roteiro:**
"Resumo de 30 segundos pra quem chegou agora: na parte 1 mostramos Go não só como linguagem de backend, mas como ferramenta de produtividade no dia a dia — scripts, automação, tooling. Se você não viu, não tem problema, essa parte 2 é praticamente independente. O que importa saber é uma coisa: Go é uma linguagem séria pra construir infraestrutura. E hoje vamos provar isso construindo infraestrutura de IA."

---

## Slide A3 — O problema que ninguém fala em voz alta

**Tela:** Frase grande, sozinha: "Toda demo de agente de IA que você viu este ano rodou na nuvem. E custou dinheiro a cada clique."

**Roteiro:**
"Quero começar com uma provocação. Quantos de vocês já usaram um agente de IA pra programar — Cursor, Copilot, Cline, Claude Code? [pausa, esperar mãos] Quantos sabem exatamente quanto cada interação custa, em tokens, em dinheiro? [pausa] É essa lacuna que essa talk ataca. Não é sobre 'IA é mágica'. É sobre entender o motor por dentro — quanto custa, onde o dinheiro vai, e se existe alternativa."

---

## Slide A4 — A resposta em uma palavra: contexto

**Tela:** Palavra "CONTEXTO" grande, com uma sub-frase: "O que você manda pro modelo é o que você paga — e o que decide se ele acerta."

**Roteiro:**
"Toda interação com um modelo de linguagem tem um custo proporcional a uma coisa: quanto texto você manda pra ele, o chamado contexto. Quanto mais contexto irrelevante, mais caro, mais lento, e às vezes pior a qualidade da resposta — o modelo se perde no meio do ruído. A pergunta que vale ouro é: como você dá pro modelo *só* o que ele precisa, no momento que ele precisa? A resposta moderna pra essa pergunta tem um nome: MCP."

---

## Slide A5 — O que é MCP, de verdade (sem hype)

**Tela:** Diagrama simples: [Agente de IA] ←→ (protocolo MCP) ←→ [Ferramentas/Dados: arquivos, APIs, bancos, vault de notas...]

**Roteiro:**
"MCP — Model Context Protocol — é um protocolo aberto, criado pela Anthropic, que padroniza como um agente de IA pede informação e executa ações fora de si mesmo. Antes do MCP, cada ferramenta de IA tinha sua própria forma proprietária de se conectar a fontes de dados. O MCP virou o 'USB' disso: você escreve um servidor MCP uma vez, e qualquer cliente compatível — Cline, Claude Code, Zed, e por aí vai — consegue conversar com ele.

O pulo do gato, e o motivo de eu trazer essa talk pra uma comunidade Go: escrever um servidor MCP é, literalmente, escrever um servidor. Coisa que Go faz muito bem, há muito tempo."

---

## Slide A6 — Mas existe um jeito certo e um jeito ingênuo de fazer isso

**Tela:** Duas caixas lado a lado, ainda sem nome de cena, só a ideia: [Caixa A: "Manda tudo, deixa o modelo se virar"] vs [Caixa B: "Manda só o que importa"]

**Roteiro:**
"E é aqui que a talk de hoje entra de verdade. Vou mostrar dois jeitos de construir um servidor MCP que resolve o MESMO problema — buscar informação num repositório de notas, tipo Obsidian. Um jeito é ingênuo: pega tudo, manda tudo. O outro usa busca semântica de verdade, rodando localmente, sem custar um centavo de nuvem. Vamos primeiro entender as peças que tornam isso possível — incluindo uma ferramenta chamada Kronk, que deixa rodar IA local em Go, com controle fino de CPU e GPU — e só no final vamos ver os dois lado a lado, rodando ao vivo, com o contador de custo na tela."

---

## Slide A7 — Roteiro de hoje

**Tela:** Lista numerada simples:
1. O motor por baixo: Kronk (conceitos — quantização, MoE, híbrido CPU/GPU)
2. Números reais de hoje, nesta máquina
3. Por dentro do código dos dois servidores MCP
4. **Ao vivo:** as duas cenas, lado a lado, com custo na tela

**Roteiro:**
"Esse é o mapa de hoje. Vamos primeiro entender o motor — o Kronk, como ele decide rodar em CPU ou GPU, o que é esse negócio de Mixture of Experts. Depois vamos abrir o código dos dois servidores MCP que eu preparei. E só no final, com tudo isso entendido, a gente vê os dois rodando ao vivo, um contra o outro, com o contador de tokens e custo do Cline na tela. Bora?"

---

## Notas de produção (não falar em voz alta)

- Slide A1/A2: ajustar com o título real do CFP e com o conteúdo de fato apresentado na parte 1, se for outra pessoa/outro contexto, adaptar o "fulano".
- Slide A3: a pausa pra perguntar "quem já usou" funciona bem pra platéia grande; se a sala for pequena/íntima, pode trocar por uma afirmação direta sem pedir resposta.
- Slide A6 é o gancho de suspense — não revelar nomes "Cena 1/Cena 2" ainda, só a ideia conceitual. Os nomes aparecem na seção de demo no final.
- Transição pro próximo bloco (deep dive Kronk): pode ser direta, "Vamos começar pelo motor: o que é o Kronk."
