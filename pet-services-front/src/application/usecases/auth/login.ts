import type { User } from "@/domain";
import type { AuthGateway } from "@/application/ports";

export interface LoginInput {
  email: string;
  password: string;
  userAgent?: string;
  ip?: string;
}

export interface LoginOutput {
  user: User;
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export class LoginUserUseCase {
  constructor(private readonly authGateway: AuthGateway) {}

  execute(input: LoginInput): Promise<LoginOutput> {
    return this.authGateway.login(input);
  }
}
