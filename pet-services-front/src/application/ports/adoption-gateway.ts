import type {
  CreateAdoptionApplicationInput,
  CreateAdoptionApplicationOutput,
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
  withdrawApplication(
    applicationId: string,
  ): Promise<WithdrawAdoptionApplicationOutput>;
}
