# TASK-012: Calcular Média de Avaliações de Produto

## Status
- [ ] Prioridade: Baixa
- [ ] Responsável: Backlog

## Contexto
A página de produto lista cada avaliação individualmente, mas não mostra uma nota média — clientes pedem um resumo rápido (ex: "4.5 de 5 estrelas").

## Critérios de Aceite (AC)
1. Calcular a média aritmética das notas (1 a 5) de um produto.
2. Produto sem avaliações retorna média `0`, sem erro.
3. A média é arredondada para 1 casa decimal.

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/reviews/`.
- Função `AverageRating(ratings []int) float64`.

## Verificação
- Execute `go test ./internal/reviews/...`.
- Casos: `[5,4,3]` → `4.0`; `[]` → `0.0`; `[5,4]` → `4.5`.
