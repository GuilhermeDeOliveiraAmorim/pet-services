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
  GetProviderOutput,
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
