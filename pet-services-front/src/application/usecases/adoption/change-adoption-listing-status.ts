import type { AdoptionGateway } from "@/application/ports/adoption-gateway";
import type { AdoptionListing } from "@/domain";

export interface ChangeAdoptionListingStatusInput {
  listingId: string;
  action: "publish" | "pause" | "archive";
}

export interface ChangeAdoptionListingStatusOutput {
  message?: string;
  listing?: AdoptionListing;
}

export class ChangeAdoptionListingStatusUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input: ChangeAdoptionListingStatusInput,
  ): Promise<ChangeAdoptionListingStatusOutput> {
    return this.adoptionGateway.changeListingStatus(input);
  }
}
