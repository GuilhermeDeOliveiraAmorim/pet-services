# Query Keys & Cache Strategy - Reference Guide

## 🎯 Estrutura Centralizada Criada

### Query Keys Centralizadas

| Domínio    | Arquivo                                                   | Status     |
| ---------- | --------------------------------------------------------- | ---------- |
| Services   | `src/application/hooks/service/service-query-keys.ts`     | ✨ NEW     |
| Providers  | `src/application/hooks/provider/provider-query-keys.ts`   | ✨ NEW     |
| Breeds     | `src/application/hooks/breed/breed-query-keys.ts`         | ✨ NEW     |
| Species    | `src/application/hooks/specie/specie-query-keys.ts`       | ✨ NEW     |
| User       | `src/application/hooks/user/user-query-keys.ts`           | ✨ NEW     |
| Reviews    | `src/application/hooks/review/use-review.ts`              | ✅ Already |
| Requests   | `src/application/hooks/request/use-request.ts`            | ✅ Already |
| Tags       | `src/application/hooks/tag/use-tag.ts`                    | ✅ Already |
| Categories | `src/application/hooks/category/use-category.ts`          | ✅ Already |
| Reference  | `src/application/hooks/reference/reference-query-keys.ts` | ✅ Already |

### Cache Configuration

| Arquivo                             | Conteúdo                     |
| ----------------------------------- | ---------------------------- |
| `src/infra/cache/cache-config.ts`   | staleTime/gcTime por domínio |
| `src/infra/cache/CACHE_PATTERNS.ts` | Exemplos de uso e patterns   |

### Documentation

| Arquivo                   | Conteúdo                 |
| ------------------------- | ------------------------ |
| `CACHE_STRATEGY.md`       | Análise completa + plano |
| `QUERY_KEYS_REFERENCE.md` | Este arquivo             |

---

## 📁 Estrutura de Query Keys

### Pattern Padrão

```typescript
const SERVICE_KEYS = {
  all: ["services"], // Base
  lists: () => [...lists()], // Todas as listas
  list: (filters?) => [...filter], // Lista com filtros
  details: () => [...details()], // Todos os detalhes
  detail: (id) => [...id], // Detalhe específico
  byProvider: (id?) => [...id], // Relacionamento
};
```

### Cada Query Key

| Chave                         | Uso                           | Exemplo                                                         |
| ----------------------------- | ----------------------------- | --------------------------------------------------------------- |
| `SERVICE_KEYS.all`            | Invalida tudo de services     | `queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.all })` |
| `SERVICE_KEYS.lists()`        | Invalida listas (sem details) | `SERVICE_KEYS.lists()` → `["services", "list"]`                 |
| `SERVICE_KEYS.list(input)`    | Lista com filtro específico   | `SERVICE_KEYS.list({ providerId: "123" })`                      |
| `SERVICE_KEYS.detail(id)`     | Detalhe de um service         | `SERVICE_KEYS.detail("service-1")`                              |
| `SERVICE_KEYS.byProvider(id)` | Services de um provider       | Para cascata de invalidação                                     |

---

## 🔄 Padrões de Invalidação

### Hook de Adição (Create)

```typescript
export const useServiceAdd = (options?: AddServiceOptions) => {
  const { addService } = useServiceUseCases();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input) => addService.execute(input),
    onSuccess: (response) => {
      // Invalidar lists (novo service aparece)
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() });
      // Cascata: invalidar services deste provider
      queryClient.invalidateQueries({
        queryKey: SERVICE_KEYS.byProvider(response.service.providerId),
      });
    },
    ...options,
  });
};
```

### Hook de Atualização (Update)

```typescript
export const useServiceUpdate = (options?: UpdateServiceOptions) => {
  const { updateService } = useServiceUseCases();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input) => updateService.execute(input),
    onSuccess: (response) => {
      // Otimismo: atualizar detail imediatamente
      queryClient.setQueryData(SERVICE_KEYS.detail(response.service.id), {
        service: response.service,
      });

      // Invalidar lists (pode ter mudado ordem/filtros)
      queryClient.invalidateQueries({
        queryKey: SERVICE_KEYS.lists(),
        refetchType: "none", // Já foi invalidado
      });
    },
    ...options,
  });
};
```

### Hook de Deleção (Delete)

