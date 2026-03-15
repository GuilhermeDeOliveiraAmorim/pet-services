import { useMemo } from "react";
import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseQueryOptions,
} from "@tanstack/react-query";
import {
  CreateReviewInput,
  ListReviewsInput,
  ListReviewsOutput,
} from "@/application/usecases/review";
import { createReviewCases } from "@/application/factories/review-usecase-factory";
import { createApiContext } from "@/infra";

const useReviewUseCases = () => {
  return useMemo(() => {
    const { reviewGateway } = createApiContext();
    return createReviewCases(reviewGateway);
  }, []);
};

const REVIEW_KEYS = {
  all: ["reviews"] as const,
  lists: () => [...REVIEW_KEYS.all, "list"] as const,
  list: (filters?: ListReviewsInput) =>
    [...REVIEW_KEYS.lists(), filters] as const,
};

type ListReviewsOptions = Omit<
  UseQueryOptions<ListReviewsOutput, Error>,
  "queryKey" | "queryFn"
>;

export function useReviewCreate() {
  const queryClient = useQueryClient();
  const { createReview } = useReviewUseCases();

  return useMutation({
    mutationFn: (input: CreateReviewInput) => createReview.execute(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: REVIEW_KEYS.lists() });
    },
  });
}

export function useReviewList(
  input?: ListReviewsInput,
  options?: ListReviewsOptions,
) {
  const { listReviews } = useReviewUseCases();

  return useQuery({
    queryKey: REVIEW_KEYS.list(input),
    queryFn: () => listReviews.execute(input),
    ...options,
  });
}
