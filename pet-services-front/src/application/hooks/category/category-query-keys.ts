import type { ListCategoriesInput } from "@/application/usecases/category";

export const CATEGORY_KEYS = {
  all: ["categories"] as const,
  lists: () => [...CATEGORY_KEYS.all, "list"] as const,
  list: (filters?: ListCategoriesInput) =>
    [...CATEGORY_KEYS.lists(), filters] as const,
} as const;
