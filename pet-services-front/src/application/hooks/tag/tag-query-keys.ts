import type { ListTagsInput } from "@/application/usecases/tag";

export const TAG_KEYS = {
  all: ["tags"] as const,
  lists: () => [...TAG_KEYS.all, "list"] as const,
  list: (filters?: ListTagsInput) => [...TAG_KEYS.lists(), filters] as const,
} as const;
