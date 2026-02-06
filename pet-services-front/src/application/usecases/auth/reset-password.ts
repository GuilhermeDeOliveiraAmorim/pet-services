import type { AuthGateway } from "@/application/ports";

export interface ResetPasswordInput {
  token: string;
  newPassword: string;
}

export interface ResetPasswordOutput {
  message?: string;
  detail?: string;
}

export class ResetPasswordUseCase {
  constructor(private readonly authGateway: AuthGateway) {}

  execute(input: ResetPasswordInput): Promise<ResetPasswordOutput> {
    return this.authGateway.resetPassword(input);
  }
}
