import type { AuthGateway } from "@/application/ports";

export interface ResendVerificationEmailInput {
  email: string;
  userAgent?: string;
  ip?: string;
}

export interface ResendVerificationEmailOutput {
  message?: string;
  detail?: string;
  verifyToken?: string;
  expiresAt?: string;
}

export class ResendVerificationEmailUseCase {
  constructor(private readonly authGateway: AuthGateway) {}

  execute(input: ResendVerificationEmailInput): Promise<ResendVerificationEmailOutput> {
    return this.authGateway.resendVerificationEmail(input);
  }
}
