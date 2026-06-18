# TASK-016: Calcular Pontos de Fidelidade por Compra

## Status
- [ ] Prioridade: Baixa
- [ ] Responsável: Backlog

## Contexto
O novo programa de fidelidade dá pontos aos clientes proporcionalmente ao valor gasto, para troca por descontos futuros.

## Critérios de Aceite (AC)
1. O cliente ganha 1 ponto para cada R$ 10,00 gastos, arredondado para baixo.
2. Compras abaixo de R$ 10,00 geram 0 pontos.
3. Valores negativos retornam erro.

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/loyalty/`.
- Função `CalculatePoints(totalCents int) (int, error)`.

## Verificação
- Execute `go test ./internal/loyalty/...`.
- Caso de teste: R$ 95,00 (9500 centavos) → 9 pontos.
