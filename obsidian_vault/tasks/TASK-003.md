# TASK-003: Busca de Produtos com Filtros de Categoria e Preço

## Status
- [ ] Prioridade: Média
- [ ] Responsável: Em andamento

## Contexto
A página de catálogo de produtos atualmente só suporta listar todos os produtos. Conforme o catálogo cresce para algumas centenas de SKUs, os clientes precisam de uma forma de filtrar os resultados por categoria e faixa de preço.

## Critérios de Aceite (AC)
1. **Filtro de categoria:** Usuários podem filtrar a lista de produtos por uma ou mais categorias.
2. **Filtro de faixa de preço:** Usuários podem filtrar por um preço mínimo e/ou máximo.
3. **Filtros combinados:** Os filtros de categoria e preço podem ser aplicados juntos.
4. **Resultados vazios:** Se nenhum produto corresponder, retornar uma lista vazia (não um erro).

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/catalog/`.
- Adicionar `SearchProducts(filter catalog.Filter) ([]Product, error)` ao serviço de catálogo.
- Struct `Filter`: `Categories []string`, `MinPrice *int`, `MaxPrice *int`.
- A filtragem deve acontecer em memória por enquanto (sem banco de dados ainda) — iterar sobre o slice de produtos existente.

## Verificação
- Execute `go test ./internal/catalog/...`.
- Verificação manual: filtrar pela categoria "electronics" com `MaxPrice=500` exclui itens de preço mais alto.
