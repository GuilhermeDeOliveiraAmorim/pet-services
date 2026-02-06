import type { AuthGateway } from "@/application/ports";

export interface RequestPasswordResetInput {
  email: string;
  userAgent?: string;
  ip?: string;
}

export interface RequestPasswordResetOutput {
  message?: string;
  detail?: string;
  resetToken?: string;
  expiresAt?: string;
}

export class RequestPasswordResetUseCase {
  constructor(private readonly authGateway: AuthGateway) {}

  execute(input: RequestPasswordResetInput): Promise<RequestPasswordResetOutput> {
    return this.authGateway.requestPasswordReset(input);
  }
}
