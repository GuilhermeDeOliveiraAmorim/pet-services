import { useMemo } from "react";
import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseMutationOptions,
  type UseQueryOptions,
} from "@tanstack/react-query";

import {
  type ChangeAdoptionListingStatusInput,
  type ChangeAdoptionListingStatusOutput,
  type CreateAdoptionGuardianProfileInput,
  type CreateAdoptionGuardianProfileOutput,
  type CreateAdoptionListingInput,
  type CreateAdoptionListingOutput,
  type ListAdoptionApplicationsByListingInput,
  type ListAdoptionApplicationsByListingOutput,
  type ListMyAdoptionApplicationsInput,
  type ListMyAdoptionApplicationsOutput,
  type ListMyAdoptionListingsInput,
  type ListMyAdoptionListingsOutput,
  type MarkAdoptionListingAsAdoptedInput,
  type MarkAdoptionListingAsAdoptedOutput,
  type CreateAdoptionApplicationInput,
  type CreateAdoptionApplicationOutput,
  type GetMyAdoptionGuardianProfileOutput,
  type ReviewAdoptionApplicationInput,
  type ReviewAdoptionApplicationOutput,
  type UpdateAdoptionGuardianProfileInput,
  type UpdateAdoptionGuardianProfileOutput,
  type UpdateAdoptionListingInput,
  type UpdateAdoptionListingOutput,
  type WithdrawAdoptionApplicationInput,
  type WithdrawAdoptionApplicationOutput,
  createAdoptionCases,
  type GetPublicAdoptionListingOutput,
  type ListPublicAdoptionListingsInput,
  type ListPublicAdoptionListingsOutput,
  AdoptionGateway,
} from "@/application";
import { createApiContext } from "@/infra";

import { ADOPTION_KEYS } from "./adoption-query-keys";

const useAdoptionUseCases = () => {
  return useMemo(() => {
    const { adoptionGateway } = createApiContext();
    const gateway = adoptionGateway as unknown as AdoptionGateway;
    return createAdoptionCases(gateway);
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

type CreateAdoptionGuardianProfileOptions = Omit<
  UseMutationOptions<
    CreateAdoptionGuardianProfileOutput,
    Error,
    CreateAdoptionGuardianProfileInput
  >,
  "mutationFn"
>;

type UpdateAdoptionGuardianProfileOptions = Omit<
  UseMutationOptions<
    UpdateAdoptionGuardianProfileOutput,
    Error,
    UpdateAdoptionGuardianProfileInput
  >,
  "mutationFn"
>;

type CreateAdoptionListingOptions = Omit<
  UseMutationOptions<
    CreateAdoptionListingOutput,
    Error,
    CreateAdoptionListingInput
  >,
  "mutationFn"
>;

type UpdateAdoptionListingOptions = Omit<
  UseMutationOptions<
    UpdateAdoptionListingOutput,
    Error,
    UpdateAdoptionListingInput
  >,
  "mutationFn"
>;

type ChangeAdoptionListingStatusOptions = Omit<
  UseMutationOptions<
    ChangeAdoptionListingStatusOutput,
    Error,
    ChangeAdoptionListingStatusInput
  >,
  "mutationFn"
>;

type ReviewAdoptionApplicationOptions = Omit<
  UseMutationOptions<
    ReviewAdoptionApplicationOutput,
    Error,
    ReviewAdoptionApplicationInput
  >,
  "mutationFn"
>;

type MarkAdoptionListingAsAdoptedOptions = Omit<
  UseMutationOptions<
    MarkAdoptionListingAsAdoptedOutput,
    Error,
    MarkAdoptionListingAsAdoptedInput
  >,
  "mutationFn"
>;

type ListMyAdoptionApplicationsOptions = Omit<
  UseQueryOptions<ListMyAdoptionApplicationsOutput, Error>,
  "queryKey" | "queryFn"
>;

type ListMyAdoptionListingsOptions = Omit<
  UseQueryOptions<ListMyAdoptionListingsOutput, Error>,
  "queryKey" | "queryFn"
>;

type ListAdoptionApplicationsByListingOptions = Omit<
  UseQueryOptions<ListAdoptionApplicationsByListingOutput, Error>,
  "queryKey" | "queryFn"
>;

type GetMyAdoptionGuardianProfileOptions = Omit<
  UseQueryOptions<GetMyAdoptionGuardianProfileOutput | null, Error>,
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

export const useAdoptionGuardianProfileCreate = (
  options?: CreateAdoptionGuardianProfileOptions,
) => {
  const queryClient = useQueryClient();
  const { createAdoptionGuardianProfile } = useAdoptionUseCases();

  return useMutation({
    mutationFn: (input) => createAdoptionGuardianProfile.execute(input),
    onSuccess: async (...args) => {
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.myGuardianProfile(),
      });
      await options?.onSuccess?.(...args);
    },
    ...options,
  });
};

export const useAdoptionGuardianProfileUpdate = (
  options?: UpdateAdoptionGuardianProfileOptions,
) => {
  const queryClient = useQueryClient();
  const { updateAdoptionGuardianProfile } = useAdoptionUseCases();

  return useMutation({
    mutationFn: (input) => updateAdoptionGuardianProfile.execute(input),
    onSuccess: async (...args) => {
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.myGuardianProfile(),
      });
      await options?.onSuccess?.(...args);
    },
    ...options,
  });
};

