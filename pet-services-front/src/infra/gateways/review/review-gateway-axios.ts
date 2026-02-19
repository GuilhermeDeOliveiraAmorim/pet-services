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

  async createReview(input: CreateReviewInput): Promise<Review> {
    const payload = {
      provider_id: input.providerId,
      rating: input.rating,
      comment: input.comment,
    };

    const response = await this.httpClient.post<Record<string, unknown>>(
      "/reviews",
      payload
    );

    return mapReviewFromApi(response.data);
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
      total?: number;
      page?: number;
      page_size?: number;
    }>(url);

    const reviews = (response.data.reviews || []).map(mapReviewFromApi);

    return {
      reviews,
      total: response.data.total || 0,
      page: response.data.page || 1,
      pageSize: response.data.page_size || 10,
    };
  }
}
