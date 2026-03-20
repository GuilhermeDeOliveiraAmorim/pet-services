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
import { mapUserFromApi } from "@/infra/mappers/user-mapper";

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
    const payload = {
      email: input.email,
      password: input.password,
      user_agent: input.userAgent,
      ip: input.ip,
    };

    const { data } = await this.http.post<{
      user: Record<string, unknown>;
      access_token: string;
      refresh_token: string;
      expires_in: number;
    }>("/auth/login", payload);

    return {
      user: mapUserFromApi(data.user),
      accessToken: data.access_token,
      refreshToken: data.refresh_token,
      expiresIn: data.expires_in,
    };
  }

  async refreshToken(input: RefreshTokenInput): Promise<RefreshTokenOutput> {
    const payload = {
      refresh_token: input.refreshToken,
      user_agent: input.userAgent,
      ip: input.ip,
    };

    const { data } = await this.http.post<{
      access_token: string;
      refresh_token: string;
      expires_in: number;
    }>("/auth/refresh", payload);

    return {
      accessToken: data.access_token,
      refreshToken: data.refresh_token,
      expiresIn: data.expires_in,
    };
  }

  async logout(input: LogoutInput): Promise<LogoutOutput> {
    const payload = {
      token_id: input.tokenId,
      revoke_all: input.revokeAll,
    };

    const { data } = await this.http.post<LogoutOutput>(
      "/auth/logout",
      payload,
    );
    return data;
  }

  async requestPasswordReset(
    input: RequestPasswordResetInput,
  ): Promise<RequestPasswordResetOutput> {
    const payload = {
      email: input.email,
      user_agent: input.userAgent,
      ip: input.ip,
    };

    const { data } = await this.http.post<{
      message?: string;
      detail?: string;
    }>("/auth/request-password-reset", payload);

    return {
      message: data.message,
      detail: data.detail,
    };
  }

  async resetPassword(input: ResetPasswordInput): Promise<ResetPasswordOutput> {
    const payload = {
      token: input.token,
      new_password: input.newPassword,
    };

    const { data } = await this.http.post<ResetPasswordOutput>(
      "/auth/reset-password",
      payload,
    );
    return data;
  }

  async resendVerificationEmail(
    input: ResendVerificationEmailInput,
  ): Promise<ResendVerificationEmailOutput> {
    const payload = {
      email: input.email,
      user_agent: input.userAgent,
      ip: input.ip,
    };

    const { data } = await this.http.post<{
      message?: string;
      detail?: string;
    }>("/auth/resend-verification-email", payload);

    return {
      message: data.message,
      detail: data.detail,
    };
  }

  async verifyEmail(input: VerifyEmailInput): Promise<VerifyEmailOutput> {
    const { data } = await this.http.post<VerifyEmailOutput>(
      "/auth/verify-email",
      input,
    );
    return data;
  }
}
