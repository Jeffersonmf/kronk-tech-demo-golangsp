# TASK-015: Mascarar Número de Cartão de Crédito em Logs

## Status
- [ ] Prioridade: Alta
- [ ] Responsável: Backlog

## Contexto
A auditoria de segurança que abrimos depois do incidente de estoque encontrou outro problema sério, sem relação com o bug original: números completos de cartão de crédito aparecem em logs de depuração do checkout que a NexCore deixou ativados em produção. Risco de compliance (PCI-DSS) real — precisa ser tratado com a mesma prioridade da TASK-001.

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
