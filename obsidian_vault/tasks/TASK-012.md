# TASK-012: Calcular Média de Avaliações de Produto

## Status
- [ ] Prioridade: Baixa
- [ ] Responsável: Backlog

## Contexto
A funcionalidade de avaliações que a NexCore entregou lista cada review individualmente, mas nunca calculou uma nota média — metade do recurso, sem o resumo que qualquer cliente espera ver de cara (ex: "4.5 de 5 estrelas"). Foi sinalizado na auditoria como "funcionalidade incompleta apresentada como concluída".

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
