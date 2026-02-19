import { ServiceGateway } from "@/application/ports";
import type { Service } from "@/domain";

export interface SearchServicesInput {
  query?: string;
  categoryId?: string;
  tagId?: string;
  latitude?: number;
  longitude?: number;
  radiusKm?: number;
  priceMin?: number;
  priceMax?: number;
  page?: number;
  pageSize?: number;
}

export interface SearchServicesOutput {
  services: Service[];
  total: number;
}

export class SearchServicesUseCase {
  constructor(private readonly serviceGateway: ServiceGateway) {}

  execute(input?: SearchServicesInput): Promise<SearchServicesOutput> {
    return this.serviceGateway.searchServices(input);
  }
}
