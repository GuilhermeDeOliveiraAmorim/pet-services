import type { AdoptionGateway } from "@/application/ports/adoption-gateway";
import type { AdoptionApplicationItem } from "./list-my-adoption-applications";

export interface ListAdoptionApplicationsByListingInput {
  listingId: string;
  page?: number;
  pageSize?: number;
}

export interface ListAdoptionApplicationsByListingOutput {
  applications: AdoptionApplicationItem[];
  pagination: {
    currentPage: number;
    perPage: number;
    totalPages: number;
    totalRecords: number;
  };
}

export class ListAdoptionApplicationsByListingUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input: ListAdoptionApplicationsByListingInput,
  ): Promise<ListAdoptionApplicationsByListingOutput> {
    return this.adoptionGateway.listApplicationsByListing(input);
  }
}
