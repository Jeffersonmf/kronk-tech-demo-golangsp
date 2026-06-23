# TASK-018: Validar e Normalizar CEP

## Status
- [ ] Prioridade: Média
- [ ] Responsável: Backlog

## Contexto
O formulário de endereço de entrega aceita o CEP em qualquer formato — outra validação que a NexCore simplesmente não escreveu — o que causa falhas silenciosas no cálculo de frete, que espera o padrão `00000-000` e quebra com qualquer variação.

## Critérios de Aceite (AC)
1. Aceitar CEP com ou sem hífen (`01310100` ou `01310-100`).
2. Normalizar a saída sempre para o formato `00000-000`.
3. CEPs com tamanho incorreto ou caracteres não numéricos retornam erro.

## Especificações Técnicas (TA)
- Pacote: `fake_project/internal/shipping/`.
- Função `NormalizeCEP(cep string) (string, error)`.

## Verificação
- Execute `go test ./internal/shipping/...`.
- Caso de teste: `"01310100"` → `"01310-100"`.
