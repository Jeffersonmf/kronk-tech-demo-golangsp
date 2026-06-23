# TASK-006: Exportação em CSV do Relatório Mensal de Vendas

## Status
- [x] Prioridade: Baixa
- [x] Responsável: Concluído (entregue em 20/05/2026)

## Contexto
Desde a entrega da NexCore, o time financeiro compilava manualmente os totais de vendas mensais, exportando dados de vários dashboards na mão — nenhuma ferramenta de relatório foi entregue, apesar de constar na proposta original. Esta foi uma das primeiras tasks que o time de plataforma conseguiu fechar depois do incidente, como um aceno rápido pro financeiro de que as coisas estavam melhorando: uma exportação única em CSV com os totais de pedidos agrupados por dia, pro fluxo de planilha que eles já usavam.

## Critérios de Aceite (AC)
1. **Intervalo de datas:** O relatório recebe uma data inicial e final e agrega os totais por dia dentro desse intervalo.
2. **Formato CSV:** As colunas de saída são `date,total_orders,total_revenue`.
3. **Moeda:** Os valores de receita são formatados com duas casas decimais, sem símbolo de moeda.
4. **Dias vazios:** Dias sem pedidos ainda devem aparecer na exportação com valores `0`.

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/reports/`.
- `reports.GenerateSalesCSV(start, end time.Time, w io.Writer) error` escreve diretamente no writer fornecido.
- Usar `encoding/csv` da biblioteca padrão — sem necessidade de dependência externa.

## Verificação
- Execute `go test ./internal/reports/...`.
- Verificação manual: exportar maio de 2026 produz 31 linhas, uma por dia.

## Observações
Encerrado no release de maio. Mantido aqui apenas para referência histórica / contexto de backlog grooming.
