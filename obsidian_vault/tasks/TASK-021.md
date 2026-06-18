# TASK-021: Verificar Estoque em Múltiplos Depósitos

## Status
- [ ] Prioridade: Média
- [ ] Responsável: Backlog

## Contexto
A empresa abriu um segundo centro de distribuição. O serviço de inventário hoje só lida com o estoque de um único depósito; precisamos somar o estoque de um produto considerando todos os depósitos.

## Critérios de Aceite (AC)
1. Receber o estoque por depósito e retornar o total somado.
2. Um depósito sem registro para o produto conta como `0`.
3. Mapa de depósitos vazio retorna `0`.

## Especificações Técnicas (TA)
- Pacote: `fake_project/internal/inventory/` (extensão de `service.go`).
- Função `TotalStockAcrossWarehouses(stockByWarehouse map[string]int) int`.

## Verificação
- Execute `go test ./internal/inventory/...`.
