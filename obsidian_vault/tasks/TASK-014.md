# TASK-014: Ordenar Produtos por Preço

## Status
- [ ] Prioridade: Baixa
- [ ] Responsável: Backlog

## Contexto
A listagem do catálogo segue a ordem de inserção dos produtos. Queremos permitir que o cliente ordene os resultados por preço, do menor para o maior ou vice-versa.

## Critérios de Aceite (AC)
1. Suporta ordenação crescente e decrescente por preço.
2. Produtos com o mesmo preço mantêm a ordem relativa original (ordenação estável).
3. Lista vazia ou com um único item não gera erro.

## Especificações Técnicas (TA)
- Pacote: `fake_project/internal/catalog/`.
- Função `SortByPrice(products []Product, ascending bool)` usando `sort.SliceStable`.

## Verificação
- Execute `go test ./internal/catalog/...`.
