import { Review } from "@/domain/entities";
import {
  CreateReviewInput,
  ListReviewsInput,
  ListReviewsOutput,
} from "@/application/usecases/review";

export interface ReviewGateway {
  createReview(input: CreateReviewInput): Promise<Review>;
  listReviews(input?: ListReviewsInput): Promise<ListReviewsOutput>;
}
