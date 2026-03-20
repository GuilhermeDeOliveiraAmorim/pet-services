export const PROVIDER_KEYS = {
  all: ["providers"] as const,

  // List queries
  lists: () => [...PROVIDER_KEYS.all, "list"] as const,
  list: () => [...PROVIDER_KEYS.lists()] as const,

  // Detail queries
  details: () => [...PROVIDER_KEYS.all, "detail"] as const,
  detail: (id: string | number) => [...PROVIDER_KEYS.details(), id] as const,

  // Photos related to a provider
  photos: () => [...PROVIDER_KEYS.all, "photos"] as const,
  photosByProvider: (providerId: string | number) =>
    [...PROVIDER_KEYS.photos(), providerId] as const,
} as const;
