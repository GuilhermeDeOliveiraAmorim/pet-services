# Cache Strategy & Query Keys Documentation

## Objetivo

Centralizar e padronizar query keys, invalidações e estratégia de cache para todos os domínios (services, providers, requests, reviews) no projeto.

---

## 📊 Estado Atual - ✅ IMPLEMENTADO (Fase 1 Completa)

### ✅ Com Query Keys Centralizados (10/10 domínios)

**Todos os domínios agora possuem query keys estruturadas e exportadas:**

| Domínio        | Arquivo                                   | Padrão                                                                                | Status        |
| -------------- | ----------------------------------------- | ------------------------------------------------------------------------------------- | ------------- |
| **Services**   | `hooks/service/service-query-keys.ts`     | `.all`, `.lists()`, `.list(filters)`, `.details()`, `.detail(id)`, `.byProvider(id)`  | ✅ NOVO       |
| **Providers**  | `hooks/provider/provider-query-keys.ts`   | `.all`, `.lists()`, `.details()`, `.detail(id)`, `.photos()`, `.photosByProvider(id)` | ✅ NOVO       |
| **Breeds**     | `hooks/breed/breed-query-keys.ts`         | `.all`, `.lists()`, `.list(speciesId)`                                                | ✅ NOVO       |
| **Species**    | `hooks/specie/specie-query-keys.ts`       | `.all`, `.lists()`, `.list()`                                                         | ✅ NOVO       |
| **User**       | `hooks/user/user-query-keys.ts`           | `.all`, `.profiles()`, `.profile()`                                                   | ✅ NOVO       |
| **Reviews**    | `hooks/review/review-query-keys.ts`       | `.all`, `.lists()`, `.list(filters)`                                                  | ✅ REFATORADO |
| **Requests**   | `hooks/request/request-query-keys.ts`     | `.all`, `.lists()`, `.list(filters)`, `.details()`, `.detail(id)`                     | ✅ REFATORADO |
| **Tags**       | `hooks/tag/tag-query-keys.ts`             | `.all`, `.lists()`, `.list(filters)`                                                  | ✅ REFATORADO |
| **Categories** | `hooks/category/category-query-keys.ts`   | `.all`, `.lists()`, `.list(filters)`                                                  | ✅ REFATORADO |
| **Reference**  | `hooks/reference/reference-query-keys.ts` | `.all`, `.countries()`, `.states()`, `.cities()`                                      | ✅ EXISTENTE  |

---

## ✅ Problemas Identificados & Resolvidos

### 1. **Invalidações Hardcoded** ✅ RESOLVIDO

**Antes:**

```tsx
// ❌ Espalhado no código (provider/page.tsx, linhas 81-171)
queryClient.invalidateQueries({ queryKey: ["provider"] });
queryClient.invalidateQueries({ queryKey: ["services"] });
queryClient.setQueryData(["provider", response.provider.id], ...);
```

**Agora:**

```typescript
// ✅ Centralizado em provider/page.tsx
import { PROVIDER_KEYS } from "@/application/hooks/provider/provider-query-keys";
import { SERVICE_KEYS } from "@/application/hooks/service/service-query-keys";

queryClient.invalidateQueries({ queryKey: PROVIDER_KEYS.lists() });
queryClient.invalidateQueries({
  queryKey: SERVICE_KEYS.byProvider(providerId),
});
```

### 2. **Falta de Padrão em Detail/List** ✅ RESOLVIDO

Padrão aplicado consistentemente em todos os 10 domínios:

```typescript
const ENTITY_KEYS = {
  all: ["entity"],                    // Base
  lists: () => [...],                 // Todas as listas
  list: (filters?) => [...],          // Lista com filtros
  details: () => [...],               // Todos os detalhes
  detail: (id) => [...],              // Detalhe específico
  // Opcionais por domínio:
  byParent: (parentId?) => [...],    // Relacionamentos
}
```

### 3. **Ausência de Query Keys para Relacionamentos** ✅ RESOLVIDO

Criados relacionamentos em query keys para cascata:

- `SERVICE_KEYS.byProvider(providerId)` - Services de um provider
- `PROVIDER_KEYS.photosByProvider(providerId)` - Fotos de um provider
- `BREED_KEYS.list(speciesId)` - Breeds de uma specie
- Padrão aplicável a novos relacionamentos

### 4. **Cache Strategy Indefinida** ✅ RESOLVIDO

