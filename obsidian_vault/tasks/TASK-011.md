# TASK-011: Formatar Valores Monetários em Reais

## Status
- [ ] Prioridade: Baixa
- [ ] Responsável: Backlog

## Contexto
Os preços são armazenados como inteiros em centavos (ex: `1050` = R$ 10,50), mas hoje são exibidos sem formatação nas telas e nos recibos.

## Critérios de Aceite (AC)
1. Converter centavos para o formato `R$ 10,50`.
2. Valores negativos exibem o sinal antes do `R$` (ex: `-R$ 5,00`).
3. Valores acima de mil usam separador de milhar (ex: `R$ 1.234,56`).

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/format/`.
- Função `FormatBRL(cents int) string`.

## Verificação
- Execute `go test ./internal/format/...`.
- Casos: `1050` → `"R$ 10,50"`; `123456` → `"R$ 1.234,56"`; `-500` → `"-R$ 5,00"`.
