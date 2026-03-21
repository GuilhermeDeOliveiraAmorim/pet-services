import type { AdoptionGateway } from "@/application/ports/adoption-gateway";

import type { AdoptionGuardianProfile } from "./get-my-adoption-guardian-profile";

export interface UpdateAdoptionGuardianProfileInput {
  displayName?: string;
  guardianType?: string;
  document?: string;
  phone?: string;
  whatsapp?: string;
  about?: string;
  cityId?: string;
  stateId?: string;
}

export interface UpdateAdoptionGuardianProfileOutput {
  message?: string;
  profile?: AdoptionGuardianProfile;
}

export class UpdateAdoptionGuardianProfileUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input: UpdateAdoptionGuardianProfileInput,
  ): Promise<UpdateAdoptionGuardianProfileOutput> {
    return this.adoptionGateway.updateGuardianProfile(input);
  }
}
