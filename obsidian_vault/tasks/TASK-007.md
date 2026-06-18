# TASK-007: Fluxo de Reembolso para Itens Devolvidos

## Status
- [ ] Prioridade: Média
- [ ] Responsável: Backlog

## Contexto
Atualmente os clientes podem solicitar uma devolução, mas não há forma automatizada de processar o reembolso quando o item devolvido chega de volta ao armazém. Os reembolsos são processados manualmente pelo suporte, o que é lento e propenso a erros.

## Critérios de Aceite (AC)
1. **Devolução recebida:** Quando uma devolução é marcada como "recebida", uma solicitação de reembolso é criada automaticamente.
2. **Reembolsos parciais:** O valor do reembolso deve suportar valores parciais do pedido (ex: um item de três).
3. **Reposição de estoque:** Quando o reembolso é processado, o estoque do item devolvido deve ser reposto.
4. **Trilha de auditoria:** Todo reembolso registra o ID do pedido original, o valor e o motivo.

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/returns/`.
- Ao receber a devolução, chamar `inventory.RestockProduct(productID, quantity)` (extensão de `fake_project/internal/inventory/service.go`) e `payments.Refund(orderID, amount)`.
- Struct de Reembolso: `OrderID`, `Amount`, `Reason`, `ProcessedAt`.

## Verificação
- Execute `go test ./internal/returns/...`.
- Verificação manual: devolver 1 de 3 unidades reembolsa 1/3 do total do item e aumenta o estoque em 1.
