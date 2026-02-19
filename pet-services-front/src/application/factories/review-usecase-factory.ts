import { ReviewGateway } from "@/application/ports/review-gateway";
import {
  CreateReviewUseCase,
  CreateReviewInput,
  ListReviewsUseCase,
  ListReviewsInput,
  ListReviewsOutput,
} from "@/application/usecases/review";

class CreateReviewUseCaseImpl implements CreateReviewUseCase {
  constructor(private reviewGateway: ReviewGateway) {}

  async execute(input: CreateReviewInput) {
    return this.reviewGateway.createReview(input);
  }
}

class ListReviewsUseCaseImpl implements ListReviewsUseCase {
  constructor(private reviewGateway: ReviewGateway) {}

  async execute(input?: ListReviewsInput): Promise<ListReviewsOutput> {
    return this.reviewGateway.listReviews(input);
  }
}

export function createReviewCases(reviewGateway: ReviewGateway) {
  return {
    createReview: new CreateReviewUseCaseImpl(reviewGateway),
    listReviews: new ListReviewsUseCaseImpl(reviewGateway),
  };
}
