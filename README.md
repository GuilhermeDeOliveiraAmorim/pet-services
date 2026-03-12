# Pet Services - Monorepo

Este repositório contém os serviços e infraestrutura do projeto Pet Services, uma plataforma para gestão de prestadores, solicitações e avaliações de serviços para pets.

## 📝 Sobre o Projeto

O **Pet Services** é um aplicativo desenvolvido para conectar donos de pets a prestadores de serviços pet locais, facilitando a vida de tutores que buscam cuidados de qualidade para seus animais de estimação.

### 🎯 Objetivo

Criar uma solução que simplifique o processo de encontrar, comparar e solicitar serviços pet na mesma região do usuário, proporcionando uma experiência simples, rápida e eficiente.

### 🚀 Funcionalidades Principais

- **Descoberta de Serviços**: Encontre facilmente prestadores de serviços pet próximos à sua localização
- **Comparação de Prestadores**: Compare avaliações, preços e especialidades dos prestadores de forma prática
- **Solicitação Simplificada**: Faça pedidos de serviços de forma rápida e direta, sem complicações
- **Serviços Locais**: Foco em conectar usuários com prestadores da mesma região, garantindo praticidade e agilidade
- **Gestão de Solicitações**: Acompanhe o status das suas solicitações e histórico de serviços
- **Sistema de Avaliações**: Compartilhe e consulte avaliações de outros tutores sobre os prestadores

## Estrutura do Projeto

```
pets-services/
├── pet-services-api/         # Código-fonte da API principal (Go)
│   ├── Dockerfile            # Dockerfile da API
│   ├── docker-compose.yml    # Compose legado (não usar, infra está em outro compose)
│   └── ...                   # Demais arquivos e pastas
└── pet-services-infra/       # Infraestrutura (Docker Compose para API, banco e MinIO)
    └── docker-compose.yml    # Compose principal para desenvolvimento
```

## Como rodar o projeto

### Pré-requisitos

- Docker e Docker Compose instalados
- Go 1.21+ (apenas para desenvolvimento local da API)

### Subindo toda a stack (API, banco e MinIO)

1. Acesse a pasta de infraestrutura:
   ```sh
   cd pet-services-infra
   ```
2. Suba os serviços:
   ```sh
   docker compose up --build
   ```

- A API estará disponível em: http://localhost:8080
- O MinIO estará disponível em: http://localhost:9001 (console)
  - Usuário: `minio`
  - Senha: `minio123`
- O banco Postgres estará em: localhost:5432
  - Usuário: `postgres`
  - Senha: `postgres`
  - Banco: `pet_services`

### Configurações de Ambiente (Infra)

No arquivo `.env` dentro de `pet-services-infra/`, mantenha pelo menos estas variáveis ajustadas:

- `ENV=development` no desenvolvimento local
- `ENV=production` em produção
- `VOLUME_PATH=/caminho/no/seu/disco` para persistência de dados do MinIO

#### Modo do Gin por ambiente

A API define o modo do Gin automaticamente com base em `ENV`:

- `ENV=production` → Gin em `release`
- qualquer outro valor (ex.: `development`) → Gin em `debug`

#### MinIO com volume configurável

O MinIO usa bind mount com `VOLUME_PATH`, permitindo escolher o disco/pasta de persistência sem alterar o compose.

Exemplo:

```sh
VOLUME_PATH=/media/seu-usuario/SeuDisco/minio
```

### Seeds automáticos no startup

Ao subir a API com migrações, são garantidos 2 users de desenvolvimento (idempotente):

- Owner:
  - Email: `owner.seed@petservices.local`
  - Senha: `Owner@123`
- Provider:
  - Email: `provider.seed@petservices.local`
  - Senha: `Provider@123`
  - Perfil de provider já criado e vinculado

Se você subir novamente os containers, os seeds são ignorados/atualizados sem duplicar registros.

### Desenvolvimento da API

1. Acesse a pasta da API:
   ```sh
   cd pet-services-api
   ```
2. Para rodar localmente (fora do Docker):
   ```sh
   go run ./cmd/api
   ```

### Rodando testes

Na pasta da API:

```sh
make test
```

## Documentação

- Acesse os arquivos `.md` na pasta `pet-services-api/` para detalhes de endpoints, análise de MVP, rate limiting e guias de implementação.

## Observações

- O compose da pasta infra é o principal para desenvolvimento local.
- O compose da API pode ser removido ou mantido apenas para referência.
- Ajuste variáveis de ambiente conforme necessário para integração com outros serviços.

---

Desenvolvido por Guilherme e colaboradores.
