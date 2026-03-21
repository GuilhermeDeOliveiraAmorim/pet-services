import type {
  CreateAdoptionApplicationInput,
  CreateAdoptionApplicationOutput,
  GetMyAdoptionGuardianProfileInput,
  GetMyAdoptionGuardianProfileOutput,
  ListAdoptionApplicationsByListingInput,
  ListAdoptionApplicationsByListingOutput,
  ListMyAdoptionListingsInput,
  ListMyAdoptionListingsOutput,
  GetPublicAdoptionListingInput,
  GetPublicAdoptionListingOutput,
  ListMyAdoptionApplicationsInput,
  ListMyAdoptionApplicationsOutput,
  ListPublicAdoptionListingsInput,
  ListPublicAdoptionListingsOutput,
  WithdrawAdoptionApplicationOutput,
} from "../usecases/adoption";

export interface AdoptionGateway {
  listPublicListings(
    input?: ListPublicAdoptionListingsInput,
  ): Promise<ListPublicAdoptionListingsOutput>;
  getPublicListing(
    input: GetPublicAdoptionListingInput,
  ): Promise<GetPublicAdoptionListingOutput>;
  createApplication(
    input: CreateAdoptionApplicationInput,
  ): Promise<CreateAdoptionApplicationOutput>;
  listMyApplications(
    input?: ListMyAdoptionApplicationsInput,
  ): Promise<ListMyAdoptionApplicationsOutput>;
  listMyListings(
    input?: ListMyAdoptionListingsInput,
  ): Promise<ListMyAdoptionListingsOutput>;
  listApplicationsByListing(
    input: ListAdoptionApplicationsByListingInput,
  ): Promise<ListAdoptionApplicationsByListingOutput>;
  getMyGuardianProfile(
    input?: GetMyAdoptionGuardianProfileInput,
  ): Promise<GetMyAdoptionGuardianProfileOutput | null>;
  withdrawApplication(
    applicationId: string,
  ): Promise<WithdrawAdoptionApplicationOutput>;
}
