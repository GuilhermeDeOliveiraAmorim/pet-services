import type { UserGateway } from "@/application/ports";

export interface DeleteUserOutput {
  message?: string;
  detail?: string;
}

export class DeleteUserUseCase {
  constructor(private readonly userGateway: UserGateway) {}

  execute(): Promise<DeleteUserOutput> {
    return this.userGateway.deleteUser();
  }
}
