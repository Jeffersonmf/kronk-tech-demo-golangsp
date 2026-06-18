# TASK-004: Notificação por E-mail quando o Pedido for Enviado

## Status
- [ ] Prioridade: Baixa
- [ ] Responsável: Backlog

## Contexto
Atualmente os clientes não têm visibilidade de quando o pedido é enviado — o suporte recebe tickets repetidos perguntando "cadê meu pedido?". O time de produto quer uma notificação simples por e-mail disparada no momento em que o status do pedido muda para "enviado".

## Critérios de Aceite (AC)
1. **Disparo:** Quando o status de um pedido muda para `shipped`, um evento de notificação é emitido.
2. **Conteúdo:** O e-mail inclui o ID do pedido, a transportadora e o código de rastreio (se disponível).
3. **Idempotência:** Reprocessar o mesmo evento de envio não deve gerar e-mails duplicados.
4. **Tratamento de falha:** Se o provedor de e-mail estiver indisponível, o evento deve ser reprocessado depois, não descartado.

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/notifications/`.
- Conectar à transição de status do pedido em `fake_project/internal/orders/` (ainda não criado) via um padrão de evento/observer.
- Interface sugerida: `notifications.NotifyShipped(orderID string, tracking TrackingInfo) error`.
- Para esta demo, o "envio" do e-mail pode ser um stub/log — a integração SMTP real está fora do escopo.

## Verificação
- Execute `go test ./internal/notifications/...`.
- Verificação manual: marcar um pedido como enviado duas vezes deve registrar apenas uma notificação no log.
