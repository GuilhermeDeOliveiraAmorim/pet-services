import type {
  CreateAdoptionApplicationInput,
  CreateAdoptionApplicationOutput,
  GetPublicAdoptionListingOutput,
  ListPublicAdoptionListingsInput,
  ListPublicAdoptionListingsOutput,
} from "../usecases/adoption";

export interface AdoptionGateway {
  listPublicListings(
    input?: ListPublicAdoptionListingsInput,
  ): Promise<ListPublicAdoptionListingsOutput>;
  getPublicListing(
    listingId: string | number,
  ): Promise<GetPublicAdoptionListingOutput>;
  createApplication(
    input: CreateAdoptionApplicationInput,
  ): Promise<CreateAdoptionApplicationOutput>;
}
