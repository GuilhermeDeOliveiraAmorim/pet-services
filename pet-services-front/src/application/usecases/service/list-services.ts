import { ServiceGateway } from "@/application/ports";
import type { Service } from "@/domain";

export interface ListServicesInput {
  providerId?: string;
  categoryId?: string;
  tagId?: string;
  priceMin?: number;
  priceMax?: number;
  page?: number;
  pageSize?: number;
}

export interface ListServicesOutput {
  services: Service[];
  total: number;
}

export class ListServicesUseCase {
  constructor(private readonly serviceGateway: ServiceGateway) {}

  execute(input?: ListServicesInput): Promise<ListServicesOutput> {
    return this.serviceGateway.listServices(input);
  }
}
