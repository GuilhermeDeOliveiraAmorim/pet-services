import { useMemo } from "react";
import {
  useMutation,
  useQuery,
  type UseMutationOptions,
  type UseQueryOptions,
} from "@tanstack/react-query";

import {
  type CreateAdoptionApplicationInput,
  type CreateAdoptionApplicationOutput,
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
    queryFn: () => getPublicAdoptionListing.execute(listingId!),
    enabled: Boolean(listingId),
    ...options,
  });
};

export const useAdoptionApplicationCreate = (
  options?: CreateAdoptionApplicationOptions,
) => {
  const { createAdoptionApplication } = useAdoptionUseCases();

  return useMutation({
    mutationFn: (input) => createAdoptionApplication.execute(input),
    ...options,
  });
};
