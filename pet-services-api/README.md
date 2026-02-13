# Pet Services API

API robusta e escalável para gerenciamento de serviços pet, desenvolvida em Go com arquitetura limpa. A plataforma conecta tutores de animais com prestadores de serviços qualificados, oferecendo funcionalidades completas de autenticação, gerenciamento de perfil, agendamento e avaliações.

## 🎯 Funcionalidades Principais

### Autenticação e Segurança

- ✅ Autenticação JWT com tokens de acesso e refresh
- ✅ Recuperação e alteração de senha
- ✅ Verificação de email com tokens de expiração
- ✅ Rate limiting para endpoints sensíveis
- ✅ Soft delete com gerenciamento de ativos/inativos
- ✅ Dois tipos de usuários: **Owner** (tutor) e **Provider** (prestador)

### Gerenciamento de Usuários

- ✅ Registro e criação de conta
- ✅ Atualização de perfil (endereço, telefone, dados pessoais)
- ✅ Upload de foto de perfil com armazenamento em MinIO
- ✅ Listagem de usuários com filtros
- ✅ Reativação e desativação de conta
- ✅ Gerenciamento de acesso por tipo de usuário

### Pets

- ✅ Cadastro e atualização de pets (nome, espécie, idade, peso, notas)
- ✅ Listagem de pets do proprietário
- ✅ Upload de múltiplas fotos por pet (até 10 fotos)
- ✅ Deleção de fotos com limpeza de armazenamento
- ✅ Soft delete de pets

### Provedores de Serviço

- ✅ Cadastro de provedor com informações comerciais
- ✅ Descrição de serviços e faixa de preço
- ✅ Upload de fotos de perfil
- ✅ Gerenciamento de endereço comercial
- ✅ Soft delete de provedores

### Serviços

- ✅ Criação de serviços com nome, descrição e precificação (fixa ou por faixa)
- ✅ Associação com categorias e tags
- ✅ Upload de fotos de serviço (até 10 fotos)
- ✅ Busca avançada com filtros por preço, categoria, tag
- ✅ Busca geoespacial por raio de proximidade (Haversine)
- ✅ Listagem com paginação
- ✅ Atualização com validação de conflitos de preço
- ✅ Soft delete de serviços

### Requisições de Serviço

- ✅ Criação de requisições (owner para provider)
- ✅ Ciclo de vida: pendente → aceita → completa/rejeitada
- ✅ Listagem com filtro por status
- ✅ Busca de requisição específica
- ✅ Ações: aceitar, rejeitar, completar (apenas provider)

### Avaliações e Reviews

- ✅ Criação de avaliações por tutores
- ✅ Rating de 1-5 estrelas
- ✅ Comentários e feedback
- ✅ Listagem de reviews por provedor

### Referências de Dados

- ✅ Listagem de países, estados e cidades (sem autenticação)
- ✅ Catálogo de espécies de animais
- ✅ Categorias de serviço (criação restrita a admins)
- ✅ Tags para serviços

### Admin

- ✅ Criação de novos usuários admins
- ✅ Criação e gerenciamento de categorias

### Infraestrutura

- ✅ Documentação automática com Swagger
- ✅ Health check da API e banco de dados
- ✅ CORS configurável
- ✅ Rate limiting customizável
- ✅ Logging estruturado
- ✅ Gerenciamento de conexões do banco
- ✅ Migrações automáticas

## 📊 Arquitetura

A API segue os princípios de **Clean Architecture** com separação clara de responsabilidades:

```
internal/
├── auth/                 # Autenticação e JWT
├── config/              # Configuração da aplicação
├── consts/              # Constantes e mensagens de erro
├── database/            # Conexão, migrações e modelos
├── entities/            # Entidades de domínio (interfaces)
├── exceptions/          # Tratamento de erros
├── factories/           # Factory pattern para instanciação
├── handlers/            # Handlers HTTP (entrada)
├── logging/             # Logger estruturado
├── mail/                # Serviço de email
├── middlewares/         # Middlewares HTTP (autenticação, rate limit)
├── models/              # Modelos GORM (saída)
├── reference/           # Serviços de referência (países, cidades)
├── repository_impl/     # Implementação de repositórios
├── routes/              # Roteamento e setup
├── storage/             # Serviço de armazenamento (MinIO)
└── usecases/            # Casos de uso (lógica de negócio)
```

### Fluxo de Requisição

