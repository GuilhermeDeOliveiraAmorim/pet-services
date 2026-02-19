import { Review } from "@/domain/entities";

export interface CreateReviewInput {
  providerId: string;
  rating: number;
  comment: string;
}

export interface CreateReviewUseCase {
  execute(input: CreateReviewInput): Promise<Review>;
}
