# TASK-020: Remover Itens Duplicados do Carrinho

## Status
- [ ] Prioridade: Média
- [ ] Responsável: Backlog

## Contexto
Um bug no frontend permite que o mesmo produto seja adicionado duas vezes como itens separados no carrinho, em vez de aumentar a quantidade de um único item.

## Critérios de Aceite (AC)
1. Itens com o mesmo `ProductID` devem ser mesclados em um único item.
2. A quantidade do item mesclado é a soma das quantidades originais.
3. A ordem dos itens restantes segue a primeira ocorrência de cada produto.

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/cart/`.
- Função `MergeDuplicateItems(items []CartItem) []CartItem`.

## Verificação
- Execute `go test ./internal/cart/...`.
