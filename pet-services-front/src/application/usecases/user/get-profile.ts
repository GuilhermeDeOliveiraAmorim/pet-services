import type { User } from "@/domain";
import type { UserGateway } from "@/application/ports";

export interface GetProfileOutput {
  user: User;
  providerId?: string;
}

export class GetProfileUseCase {
  constructor(private readonly userGateway: UserGateway) {}

  execute(): Promise<GetProfileOutput> {
    return this.userGateway.getProfile();
  }
}
