import type { AdoptionGateway } from "@/application/ports/adoption-gateway";

export interface WithdrawAdoptionApplicationInput {
  applicationId: string;
}

export interface WithdrawAdoptionApplicationOutput {
  id: string;
  status: string;
}

export class WithdrawAdoptionApplicationUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input: WithdrawAdoptionApplicationInput,
  ): Promise<WithdrawAdoptionApplicationOutput> {
    return this.adoptionGateway.withdrawApplication(input.applicationId);
  }
}
