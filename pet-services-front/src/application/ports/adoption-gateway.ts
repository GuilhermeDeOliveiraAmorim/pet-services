import type {
  CreateAdoptionApplicationInput,
  CreateAdoptionApplicationOutput,
  CreateAdoptionGuardianProfileInput,
  CreateAdoptionGuardianProfileOutput,
  CreateAdoptionListingInput,
  CreateAdoptionListingOutput,
  ChangeAdoptionListingStatusInput,
  ChangeAdoptionListingStatusOutput,
  GetMyAdoptionGuardianProfileInput,
  GetMyAdoptionGuardianProfileOutput,
  ListAdoptionApplicationsByListingInput,
  ListAdoptionApplicationsByListingOutput,
  MarkAdoptionListingAsAdoptedInput,
  MarkAdoptionListingAsAdoptedOutput,
  ReviewAdoptionApplicationInput,
  ReviewAdoptionApplicationOutput,
  ListMyAdoptionListingsInput,
  ListMyAdoptionListingsOutput,
  GetPublicAdoptionListingInput,
  GetPublicAdoptionListingOutput,
  ListMyAdoptionApplicationsInput,
  ListMyAdoptionApplicationsOutput,
  ListPublicAdoptionListingsInput,
  ListPublicAdoptionListingsOutput,
  UpdateAdoptionGuardianProfileInput,
  UpdateAdoptionGuardianProfileOutput,
  UpdateAdoptionListingInput,
  UpdateAdoptionListingOutput,
  WithdrawAdoptionApplicationOutput,
} from "../usecases/adoption";

export interface AdoptionGateway {
  createGuardianProfile(
    input: CreateAdoptionGuardianProfileInput,
  ): Promise<CreateAdoptionGuardianProfileOutput>;
  updateGuardianProfile(
    input: UpdateAdoptionGuardianProfileInput,
  ): Promise<UpdateAdoptionGuardianProfileOutput>;
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
  createListing(
    input: CreateAdoptionListingInput,
  ): Promise<CreateAdoptionListingOutput>;
  updateListing(
    input: UpdateAdoptionListingInput,
  ): Promise<UpdateAdoptionListingOutput>;
  changeListingStatus(
    input: ChangeAdoptionListingStatusInput,
  ): Promise<ChangeAdoptionListingStatusOutput>;
  reviewApplication(
    input: ReviewAdoptionApplicationInput,
  ): Promise<ReviewAdoptionApplicationOutput>;
  markListingAsAdopted(
    input: MarkAdoptionListingAsAdoptedInput,
  ): Promise<MarkAdoptionListingAsAdoptedOutput>;
  withdrawApplication(
    applicationId: string,
  ): Promise<WithdrawAdoptionApplicationOutput>;
}