```
HTTP Request
    ↓
Routes
    ↓
Middlewares (Auth, RateLimit, ProfileComplete)
    ↓
Handlers (Validação de entrada)
    ↓
Use Cases (Lógica de negócio)
    ↓
Repositories (Acesso a dados)
    ↓
Database (PostgreSQL)
```

## 🛠️ Stack Técnico

| Componente                 | Tecnologia       | Versão |
| -------------------------- | ---------------- | ------ |
| **Linguagem**              | Go               | 1.24+  |
| **Framework Web**          | Gin              | 1.11+  |
| **ORM**                    | GORM             | 1.31+  |
| **Banco de Dados**         | PostgreSQL       | 12+    |
| **Armazenamento**          | MinIO            | 7.0+   |
| **Autenticação**           | JWT (golang-jwt) | -      |
| **Documentação**           | Swagger/OpenAPI  | -      |
| **Containerização**        | Docker           | -      |
| **Gerenciador de Pacotes** | Go Modules       | -      |

## 📦 Dependências Principais

```go
github.com/gin-gonic/gin                    # Framework web
github.com/gin-contrib/cors                 # CORS middleware
gorm.io/gorm                               # ORM
gorm.io/driver/postgres                    # Driver PostgreSQL
github.com/minio/minio-go                  # Cliente MinIO
golang.org/x/crypto                        # Hashing de senha
github.com/swaggo/swag                     # Swagger
github.com/oklog/ulid                      # Geração de IDs
```

## 🚀 Início Rápido

### Pré-requisitos

- **Go 1.24+** instalado
- **Docker e Docker Compose** (opcional, recomendado)
- **PostgreSQL 12+** ou container Docker
- **MinIO** (para armazenamento de fotos)

### 1. Clone o Repositório

```bash
git clone https://github.com/GuilhermeDeOliveiraAmorim/pet-services.git
cd pet-services/pet-services-api
```

### 2. Configure Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```bash
# Banco de Dados
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=sua_senha_postgres
DB_PORT=5432
DB_NAME=pet_services

# Segurança (JWT)
JWT_SECRET=sua_chave_secreta_super_segura_aqui
JWT_ACCESS_SECRET=sua_chave_access_token
JWT_REFRESH_SECRET=sua_chave_refresh_token

# Servidor
SERVER_PORT=8080
SWAGGER_HOST=localhost:8080

# Frontend
FRONT_END_URL_DEV=http://localhost:3000
FRONT_END_URL_PROD=https://seu-dominio.com

# MinIO (Armazenamento)
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
IMAGE_BUCKET_NAME=pet-services-images
URL_BUCKET_NAME=pet-services-urls

# Email
EMAIL_SERVICE_API_KEY=sua_chave_api_email
EMAIL_VERIFY=seu_email_verificacao

# Segurança
RESET_PASSWORD_EXPIRATION_TIME=3600
EMAIL_CHANGE_EXPIRATION_TIME=86400
VERIFY_EMAIL_EXPIRATION_TIME=86400
MAX_CHANGE_EMAIL_ATTEMPTS=3
```

### 3. Suba a Infraestrutura (Docker)

```bash
cd ../pet-services-infra
docker-compose up -d
```

Este comando sobe:

- PostgreSQL na porta 5432
- MinIO na porta 9000
- Mailpit (teste de email) na porta 1025

### 4. Instale Dependências

```bash
cd ../pet-services-api
go mod tidy
go mod download
```

### 5. Execute as Migrações

```bash
make migrate
```

### 6. Inicie a API

```bash
make run
```

A API estará disponível em `http://localhost:8080`

Acesse a documentação Swagger em `http://localhost:8080/swagger/index.html`

## 🔄 Desenvolvimento

### Hot Reload com Air

Para desenvolvimento com reload automático:

```bash
make dev
```

Certifique-se de ter o `air` instalado:

```bash
go install github.com/cosmtrek/air@latest
```

### Testes

Execute todos os testes:

```bash
make test
```

Gere relatório de cobertura:

```bash
make test-coverage
```

### Lint e Formatação

```bash
# Formatar código
make fmt

# Lint
make lint
```

### Gerar Documentação Swagger

```bash
make doc
```

Isso gera/atualiza a documentação Swagger baseada nos comentários das handlers.

## 📡 Endpoints Principais

### Autenticação

