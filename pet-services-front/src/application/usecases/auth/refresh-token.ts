import type { AuthGateway } from "@/application/ports";

export interface RefreshTokenInput {
  refreshToken: string;
  userAgent?: string;
  ip?: string;
}

export interface RefreshTokenOutput {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export class RefreshTokenUseCase {
  constructor(private readonly authGateway: AuthGateway) {}

  execute(input: RefreshTokenInput): Promise<RefreshTokenOutput> {
    return this.authGateway.refreshToken(input);
  }
}
