# TASK-015: Mascarar Número de Cartão de Crédito em Logs

## Status
- [ ] Prioridade: Alta
- [ ] Responsável: Backlog

## Contexto
Uma auditoria de segurança identificou que números completos de cartão de crédito aparecem em logs de depuração do checkout — risco de compliance (PCI-DSS).

## Critérios de Aceite (AC)
1. Apenas os últimos 4 dígitos do cartão permanecem visíveis.
2. Os demais dígitos são substituídos por `*`.
3. Entradas com menos de 4 dígitos retornam erro de entrada inválida.

## Especificações Técnicas (TA)
- Pacote: `fake_project/internal/payments/`.
- Função `MaskCardNumber(number string) (string, error)`.

## Verificação
- Execute `go test ./internal/payments/...`.
- Caso de teste: `"4111111111111234"` → `"************1234"`.