```typescript
export const useServiceDelete = (options?: DeleteServiceOptions) => {
  const { deleteService } = useServiceUseCases();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (serviceId) => deleteService.execute(serviceId),
    onSuccess: (_, serviceId) => {
      // Remover do cache (já não existe)
      queryClient.removeQueries({
        queryKey: SERVICE_KEYS.detail(serviceId),
      });

      // Invalidar lists (service desapareceu)
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() });
    },
    ...options,
  });
};
```

---

## 📊 Cache Strategy por Domínio

```typescript
// Static data - never changes
reference:     staleTime: ∞,         gcTime: 30 min
species:       staleTime: 1 hour,    gcTime: 1 hour
breeds:        staleTime: 1 hour,    gcTime: 1 hour

// Semi-static - rarely changes
categories:    staleTime: 30 min,    gcTime: 1 hour
tags:          staleTime: 30 min,    gcTime: 1 hour

// Dynamic data
user:          staleTime: 5 min,     gcTime: 30 min    (no refetch on focus)
services:      staleTime: 1 min,     gcTime: 10 min    (no refetch on focus)
providers:     staleTime: 1 min,     gcTime: 10 min    (no refetch on focus)

// Very dynamic
reviews:       staleTime: 5 min,     gcTime: 30 min    (refetch on focus!)
requests:      staleTime: 30 sec,    gcTime: 10 min    (refetch on focus!)
```

**Legenda:**

- `staleTime`: Quanto tempo dados são considerados "fresh" (não refetch)
- `gcTime`: Quanto tempo manter dados em cache após unmount
- `refetch on focus`: Refetch quando usuário retorna para aba

---

## 🔗 Cascata de Invalidação

```
┌─────────────────┐
│  updateProvider │
└────────┬────────┘
         │
         ├─→ PROVIDER_KEYS.detail(id)      [setQueryData]
         ├─→ PROVIDER_KEYS.lists()         [invalidate]
         └─→ SERVICE_KEYS.byProvider(id)   [invalidate] ← CASCATA

┌──────────────────┐
│   updateService  │
└────────┬─────────┘
         │
         ├─→ SERVICE_KEYS.detail(id)        [setQueryData]
         ├─→ SERVICE_KEYS.lists()           [invalidate]
         └─→ SERVICE_KEYS.byProvider(pId)   [invalidate]

┌──────────────────┐
│   createReview   │
└────────┬─────────┘
         │
         ├─→ REVIEW_KEYS.lists()            [invalidate]
         └─→ PROVIDER_KEYS.detail(pId)     [invalidate] ← CASCATA (rating mudou!)

┌──────────────────┐
│   acceptRequest  │
└────────┬─────────┘
         │
         ├─→ REQUEST_KEYS.lists()           [invalidate]
         └─→ REQUEST_KEYS.detail(id)       [invalidate]
```

---

## 💻 Exemplos Práticos

### ✅ Correto (Usar Query Keys)

```typescript
// Em um hook
import { SERVICE_KEYS } from '@/application/hooks/service/service-query-keys';

export const useServiceUpdate = (options?: UpdateServiceOptions) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input) => updateService.execute(input),
    onSuccess: (response) => {
      // ✅ Usar query keys centralizados
      queryClient.setQueryData(SERVICE_KEYS.detail(response.service.id), ...);
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() });
    },
    ...options,
  });
};

// Em um componente
const { data: service } = useServiceGet(serviceId);
const { mutateAsync: update } = useServiceUpdate();

const handleUpdate = async () => {
  await update({ id: serviceId, name: "new name" });
  // Cache é invalidado automaticamente via hook
};
```

### ❌ Evitar (Hardcoded)

```typescript
// ❌ NÃO FAZER
export const useServiceUpdate = (options?: UpdateServiceOptions) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input) => updateService.execute(input),
    onSuccess: (response) => {
      // ❌ Hardcoded - difícil manter, sem centralização
      queryClient.invalidateQueries({ queryKey: ["services"] });
      queryClient.invalidateQueries({ queryKey: ["services", "list"] });
    },
    ...options,
  });
};
```

---

## 🚀 Próximos Passos (TODO)

### Fase 1: Integração Imediata ✨ NOW

- [ ] Atualizar `use-service.ts` para usar `SERVICE_KEYS`
- [ ] Atualizar `use-provider.ts` para usar `PROVIDER_KEYS`
- [ ] Atualizar `use-breed.ts` para usar `BREED_KEYS`
- [ ] Atualizar `use-specie.ts` para usar `SPECIE_KEYS`
- [ ] Atualizar `use-user.ts` para usar `USER_KEYS`
- [ ] Remover invalidações hardcoded em `provider/page.tsx`
- [ ] Executar testes e lint

