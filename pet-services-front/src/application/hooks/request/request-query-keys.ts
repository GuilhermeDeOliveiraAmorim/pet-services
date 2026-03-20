import type { ListRequestsInput } from "@/application/usecases/request";

export const REQUEST_KEYS = {
  all: ["requests"] as const,
  lists: () => [...REQUEST_KEYS.all, "list"] as const,
  list: (filters?: ListRequestsInput) =>
    [...REQUEST_KEYS.lists(), filters] as const,
  details: () => [...REQUEST_KEYS.all, "detail"] as const,
  detail: (id: string) => [...REQUEST_KEYS.details(), id] as const,
} as const;
