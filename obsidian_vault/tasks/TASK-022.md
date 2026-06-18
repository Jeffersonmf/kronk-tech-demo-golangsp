# TASK-022: Gerar Código de Rastreio Fake para Testes

## Status
- [ ] Prioridade: Baixa
- [ ] Responsável: Backlog

## Contexto
Os testes de integração do fluxo de envio dependem de chamar a API real da transportadora para obter um código de rastreio. Precisamos de um gerador local para usar em ambientes de teste.

## Critérios de Aceite (AC)
1. Gerar um código no formato `BR` + 9 dígitos + `BR` (padrão dos Correios), ex: `BR123456789BR`.
2. Os 9 dígitos devem ser aleatórios a cada chamada.
3. O código gerado deve sempre ter exatamente 13 caracteres.

## Especificações Técnicas (TA)
- Pacote: `fake_project/internal/shipping/`.
- Função `GenerateFakeTrackingCode() string` usando `math/rand`.

## Verificação
- Execute `go test ./internal/shipping/...`.
- Validar o formato com regex e o tamanho fixo de 13 caracteres.
