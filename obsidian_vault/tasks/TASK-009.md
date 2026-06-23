# TASK-009: Calcular Frete por Peso

## Status
- [ ] Prioridade: Média
- [ ] Responsável: Backlog

## Contexto
O frete hoje é um valor fixo de R$ 20,00 para qualquer pedido, independente do peso — um comentário no código da NexCore chama isso de "placeholder pro MVP, substituir antes do lançamento". O lançamento aconteceu há mais de um ano. Precisamos calcular o frete de verdade, com base no peso total dos itens do carrinho.

## Critérios de Aceite (AC)
1. Pedidos de até 1kg custam R$ 10,00.
2. Acima de 1kg, soma-se R$ 5,00 por kg adicional, arredondado para cima.
3. Peso zero ou negativo retorna erro.

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/shipping/`.
- Função `CalculateShippingCost(weightKg float64) (int, error)` retorna o valor em centavos.

## Verificação
- Execute `go test ./internal/shipping/...`.
- Caso de teste: 2.5kg deve custar R$ 25,00 (2500 centavos).
