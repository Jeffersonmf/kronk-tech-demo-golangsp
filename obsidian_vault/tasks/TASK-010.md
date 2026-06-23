# TASK-010: Gerar Número de Pedido Sequencial

## Status
- [ ] Prioridade: Baixa
- [ ] Responsável: Backlog

## Contexto
Os pedidos hoje só têm o UUID interno que a NexCore usou como identificador único — ótimo para o banco de dados, péssimo para um atendente de suporte ditar por telefone. O time de suporte pediu um número curto e sequencial (ex: `#00042`) pra facilitar o atendimento ao cliente, algo que devia ter sido óbvio desde o desenho original do sistema de pedidos.

## Critérios de Aceite (AC)
1. O primeiro pedido gera o número `#00001`.
2. Números seguintes incrementam em 1, sem repetição.
3. A geração deve ser segura para chamadas concorrentes.

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/orders/`.
- Struct `OrderNumberGenerator` com contador interno protegido (mutex ou `atomic.Int64`).
- Método `Next() string` retornando o número formatado.

## Verificação
- Execute `go test ./internal/orders/... -race`.
- Gerar 100 números a partir de goroutines concorrentes e verificar que são únicos e sequenciais.
