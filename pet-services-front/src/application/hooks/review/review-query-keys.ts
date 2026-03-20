import type { ListReviewsInput } from "@/application/usecases/review";

export const REVIEW_KEYS = {
  all: ["reviews"] as const,
  lists: () => [...REVIEW_KEYS.all, "list"] as const,
  list: (filters?: ListReviewsInput) =>
    [...REVIEW_KEYS.lists(), filters] as const,
} as const;