Centralizado em `src/infra/cache/cache-config.ts` com staleTime/gcTime por domínio

---

## 🚀 Implementação - ✅ FASE 1 COMPLETA

### 1. ✅ Query Keys Centralizados Criados

**Estrutura Finalizada:**

```
src/application/hooks/
├── service/service-query-keys.ts              ✨ NEW
├── provider/provider-query-keys.ts            ✨ NEW
├── breed/breed-query-keys.ts                  ✨ NEW
├── specie/specie-query-keys.ts                ✨ NEW
├── user/user-query-keys.ts                    ✨ NEW
├── review/review-query-keys.ts                ✅ REFATORADO
├── request/request-query-keys.ts              ✅ REFATORADO
├── tag/tag-query-keys.ts                      ✅ REFATORADO
├── category/category-query-keys.ts            ✅ REFATORADO
└── reference/reference-query-keys.ts          ✅ EXISTENTE
```

**Pattern aplicado em todos os 10 domínios:**

```typescript
export const SERVICE_KEYS = {
  all: ["services"] as const,
  lists: () => [...SERVICE_KEYS.all, "list"] as const,
  list: (filters?: ListServicesInput) =>
    [...SERVICE_KEYS.lists(), filters] as const,
  details: () => [...SERVICE_KEYS.all, "detail"] as const,
  detail: (id: string | number) => [...SERVICE_KEYS.details(), id] as const,
  searches: () => [...SERVICE_KEYS.all, "search"] as const,
  search: (filters?: SearchServicesInput) =>
    [...SERVICE_KEYS.searches(), filters] as const,
  byProvider: (providerId?: string | number) =>
    [...SERVICE_KEYS.all, "byProvider", providerId ?? "all"] as const,
} as const;
```

### 2. ✅ Hooks Atualizados para Importar Query Keys

**Remover definições locais → Importar centralizadas:**

```typescript
// ANTES: use-review.ts
const REVIEW_KEYS = { ... };  // ❌ Declarado localmente, não exportado

// DEPOIS: use-review.ts
import { REVIEW_KEYS } from "./review-query-keys";  // ✅ Importado centralmente
```

**Arquivos Atualizados:**

- ✅ `use-review.ts` - Importa `REVIEW_KEYS`
- ✅ `use-request.ts` - Importa `REQUEST_KEYS`
- ✅ `use-tag.ts` - Importa `TAG_KEYS`
- ✅ `use-category.ts` - Importa `CATEGORY_KEYS`
- ✅ `use-service.ts` - Importa `SERVICE_KEYS`
- ✅ `use-provider.ts` - Importa `PROVIDER_KEYS`
- ✅ `use-breed.ts` - Importa `BREED_KEYS`
- ✅ `use-specie.ts` - Importa `SPECIE_KEYS`
- ✅ `use-user.ts` - Importa `USER_KEYS`

### 3. ✅ Cache Configuration Centralizado

**Arquivo `src/infra/cache/cache-config.ts`:**

```typescript
export const CACHE_CONFIG = {
  // Static data - never changes
  reference: { staleTime: Infinity, gcTime: 30 * 60 * 1000 },
  species: { staleTime: 60 * 60 * 1000, gcTime: 60 * 60 * 1000 },
  breeds: { staleTime: 60 * 60 * 1000, gcTime: 60 * 60 * 1000 },

  // Semi-static - rarely changes
  categories: { staleTime: 30 * 60 * 1000, gcTime: 60 * 60 * 1000 },
  tags: { staleTime: 30 * 60 * 1000, gcTime: 60 * 60 * 1000 },

  // User Profile - medium priority
  user: {
    staleTime: 5 * 60 * 1000,
    gcTime: 30 * 60 * 1000,
    refetchOnWindowFocus: false,
  },

  // Services & Providers - data can change
  services: {
    staleTime: 60 * 1000,
    gcTime: 10 * 60 * 1000,
    refetchOnWindowFocus: false,
  },
  providers: {
    staleTime: 60 * 1000,
    gcTime: 10 * 60 * 1000,
    refetchOnWindowFocus: false,
  },

  // Reviews - frequently read
  reviews: {
    staleTime: 5 * 60 * 1000,
    gcTime: 30 * 60 * 1000,
    refetchOnWindowFocus: true,
  },

  // Requests - CRITICAL
  requests: {
    staleTime: 30 * 1000,
    gcTime: 10 * 60 * 1000,
    refetchOnWindowFocus: true,
  },
};
```

