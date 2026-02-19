import { useMemo } from "react";
import {
  useMutation,
  useQuery,
  type UseMutationOptions,
  type UseQueryOptions,
} from "@tanstack/react-query";

import type {
  AddUserPhotoInput,
  AddUserPhotoOutput,
  ChangePasswordInput,
  ChangePasswordOutput,
  DeactivateUserOutput,
  DeleteUserOutput,
  GetProfileOutput,
  ReactivateUserOutput,
  RegisterUserInput,
  RegisterUserOutput,
  UpdateUserInput,
  UpdateUserOutput,
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

type GetProfileOptions = Omit<
  UseQueryOptions<GetProfileOutput, Error>,
  "queryKey" | "queryFn"
>;

type UpdateUserOptions = Omit<
  UseMutationOptions<UpdateUserOutput, Error, UpdateUserInput>,
  "mutationFn"
>;

type DeactivateUserOptions = Omit<
  UseMutationOptions<DeactivateUserOutput, Error, void>,
  "mutationFn"
>;

type ReactivateUserOptions = Omit<
  UseMutationOptions<ReactivateUserOutput, Error, void>,
  "mutationFn"
>;

type AddUserPhotoOptions = Omit<
  UseMutationOptions<AddUserPhotoOutput, Error, AddUserPhotoInput>,
  "mutationFn"
>;

type DeleteUserOptions = Omit<
  UseMutationOptions<DeleteUserOutput, Error, void>,
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

export const useUserProfile = (options?: GetProfileOptions) => {
  const { getProfileUseCase } = useUserUseCases();

  return useQuery({
    queryKey: ["user-profile"],
    queryFn: () => getProfileUseCase.execute(),
    ...options,
  });
};

export const useUserUpdate = (options?: UpdateUserOptions) => {
  const { updateUserUseCase } = useUserUseCases();

  return useMutation({
    mutationFn: (input) => updateUserUseCase.execute(input),
    ...options,
  });
};

export const useUserDeactivate = (options?: DeactivateUserOptions) => {
  const { deactivateUserUseCase } = useUserUseCases();

  return useMutation({
    mutationFn: () => deactivateUserUseCase.execute(),
    ...options,
  });
};

export const useUserReactivate = (options?: ReactivateUserOptions) => {
  const { reactivateUserUseCase } = useUserUseCases();

  return useMutation({
    mutationFn: () => reactivateUserUseCase.execute(),
    ...options,
  });
};

export const useUserAddPhoto = (options?: AddUserPhotoOptions) => {
  const { addUserPhotoUseCase } = useUserUseCases();

  return useMutation({
    mutationFn: (input) => addUserPhotoUseCase.execute(input),
    ...options,
  });
};

export const useUserDelete = (options?: DeleteUserOptions) => {
  const { deleteUserUseCase } = useUserUseCases();

  return useMutation({
    mutationFn: () => deleteUserUseCase.execute(),
    ...options,
  });
};
