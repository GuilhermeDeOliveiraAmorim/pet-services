import type { AdoptionGateway } from "@/application/ports/adoption-gateway";

export interface CreateAdoptionApplicationInput {
  listingId: string;
  motivation: string;
  housingType?: string;
  petExperience?: string;
  contactPhone?: string;
  familyMembers?: number;
  agreesHomeVisit?: boolean;
  hasOtherPets?: boolean;
}

export interface CreateAdoptionApplicationOutput {
  id: string;
  listingId: string;
  applicantUserId: string;
  status: string;
}

export class CreateAdoptionApplicationUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input: CreateAdoptionApplicationInput,
  ): Promise<CreateAdoptionApplicationOutput> {
    return this.adoptionGateway.createApplication(input);
  }
}
