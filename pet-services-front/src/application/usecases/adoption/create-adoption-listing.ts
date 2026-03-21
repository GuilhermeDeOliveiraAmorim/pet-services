import type { AdoptionGateway } from "@/application/ports/adoption-gateway";
import type { AdoptionListing } from "@/domain";

export interface CreateAdoptionListingInput {
  petId: string;
  title: string;
  description: string;
  adoptionReason: string;
  sex: string;
  size: string;
  ageGroup: string;
  cityId: string;
  stateId: string;
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

export interface CreateAdoptionListingOutput {
  message?: string;
  listing?: AdoptionListing;
}

export class CreateAdoptionListingUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input: CreateAdoptionListingInput,
  ): Promise<CreateAdoptionListingOutput> {
    return this.adoptionGateway.createListing(input);
  }
}
