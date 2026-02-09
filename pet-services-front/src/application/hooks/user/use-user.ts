import { useMemo } from "react";
import { useMutation, type UseMutationOptions } from "@tanstack/react-query";

import type {
  ChangePasswordInput,
  ChangePasswordOutput,
  RegisterUserInput,
  RegisterUserOutput,
} from "@/application";
import { createUserUseCases } from "@/application/factories/user-usecase-factory";
import { createApiContext } from "@/infra";

const useUserUseCases = () => {
  return useMemo(() => {
    const { userGateway } = createApiContext();
    return createUserUseCases(userGateway);
  }, []);
};

type RegisterUserOptions = Omit<
  UseMutationOptions<RegisterUserOutput, Error, RegisterUserInput>,
  "mutationFn"
>;

type ChangePasswordOptions = Omit<
  UseMutationOptions<ChangePasswordOutput, Error, ChangePasswordInput>,
  "mutationFn"
>;

export const useUserRegister = (options?: RegisterUserOptions) => {
  const { registerUserUseCase } = useUserUseCases();

  return useMutation({
    mutationFn: (input) => registerUserUseCase.execute(input),
    ...options,
  });
};

export const useUserChangePassword = (options?: ChangePasswordOptions) => {
  const { changePasswordUseCase } = useUserUseCases();

  return useMutation({
    mutationFn: (input) => changePasswordUseCase.execute(input),
    ...options,
  });
};
