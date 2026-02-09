import type { UserGateway } from "@/application/ports";

export interface ReactivateUserOutput {
  message?: string;
  detail?: string;
}

export class ReactivateUserUseCase {
  constructor(private readonly userGateway: UserGateway) {}

  execute(): Promise<ReactivateUserOutput> {
    return this.userGateway.reactivateUser();
  }
}
