# TASK-025: Endpoint de Health Check

## Status
- [ ] Prioridade: Baixa
- [ ] Responsável: Backlog

## Contexto
A NexCore nunca deixou nem o básico de observabilidade pronto — sem health check, a infraestrutura não tinha como saber com confiança se uma instância nova estava de fato no ar antes de liberar tráfego pra ela. O time de infra está retomando esses fundamentos agora que o incêndio principal foi apagado.

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
