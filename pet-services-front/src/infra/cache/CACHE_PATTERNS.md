# Cache Patterns & Implementation Guide

## Importações Corretas

Sempre use os arquivos `-query-keys.ts` para importar:

```typescript
import { SERVICE_KEYS } from "@/application/hooks/service/service-query-keys";
import { PROVIDER_KEYS } from "@/application/hooks/provider/provider-query-keys";
import { REVIEW_KEYS } from "@/application/hooks/review/review-query-keys";
import { REQUEST_KEYS } from "@/application/hooks/request/request-query-keys";
import { TAG_KEYS } from "@/application/hooks/tag/tag-query-keys";
import { CATEGORY_KEYS } from "@/application/hooks/category/category-query-keys";
import { BREED_KEYS } from "@/application/hooks/breed/breed-query-keys";
import { SPECIE_KEYS } from "@/application/hooks/specie/specie-query-keys";
import { USER_KEYS } from "@/application/hooks/user/user-query-keys";
```

---

## Padrão 1: Create (Addição)

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
        queryKey: SERVICE_KEYS.byProvider(response.service.providerId) 
      });
    },
    ...options,
  });
};
```

---

## Padrão 2: Update (Modificação)

```typescript
export const useServiceUpdate = (options?: UpdateServiceOptions) => {
  const { updateService } = useServiceUseCases();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input) => updateService.execute(input),
    onSuccess: (response) => {
      // Otimismo: atualizar detail imediatamente
      queryClient.setQueryData(
        SERVICE_KEYS.detail(response.service.id),
        { service: response.service }
      );
      
      // Invalidar lists (pode ter mudado ordem/filtros)
      queryClient.invalidateQueries({ 
        queryKey: SERVICE_KEYS.lists(),
        refetchType: 'none'  // Já foi invalidado
      });
    },
    ...options,
  });
};
```

---

## Padrão 3: Delete (Remoção)

```typescript
export const useServiceDelete = (options?: DeleteServiceOptions) => {
  const { deleteService } = useServiceUseCases();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (serviceId) => deleteService.execute(serviceId),
    onSuccess: (_, serviceId) => {
      // Remover do cache (já não existe)
      queryClient.removeQueries({ 
        queryKey: SERVICE_KEYS.detail(serviceId) 
      });
      
      // Invalidar lists (service desapareceu)
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() });
    },
    ...options,
  });
};
```

---

## Cascata de Invalidação

### Provider Update → Service Invalidation

```typescript
export const useProviderUpdate = (options?: UpdateProviderOptions) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input) => updateProvider.execute(input),
    onSuccess: (response, variables) => {
      const providerId = response.provider.id;

      // 1. Atualizar provider detail
      queryClient.setQueryData(PROVIDER_KEYS.detail(providerId), {
        provider: response.provider,
      });

      // 2. Invalidar provider lists
      queryClient.invalidateQueries({ queryKey: PROVIDER_KEYS.lists() });

      // 3. Invalidar services deste provider (CASCATA)
      queryClient.invalidateQueries({
        queryKey: SERVICE_KEYS.byProvider(providerId),
      });
    },
  });
};
```

### Review Create → Provider Rating Update

```typescript
export const useReviewCreate = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input) => createReview.execute(input),
    onSuccess: (response) => {
      // 1. Invalidar reviews lists
      queryClient.invalidateQueries({ queryKey: REVIEW_KEYS.lists() });

      // 2. Invalidar provider (rating pode ter mudado)
      if (response.review.providerId) {
        queryClient.invalidateQueries({
          queryKey: PROVIDER_KEYS.detail(response.review.providerId),
        });
      }
    },
  });
};
```

---

## Prefetch Strategies

### Prefetch Detail (Hover)

```typescript
const prefetchServiceDetails = (queryClient, serviceId: string) => {
  queryClient.prefetchQuery({
    queryKey: SERVICE_KEYS.detail(serviceId),
    queryFn: () => getService(serviceId),
  });
};
```

### Prefetch Lists (Navigation)

```typescript
const prefetchProviderServices = (queryClient, providerId: string) => {
  queryClient.prefetchQuery({
    queryKey: SERVICE_KEYS.byProvider(providerId),
    queryFn: () => listServices({ providerId }),
  });
};
```

---

## Checklist para Novos Hooks

- [ ] Criar arquivo `*-query-keys.ts` com query keys estruturados
- [ ] Importar query keys no hook: `import { KEYS } from './query-keys'`
- [ ] Usar `KEYS.list()` para queries
- [ ] Usar `KEYS.detail(id)` para detalhes
- [ ] Invalidar no `onSuccess` de mutations
- [ ] Documentar cascatas de invalidação
- [ ] Testar com React Query DevTools
- [ ] Validar TypeScript: `npm run type-check`