export const useAdoptionListingCreate = (
  options?: CreateAdoptionListingOptions,
) => {
  const queryClient = useQueryClient();
  const { createAdoptionListing } = useAdoptionUseCases();

  return useMutation({
    mutationFn: (input) => createAdoptionListing.execute(input),
    onSuccess: async (...args) => {
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.myListingsLists(),
      });
      await options?.onSuccess?.(...args);
    },
    ...options,
  });
};

export const useAdoptionListingUpdate = (
  options?: UpdateAdoptionListingOptions,
) => {
  const queryClient = useQueryClient();
  const { updateAdoptionListing } = useAdoptionUseCases();

  return useMutation({
    mutationFn: (input) => updateAdoptionListing.execute(input),
    onSuccess: async (...args) => {
      const [, variables] = args;
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.myListingsLists(),
      });
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.detail(variables.listingId),
      });
      await options?.onSuccess?.(...args);
    },
    ...options,
  });
};

export const useAdoptionListingStatusChange = (
  options?: ChangeAdoptionListingStatusOptions,
) => {
  const queryClient = useQueryClient();
  const { changeAdoptionListingStatus } = useAdoptionUseCases();

  return useMutation({
    mutationFn: (input) => changeAdoptionListingStatus.execute(input),
    onSuccess: async (...args) => {
      const [, variables] = args;
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.myListingsLists(),
      });
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.detail(variables.listingId),
      });
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.lists(),
      });
      await options?.onSuccess?.(...args);
    },
    ...options,
  });
};

export const useAdoptionApplicationReview = (
  options?: ReviewAdoptionApplicationOptions,
) => {
  const queryClient = useQueryClient();
  const { reviewAdoptionApplication } = useAdoptionUseCases();

  return useMutation({
    mutationFn: (input) => reviewAdoptionApplication.execute(input),
    onSuccess: async (...args) => {
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.listingApplicationsLists(),
      });
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.myListingsLists(),
      });
      await options?.onSuccess?.(...args);
    },
    ...options,
  });
};

export const useAdoptionListingMarkAsAdopted = (
  options?: MarkAdoptionListingAsAdoptedOptions,
) => {
  const queryClient = useQueryClient();
  const { markAdoptionListingAsAdopted } = useAdoptionUseCases();

  return useMutation({
    mutationFn: (input) => markAdoptionListingAsAdopted.execute(input),
    onSuccess: async (...args) => {
      const [, variables] = args;
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.myListingsLists(),
      });
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.detail(variables.listingId),
      });
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.listingApplicationsLists(),
      });
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.lists(),
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

export const useMyAdoptionListings = (
  input?: ListMyAdoptionListingsInput,
  options?: ListMyAdoptionListingsOptions,
) => {
  const { listMyAdoptionListings } = useAdoptionUseCases();

  return useQuery({
    queryKey: ADOPTION_KEYS.myListingsList(input),
    queryFn: () => listMyAdoptionListings.execute(input),
    ...options,
  });
};

export const useAdoptionApplicationsByListing = (
  input?: ListAdoptionApplicationsByListingInput,
  options?: ListAdoptionApplicationsByListingOptions,
) => {
  const { listAdoptionApplicationsByListing } = useAdoptionUseCases();

  return useQuery({
    queryKey: ADOPTION_KEYS.listingApplicationsList(
      input ?? { listingId: "unknown" },
    ),
    queryFn: () => listAdoptionApplicationsByListing.execute(input!),
    enabled: Boolean(input?.listingId),
    ...options,
  });
};

export const useMyAdoptionGuardianProfile = (
  options?: GetMyAdoptionGuardianProfileOptions,
) => {
  const { getMyAdoptionGuardianProfile } = useAdoptionUseCases();

  return useQuery({
    queryKey: ADOPTION_KEYS.myGuardianProfile(),
    queryFn: () => getMyAdoptionGuardianProfile.execute(),
    retry: false,
    ...options,
  });
};

export const useGuardianStatus = (
  options?: GetMyAdoptionGuardianProfileOptions,
) => {
  const query = useMyAdoptionGuardianProfile(options);

  return {
    ...query,
    isApprovedGuardian: query.data?.profile?.approvalStatus === "approved",
  };
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
      await queryClient.invalidateQueries({
        queryKey: ADOPTION_KEYS.listingApplicationsLists(),
      });
      await options?.onSuccess?.(...args);
    },
    ...options,
  });
};
