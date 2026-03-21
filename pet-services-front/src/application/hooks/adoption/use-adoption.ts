import { useMemo } from "react";
import { useQuery, type UseQueryOptions } from "@tanstack/react-query";

import {
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
