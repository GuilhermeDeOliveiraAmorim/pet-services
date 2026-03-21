import type { AdoptionGateway } from "@/application/ports/adoption-gateway";

import type { AdoptionGuardianProfile } from "./get-my-adoption-guardian-profile";

export interface CreateAdoptionGuardianProfileInput {
  displayName: string;
  guardianType: string;
  document: string;
  phone: string;
  whatsapp?: string;
  about?: string;
  cityId: string;
  stateId: string;
}

export interface CreateAdoptionGuardianProfileOutput {
  message?: string;
  detail?: string;
  profile?: AdoptionGuardianProfile;
}

export class CreateAdoptionGuardianProfileUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input: CreateAdoptionGuardianProfileInput,
  ): Promise<CreateAdoptionGuardianProfileOutput> {
    return this.adoptionGateway.createGuardianProfile(input);
  }
}
