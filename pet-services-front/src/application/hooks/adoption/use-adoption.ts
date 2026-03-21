import { useMemo } from "react";
import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseMutationOptions,
  type UseQueryOptions,
} from "@tanstack/react-query";

import {
  type ListMyAdoptionApplicationsInput,
  type ListMyAdoptionApplicationsOutput,
  type CreateAdoptionApplicationInput,
  type CreateAdoptionApplicationOutput,
  type WithdrawAdoptionApplicationInput,
  type WithdrawAdoptionApplicationOutput,
  createAdoptionCases,
  type GetPublicAdoptionListingOutput,
  type ListPublicAdoptionListingsInput,
  type ListPublicAdoptionListingsOutput,
} from "@/application";
import { createApiContext } from "@/infra";

import { ADOPTION_KEYS } from "./adoption-query-keys";

const useAdoptionUseCases = () => {
  return useMemo(() => {
    const { adoptionGateway } = createApiContext();
    return createAdoptionCases(adoptionGateway);
  }, []);
};

type ListPublicAdoptionListingsOptions = Omit<
  UseQueryOptions<ListPublicAdoptionListingsOutput, Error>,
  "queryKey" | "queryFn"
>;

type GetPublicAdoptionListingOptions = Omit<
  UseQueryOptions<GetPublicAdoptionListingOutput, Error>,
  "queryKey" | "queryFn"
>;

type CreateAdoptionApplicationOptions = Omit<
  UseMutationOptions<
    CreateAdoptionApplicationOutput,
    Error,
    CreateAdoptionApplicationInput
  >,
  "mutationFn"
>;

type ListMyAdoptionApplicationsOptions = Omit<
  UseQueryOptions<ListMyAdoptionApplicationsOutput, Error>,
  "queryKey" | "queryFn"
>;

type WithdrawAdoptionApplicationOptions = Omit<
  UseMutationOptions<
    WithdrawAdoptionApplicationOutput,
    Error,
    WithdrawAdoptionApplicationInput
  >,
  "mutationFn"
>;

export const usePublicAdoptionListings = (
  input?: ListPublicAdoptionListingsInput,
  options?: ListPublicAdoptionListingsOptions,
) => {
  const { listPublicAdoptionListings } = useAdoptionUseCases();

  return useQuery({
    queryKey: ADOPTION_KEYS.list(input),
    queryFn: () => listPublicAdoptionListings.execute(input),
    ...options,
  });
};

export const usePublicAdoptionListing = (
  listingId?: string | number,
  options?: GetPublicAdoptionListingOptions,
) => {
  const { getPublicAdoptionListing } = useAdoptionUseCases();

  return useQuery({
    queryKey: ADOPTION_KEYS.detail(listingId ?? "unknown"),
    queryFn: () =>
      getPublicAdoptionListing.execute({
        listingId: listingId!,
      }),
    enabled: Boolean(listingId),
    ...options,
  });
};

export const useAdoptionApplicationCreate = (
  options?: CreateAdoptionApplicationOptions,
) => {
  const queryClient = useQueryClient();
  const { createAdoptionApplication } = useAdoptionUseCases();

  return useMutation({
    mutationFn: (input) => createAdoptionApplication.execute(input),
    onSuccess: async (...args) => {
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.myApplicationsLists(),
      });
      await options?.onSuccess?.(...args);
    },
    ...options,
  });
};

export const useMyAdoptionApplications = (
  input?: ListMyAdoptionApplicationsInput,
  options?: ListMyAdoptionApplicationsOptions,
) => {
  const { listMyAdoptionApplications } = useAdoptionUseCases();

  return useQuery({
    queryKey: ADOPTION_KEYS.myApplicationsList(input),
    queryFn: () => listMyAdoptionApplications.execute(input),
    ...options,
  });
};

export const useAdoptionApplicationWithdraw = (
  options?: WithdrawAdoptionApplicationOptions,
) => {
  const queryClient = useQueryClient();
  const { withdrawAdoptionApplication } = useAdoptionUseCases();

  return useMutation({
    mutationFn: (input) => withdrawAdoptionApplication.execute(input),
    onSuccess: async (...args) => {
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.myApplicationsLists(),
      });
      await options?.onSuccess?.(...args);
    },
    ...options,
  });
};
