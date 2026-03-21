import type { AdoptionGateway } from "@/application/ports/adoption-gateway";
import type { AdoptionListing } from "@/domain";

export interface ListMyAdoptionListingsInput {
  page?: number;
  pageSize?: number;
}

export interface ListMyAdoptionListingsOutput {
  listings: AdoptionListing[];
  pagination: {
    currentPage: number;
    perPage: number;
    totalPages: number;
    totalRecords: number;
  };
}

export class ListMyAdoptionListingsUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input?: ListMyAdoptionListingsInput,
  ): Promise<ListMyAdoptionListingsOutput> {
    return this.adoptionGateway.listMyListings(input);
  }
}
