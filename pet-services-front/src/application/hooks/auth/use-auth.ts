import { useMemo } from "react";
import { useMutation, type UseMutationOptions } from "@tanstack/react-query";

import {
  type LoginInput,
  type LoginOutput,
  type RefreshTokenInput,
  type RefreshTokenOutput,
  type LogoutInput,
  type LogoutOutput,
  type RequestPasswordResetInput,
  type RequestPasswordResetOutput,
  type ResetPasswordInput,
  type ResetPasswordOutput,
  type ResendVerificationEmailInput,
  type ResendVerificationEmailOutput,
  type VerifyEmailInput,
  type VerifyEmailOutput,
} from "@/application/usecases/auth";
import { createAuthUseCases } from "@/application/factories/auth-usecase-factory";
import { createApiContext } from "@/infra";
import { useAuthSession } from "./use-auth-session";

const useAuthUseCases = () => {
  return useMemo(() => {
    const { authGateway } = createApiContext();
    return createAuthUseCases(authGateway);
  }, []);
};

type LoginOptions = Omit<
  UseMutationOptions<LoginOutput, Error, LoginInput>,
  "mutationFn"
>;

type RefreshTokenOptions = Omit<
  UseMutationOptions<RefreshTokenOutput, Error, RefreshTokenInput>,
  "mutationFn"
>;

type LogoutOptions = Omit<
  UseMutationOptions<LogoutOutput, Error, LogoutInput>,
  "mutationFn"
>;

type RequestPasswordResetOptions = Omit<
  UseMutationOptions<
    RequestPasswordResetOutput,
    Error,
    RequestPasswordResetInput
  >,
  "mutationFn"
>;

type ResetPasswordOptions = Omit<
  UseMutationOptions<ResetPasswordOutput, Error, ResetPasswordInput>,
  "mutationFn"
>;

type ResendVerificationEmailOptions = Omit<
  UseMutationOptions<
    ResendVerificationEmailOutput,
    Error,
    ResendVerificationEmailInput
  >,
  "mutationFn"
>;

type VerifyEmailOptions = Omit<
  UseMutationOptions<VerifyEmailOutput, Error, VerifyEmailInput>,
  "mutationFn"
>;

export const useAuthLogin = (options?: LoginOptions) => {
  const { loginUseCase } = useAuthUseCases();

  return useMutation({
    mutationFn: (input) => loginUseCase.execute(input),
    ...options,
  });
};

export const useAuthRefreshToken = (options?: RefreshTokenOptions) => {
  const { refreshTokenUseCase } = useAuthUseCases();
  const { setSession } = useAuthSession();

  return useMutation({
    mutationFn: (input) => refreshTokenUseCase.execute(input),
    onSuccess: (data, variables, context, mutation) => {
      const expiresAt = Date.now() + data.expiresIn * 1000;
      setSession({
        accessToken: data.accessToken,
        refreshToken: data.refreshToken,
        expiresAt,
      });

      options?.onSuccess?.(data, variables, context, mutation);
    },
    ...options,
  });
};

export const useAuthLogout = (options?: LogoutOptions) => {
  const { logoutUseCase } = useAuthUseCases();

  return useMutation({
    mutationFn: (input) => logoutUseCase.execute(input),
    ...options,
  });
};

export const useAuthRequestPasswordReset = (
  options?: RequestPasswordResetOptions,
) => {
  const { requestPasswordResetUseCase } = useAuthUseCases();

  return useMutation({
    mutationFn: (input) => requestPasswordResetUseCase.execute(input),
    ...options,
  });
};

export const useAuthResetPassword = (options?: ResetPasswordOptions) => {
  const { resetPasswordUseCase } = useAuthUseCases();

  return useMutation({
    mutationFn: (input) => resetPasswordUseCase.execute(input),
    ...options,
  });
};

export const useAuthResendVerificationEmail = (
  options?: ResendVerificationEmailOptions,
) => {
  const { resendVerificationEmailUseCase } = useAuthUseCases();

  return useMutation({
    mutationFn: (input) => resendVerificationEmailUseCase.execute(input),
    ...options,
  });
};

export const useAuthVerifyEmail = (options?: VerifyEmailOptions) => {
  const { verifyEmailUseCase } = useAuthUseCases();

  return useMutation({
    mutationFn: (input) => verifyEmailUseCase.execute(input),
    ...options,
  });
};
