import type {
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
} from "../usecases/auth";

export interface AuthGateway {
  login(input: LoginInput): Promise<LoginOutput>;
  refreshToken(input: RefreshTokenInput): Promise<RefreshTokenOutput>;
  logout(input: LogoutInput): Promise<LogoutOutput>;
  requestPasswordReset(input: RequestPasswordResetInput): Promise<RequestPasswordResetOutput>;
  resetPassword(input: ResetPasswordInput): Promise<ResetPasswordOutput>;
  resendVerificationEmail(input: ResendVerificationEmailInput): Promise<ResendVerificationEmailOutput>;
  verifyEmail(input: VerifyEmailInput): Promise<VerifyEmailOutput>;
}
