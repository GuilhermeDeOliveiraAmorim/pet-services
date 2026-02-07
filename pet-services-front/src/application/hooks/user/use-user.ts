import { useMemo } from "react";
import { useMutation, type UseMutationOptions } from "@tanstack/react-query";

import type { RegisterUserInput, RegisterUserOutput } from "@/application";
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

export const useUserRegister = (options?: RegisterUserOptions) => {
  const { registerUserUseCase } = useUserUseCases();

  return useMutation({
    mutationFn: (input) => registerUserUseCase.execute(input),
    ...options,
  });
};
