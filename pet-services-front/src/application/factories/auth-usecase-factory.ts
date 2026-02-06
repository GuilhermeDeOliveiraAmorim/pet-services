import type { AuthGateway } from "../ports";
import {
  LoginUserUseCase,
  LogoutUseCase,
  RefreshTokenUseCase,
  RequestPasswordResetUseCase,
  ResetPasswordUseCase,
  ResendVerificationEmailUseCase,
  VerifyEmailUseCase,
} from "../usecases/auth";

export const createAuthUseCases = (gateway: AuthGateway) => {
  return {
    loginUseCase: new LoginUserUseCase(gateway),
    refreshTokenUseCase: new RefreshTokenUseCase(gateway),
    logoutUseCase: new LogoutUseCase(gateway),
    requestPasswordResetUseCase: new RequestPasswordResetUseCase(gateway),
    resetPasswordUseCase: new ResetPasswordUseCase(gateway),
    verifyEmailUseCase: new VerifyEmailUseCase(gateway),
    resendVerificationEmailUseCase: new ResendVerificationEmailUseCase(gateway),
  };
};
