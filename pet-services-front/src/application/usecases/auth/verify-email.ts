import type { AuthGateway } from "@/application/ports";

export interface VerifyEmailInput {
  token: string;
}

export interface VerifyEmailOutput {
  message?: string;
  detail?: string;
}

export class VerifyEmailUseCase {
  constructor(private readonly authGateway: AuthGateway) {}

  execute(input: VerifyEmailInput): Promise<VerifyEmailOutput> {
    return this.authGateway.verifyEmail(input);
  }
}
