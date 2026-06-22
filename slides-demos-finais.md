# Roteiro — Demos Finais: Cena 1 vs Cena 2

Último bloco da talk. Numeração continua de `slides-kronk-deepdive.md` (que terminou no slide 14). Pré-requisito mental do público nesse ponto: já sabem o que é Kronk, MoE, híbrido CPU/GPU, e já viram o código dos dois servidores MCP. Agora é a hora da prova real.

**Setup técnico que precisa estar pronto antes de subir esse bloco** (ver notas de produção no final): Cline configurado com dois providers prontos pra trocar — o provider de nuvem paga ("Cline", conta com créditos, pra mostrar custo real em $) e o provider local (Kronk, `http://localhost:11435/v1`, custo $0). Servidores `mcp_vault` (`vault_dumb` :9001, `vault_smart` :9002) rodando.

---

## Slide 15 — Chegou a hora

**Tela:** Texto simples: "Mesma pergunta. Mesmo agente. Duas formas de buscar contexto." + os nomes revelados agora: **Cena 1** e **Cena 2**.

**Roteiro:**
"Lembram do gancho do começo — 'tem um jeito ingênuo e um jeito inteligente de fazer isso'? Chegou a hora de ver os dois, ao vivo, na mesma tela. Vou rodar exatamente a mesma tarefa duas vezes: pedir pro Cline ler uma tarefa do nosso vault de notas e corrigir um bug de concorrência num código Go real. A única coisa que muda entre as duas rodadas é qual servidor MCP está respondendo, e qual modelo está por trás. Chamo a primeira de Cena 1, a segunda de Cena 2."

---

## Slide 16 — Apresentando a Cena 1

**Tela:** print da tela de configuração do Cline mostrando: Provider = conta paga na nuvem, MCP ativo = `vault_dumb`.

**Roteiro:**
"Cena 1 representa o jeito mais comum de montar isso hoje: um modelo na nuvem — pago, por token — conectado a um servidor MCP que não filtra nada. É o que a maioria dos times monta quando está com pressa: 'deixa o modelo decidir o que é relevante, a gente só manda tudo'. Vou deixar o contador de tokens e o contador de custo do Cline visíveis na tela o tempo todo. Olhem esses números antes de eu começar: zero."

---

## Slide 17 — Cena 1 ao vivo

**Tela:** A tela real do VS Code/Cline, sem slide — é a hora da demo de fato. Ter o prompt já preparado pra colar (ver `prompt1.txt`/`prompt2.txt`/`prompt3.txt` no repo, ou o prompt da TASK-001 já validado).

**Roteiro (falar enquanto a chamada da ferramenta acontece):**
"Vou colar o prompt: 'Leia a TASK-001 no vault do Obsidian e resolva o problema descrito nela.' [colar e enviar] Reparem: o Cline está decidindo chamar a ferramenta `read_vault`. Não fui eu que disse pra ele chamar essa ferramenta especificamente — ele só tem uma ferramenta disponível nessa cena, então é a única opção. [esperar a resposta da ferramenta aparecer] Olha o contador de tokens subindo... e o custo junto. Essa resposta sozinha, só de buscar contexto, já consumiu [X] mil tokens."

**Contingência (ler antes, não falar pro público):** Se a geração do modelo de nuvem demorar ou travar em algum tool-call subsequente, o ponto importante já foi feito no instante em que a ferramenta `read_vault` devolve a resposta — é aí que o contador salta. Pode pausar a narrativa aí, comentar o número, e seguir pra Cena 2 sem esperar o fix terminar de verdade nesta cena.

---

## Slide 18 — O número da Cena 1

**Tela:** Print/captura de tela real do contador do Cline ao final dessa rodada — tokens e custo em destaque, fonte grande.

**Roteiro:**
"Aqui está o resultado da Cena 1: [ler os números reais capturados]. Token e dinheiro de verdade, gastos só porque a ferramenta de busca não soube filtrar nada. Guardem esse número na cabeça. Agora vamos trocar duas coisas — só duas — e ver o que muda."

---

## Slide 19 — Apresentando a Cena 2

**Tela:** print da configuração trocada: Provider = Kronk local (`http://localhost:11435/v1`), MCP ativo = `vault_smart`.