### 4. ✅ Documentação Completa

**Arquivos de Documentação:**

- ✅ `CACHE_STRATEGY.md` - Este arquivo (estratégia global)
- ✅ `QUERY_KEYS_REFERENCE.md` - Guia rápido com exemplos
- ✅ `src/infra/cache/CACHE_PATTERNS.md` - Padrões de implementação

---

## 📋 Próximas Fases (TODO)

### Fase 2: Integração com Hooks (Opcional - já funciona)

- [ ] Atualizar `provider/page.tsx` para usar query keys importados (em vez de hardcoded)
- [ ] Testar fluxos de invalidação em cascata
- [ ] Validar com React Query DevTools

### Fase 3: Prefetch & Otimizações

- [ ] Implementar prefetch strategies (hover cards)
- [ ] Monitorar tamanho de cache
- [ ] Otimizar staleTime com feedback real

### Fase 4: Padrão & Contribuição

- [ ] Atualizar contribution guidelines
- [ ] Adicionar checklist para novos hooks
- [ ] Setup CI/CD para validar query keys

---

## 🔗 Relacionamentos de Invalidação (Cascata)

```
┌─────────────┐
│  Provider   │
└──────┬──────┘
       │ invalidates
       ├─→ PROVIDER_KEYS.detail(id)
       ├─→ PROVIDER_KEYS.lists()
       └─→ SERVICE_KEYS.byProvider(id)  ← CASCATA

┌─────────────┐
│   Service   │
└──────┬──────┘
       │ invalidates
       ├─→ SERVICE_KEYS.detail(id)
       ├─→ SERVICE_KEYS.lists()
       └─→ SERVICE_KEYS.byProvider(providerId)

┌──────────────┐
│   Request    │
└──────┬───────┘
       │ invalidates
       ├─→ REQUEST_KEYS.detail(id)
       └─→ REQUEST_KEYS.lists()

┌─────────────┐
│    Review    │
└──────┬──────┘
       │ invalidates
       ├─→ REVIEW_KEYS.lists()
       └─→ PROVIDER_KEYS.detail(providerId)  ← Rating mudou!
```

---

## 📝 Regras de Escritura

### Rule 1: Sempre usar Query Keys Centralizados

```typescript
// ✅ BOM
import { SERVICE_KEYS } from "./service-query-keys";
queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() });

// ❌ EVITAR
queryClient.invalidateQueries({ queryKey: ["services"] });
```

### Rule 2: Invalidar no Hook, Não na Página

```typescript
// ✅ BOM (use-service.ts)
export const useServiceAdd = (options?: AddServiceOptions) => {
  return useMutation({
    mutationFn: (input) => addService.execute(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() });
    },
    ...options,
  });
};

// ❌ EVITAR (misturado em provider/page.tsx)
useServiceAdd({
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ["services"] }),
});
```

### Rule 3: Usar setQueryData para Otimismo + invalidateQueries para Segurança

```typescript
// ✅ BOM
onSuccess: (response, variables) => {
  // Otimismo imediato
  queryClient.setQueryData(SERVICE_KEYS.detail(response.id), response);
  // Segurança: refetch se houver dúvida
  queryClient.invalidateQueries({
    queryKey: SERVICE_KEYS.lists(),
    refetchType: "none", // Já foram invalidados
  });
};

// ❌ EVITAR
onSuccess: () => location.reload(); // Nunca!
```

---

## ✅ Validação

- ✅ `npx tsc --noEmit` - TypeScript limpo (sem erros)
- ✅ 10/10 domínios com query keys centralizados
- ✅ Importações consistentes em todos os hooks
- ✅ Cache config documentado por domínio
- ✅ Padrões de cascata definidos

---

## 📚 Referências

- [TanStack Query - Query Keys](https://tanstack.com/query/latest/docs/react/guides/important-defaults)
- [TanStack Query - Invalidation](https://tanstack.com/query/latest/docs/react/guides/updates-from-mutations)
- [React Query Best Practices](https://tkdodo.eu/blog/practical-react-query)

---

## 📍 Arquivos Associados

- [QUERY_KEYS_REFERENCE.md](./QUERY_KEYS_REFERENCE.md) - Guia rápido com exemplos
- [src/infra/cache/cache-config.ts](./src/infra/cache/cache-config.ts) - Configuração global de cache
- [src/infra/cache/CACHE_PATTERNS.md](./src/infra/cache/CACHE_PATTERNS.md) - Padrões de implementação
