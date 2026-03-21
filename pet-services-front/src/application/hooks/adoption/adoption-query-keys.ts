import type { ListPublicAdoptionListingsInput } from "@/application/usecases/adoption";

export const ADOPTION_KEYS = {
  all: ["adoption"] as const,
  lists: () => [...ADOPTION_KEYS.all, "list"] as const,
  list: (filters?: ListPublicAdoptionListingsInput) =>
    [...ADOPTION_KEYS.lists(), filters] as const,
  details: () => [...ADOPTION_KEYS.all, "detail"] as const,
  detail: (listingId: string | number) =>
    [...ADOPTION_KEYS.details(), listingId] as const,
} as const;
