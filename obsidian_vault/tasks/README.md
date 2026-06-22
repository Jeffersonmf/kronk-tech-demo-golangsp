# Vault de Tasks — propósito

Este vault simula o backlog de um pequeno time de e-commerce no Obsidian, usado como "fonte da verdade" para a demo do GolangSP. São 25 tasks (TASK-001 a TASK-025), todas no mesmo formato: **Status / Contexto / Critérios de Aceite (AC) / Especificações Técnicas (TA) / Verificação**.

## A task principal

**TASK-001** é a "carro-chefe" da demo: descreve uma race condition real no inventário, em `fake_project/internal/inventory/service.go`, verificada com `go test -race`. É a única task que o agente de IA precisa de fato resolver ao vivo — é a mais densa do vault e a que está conectada a código real no repositório.

## As demais 24 (ruído de backlog)

**TASK-002 a TASK-025** são tarefas Go pequenas e autocontidas, cada uma cabendo numa única função/arquivo novo. Qualquer uma pode ser implementada do zero ao verde (`go test ./...`) em poucos minutos — boas candidatas extras pra implementar ao vivo, mas não fazem parte do roteiro principal da demo.

| Task | Título |
|---|---|
| TASK-002 | Validação de Cupom de Desconto no Checkout |
| TASK-003 | Busca de Produtos com Filtros de Categoria e Preço |
| TASK-004 | Notificação por E-mail quando o Pedido for Enviado |
| TASK-005 | Migrar Autenticação para OAuth2 |
| TASK-006 | Exportação em CSV do Relatório Mensal de Vendas |
| TASK-007 | Fluxo de Reembolso para Itens Devolvidos |
| TASK-008 | Validar CPF no Cadastro de Cliente |
| TASK-009 | Calcular Frete por Peso |
| TASK-010 | Gerar Número de Pedido Sequencial |
| TASK-011 | Formatar Valores Monetários em Reais |
| TASK-012 | Calcular Média de Avaliações de Produto |
| TASK-013 | Paginar Lista de Produtos |
| TASK-014 | Ordenar Produtos por Preço |
| TASK-015 | Mascarar Número de Cartão de Crédito em Logs |
| TASK-016 | Calcular Pontos de Fidelidade por Compra |
| TASK-017 | Gerar Slug a partir do Nome do Produto |
| TASK-018 | Validar e Normalizar CEP |
| TASK-019 | Calcular Imposto sobre o Pedido |
| TASK-020 | Remover Itens Duplicados do Carrinho |
| TASK-021 | Verificar Estoque em Múltiplos Depósitos |
| TASK-022 | Gerar Código de Rastreio Fake para Testes |
| TASK-023 | Calcular Desconto Progressivo por Quantidade |
| TASK-024 | Validar Força de Senha no Cadastro |
| TASK-025 | Endpoint de Health Check |

## Por que o vault tem esse tamanho

A demo compara duas formas de um agente de IA consumir este vault como contexto:

- **MCP Tradicional (Cena 1)** carrega o vault *inteiro* como contexto bruto — o custo em tokens escala com o tamanho do backlog (25 tasks).
- **Kronk + RAG (Cena 2)** recupera apenas o(s) documento(s) *semanticamente relevantes* para a tarefa pedida — idealmente 1-2 tasks, não as 25.

Com poucas tasks, as duas abordagens leriam praticamente o mesmo conteúdo e a diferença de custo não seria visível. Com 25, o contraste de tokens/custo entre "ler tudo" e "recuperar só o relevante" aparece de forma clara na comparação ao vivo — e, como cada task de ruído é pequena e independente, qualquer uma serve como tarefa real de implementação caso a demo precise de uma segunda rodada.
