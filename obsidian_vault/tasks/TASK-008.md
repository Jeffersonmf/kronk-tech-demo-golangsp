# TASK-008: Validar CPF no Cadastro de Cliente

## Status
- [ ] Prioridade: Média
- [ ] Responsável: Backlog

## Contexto
Na auditoria pós-incidente encontramos que o cadastro de clientes aceita literalmente qualquer string no campo CPF — a NexCore nunca implementou validação nenhuma, nem de formato. Isso já causou falhas na integração com a transportadora quando o documento cadastrado é inválido, e é o tipo de corte de canto que só aparece quando alguém finalmente vai ler o código com atenção.

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
