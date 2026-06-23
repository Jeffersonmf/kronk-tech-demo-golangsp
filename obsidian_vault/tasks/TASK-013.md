# TASK-013: Paginar Lista de Produtos

## Status
- [ ] Prioridade: Média
- [ ] Responsável: Em andamento

## Contexto
A listagem de produtos retorna todos os itens de uma vez — outro ponto que a NexCore justificou como "suficiente pro MVP" e nunca revisitou. Com o catálogo real crescendo bem além do que foi testado na entrega, isso já fica visivelmente lento, e a interface nem tem como exibir "página 2", "página 3", etc.

## Critérios de Aceite (AC)
1. A função recebe a lista completa, o número da página (começando em 1) e o tamanho da página.
2. Retorna apenas os itens correspondentes à página solicitada.
3. Uma página além do total de itens retorna uma lista vazia, sem erro.

## Especificações Técnicas (TA)
- Pacote: `fake_project/internal/catalog/`.
- Função `Paginate(products []Product, page, pageSize int) []Product`.

## Verificação
- Execute `go test ./internal/catalog/...`.
- Caso de teste: 10 produtos, `pageSize=3`, `page=2` → itens 4 a 6.
