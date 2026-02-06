import type {
  AuthGateway,
  LoginInput,
  LoginOutput,
  LogoutInput,
  LogoutOutput,
  RefreshTokenInput,
  RefreshTokenOutput,
  RequestPasswordResetInput,
  RequestPasswordResetOutput,
  ResendVerificationEmailInput,
  ResendVerificationEmailOutput,
  ResetPasswordInput,
  ResetPasswordOutput,
  VerifyEmailInput,
  VerifyEmailOutput,
} from "@/application";
import type { AxiosInstance } from "axios";

export class AuthGatewayAxios implements AuthGateway {
  constructor(private readonly http: AxiosInstance) {}

  setAccessToken(token?: string) {
    if (token) {
      this.http.defaults.headers.common.Authorization = `Bearer ${token}`;
      return;
    }

    delete this.http.defaults.headers.common.Authorization;
  }

  async login(input: LoginInput): Promise<LoginOutput> {
    const { data } = await this.http.post<LoginOutput>("/auth/login", input);
    return data;
  }

  async refreshToken(input: RefreshTokenInput): Promise<RefreshTokenOutput> {
    const { data } = await this.http.post<RefreshTokenOutput>("/auth/refresh", input);
    return data;
  }

  async logout(input: LogoutInput): Promise<LogoutOutput> {
    const { data } = await this.http.post<LogoutOutput>("/auth/logout", input);
    return data;
  }

  async requestPasswordReset(input: RequestPasswordResetInput): Promise<RequestPasswordResetOutput> {
    const { data } = await this.http.post<RequestPasswordResetOutput>("/auth/request-password-reset", input);
    return data;
  }

  async resetPassword(input: ResetPasswordInput): Promise<ResetPasswordOutput> {
    const { data } = await this.http.post<ResetPasswordOutput>("/auth/reset-password", input);
    return data;
  }

  async resendVerificationEmail(input: ResendVerificationEmailInput): Promise<ResendVerificationEmailOutput> {
    const { data } = await this.http.post<ResendVerificationEmailOutput>(
      "/auth/resend-verification-email",
      input,
    );
    return data;
  }

  async verifyEmail(input: VerifyEmailInput): Promise<VerifyEmailOutput> {
    const { data } = await this.http.post<VerifyEmailOutput>("/auth/verify-email", input);
    return data;
  }
}
