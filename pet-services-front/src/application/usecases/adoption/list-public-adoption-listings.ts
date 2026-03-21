import type { AdoptionGateway } from "@/application/ports/adoption-gateway";
import type { AdoptionListing } from "@/domain";

export interface ListPublicAdoptionListingsInput {
  page?: number;
  pageSize?: number;
  sex?: string;
  size?: string;
  ageGroup?: string;
  cityId?: string;
  stateId?: string;
}

export interface ListPublicAdoptionListingsOutput {
  listings: AdoptionListing[];
  pagination: {
    currentPage: number;
    perPage: number;
    totalPages: number;
    totalRecords: number;
  };
}

export class ListPublicAdoptionListingsUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input?: ListPublicAdoptionListingsInput,
  ): Promise<ListPublicAdoptionListingsOutput> {
    return this.adoptionGateway.listPublicListings(input);
  }
}
