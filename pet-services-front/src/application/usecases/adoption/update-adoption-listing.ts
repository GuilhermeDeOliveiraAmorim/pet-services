import type { AdoptionGateway } from "@/application/ports/adoption-gateway";
import type { AdoptionListing } from "@/domain";

export interface UpdateAdoptionListingInput {
  listingId: string;
  title?: string;
  description?: string;
  adoptionReason?: string;
  sex?: string;
  size?: string;
  ageGroup?: string;
  cityId?: string;
  stateId?: string;
  latitude?: number;
  longitude?: number;
  vaccinated?: boolean;
  neutered?: boolean;
  dewormed?: boolean;
  specialNeeds?: boolean;
  goodWithChildren?: boolean;
  goodWithDogs?: boolean;
  goodWithCats?: boolean;
  requiresHouseScreening?: boolean;
}

export interface UpdateAdoptionListingOutput {
  message?: string;
  listing?: AdoptionListing;
}

export class UpdateAdoptionListingUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input: UpdateAdoptionListingInput,
  ): Promise<UpdateAdoptionListingOutput> {
    return this.adoptionGateway.updateListing(input);
  }
}
