import type { AdoptionGateway } from "@/application/ports/adoption-gateway";

export interface AdoptionApplicationItem {
  id: string;
  listingId: string;
  status: string;
  motivation: string;
  contactPhone: string;
  reviewedAt?: string;
  createdAt: string;
}

export interface AdoptionApplicationsPagination {
  currentPage: number;
  perPage: number;
  totalPages: number;
  totalRecords: number;
}

export interface ListMyAdoptionApplicationsInput {
  page?: number;
  pageSize?: number;
}

export interface ListMyAdoptionApplicationsOutput {
  applications: AdoptionApplicationItem[];
  pagination: AdoptionApplicationsPagination;
}

export class ListMyAdoptionApplicationsUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input?: ListMyAdoptionApplicationsInput,
  ): Promise<ListMyAdoptionApplicationsOutput> {
    return this.adoptionGateway.listMyApplications(input);
  }
}
