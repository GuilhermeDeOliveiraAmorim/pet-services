This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

## Testes locais (Auth)

Pré-requisitos:

- API rodando em `http://localhost:8080` (ou defina `API_URL`)
- Mailpit opcional para verificação de email: `http://localhost:8025`

Variáveis usadas pelos scripts:

- `API_URL` (opcional, padrão: `http://localhost:8080`)
- `TEST_EMAIL` (obrigatório para alguns testes)
- `TEST_PASSWORD` (opcional; padrão: `123QWEasd@`)
- `NEW_PASSWORD` (opcional; usado no change-password)

Scripts disponíveis:

```bash
npm run test:auth
```

Login → refresh → logout.

```bash
npm run test:auth-flow
```

Registro → reenvio de verificação → verificação → login.

```bash
npm run test:verify-email
```

Reenvia verificação e valida o token.

```bash
npm run test:change-password
```

Login → change-password.

Exemplo:

```bash
API_URL=http://localhost:8080 \
TEST_EMAIL=guilherme.o.a.ufal@gmail.com \
TEST_PASSWORD=123QWEasd@ \
NEW_PASSWORD=NovaSenha@123 \
npm run test:change-password
```

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.
