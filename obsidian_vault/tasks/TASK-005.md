# TASK-005: Migrar Autenticação para OAuth2

## Status
- [ ] Prioridade: Alta
- [ ] Responsável: Backlog

## Contexto
A plataforma usa o login caseiro de usuário/senha + token de sessão que a NexCore implementou — funcional, mas sem nenhuma opção de SSO, porque "não estava no MVP". Agora que a ZYX está voltando a fechar contratos corporativos, vários clientes grandes estão pedindo "Login com Google" / "Login com Microsoft" para que seus funcionários possam usar SSO ao fazer pedidos em volume. É a primeira vez desde o incidente que estamos construindo algo novo em vez de só consertar herança.

## Critérios de Aceite (AC)
1. **Login OAuth2:** Usuários podem se autenticar via pelo menos um provedor OAuth2 externo (Google).
2. **Vinculação de conta:** Um login OAuth2 pode ser vinculado a uma conta existente baseada em e-mail.
3. **Fallback:** O login por usuário/senha existente deve continuar funcionando durante a migração.
4. **Expiração de token:** Os tokens OAuth2 são renovados de forma transparente antes de expirar.

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/auth/`.
- Usar `golang.org/x/oauth2` para o fluxo do provedor.
- Armazenar os tokens do provedor separadamente do armazenamento de sessão existente; não sobrescrever as credenciais legadas.
- Adicionar `auth.ExchangeCode(provider, code string) (Session, error)`.

## Verificação
- Execute `go test ./internal/auth/...`.
- Verificação manual: um usuário que já se registrou com e-mail/senha consegue vincular uma conta Google e fazer login por qualquer um dos dois métodos depois.
