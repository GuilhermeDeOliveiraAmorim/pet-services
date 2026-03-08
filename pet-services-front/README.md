# pet-services-front

Frontend do projeto Pet Services, construído com Next.js (App Router), React, TypeScript e Chakra UI.

## Stack principal

- Next.js 16
- React 19
- TypeScript
- Chakra UI v3
- TanStack Query
- Axios

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

## Scripts úteis

- `npm run dev`: sobe ambiente de desenvolvimento (webpack)
- `npm run build`: build de produção
- `npm run start`: inicia build de produção
- `npm run lint`: validação ESLint

## Fluxo de cadastro (atual)

O cadastro segue o contrato atual da API em 2 etapas:

1. `POST /users/register`: criação inicial com dados básicos (`name`, `login`, `phone`)
2. `PUT /users`: complementação de perfil após autenticação

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
