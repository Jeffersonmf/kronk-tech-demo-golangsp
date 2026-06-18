# TASK-023: Calcular Desconto Progressivo por Quantidade

## Status
- [ ] Prioridade: Média
- [ ] Responsável: Backlog

## Contexto
Para incentivar compras em volume, queremos aplicar um desconto automático conforme a quantidade de unidades de um mesmo item no carrinho.

## Critérios de Aceite (AC)
1. 1 a 2 unidades: sem desconto.
2. 3 a 5 unidades: 5% de desconto.
3. 6 ou mais unidades: 10% de desconto.
4. O desconto é calculado sobre o valor total do item (preço unitário × quantidade).

## Especificações Técnicas (TA)
- Pacote: `fake_project/internal/promotions/`.
- Função `QuantityDiscount(unitPriceCents, quantity int) int` retorna o valor do desconto em centavos.

## Verificação
- Execute `go test ./internal/promotions/...`.
- Caso de teste: 4 unidades de R$ 100,00 → desconto de R$ 20,00 (5% de R$ 400,00).
