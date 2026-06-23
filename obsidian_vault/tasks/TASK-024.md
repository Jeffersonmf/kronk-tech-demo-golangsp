# TASK-024: Validar Força de Senha no Cadastro

## Status
- [ ] Prioridade: Alta
- [ ] Responsável: Backlog

## Contexto
A mesma auditoria de segurança que encontrou o vazamento de número de cartão (TASK-015) também encontrou isto: o cadastro de usuários aceita senhas como `"123"` ou `"abc"` sem nenhuma validação de força. Mais um item que devia ter sido coberto desde o primeiro dia e a NexCore deixou passar.

## Critérios de Aceite (AC)
1. A senha deve ter no mínimo 8 caracteres.
2. Deve conter ao menos uma letra e um número.
3. Senhas que não atendem aos critérios retornam um erro descrevendo o motivo.

## Especificações Técnicas (TA)
- Pacote: `fake_project/internal/auth/`.
- Função `ValidatePasswordStrength(password string) error`.

## Verificação
- Execute `go test ./internal/auth/...`.
- Testar senhas válidas e inválidas (curta, somente números, somente letras).
