# TASK-019: Calcular Imposto sobre o Pedido

## Status
- [ ] Prioridade: Média
- [ ] Responsável: Backlog

## Contexto
Os preços exibidos não incluem impostos — o time financeiro só descobriu isso quando começou a tentar reconciliar as notas fiscais com o que a plataforma registrava, e percebeu que o checkout da NexCore nunca calculou imposto nenhum, em nenhum lugar. Para o resumo do checkout, precisamos calcular o valor do imposto separadamente, usando uma alíquota fixa para simplificar por enquanto.

## Critérios de Aceite (AC)
1. Aplicar uma alíquota fixa de 10% sobre o subtotal do pedido.
2. O valor do imposto é arredondado para o centavo mais próximo.
3. Subtotal igual a zero resulta em imposto zero.

## Especificações Técnicas (TA)
- Pacote: `fake_project/internal/checkout/`.
- Função `CalculateTax(subtotalCents int) int`.

## Verificação
- Execute `go test ./internal/checkout/...`.
- Caso de teste: subtotal de R$ 100,00 → imposto de R$ 10,00.
