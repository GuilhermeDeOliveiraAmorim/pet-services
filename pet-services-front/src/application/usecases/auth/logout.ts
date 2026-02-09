import type { AuthGateway } from "@/application/ports";

export interface LogoutInput {
  userId?: string;
  tokenId?: string;
  revokeAll: boolean;
}

export interface LogoutInputBody {
  tokenId?: string;
  revokeAll: boolean;
}

export interface LogoutOutput {
  message?: string;
  detail?: string;
}

export class LogoutUseCase {
  constructor(private readonly authGateway: AuthGateway) {}

  execute(input: LogoutInput): Promise<LogoutOutput> {
    return this.authGateway.logout(input);
  }
}
