# TASK-001: Corrigir Race Condition no Inventário

## Status
- [ ] Prioridade: Alta
- [ ] Responsável: Agente de IA

## Contexto
Esse é o incidente que deu origem a todo o backlog deste vault (ver `obsidian_vault/tasks/README.md`): a plataforma reescrita pela NexCore foi pra produção sem testar concorrência de verdade, e durante a primeira campanha promocional grande depois da entrega, vendemos produtos que não existiam mais no estoque — pedidos com **estoque negativo** sob alto tráfego (promoções/flash sales).

## Critérios de Aceite (AC)
1.  **Atomicidade:** A redução de estoque deve ser atômica.
2.  **Validação:** O sistema deve verificar se `estoque_disponivel >= quantidade_solicitada` antes de deduzir.
3.  **Sem Negativos:** O estoque do produto nunca deve ficar abaixo de zero.
4.  **Tratamento de Erro:** Retornar um erro claro `ErrInsufficientStock` (estoque insuficiente) quando a validação falhar.

## Especificações Técnicas (TA)
- A lógica de inventário está em `fake_project/internal/inventory/service.go`.
- Use primitivas de concorrência do Go (mutexes ou channels) ou operações atômicas.
- A implementação atual usa uma simples verificação `if` sem lock.
- Caso de teste para verificação: `fake_project/internal/inventory/service_test.go`.

## Verificação
- Execute `go test ./internal/inventory/...`
- O teste deve passar com a flag `-race` habilitada.

## Nota da Demo
Esta é a **única** tarefa neste vault que o agente de IA precisa de fato resolver durante a demo ao vivo. As tarefas TASK-002 a TASK-025 são débito técnico herdado da NexCore, mantidas neste vault como "ruído" realista — veja `obsidian_vault/tasks/README.md` para entender por que existem. Um agente baseado em retrieval (Kronk + RAG) deve trazer esta tarefa (e `fake_project/internal/inventory/service.go`) como contexto relevante e ignorar o restante; um agente que carrega o vault inteiro como contexto bruto vai trazer todas elas.