- `POST /auth/login` - Login com email/senha
- `POST /auth/refresh` - Renovar token JWT
- `POST /auth/logout` - Logout (revoga tokens)
- `POST /auth/request-password-reset` - Solicitar reset de senha
- `POST /auth/reset-password` - Resetar senha
- `POST /auth/verify-email` - Verificar email
- `POST /auth/resend-verification-email` - Reenviar email de verificação

### Usuários

- `POST /users/register` - Registrar novo usuário
- `GET /users/profile` - Obter perfil autenticado
- `GET /users/:user_id` - Obter dados de usuário específico
- `GET /users` - Listar usuários com paginação
- `PUT /users` - Atualizar perfil
- `DELETE /users` - Deletar conta
- `POST /users/reactivate` - Reativar conta desativada
- `POST /users/deactivate` - Desativar conta
- `POST /users/change-password` - Alterar senha
- `POST /users/photos` - Upload de foto de perfil

### Pets

- `GET /pets` - Listar pets do proprietário
- `GET /pets/:pet_id` - Obter detalhes do pet
- `POST /pets` - Criar novo pet
- `PUT /pets/:pet_id` - Atualizar pet
- `DELETE /pets/:pet_id` - Deletar pet (soft delete)
- `POST /pets/:pet_id/photos` - Upload de foto do pet
- `DELETE /pets/:pet_id/photos/:photo_id` - Deletar foto do pet

### Provedores

- `POST /providers` - Criar provedor
- `GET /providers/:provider_id` - Obter detalhes do provedor
- `POST /providers/photos` - Upload de foto do provedor
- `DELETE /providers/:provider_id/photos/:photo_id` - Deletar foto
- `DELETE /providers/:provider_id` - Deletar provedor (soft delete)

### Serviços

- `GET /services` - Listar serviços com paginação
- `GET /services/:service_id` - Obter detalhes do serviço
- `GET /services/search` - Buscar serviços (filtros + geolocalização)
- `POST /services` - Criar novo serviço
- `PUT /services/:service_id` - Atualizar serviço
- `DELETE /services/:service_id` - Deletar serviço (soft delete)
- `POST /services/:service_id/photos` - Upload de foto do serviço
- `DELETE /services/:service_id/photos/:photo_id` - Deletar foto
- `POST /services/:service_id/tags` - Adicionar tag ao serviço
- `POST /services/:service_id/categories` - Adicionar categoria ao serviço

### Requisições

- `GET /requests` - Listar requisições
- `GET /requests/:request_id` - Obter detalhes da requisição
- `POST /requests` - Criar requisição (owner)
- `PATCH /requests/:request_id/accept` - Aceitar requisição (provider)
- `PATCH /requests/:request_id/reject` - Rejeitar requisição (provider)
- `PATCH /requests/:request_id/complete` - Completar requisição (provider)

### Reviews

- `GET /reviews` - Listar reviews
- `POST /providers/:provider_id/reviews` - Criar review (owner)

### Referências

- `GET /reference/countries` - Listar países
- `GET /reference/states` - Listar estados
- `GET /reference/cities` - Listar cidades
- `GET /species` - Listar espécies de animais
- `GET /categories` - Listar categorias
- `GET /tags` - Listar tags

### Admin

- `POST /admin` - Criar novo admin
- `POST /admin/categories` - Criar categoria

### Health Check

- `GET /health` - Health check da API e banco de dados

## 🔐 Autenticação e Autorização

### Fluxo de Autenticação

```
1. POST /auth/login { email, password }
   ↓
2. Retorna { accessToken, refreshToken }
   ↓
3. Inclua em headers: Authorization: Bearer {accessToken}
```

### Tipos de Usuário

| Tipo         | Descrição             | Acesso                                                         |
| ------------ | --------------------- | -------------------------------------------------------------- |
| **Owner**    | Tutor de animais      | Gerenciar pets, criar requisições, avaliar provedores          |
| **Provider** | Prestador de serviços | Gerenciar provedores e serviços, aceitar/completar requisições |
| **Admin**    | Administrador         | Criar admins, categorias, gerenciar plataforma                 |

### Middlewares

- **AuthMiddleware**: Valida JWT e extrai user_id
- **ProfileCompleteMiddleware**: Verifica se perfil está completo (bloqueia algumas rotas)
- **OwnerOnlyMiddleware**: Restringe acesso a owners
- **ProviderOnlyMiddleware**: Restringe acesso a providers
- **AdminOnlyMiddleware**: Restringe acesso a admins
- **DefaultRateLimitMiddleware**: Limita requisições gerais (100 req/min)
- **StrictRateLimitMiddleware**: Limita endpoints sensíveis (10 req/min)

