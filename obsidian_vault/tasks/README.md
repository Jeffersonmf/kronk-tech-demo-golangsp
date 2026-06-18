# Vault de Tasks — propósito

Este vault simula o backlog de um pequeno time de e-commerce no Obsidian, usado como "fonte da verdade" para a demo do GolangSP. São 25 tasks (TASK-001 a TASK-025), todas no mesmo formato (Status / Contexto / Critérios de Aceite (AC) / Especificações Técnicas (TA) / Verificação).

- **TASK-001** é a tarefa "carro-chefe": descreve a race condition de inventário em `fake_project/internal/inventory/service.go`, verificada com `go test -race`. É a mais densa do vault.
- **TASK-002 a TASK-025** são tarefas Go pequenas e autocontidas, cada uma cabendo numa única função/arquivo novo (validação de CPF/CEP, formatação de moeda, paginação, cálculo de frete/desconto/imposto, health check, etc.). Qualquer uma delas pode ser escolhida e implementada do zero ao verde (`go test ./...`) em poucos minutos — boas candidatas para o agente de IA implementar ao vivo na demo.

**Por que o vault tem esse tamanho:** a demo compara duas formas de um agente de IA consumir este vault como contexto:
- MCP Tradicional (Cena 1) carrega o vault *inteiro* como contexto bruto — o custo em tokens escala com o tamanho do backlog (25 tasks).
- Kronk + RAG (Cena 2) recupera apenas o(s) documento(s) *semanticamente relevantes* para a tarefa pedida — idealmente 1-2 tasks, não as 25.

Com poucas tasks, as duas abordagens leriam praticamente o mesmo conteúdo e a diferença de custo não seria visível. Com 25, o contraste de tokens/custo entre "ler tudo" e "recuperar só o relevante" aparece de forma clara na comparação ao vivo — e, como cada task é pequena e independente, qualquer uma serve como tarefa real de implementação na demo.
