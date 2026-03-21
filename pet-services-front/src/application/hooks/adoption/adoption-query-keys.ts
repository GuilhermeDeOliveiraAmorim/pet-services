import type { ListPublicAdoptionListingsInput } from "@/application/usecases/adoption";
import type {
  ListAdoptionApplicationsByListingInput,
  ListMyAdoptionListingsInput,
} from "@/application/usecases/adoption";
import type { ListMyAdoptionApplicationsInput } from "@/application/usecases/adoption";

export const ADOPTION_KEYS = {
  all: ["adoption"] as const,
  lists: () => [...ADOPTION_KEYS.all, "list"] as const,
  list: (filters?: ListPublicAdoptionListingsInput) =>
    [...ADOPTION_KEYS.lists(), filters] as const,
  details: () => [...ADOPTION_KEYS.all, "detail"] as const,
  detail: (listingId: string | number) =>
    [...ADOPTION_KEYS.details(), listingId] as const,
  myApplicationsLists: () => [...ADOPTION_KEYS.all, "my-applications"] as const,
  myApplicationsList: (filters?: ListMyAdoptionApplicationsInput) =>
    [...ADOPTION_KEYS.myApplicationsLists(), filters] as const,
  myListingsLists: () => [...ADOPTION_KEYS.all, "my-listings"] as const,
  myListingsList: (filters?: ListMyAdoptionListingsInput) =>
    [...ADOPTION_KEYS.myListingsLists(), filters] as const,
  listingApplicationsLists: () =>
    [...ADOPTION_KEYS.all, "listing-applications"] as const,
  listingApplicationsList: (input: ListAdoptionApplicationsByListingInput) =>
    [...ADOPTION_KEYS.listingApplicationsLists(), input] as const,
  myGuardianProfile: () =>
    [...ADOPTION_KEYS.all, "my-guardian-profile"] as const,
} as const;
