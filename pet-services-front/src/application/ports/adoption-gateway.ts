import type {
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
}
