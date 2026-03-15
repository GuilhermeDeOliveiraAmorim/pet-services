import { AxiosInstance } from "axios";
import { Request } from "@/domain/entities";
import { RequestGateway } from "@/application/ports/request-gateway";
import {
  AddRequestInput,
  GetRequestInput,
  ListRequestsInput,
  ListRequestsOutput,
  AcceptRequestInput,
  RejectRequestInput,
  CompleteRequestInput,
} from "@/application/usecases/request";
import { mapRequestFromApi } from "@/infra/mappers/request-mapper";

export class RequestGatewayAxios implements RequestGateway {
  constructor(private readonly httpClient: AxiosInstance) {}

  private extractRequestFromActionResponse(
    data: Record<string, unknown>,
  ): Record<string, unknown> {
    const request = data.request;

    if (request && typeof request === "object") {
      return request as Record<string, unknown>;
    }

    return data;
  }

  async addRequest(input: AddRequestInput): Promise<Request> {
    const payload = {
      provider_id: input.providerId,
      service_id: input.serviceId,
      pet_id: input.petId,
      notes: input.notes,
    };

    const response = await this.httpClient.post<Record<string, unknown>>(
      "/requests",
      payload,
    );

    return mapRequestFromApi(response.data);
  }

  async getRequest(input: GetRequestInput): Promise<Request> {
    const response = await this.httpClient.get<Record<string, unknown>>(
      `/requests/${input.id}`,
    );

    return mapRequestFromApi(response.data);
  }

  async listRequests(input?: ListRequestsInput): Promise<ListRequestsOutput> {
    const params = new URLSearchParams();

    if (input?.userId) params.append("user_id", input.userId);
    if (input?.providerId) params.append("provider_id", input.providerId);
    if (input?.serviceId) params.append("service_id", input.serviceId);
    if (input?.status) params.append("status", input.status);
    if (input?.page) params.append("page", input.page.toString());
    if (input?.pageSize) params.append("page_size", input.pageSize.toString());

    const queryString = params.toString();
    const url = queryString ? `/requests?${queryString}` : "/requests";

    const response = await this.httpClient.get<{
      requests?: Record<string, unknown>[];
      total?: number;
      page?: number;
      page_size?: number;
    }>(url);

    const requests = (response.data.requests || []).map(mapRequestFromApi);

    return {
      requests,
      total: response.data.total || 0,
      page: response.data.page || 1,
      pageSize: response.data.page_size || 10,
    };
  }

  async acceptRequest(input: AcceptRequestInput): Promise<Request> {
    const response = await this.httpClient.patch<Record<string, unknown>>(
      `/requests/${input.id}/accept`,
    );

    return mapRequestFromApi(
      this.extractRequestFromActionResponse(response.data),
    );
  }

  async rejectRequest(input: RejectRequestInput): Promise<Request> {
    const payload = {
      reason: input.rejectReason,
    };

    const response = await this.httpClient.patch<Record<string, unknown>>(
      `/requests/${input.id}/reject`,
      payload,
    );

    return mapRequestFromApi(
      this.extractRequestFromActionResponse(response.data),
    );
  }

  async completeRequest(input: CompleteRequestInput): Promise<Request> {
    const response = await this.httpClient.patch<Record<string, unknown>>(
      `/requests/${input.id}/complete`,
    );

    return mapRequestFromApi(
      this.extractRequestFromActionResponse(response.data),
    );
  }
}
