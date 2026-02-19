import { AxiosInstance } from "axios";
import { CategoryGateway } from "@/application/ports/category-gateway";
import {
  ListCategoriesInput,
  ListCategoriesOutput,
} from "@/application/usecases/category";
import { mapCategoryFromApi } from "@/infra/mappers/category-mapper";

export class CategoryGatewayAxios implements CategoryGateway {
  constructor(private readonly httpClient: AxiosInstance) {}

  async listCategories(
    input?: ListCategoriesInput,
  ): Promise<ListCategoriesOutput> {
    const params = new URLSearchParams();

    if (input?.name) params.append("name", input.name);
    if (input?.page) params.append("page", input.page.toString());
    if (input?.pageSize) params.append("page_size", input.pageSize.toString());

    const queryString = params.toString();
    const url = queryString
      ? `/util/categories?${queryString}`
      : "/util/categories";

    const response = await this.httpClient.get<{
      categories?: Record<string, unknown>[];
      total?: number;
    }>(url);

    const categories = (response.data.categories || []).map(mapCategoryFromApi);

    return {
      categories,
      total: response.data.total || 0,
      page: input?.page || 1,
      pageSize: input?.pageSize || 10,
    };
  }
}
