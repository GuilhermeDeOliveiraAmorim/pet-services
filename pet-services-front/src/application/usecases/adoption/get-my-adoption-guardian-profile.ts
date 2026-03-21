import type { AdoptionGateway } from "@/application/ports/adoption-gateway";

export interface GetMyAdoptionGuardianProfileInput {
  forceRefresh?: boolean;
}

export interface AdoptionGuardianProfile {
  id: string;
  userId: string;
  displayName: string;
  guardianType: string;
  document: string;
  phone: string;
  whatsapp: string;
  about: string;
  cityId: string;
  stateId: string;
  approvalStatus: string;
  approvedBy?: string;
  approvedAt?: string;
}

export interface GetMyAdoptionGuardianProfileOutput {
  profile: AdoptionGuardianProfile;
}

export class GetMyAdoptionGuardianProfileUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input?: GetMyAdoptionGuardianProfileInput,
  ): Promise<GetMyAdoptionGuardianProfileOutput | null> {
    return this.adoptionGateway.getMyGuardianProfile(input);
  }
}
