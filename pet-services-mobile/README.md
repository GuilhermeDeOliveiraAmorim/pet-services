# pet-services-mobile

Aplicativo React Native (Expo) para o ecossistema do Pet Services.

## Stack adotada

- Expo SDK 55 + React Native + TypeScript
- Expo Router para navegacao
- Zustand para sessao/autenticacao local
- Axios para cliente HTTP
- TanStack React Query para cache e sincronizacao de dados

## Estrutura inicial

- app/: rotas e telas (expo-router)
- src/features/: regras por dominio (auth, pets, providers)
- src/lib/http: cliente HTTP
- src/lib/config: configuracoes de ambiente
- src/providers: providers globais (React Query)

## Fluxo inicial implementado

- /(auth)/login
- /(auth)/register?user_type=provider
- /(tabs)/owner
- /(tabs)/partner

## Rodando localmente

1. Instale dependencias:
   npm install
2. Inicie o app:
   npm run start
3. Execute no Android:
   npm run android
4. Execute no iOS (macOS):
   npm run ios

## Configurando API

Defina `expo.extra.apiUrl` em app.json para apontar para sua API:

```json
{
  "expo": {
    "extra": {
      "apiUrl": "http://SEU_IP:8080"
    }
  }
}
```

Observacao para dispositivo fisico: `localhost` do celular nao aponta para seu computador.
Use o IP local da sua maquina na rede.

## Proximos passos recomendados

1. Integrar login/register reais em src/features/auth/api.ts
2. Persistir sessao com storage seguro (expo-secure-store)
3. Criar telas de listagem de servicos e detalhes
4. Adicionar validacao de formularios (zod + react-hook-form)
5. Configurar EAS Build para Play Store/App Store
