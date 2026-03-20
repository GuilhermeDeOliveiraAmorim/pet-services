import type {
  ListServicesInput,
  SearchServicesInput,
} from "@/application/usecases/service";

export const SERVICE_KEYS = {
  all: ["services"] as const,

  // List queries
  lists: () => [...SERVICE_KEYS.all, "list"] as const,
  list: (filters?: ListServicesInput) =>
    [...SERVICE_KEYS.lists(), filters] as const,

  // Detail queries
  details: () => [...SERVICE_KEYS.all, "detail"] as const,
  detail: (id: string | number) => [...SERVICE_KEYS.details(), id] as const,

  // Search queries (different from list)
  searches: () => [...SERVICE_KEYS.all, "search"] as const,
  search: (filters?: SearchServicesInput) =>
    [...SERVICE_KEYS.searches(), filters] as const,

  // By provider (for cascade invalidation)
  byProvider: (providerId?: string | number) =>
    [...SERVICE_KEYS.all, "byProvider", providerId ?? "all"] as const,
} as const;
