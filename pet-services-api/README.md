# Pet Services API

API RESTful para gerenciamento de serviços pet, desenvolvida em Go utilizando o framework Gin, GORM para acesso ao banco de dados e autenticação JWT. Documentação automática via Swagger.

## Funcionalidades

- Cadastro, autenticação e gerenciamento de usuários (admin e cliente)
- Gerenciamento de pets, raças, espécies, categorias e serviços
- Sistema de avaliação e reviews
- Upload de fotos
- Recuperação e alteração de senha
- Health check da API e do banco de dados
- Rate limiting e middlewares customizados

## Stack Principal

- **Go** 1.24+
- **Gin** (HTTP framework)
- **GORM** (ORM)
- **PostgreSQL** (recomendado)
- **Swagger** (documentação automática)
- **Docker** (deploy containerizado)

## Como rodar localmente

### Pré-requisitos

- Go 1.24+
- Docker e Docker Compose (opcional, recomendado para ambiente completo)

### 1. Clonar o repositório

```bash
git clone https://github.com/GuilhermeDeOliveiraAmorim/pet-services.git
cd pet-services/pet-services-api
```

### 2. Configurar variáveis de ambiente

Crie um arquivo `.env` baseado no `.env.example` (se existir) ou configure as variáveis necessárias para conexão com o banco e JWT.

### 3. Subir banco de dados (opcional)

```bash
cd ../pet-services-infra
sudo docker-compose up -d
```

### 4. Rodar as migrações

```bash
cd ../pet-services-api
make migrate
```

### 5. Rodar a API

```bash
make run
```

A API estará disponível em `http://localhost:8080`.

## Documentação Swagger

Acesse `http://localhost:8080/swagger/index.html` para visualizar e testar os endpoints.

## Principais Endpoints

- `POST /users/register` — Cadastro de usuário
- `POST /auth/login` — Login
- `GET /user/profile` — Perfil do usuário autenticado
- `GET /health` — Health check da API e banco

## Testes

```bash
make test
```

## Build e Deploy com Docker

```bash
docker build -t pet-services-api .
docker run -p 8080:8080 pet-services-api
```

## Estrutura de Pastas

- `cmd/api/` — Entrypoint da aplicação
- `internal/` — Domínio, usecases, handlers, middlewares, entidades, repositórios
- `docs/` — Documentação Swagger
- `bin/` — Binários gerados

## Contribuição

Pull requests são bem-vindos! Abra uma issue para discutir melhorias ou reportar bugs.

## Licença

MIT
