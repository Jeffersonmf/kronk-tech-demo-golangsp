# TASK-024: Validar Força de Senha no Cadastro

## Status
- [ ] Prioridade: Alta
- [ ] Responsável: Backlog

## Contexto
O cadastro de usuários hoje aceita senhas como `"123"` ou `"abc"`, o que foi apontado como risco de segurança em uma auditoria recente.

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
