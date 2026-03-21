import type { AdoptionGateway } from "@/application/ports/adoption-gateway";
import type { AdoptionListing } from "@/domain";

export interface GetPublicAdoptionListingInput {
  listingId: string | number;
}

export interface GetPublicAdoptionListingOutput {
  listing: AdoptionListing;
}

export class GetPublicAdoptionListingUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input: GetPublicAdoptionListingInput,
  ): Promise<GetPublicAdoptionListingOutput> {
    return this.adoptionGateway.getPublicListing(input);
  }
}
