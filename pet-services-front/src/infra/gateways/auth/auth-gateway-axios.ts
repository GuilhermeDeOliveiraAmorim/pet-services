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
import type { User } from "@/domain";
import type { AxiosInstance } from "axios";

const mapUserFromApi = (user?: Record<string, any> | null): User => {
  const raw = user ?? {};

  return {
    id: raw.id ?? "",
    active: raw.active ?? false,
    createdAt: raw.created_at ?? raw.createdAt ?? "",
    updatedAt: raw.updated_at ?? raw.updatedAt ?? null,
    deactivatedAt: raw.deactivated_at ?? raw.deactivatedAt ?? null,
    name: raw.name ?? "",
    userType: raw.userType ?? raw.user_type ?? "owner",
    login: {
      email: raw.login?.email ?? "",
      password: raw.login?.password ?? "",
    },
    phone: {
      countryCode: raw.phone?.countryCode ?? raw.phone?.country_code ?? "",
      areaCode: raw.phone?.areaCode ?? raw.phone?.area_code ?? "",
      number: raw.phone?.number ?? "",
    },
    address: {
      street: raw.address?.street ?? "",
      number: raw.address?.number ?? "",
      neighborhood: raw.address?.neighborhood ?? "",
      city: raw.address?.city ?? "",
      zipCode: raw.address?.zipCode ?? raw.address?.zip_code ?? "",
      state: raw.address?.state ?? "",
      country: raw.address?.country ?? "",
      complement: raw.address?.complement ?? "",
      location: {
        latitude:
          raw.address?.location?.latitude ?? raw.address?.location?.lat ?? 0,
        longitude:
          raw.address?.location?.longitude ?? raw.address?.location?.lng ?? 0,
      },
    },
    emailVerified: raw.emailVerified ?? raw.email_verified ?? false,
    photos: raw.photos ?? [],
    pets: raw.pets ?? [],
  } as User;
};

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
      user: Record<string, any>;
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
      reset_token?: string;
      expires_at?: string;
    }>("/auth/request-password-reset", payload);

    return {
      message: data.message,
      detail: data.detail,
      resetToken: data.reset_token,
      expiresAt: data.expires_at,
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
      verify_token?: string;
      expires_at?: string;
    }>("/auth/resend-verification-email", payload);

    return {
      message: data.message,
      detail: data.detail,
      verifyToken: data.verify_token,
      expiresAt: data.expires_at,
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
