import { AxiosInstance } from "axios";
import { TagGateway } from "@/application/ports/tag-gateway";
import { ListTagsInput, ListTagsOutput } from "@/application/usecases/tag";
import { mapTagFromApi } from "@/infra/mappers/tag-mapper";

export class TagGatewayAxios implements TagGateway {
  constructor(private readonly httpClient: AxiosInstance) {}

  async listTags(input?: ListTagsInput): Promise<ListTagsOutput> {
    const params = new URLSearchParams();

    if (input?.name) params.append("name", input.name);
    if (input?.page) params.append("page", input.page.toString());
    if (input?.pageSize) params.append("page_size", input.pageSize.toString());

    const queryString = params.toString();
    const url = queryString ? `/tags?${queryString}` : "/tags";

    const response = await this.httpClient.get<{
      tags?: Record<string, unknown>[];
      Tags?: Record<string, unknown>[];
      total?: number;
      Total?: number;
    }>(url);

    const rawTags = response.data.tags ?? response.data.Tags ?? [];
    const tags = rawTags.map(mapTagFromApi);

    return {
      tags,
      total: response.data.total ?? response.data.Total ?? 0,
      page: input?.page || 1,
      pageSize: input?.pageSize || 10,
    };
  }
}
