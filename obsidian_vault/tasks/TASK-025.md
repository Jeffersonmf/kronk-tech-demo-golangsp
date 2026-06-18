# TASK-025: Endpoint de Health Check

## Status
- [ ] Prioridade: Baixa
- [ ] Responsável: Backlog

## Contexto
A equipe de infraestrutura precisa de um endpoint simples para que o load balancer verifique se o serviço está no ar antes de liberar tráfego para uma nova instância.

## Critérios de Aceite (AC)
1. `GET /health` retorna status HTTP 200.
2. O corpo da resposta é o JSON `{"status":"ok"}`.
3. O endpoint não depende de banco de dados ou serviços externos.

## Especificações Técnicas (TA)
- Novo pacote: `fake_project/internal/health/`.
- Handler `HealthCheckHandler(w http.ResponseWriter, r *http.Request)`.
- Registrar a rota em `fake_project/cmd/server/` (diretório já existe, hoje vazio).

## Verificação
- Execute `go test ./internal/health/...`.
- Requisição para `/health` retorna status 200 e o corpo esperado.
