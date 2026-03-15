import { AxiosInstance } from "axios";
import { Review } from "@/domain/entities";
import { ReviewGateway } from "@/application/ports/review-gateway";
import {
  CreateReviewInput,
  ListReviewsInput,
  ListReviewsOutput,
} from "@/application/usecases/review";
import { mapReviewFromApi } from "@/infra/mappers/review-mapper";

export class ReviewGatewayAxios implements ReviewGateway {
  constructor(private readonly httpClient: AxiosInstance) {}

  private extractReviewFromActionResponse(
    data: Record<string, unknown>,
  ): Record<string, unknown> {
    const review = data.review;

    if (review && typeof review === "object") {
      return review as Record<string, unknown>;
    }

    return data;
  }

  async createReview(input: CreateReviewInput): Promise<Review> {
    const payload = {
      rating: input.rating,
      comment: input.comment,
    };

    const response = await this.httpClient.post<Record<string, unknown>>(
      `/providers/${input.providerId}/reviews`,
      payload,
    );

    return mapReviewFromApi(
      this.extractReviewFromActionResponse(response.data),
    );
  }

  async listReviews(input?: ListReviewsInput): Promise<ListReviewsOutput> {
    const params = new URLSearchParams();

    if (input?.providerId) params.append("provider_id", input.providerId);
    if (input?.userId) params.append("user_id", input.userId);
    if (input?.page) params.append("page", input.page.toString());
    if (input?.pageSize) params.append("page_size", input.pageSize.toString());

    const queryString = params.toString();
    const url = queryString ? `/reviews?${queryString}` : "/reviews";

    const response = await this.httpClient.get<{
      reviews?: Record<string, unknown>[];
      total_items?: number;
      current_page?: number;
      items_per_page?: number;
      total?: number;
      page?: number;
      page_size?: number;
    }>(url);

    const reviews = (response.data.reviews || []).map(mapReviewFromApi);

    return {
      reviews,
      total: response.data.total_items || 0,
      page: response.data.current_page || 1,
      pageSize: response.data.items_per_page || 10,
    };
  }
}
