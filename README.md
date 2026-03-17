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

## 📦 Atualizações Recentes (Mar/2026)

- **Catálogo de serviços evoluído**:
  - Paginação no catálogo e na busca
  - Filtros por categoria, tag e faixa de preço
  - Busca textual com ordenação por relevância, preço e duração
  - Busca por localização (CEP, geolocalização do navegador e raio em km)
- **API de serviços enriquecida**:
  - Resposta padronizada com `total_items`
  - Campos de avaliação no item de serviço (`average_rating`, `review_count`)
  - Distância aproximada (`distance_km`) quando há coordenadas na busca
- **Fluxo completo de solicitações**:
  - Página de solicitações com filtro por status
  - Ações do provider: aceitar, recusar (com motivo) e concluir
  - Fluxo do owner para avaliar provider após solicitação concluída
- **Detalhes de provider e serviço aprimorados**:
  - Exibição de avaliação média e volume de avaliações
  - Melhorias na apresentação de fotos e dados comerciais
- **Seeds de desenvolvimento expandidos**:
  - Base com providers, owners, pets, requests e reviews em cenário realista
  - Dados focados em Aracaju/SE para facilitar testes de localização e catálogo
- **Operação da stack simplificada**:
  - Novo comando `make rebuild-up` no Makefile raiz
  - Compose da infra com API, Front, Postgres, MinIO e Mailpit

## Estrutura do Projeto

```
pet-services/
├── pet-services-api/         # API principal (Go)
├── pet-services-front/       # Frontend web (Next.js + React + TypeScript)
├── pet-services-mobile/      # App mobile
├── pet-services-infra/       # Infra local (Docker Compose)
└── Makefile                  # Atalhos para subir/parar/rebuildar a stack
```

## Como rodar o projeto

### Pré-requisitos

- Docker e Docker Compose instalados
- Go 1.21+ (apenas para desenvolvimento local da API)
- Node.js 20+ (apenas para desenvolvimento local do front)

### Comandos rápidos (raiz do monorepo)

```sh
make up         # sobe a stack em background
make logs       # acompanha logs
make down       # derruba a stack
make rebuild-up # derruba, rebuilda e sobe novamente
```

### Subindo toda a stack (API, Front, banco, MinIO e Mailpit)

1. Acesse a pasta de infraestrutura:
   ```sh
   cd pet-services-infra
   ```
2. Suba os serviços:
   ```sh
   docker compose up --build
   ```

- A API estará disponível em: http://localhost:8080
- O Front estará disponível em: http://localhost:3000
- O MinIO estará disponível em:
  - Console: http://localhost:9001
  - API S3: http://localhost:9002
  - Credenciais: definidas por `MINIO_ROOT_USER` e `MINIO_ROOT_PASSWORD`
- O banco Postgres estará em: localhost:5433 (container expõe 5432 internamente)
  - Credenciais: definidas por `DB_USER`, `DB_PASS` e `DB_NAME`
- O Mailpit (SMTP dev + inbox) estará em:
  - SMTP: localhost:1025
  - UI: http://localhost:8025

### Configurações de Ambiente (Infra)

No arquivo `.env` dentro de `pet-services-infra/`, mantenha pelo menos estas variáveis ajustadas:

- `ENV=development` no desenvolvimento local
- `ENV=production` em produção
- `MINIO_DATA_PATH=/caminho/no/seu/disco` para persistência dos dados do MinIO

#### Modo do Gin por ambiente

A API define o modo do Gin automaticamente com base em `ENV`:

- `ENV=production` → Gin em `release`
- qualquer outro valor (ex.: `development`) → Gin em `debug`

#### MinIO com volume configurável

O MinIO usa bind mount com `MINIO_DATA_PATH`, permitindo escolher o disco/pasta de persistência sem alterar o compose.

Exemplo:

```sh
MINIO_DATA_PATH=/media/seu-usuario/SeuDisco/minio
```

### Seeds automáticos no startup

Ao subir a API com migrações, são garantidos seeds idempotentes para acelerar testes de QA local:

- Usuários básicos seed:
  - Owner: `owner.seed@petservices.local` / `Owner@123`
  - Provider: `provider.seed@petservices.local` / `Provider@123`
- Cenário expandido:
  - 5 providers com perfis e serviços em Aracaju/SE
  - 4 owners adicionais com pets vinculados
  - Solicitações em múltiplos status e avaliações já criadas

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
- O compose da API é legado para referência local isolada.
- Ajuste variáveis de ambiente conforme necessário para integração com outros serviços.

---

Desenvolvido por Guilherme e colaboradores.
