# TASK-008: Validar CPF no Cadastro de Cliente

## Status
- [ ] Prioridade: Média
- [ ] Responsável: Backlog

## Contexto
O cadastro de clientes aceita qualquer string no campo CPF, o que tem causado falhas na integração com a transportadora quando o documento é inválido.

## Critérios de Aceite (AC)
1. CPF deve ter 11 dígitos numéricos.
2. Os dígitos verificadores devem ser validados conforme o algoritmo oficial.
3. CPFs com todos os dígitos iguais (ex: `00000000000`) devem ser rejeitados.

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/customers/`.
- Função `ValidateCPF(cpf string) error`.
- Retornar `ErrInvalidCPF` quando inválido.

## Verificação
- Execute `go test ./internal/customers/...`.
- Testar com pelo menos um CPF válido conhecido e variações inválidas (tamanho errado, dígitos repetidos, dígito verificador incorreto).