**Roteiro:**
"Troquei o provider do Cline pra apontar pro Kronk, rodando aqui nesta máquina, neste momento — sem internet envolvida nessa chamada. E troquei o MCP ativo pra `vault_smart`, o servidor com busca semântica que abrimos o código mais cedo. Mesmo prompt, exatamente o mesmo texto que usei na Cena 1. Vamos ver."

---

## Slide 20 — Cena 2 ao vivo

**Tela:** Tela real do VS Code/Cline de novo.

**Roteiro (falar enquanto a chamada da ferramenta acontece):**
"[colar o mesmo prompt e enviar] Reparem que agora a ferramenta disponível é outra: `search_vault`. O modelo vai gerar a query de busca e chamar essa ferramenta. [esperar resposta] E voltou só com os documentos relevantes — a TASK-001, possivelmente mais um ou dois relacionados — em vez do vault inteiro. Olhem o contador: [ler o número]. E o custo: zero. Porque está rodando local."

**Contingência (não falar pro público):** O modelo local (`gpt-oss-20b-Q8_0`, MoE, rodando via Metal no M4 Pro) com raciocínio em `medium`/`high` pode demorar mais que uma resposta de nuvem pra completar a correção de código de ponta a ponta — nos testes de hoje, um prompt do tamanho real do Cline (~2k tokens de contexto) já leva ~2s só de TTFT, e o raciocínio + tool calls subsequentes somam tempo. O ponto da demo é o contraste de tokens/custo na chamada da ferramenta, que aparece rápido. Se o fix completo não terminar em tempo de palco, está OK comentar isso abertamente ("e aqui ele continua trabalhando na correção, localmente, sem custar nada — vamos deixar rodando e voltar no final") e seguir pro fechamento, ou cortar pra um clipe pré-gravado do resultado final (`go test -race` passando) se for vídeo gravado.

---

## Slide 21 — Lado a lado: o contraste final

**Tela:** Os dois prints (Slide 18 e o número real da Cena 2) um do lado do outro, com uma seta ou X grande mostrando a diferença.

| | Cena 1 (nuvem + dumb) | Cena 2 (local + smart) |
|---|---|---|
| Tokens consumidos | [número real capturado] | [número real capturado] |
| Custo | [$ real capturado] | $0,00 |
| Onde rodou | Servidor de terceiro | Esta máquina, agora |
| Internet necessária | Sim | Não |

**Roteiro:**
"Essa é a imagem que eu queria que vocês levassem pra casa. Mesma tarefa, mesmo agente, mesmo código pra corrigir. A diferença inteira está em como a informação chega até o modelo, e onde o modelo roda. E o motor que tornou a Cena 2 possível, rodando aqui, em Go, sem nuvem — foi tudo que vimos no começo da talk: Kronk, com offload de GPU bem ajustado, cache de contexto pensado pra agentes, e um servidor MCP que faz busca de verdade em vez de só empilhar tudo."

---

## Slide 22 — Fechamento

**Tela:** Slide final simples: "Go não é só a linguagem do seu backend. É a linguagem da sua infraestrutura de IA também." + links/contato/repo do GitHub.

**Roteiro:**
"Pra fechar: tudo que vocês viram hoje — o servidor MCP, o motor de inferência, a busca vetorial — é Go, de ponta a ponta. Não precisamos de Python, não precisamos de nuvem pra ter um agente de IA útil. O código está aberto no repositório [mostrar link]. Obrigado!"

---

## Notas de produção (não falar em voz alta)

- **Pré-requisito crítico:** ensaiar a troca de provider e de MCP ativo no Cline ANTES de subir no palco — são poucos cliques (Settings → API Provider, e toggle `disabled` em `cline_mcp_settings.json` ou na UI de MCP Servers), mas trocar ao vivo sem ensaio é risco de travar a narrativa.
- **Os números dos slides 18/21 são placeholders** — substituir pelos valores reais capturados num ensaio completo antes da talk. Não inventar números.
- **Confirmado:** a Cena 1 usa a conta paga na nuvem (API key "GolangSP" criada hoje, com créditos) — custo real em $ vs $0 da Cena 2. Cuidado pra não deixar a tela exibir a API key em nenhum momento (cobrir esse campo nas configurações antes de compartilhar tela, ou trocar a key depois da talk já que ela foi colada em texto puro numa sessão de chat antes).
- Considerar gravar essa seção de demo previamente (igual o restante da execução local), dado o que vimos hoje sobre tempo de resposta do modelo local — reduz risco de travar a apresentação ao vivo.
