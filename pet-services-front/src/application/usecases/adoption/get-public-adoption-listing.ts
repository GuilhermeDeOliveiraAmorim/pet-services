import type { AdoptionGateway } from "@/application/ports/adoption-gateway";
import type { AdoptionListing } from "@/domain";

export interface GetPublicAdoptionListingOutput {
  listing: AdoptionListing;
}

export class GetPublicAdoptionListingUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(listingId: string | number): Promise<GetPublicAdoptionListingOutput> {
    return this.adoptionGateway.getPublicListing(listingId);
  }
}
