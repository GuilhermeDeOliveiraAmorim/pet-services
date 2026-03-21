# pet-services-front

Frontend do projeto Pet Services, construído com Next.js (App Router), React, TypeScript e Chakra UI.

<details>
<summary><strong>Sumário</strong> (clique para expandir)</summary>

- [Visão Geral](#pet-services-front)
- [Stack Principal](#stack-principal)
- [Arquitetura & Padrões](#arquitetura--padrões)
- [Fluxos Principais](#fluxos-principais)
- [Estratégia de Cache & Query Keys](#estratégia-de-cache--query-keys)
- [Como Rodar](#como-rodar)
- [Scripts Úteis](#scripts-úteis)
- [Testes Locais](#testes-locais-auth)
- [Dicas de Desenvolvimento](#dicas-de-desenvolvimento)
- [Documentação de Cache & Query Keys](#documentação-de-cache--query-keys)

</details>

---

## Stack principal

- Next.js 16 (App Router)
- React 19
- TypeScript
- Chakra UI v3
- TanStack Query (React Query)
- Axios

---

## Arquitetura & Padrões

- **Estrutura modular:**
  - `src/app/`: rotas protegidas e públicas (Next.js App Router)
  - `src/application/`: camada de aplicação (usecases, gateways, hooks)
  - `src/domain/`: entidades e tipos de domínio
  - `src/components/`: componentes reutilizáveis (UI, formulários, navegação)
  - `src/infra/`: gateways HTTP, cache, providers
- **Padrões:**
  - Application Layer (usecases, gateways, hooks)
  - Query Keys centralizadas para cache e invalidação
  - Componentização e reutilização máxima
  - Navegação protegida por autenticação e roles
  - Testes automatizados para flows críticos

---

## Fluxos Principais

- **Autenticação:** login, registro, refresh, troca de senha, verificação de e-mail
- **Usuário:** atualização de perfil, upload de foto, roles (owner, provider, admin)
- **Pets:** cadastro, edição, upload de fotos, deleção
- **Serviços:** busca, filtros, detalhes, criação/edição (provider)
- **Adoção & Guardian:**
  - Cadastro e edição de perfil de responsável (guardian)
  - Criação e edição de listagens de adoção
  - Candidatura, revisão e status
  - Aprovação de perfis e candidaturas (admin)
- **Admin:** gerenciamento de categorias, aprovação/rejeição de perfis sensíveis
- **Avaliações:** reviews, ratings, lembretes

---

## Estratégia de Cache & Query Keys

- Query Keys centralizadas por domínio (ver arquivos em `src/application/hooks/*-query-keys.ts`)
- Invalidação e atualização de cache padronizadas
- Documentação detalhada em [CACHE_STRATEGY.md](CACHE_STRATEGY.md) e [QUERY_KEYS_REFERENCE.md](QUERY_KEYS_REFERENCE.md)
- Exemplos de uso e padrões para novos domínios

---

## Como rodar

Pré-requisitos:

- Node.js 20+
- API backend rodando (padrão: `http://localhost:8080`)

Instalação e execução:

```bash
npm install
npm run dev
```

Aplicação local: `http://localhost:3000`

---

## Scripts úteis

- `npm run dev`: ambiente de desenvolvimento
- `npm run build`: build de produção
- `npm run start`: inicia build de produção
- `npm run lint`: validação ESLint

---

## Testes locais (Auth)

Variáveis suportadas:

- `API_URL` (opcional, padrão: `http://localhost:8080`)
- `TEST_EMAIL` (obrigatório em parte dos scripts)
- `TEST_PASSWORD` (opcional, padrão interno dos scripts)
- `NEW_PASSWORD` (opcional, usado no fluxo de troca de senha)

Scripts:

- `npm run test:auth`: login → refresh → logout
- `npm run test:auth-flow`: registro → verificação → login
- `npm run test:verify-email`: reenvio e validação de token de verificação
- `npm run test:change-password`: troca de senha autenticado
- `npm run test:update-user`: atualização de perfil

---

## Dicas de Desenvolvimento

- Use os hooks e usecases da camada de aplicação para acessar dados e mutações
- Sempre utilize as query keys centralizadas para cache/invalidação
- Consulte [CACHE_STRATEGY.md](CACHE_STRATEGY.md) e [QUERY_KEYS_REFERENCE.md](QUERY_KEYS_REFERENCE.md) para padrões de cache
- Utilize componentes reutilizáveis para formulários e navegação
- Siga o padrão de roles para proteger rotas e ações sensíveis
- Rode `npm run lint` e `npm run build` antes de commitar

---

## Documentação de Cache & Query Keys

- [CACHE_STRATEGY.md](CACHE_STRATEGY.md): Estratégia completa de cache, problemas resolvidos, exemplos
- [QUERY_KEYS_REFERENCE.md](QUERY_KEYS_REFERENCE.md): Guia de referência das query keys, exemplos e padrões
