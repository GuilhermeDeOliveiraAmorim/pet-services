import { Review } from "@/domain/entities";

export interface ListReviewsInput {
  providerId?: string;
  userId?: string;
  page?: number;
  pageSize?: number;
}

export interface ListReviewsOutput {
  reviews: Review[];
  total: number;
  page: number;
  pageSize: number;
}

export interface ListReviewsUseCase {
  execute(input?: ListReviewsInput): Promise<ListReviewsOutput>;
}