## 💾 Modelos de Dados

### User

```go
type User struct {
    ID               string    // ULID
    Email            string    // Único
    Password         string    // Hashed
    Name             string
    UserType         string    // "owner" ou "provider"
    Phone            Phone
    Address          Address
    ProfileComplete  bool
    Active           bool      // Soft delete
    CreatedAt        time.Time
    UpdatedAt        *time.Time
    DeactivatedAt    *time.Time
}
```

### Pet

```go
type Pet struct {
    ID          string
    UserID      string        // FK User
    SpecieID    string        // FK Species
    Name        string
    Age         int
    Weight      float64
    Notes       string
    Photos      []Photo       // Até 10 fotos
    Active      bool          // Soft delete
    CreatedAt   time.Time
    UpdatedAt   *time.Time
    DeactivatedAt *time.Time
}
```

### Service

```go
type Service struct {
    ID           string
    ProviderID   string        // FK Provider
    Name         string
    Description  string
    Price        float64       // Preço fixo (opcional)
    PriceMinimum float64       // Preço mínimo (opcional)
    PriceMaximum float64       // Preço máximo (opcional)
    Duration     int           // Em minutos
    Photos       []Photo       // Até 10 fotos
    Categories   []Category
    Tags         []Tag
    Active       bool          // Soft delete
    CreatedAt    time.Time
    UpdatedAt    *time.Time
    DeactivatedAt *time.Time
}
```

### Request

```go
type Request struct {
    ID          string
    OwnerID     string        // FK User (Owner)
    ProviderID  string        // FK Provider
    ServiceID   string        // FK Service
    PetID       string        // FK Pet
    Status      string        // "pending", "accepted", "completed", "rejected"
    Active      bool          // Soft delete
    CreatedAt   time.Time
    UpdatedAt   *time.Time
    DeactivatedAt *time.Time
}
```

### Review

```go
type Review struct {
    ID         string
    OwnerID    string        // FK User (Owner)
    ProviderID string        // FK Provider
    RequestID  string        // FK Request
    Rating     int           // 1-5
    Comment    string
    CreatedAt  time.Time
}
```

## 🖼️ Gerenciamento de Fotos

### Características

- **Armazenamento**: MinIO S3-compatible
- **Limite por usuário**: 1 foto de perfil
- **Limite por pet**: 10 fotos
- **Limite por serviço**: 10 fotos
- **Limite por provedor**: 10 fotos
- **URLs Assinadas**: 15 minutos de validade
- **Limpeza**: Ao deletar foto, arquivo é removido do armazenamento

### Exemplo de Upload

```bash
curl -X POST http://localhost:8080/pets/123/photos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "photo=@/path/to/photo.jpg"
```

## 🔍 Busca e Filtros

### Busca de Serviços Simples

```bash
GET /services?categoryId=xxx&tagId=yyy&priceMin=100&priceMax=500&page=1&pageSize=10
```

### Busca Geoespacial

```bash
GET /services/search?query=grooming&latitude=-9.5&longitude=-35.7&radiusKm=10&page=1
```

Usa a **fórmula de Haversine** para calcular distância em quilômetros.

## 📝 Soft Delete

Todos os registros principais usam soft delete:

```
User, Pet, Service, Provider, Request
```

Quando um registro é deletado:

- `active` é definido como `false`
- `deactivated_at` é preenchido com o timestamp
- O registro permanece no banco (audit trail)
- Listagens filtram apenas registros `active = true` por padrão

## 🛡️ Segurança

### Senhas

- Hashing com bcrypt
- Mínimo de 8 caracteres
- Suporta caracteres especiais

### JWT

- Access token: Curta duração (~15 min)
- Refresh token: Longa duração (~7 dias)
- Revogação ao logout
- Validação de signature

### Rate Limiting

- **Geral**: 100 requisições por minuto
- **Autenticação**: 10 requisições por minuto (login, reset de senha)
- Baseado em IP

### CORS

- Configurável por variáveis de ambiente
- Suporta múltiplos domínios
- Credenciais seguras em produção

### Validações

- Email e telefone únicos
- Validação de campos obrigatórios
- Validação de comprimento de string
- Validação de valores numéricos (idade, peso, preço)
- Validação de ranges (rating 1-5, preço mínimo ≤ máximo)

