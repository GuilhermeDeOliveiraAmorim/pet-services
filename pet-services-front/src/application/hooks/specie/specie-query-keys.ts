export const SPECIE_KEYS = {
  all: ["species"] as const,

  lists: () => [...SPECIE_KEYS.all, "list"] as const,
  list: () => [...SPECIE_KEYS.lists()] as const,
} as const;
