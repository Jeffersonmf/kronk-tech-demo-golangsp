# Backlog — Plataforma de E-commerce (ZYX Commerce)

## Contexto (leia antes de pegar qualquer task)

Em 2025 a ZYX Commerce contratou a **NexCore Consultoria** pra reescrever a plataforma de vendas do zero, prometendo "arquitetura moderna, escalável, pronta pra Black Friday". Entregaram em cima do prazo, fizeram a festa de encerramento de contrato, e foram embora.

Três semanas depois da entrega, durante a primeira campanha promocional grande, o time de operações começou a receber chamados de clientes reportando pedidos com **estoque negativo** — produtos vendidos que não existiam mais no depósito. Faturamento bloqueado, devolução em massa, e o time de suporte trabalhando no escuro porque a NexCore não documentou praticamente nada além do código em si.

Desde então o time de plataforma (nós, internos, três pessoas) ficou com a missão de **estabilizar o que foi herdado**. Conforme fomos auditando o código, fomos encontrando mais problemas — nada tão grave quanto o do estoque, mas o suficiente pra preencher um backlog inteiro de débito técnico que devia ter sido pego em code review e não foi.

Esse vault é esse backlog. Task por task, na ordem que fomos encontrando os problemas.

## Prioridade máxima — bloqueando o negócio

- **[[TASK-001]]** — Race condition no módulo de inventário (`internal/inventory`). É a causa raiz do incidente de estoque negativo. Confirmado com `go test -race`: múltiplas goroutines conseguem passar pela checagem de estoque disponível antes da subtração ser aplicada. **Ninguém pega outra task até essa estar fechada e validada em produção.**

## Backlog geral — débito técnico herdado da NexCore

As tasks **TASK-002 a TASK-025** são o restante do que a auditoria encontrou: validações que nunca foram implementadas (CPF, CEP, força de senha), cálculos espalhados sem teste (frete, imposto, desconto progressivo, pontos de fidelidade), funcionalidades que o contrato dizia "concluídas" e na prática eram só um TODO com nome bonito (notificação de envio, exportação de relatório, fluxo de reembolso), e um punhado de coisas que qualquer revisor sênior teria pego (paginação que não pagina, busca sem filtro de fato, log vazando número de cartão sem máscara).

Nenhuma delas é urgente como a TASK-001, mas o acúmulo é o motivo pelo qual qualquer mudança nessa base hoje dá medo. Prioridade de cada uma está marcada no campo `Status` do próprio arquivo — segue aproximadamente a ordem em que a auditoria foi encontrando.

## Como trabalhar nesse backlog

- Toda task nova entra com `Responsável: Backlog` até alguém puxar pra sprint atual.
- Task sem `Critérios de Aceite` claros volta pro backlog até ser detalhada — já tivemos problema demais entregando "no escuro" pra confiar de novo nesse estilo.
- `Verificação` deve sempre citar o comando de teste exato, não "testar manualmente". A NexCore testava manualmente. Olha onde isso nos trouxe.
- Se encontrar mais sujeira no código que não está catalogada aqui, abre uma task nova — não conserta "de passagem" sem registrar, senão perdemos o histórico de novo.
