import type { UserGateway } from "@/application/ports";

export interface DeactivateUserOutput {
  message?: string;
  detail?: string;
}

export class DeactivateUserUseCase {
  constructor(private readonly userGateway: UserGateway) {}

  execute(): Promise<DeactivateUserOutput> {
    return this.userGateway.deactivateUser();
  }
}
