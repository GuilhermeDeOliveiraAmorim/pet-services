import { useMemo } from "react";
import {
  useMutation,
  useQuery,
  type UseMutationOptions,
  type UseQueryOptions,
} from "@tanstack/react-query";

import type {
  AddProviderInput,
  AddProviderOutput,
  AddProviderPhotoOutput,
  DeleteProviderOutput,
  DeleteProviderPhotoOutput,
  GetProviderOutput,
  ListProvidersOutput,
  UpdateProviderInput,
  UpdateProviderOutput,
} from "@/application";
import { createProviderCases } from "@/application/factories/provider-usecase-factory";
import { createApiContext } from "@/infra";

const useProviderUseCases = () => {
  return useMemo(() => {
    const { providerGateway } = createApiContext();
    return createProviderCases(providerGateway);
  }, []);
};

type AddProviderOptions = Omit<
  UseMutationOptions<AddProviderOutput, Error, AddProviderInput>,
  "mutationFn"
>;

type GetProviderOptions = Omit<
  UseQueryOptions<GetProviderOutput, Error>,
  "queryKey" | "queryFn"
>;

type UpdateProviderOptions = Omit<
  UseMutationOptions<UpdateProviderOutput, Error, UpdateProviderInput>,
  "mutationFn"
>;

type DeleteProviderOptions = Omit<
  UseMutationOptions<DeleteProviderOutput, Error, string | number>,
  "mutationFn"
>;

type ListProvidersOptions = Omit<
  UseQueryOptions<ListProvidersOutput, Error>,
  "queryKey" | "queryFn"
>;

type AddProviderPhotoOptions = Omit<
  UseMutationOptions<
    AddProviderPhotoOutput,
    Error,
    { providerId: string | number; photo: File }
  >,
  "mutationFn"
>;

type DeleteProviderPhotoOptions = Omit<
  UseMutationOptions<
    DeleteProviderPhotoOutput,
    Error,
    { providerId: string | number; photoId: string | number }
  >,
  "mutationFn"
>;

export const useProviderAdd = (options?: AddProviderOptions) => {
  const { addProvider } = useProviderUseCases();

  return useMutation({
    mutationFn: (input) => addProvider.execute(input),
    ...options,
  });
};

export const useProviderGet = (
  providerId?: string | number,
  options?: GetProviderOptions,
) => {
  const { getProvider } = useProviderUseCases();

  return useQuery({
    queryKey: ["provider", providerId],
    queryFn: () => getProvider.execute(providerId!),
    enabled: Boolean(providerId),
    ...options,
  });
};

export const useProviderUpdate = (options?: UpdateProviderOptions) => {
  const { updateProvider } = useProviderUseCases();

  return useMutation({
    mutationFn: (input) => updateProvider.execute(input),
    ...options,
  });
};

export const useProviderDelete = (options?: DeleteProviderOptions) => {
  const { deleteProvider } = useProviderUseCases();

  return useMutation({
    mutationFn: (providerId) => deleteProvider.execute(providerId),
    ...options,
  });
};

export const useProviderList = (options?: ListProvidersOptions) => {
  const { listProviders } = useProviderUseCases();

  return useQuery({
    queryKey: ["providers"],
    queryFn: () => listProviders.execute(),
    ...options,
  });
};

export const useProviderAddPhoto = (options?: AddProviderPhotoOptions) => {
  const { addProviderPhoto } = useProviderUseCases();

  return useMutation({
    mutationFn: (input) => addProviderPhoto.execute(input),
    ...options,
  });
};

export const useProviderDeletePhoto = (
  options?: DeleteProviderPhotoOptions,
) => {
  const { deleteProviderPhoto } = useProviderUseCases();

  return useMutation({
    mutationFn: (input) =>
      deleteProviderPhoto.execute(input.providerId, input.photoId),
    ...options,
  });
};