### Fase 2: Cascata e Relacionamentos 🔗

- [ ] Implementar cascata de invalidação (provider → services)
- [ ] Implementar cascata de invalidação (review → provider rating)
- [ ] Testar fluxos: criar/atualizar/deletar

### Fase 3: Performance 📈

- [ ] Implementar prefetch strategies (hover cards, navegação)
- [ ] Otimizar staleTime por domínio com feedback do usuário
- [ ] Monitorar tamanho de cache (DevTools)

### Fase 4: Documentação & Padrão 📚

- [ ] Atualizar contribution guidelines
- [ ] Criar checklist para novos hooks
- [ ] Adicionar exemplos em codebase

---

## 🔍 Quick Reference - Query Keys Map

```typescript
// Services
SERVICE_KEYS.all; // ["services"]
SERVICE_KEYS.lists(); // ["services", "list"]
SERVICE_KEYS.list({ providerId }); // ["services", "list", { providerId }]
SERVICE_KEYS.details(); // ["services", "detail"]
SERVICE_KEYS.detail("id"); // ["services", "detail", "id"]
SERVICE_KEYS.searches(); // ["services", "search"]
SERVICE_KEYS.search({ query }); // ["services", "search", { query }]
SERVICE_KEYS.byProvider("providerId"); // ["services", "byProvider", "providerId"]

// Providers
PROVIDER_KEYS.all; // ["providers"]
PROVIDER_KEYS.lists(); // ["providers", "list"]
PROVIDER_KEYS.details(); // ["providers", "detail"]
PROVIDER_KEYS.detail("id"); // ["providers", "detail", "id"]

// Adoption (Adoção & Guardian)
ADOPTION_KEYS.all; // ["adoption"]
ADOPTION_KEYS.lists(); // ["adoption", "list"]
ADOPTION_KEYS.list(filters); // ["adoption", "list", filters]
ADOPTION_KEYS.details(); // ["adoption", "detail"]
ADOPTION_KEYS.detail(listingId); // ["adoption", "detail", listingId]
ADOPTION_KEYS.myApplicationsLists(); // ["adoption", "my-applications"]
ADOPTION_KEYS.myApplicationsList(filters); // ["adoption", "my-applications", filters]
ADOPTION_KEYS.myListingsLists(); // ["adoption", "my-listings"]
ADOPTION_KEYS.myListingsList(filters); // ["adoption", "my-listings", filters]
ADOPTION_KEYS.listingApplicationsLists(); // ["adoption", "listing-applications"]
ADOPTION_KEYS.listingApplicationsList(input); // ["adoption", "listing-applications", input]
ADOPTION_KEYS.myGuardianProfile(); // ["adoption", "my-guardian-profile"]
ADOPTION_KEYS.management(); // ["adoption", "management"]
```

---

## 🐾 Exemplos de Uso — Adoção & Guardian

### Consulta de listagens públicas de adoção

```typescript
const { data } = usePublicAdoptionListings({ specieId: "dog" });
// queryKey: ADOPTION_KEYS.list({ specieId: "dog" })
```

### Consulta do perfil de responsável (guardian)

```typescript
const { data } = useMyAdoptionGuardianProfile();
// queryKey: ADOPTION_KEYS.myGuardianProfile()
```

### Criação de candidatura à adoção

```typescript
const mutation = useAdoptionApplicationCreate();
mutation.mutate({ listingId, message });
// onSuccess: invalidateQueries(ADOPTION_KEYS.myApplicationsLists())
```

### Atualização de listagem de adoção

```typescript
const mutation = useAdoptionListingUpdate();
mutation.mutate({ listingId, ... });
// onSuccess: invalidateQueries(ADOPTION_KEYS.myListingsLists())
//            invalidateQueries(ADOPTION_KEYS.detail(listingId))
```

### Padrão de Invalidação/Cascata

- Sempre use as keys do domínio para garantir atualização correta dos dados do usuário logado e dos detalhes/listas públicas.
- Exemplo: ao criar/editar uma listagem, invalide tanto as listas do usuário quanto o detalhe da listagem.

---

---

## 📞 Support

Dúvidas? Ver:

- `CACHE_STRATEGY.md` - Análise e plano completo
- `src/infra/cache/CACHE_PATTERNS.ts` - Exemplos detalhados
- TanStack Query Docs: https://tanstack.com/query/latest
