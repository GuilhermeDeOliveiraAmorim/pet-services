export const BREED_KEYS = {
  all: ["breeds"] as const,

  lists: () => [...BREED_KEYS.all, "list"] as const,
  list: (speciesId?: string) =>
    [...BREED_KEYS.lists(), speciesId ?? "all"] as const,
} as const;
