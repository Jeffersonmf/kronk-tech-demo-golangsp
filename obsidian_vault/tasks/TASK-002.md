# TASK-002: Adicionar Validação de Cupom de Desconto no Checkout

## Status
- [ ] Prioridade: Média
- [ ] Responsável: Backlog

## Contexto
O time de marketing quer começar a rodar campanhas promocionais com códigos de desconto (ex: `SUMMER10`, `WELCOME20`). Atualmente o fluxo de checkout não tem o conceito de cupom — os totais são calculados apenas a partir do preço e quantidade dos produtos.

## Critérios de Aceite (AC)
1. **Busca do código:** O sistema deve validar um código de cupom contra uma lista de promoções ativas.
2. **Expiração:** Cupons expirados ou ainda não ativos devem ser rejeitados com um erro claro.
3. **Aplicação do desconto:** Cupons válidos aplicam um desconto percentual ou de valor fixo sobre o subtotal do pedido.
4. **Uso único por pedido:** Apenas um código de cupom pode ser aplicado por pedido.

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/promotions/`.
- A struct de Cupom deve incluir `Code`, `DiscountType` (percentage|fixed), `Value`, `ValidFrom`, `ValidUntil`.
- O fluxo de checkout (`fake_project/internal/checkout/`) chama `promotions.Validate(code)` antes de calcular o total final.
- Retornar erros dedicados `ErrInvalidCoupon` / `ErrCouponExpired` para tratamento claro de erros.

## Verificação
- Execute `go test ./internal/promotions/...`.
- Verificação manual: aplicar `SUMMER10` em um carrinho de $100 deve resultar em um total de $90.
