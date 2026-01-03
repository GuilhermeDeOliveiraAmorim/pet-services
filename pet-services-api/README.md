# Pet Services API

Backend Go (DDD) para marketplace de serviços pet, com autenticação JWT, gestão de prestadores, solicitações e avaliações.

## Visão Geral

- Linguagem: Go 1.25.5
- Dependências: github.com/google/uuid
- Estilo: DDD, casos de uso isolados na camada `internal/application`, regras de domínio em `internal/domain`
- Autenticação: JWT (access + refresh), senhas com hasher plugável (ex.: bcrypt)
- Observabilidade: logging estruturado com `log/slog` via helper `internal/application/logging`

## Domínios e Casos de Uso

- Auth: signup, login, refresh, logout
- User: perfil (get/update), trocar senha, reset de senha, verificação de email (request/confirm), delete account
- Provider: cadastro, atualizar perfil, serviços (add/update/remove), disponibilidade (semanal e por dia), fotos (add/remove/reordenar), ativar/desativar, moderação (aprovar/rejeitar), listar por localização
- Request: criar solicitação, aceitar/rejeitar/concluir/cancelar, rastrear/listar por dono, prestador ou status
- Review: submeter avaliação (apenas após conclusão), listar avaliações do prestador

## Arquitetura

- `internal/domain`: entidades, value objects, regras e erros de domínio, interfaces de repositório e serviços
- `internal/application`: casos de uso orquestram repositórios/serviços; sem dependência de detalhes de infra
- `internal/application/logging`: helper `UseCase` para log de início/fim com duração e status
- Infra (repositórios, HTTP, JWT, email, bcrypt) não incluída neste repositório

## Configuração

1. Instale Go 1.25.5
2. `go mod tidy`
3. (Opcional) Crie `.env` para credenciais de infra quando implementar camada externa

## Build e Testes

- Build: `go build ./...`
- (Quando houver testes) `go test ./...`

## Logging

- Injeção de `*slog.Logger` nos casos de uso (passe `nil` para usar `slog.Default`)
- `UseCase(ctx, logger, "UseCaseName", attrs...)` registra `start` e `end` com `duration` e `status` (`ok` ou `error`)

## Próximos Passos (sugestão)

- Implementar infra: repositórios (ex.: Postgres), token service JWT, password hasher (bcrypt), email service
- Expor HTTP/REST ou gRPC com middlewares de autenticação
- Criar migrações e seeds iniciais
- Adicionar testes unitários e integração
