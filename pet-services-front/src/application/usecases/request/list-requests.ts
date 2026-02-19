import { Request } from "@/domain/entities";

export interface ListRequestsInput {
  userId?: string;
  providerId?: string;
  serviceId?: string;
  status?: string;
  page?: number;
  pageSize?: number;
}

export interface ListRequestsOutput {
  requests: Request[];
  total: number;
  page: number;
  pageSize: number;
}

export interface ListRequestsUseCase {
  execute(input?: ListRequestsInput): Promise<ListRequestsOutput>;
}
