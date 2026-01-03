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
3. Copie `.env.example` para `.env` e configure credenciais:
   ```bash
   cp .env.example .env
   # Edit .env with your local DB credentials
   ```
4. (Opcional) Use Docker para PostgreSQL:
   ```bash
   docker-compose up -d postgres
   ```

## Build e Testes

- Build: `make build` ou `go build ./cmd/api`
- Run API: `make run` (compila e executa cmd/api/main.go)
- Testes: `make test` ou `go test ./...`
- Migrations: `make migrate`

## Quick Start

1. **Setup DB**:

   ```bash
   make setup  # Instala deps, sobe PostgreSQL, roda migrações
   ```

2. **Run API**:

   ```bash
   make run
   # Ou manualmente:
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/pet_services?sslmode=disable"
   go run ./cmd/api/main.go
   ```

3. **Verify**:
   - Health check: `curl http://localhost:8080/health`
   - Readiness: `curl http://localhost:8080/ready`
   - Swagger UI: `http://localhost:8080/api/v1/swagger/index.html`

## Logging

- Injeção de `*slog.Logger` nos casos de uso (passe `nil` para usar `slog.Default`)
- `UseCase(ctx, logger, "UseCaseName", attrs...)` registra `start` e `end` com `duration` e `status` (`ok` ou `error`)

## Próximos Passos (sugestão)

- Implementar infra: repositórios (ex.: Postgres), token service JWT, password hasher (bcrypt), email service
- Expor HTTP/REST ou gRPC com middlewares de autenticação
- Criar migrações e seeds iniciais
- Adicionar testes unitários e integração

## API Server (`cmd/api/main.go`)

Startup sequence:

1. Load environment (`.env`)
2. Connect to Postgres via `infra/database.Open()`
3. Run migrations
4. Build use case factory with real/stub implementations
5. Register HTTP routes + middlewares (CORS, request ID, logging)
6. Start Gin server with graceful shutdown on SIGINT/SIGTERM

**Endpoints**:

- `GET /health` - Liveness check
- `GET /ready` - Readiness check (DB connectivity)
- `GET /api/v1/swagger/*any` - Swagger UI
- `POST /api/v1/auth/signup` - User signup
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - Logout
- (Protected) `/api/v1/users/*` - User profile/settings
- (Protected) `/api/v1/providers/*` - Provider operations
- (Protected) `/api/v1/requests/*` - Service requests
- (Protected) `/api/v1/reviews/*` - Reviews

**Middleware Stack**:

- Recovery (panic handler)
- Request ID (generated or from header)
- Structured logging (slog with request context)
- CORS (configurable origins)
- Auth (Bearer token on protected routes)

**Configuration** (via environment):

- `DATABASE_URL` or `DB_USER/DB_PASS/DB_NAME/DB_HOST/DB_PORT/DB_SSLMODE`
- `PORT` (default: 8080)
- `CORS_ORIGINS` (default: \*)
- `PASSWORD_RESET_BASE_URL`, `EMAIL_VERIFY_BASE_URL`

**Stub Implementations** (TODO: replace):

- `TokenService` - Returns dummy tokens
- `PasswordHasher` - Naive hashing (use bcrypt in prod)
- `EmailService` - Prints to console (use SMTP/SendGrid in prod)