## 📊 Migrações

As migrações são executadas automaticamente na inicialização:

### Migrações Aplicadas

| Versão         | Descrição                                                |
| -------------- | -------------------------------------------------------- |
| 20260110000000 | Criar esquema inicial (users, pets, services, etc)       |
| 20260204000000 | Criar tabela para tokens de reset de senha e verificação |
| 20260213000000 | Adicionar campo `profile_complete` na tabela users       |
| 20260215000000 | Criar tabela de refresh tokens                           |

Para adicionar nova migração, crie função em `internal/database/migrations.go` e registre em `getMigrations()`.

## 🧪 Testes

### Estrutura

```bash
go test -v -race -coverprofile=coverage.out ./...
```

### Cobertura

```bash
go tool cover -html=coverage.out
```

### Exemplo de Teste

```go
func TestAddPetUseCase_Execute(t *testing.T) {
    // Arrange
    mockUserRepo := &MockUserRepository{}
    mockPetRepo := &MockPetRepository{}

    // Act
    uc := usecases.NewAddPetUseCase(mockUserRepo, mockPetRepo, logger)

    // Assert
    // ...
}
```

## 📈 Performance

### Otimizações

- **Pool de conexões**: Max 20 conexões abertas, 10 ociosas
- **Índices de banco**: Email, phone, active, createdAt
- **Paginação**: Padrão 10 itens/página
- **Lazy loading**: Associações carregadas sob demanda com `Preload`
- **Caching**: URLs assinadas em memória (15 min)

### Benchmarks

```bash
go test -bench=. -benchmem ./...
```

## 🚢 Deploy

### Docker

```bash
# Build image
docker build -t pet-services-api:latest .

# Run container
docker run -e DB_HOST=host.docker.internal \
           -e JWT_SECRET=seu_secret \
           -p 8080:8080 \
           pet-services-api:latest
```

### Docker Compose

```bash
cd pet-services-infra
docker-compose up -d
```

### Variáveis de Produção

- Usar secrets manager (AWS Secrets Manager, HashiCorp Vault)
- JWT_SECRET com mínimo 32 caracteres
- Desabilitar CORS wildcard
- Usar HTTPS
- MinIO com SSL
- PostgreSQL com conexão SSL

## 🐛 Logging

Logs estruturados com `slog`:

```
[Start] Migrações do banco concluídas com sucesso
[Migration] aplicando migração version=20260110000000
[Api.GetPetUseCase] Erro ao buscar pet error="record not found"
```

Níveis:

- **INFO**: Eventos importantes
- **WARN**: Avisos
- **ERROR**: Erros não fatais
- **DEBUG**: Informações detalhadas (em desenvolvimento)

## 📚 Estrutura de Respostas

### Sucesso

```json
{
  "message": "Pet adicionado com sucesso",
  "detail": "O pet foi registrado no sistema",
  "pet": {
    "id": "01ARZ3NDEKTSV4RRFFQ69G5FAV",
    "name": "Rex",
    "species_id": "001",
    "age": 3,
    "weight": 25.5,
    "active": true,
    "created_at": "2026-02-13T10:30:00Z"
  }
}
```

### Erro

```json
{
  "type": "https://example.com/errors/bad-request",
  "title": "Email já cadastrado",
  "detail": "O email informado já está associado a outra conta",
  "status": 400,
  "instance": "/users/register"
}
```

## 🤝 Contribuindo

### Setup de Desenvolvimento

```bash
# Clone e configure
git clone ...
cd pet-services-api
cp .env.example .env
make setup

# Inicie desenvolvimento
make dev
```

### Padrões de Código

- Follow Go idioms: `gofmt`, `go vet`
- Use context para operações assíncronas
- Tratar erros explicitamente
- Documentar funções públicas
- Use interfaces para dependências

### Commit Messages

```
feat: adicionar endpoint de busca geoespacial
fix: corrigir validação de email duplicado
refactor: simplificar lógica de soft delete
docs: atualizar documentação de autenticação
test: adicionar testes de rate limiting
```

## 📞 Suporte e Contato

- **Autor**: Guilherme de Oliveira Amorim
- **Email**: guilherme@example.com
- **Issues**: GitHub Issues
- **Documentação**: Swagger em `/swagger/index.html`

## 📄 Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.

---

**Última atualização**: 13 de fevereiro de 2026  
**Versão da API**: 1.0.0  
**Versão do Go**: 1.24+
