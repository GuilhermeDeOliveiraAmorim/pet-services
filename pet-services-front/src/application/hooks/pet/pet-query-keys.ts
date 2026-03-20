export const PET_KEYS = {
  all: ["pets"] as const,

  // List queries
  lists: () => [...PET_KEYS.all, "list"] as const,
  list: () => [...PET_KEYS.lists()] as const,

  // Owner queries
  byOwner: () => [...PET_KEYS.all, "owner"] as const,
  byOwnerId: (ownerId?: string) =>
    [...PET_KEYS.byOwner(), ownerId ?? "all"] as const,

  // Detail queries
  details: () => [...PET_KEYS.all, "detail"] as const,
  detail: (id: string | number) => [...PET_KEYS.details(), id] as const,

  // Photos
  photos: () => [...PET_KEYS.all, "photos"] as const,
  photosByPet: (petId: string | number) =>
    [...PET_KEYS.photos(), petId] as const,
} as const;
