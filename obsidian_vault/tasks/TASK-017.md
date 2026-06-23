# TASK-017: Gerar Slug a partir do Nome do Produto

## Status
- [ ] Prioridade: Baixa
- [ ] Responsável: Backlog

## Contexto
As URLs de produto hoje usam o ID interno (ex: `/produtos/a1b2c3`) — a NexCore nunca pensou em SEO na hora de desenhar as rotas do catálogo. Agora que o time de marketing está retomando investimento em tráfego orgânico (depois de um ano evitando expor a plataforma demais), eles querem URLs amigáveis baseadas no nome do produto (ex: `/produtos/tenis-de-corrida-azul`).

## Critérios de Aceite (AC)
1. Converter o texto para minúsculas e remover acentos.
2. Substituir espaços e caracteres especiais por hífen.
3. Remover hífens duplicados e hífens nas extremidades.

## Especificações Técnicas (TA)
- Pacote: `fake_project/internal/catalog/`.
- Função `Slugify(name string) string`.

## Verificação
- Execute `go test ./internal/catalog/...`.
- Caso de teste: `"Tênis de Corrida Azul!"` → `"tenis-de-corrida-azul"`.
