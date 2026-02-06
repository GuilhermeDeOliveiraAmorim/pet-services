import type { UserGateway } from "@/application/ports";

export interface ChangePasswordInput {
  userId: string;
  oldPassword: string;
  newPassword: string;
}

export interface ChangePasswordInputBody {
  oldPassword: string;
  newPassword: string;
}

export interface ChangePasswordOutput {
  message?: string;
  detail?: string;
}

export class ChangePasswordUseCase {
  constructor(private readonly userGateway: UserGateway) {}

  execute(input: ChangePasswordInput): Promise<ChangePasswordOutput> {
    return this.userGateway.changePassword(input);
  }
}
