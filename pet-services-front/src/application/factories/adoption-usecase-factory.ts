import type { AdoptionGateway } from "../ports/adoption-gateway";
import {
  GetPublicAdoptionListingUseCase,
  ListPublicAdoptionListingsUseCase,
} from "../usecases/adoption";

export const createAdoptionCases = (gateway: AdoptionGateway) => {
  return {
    listPublicAdoptionListings: new ListPublicAdoptionListingsUseCase(gateway),
    getPublicAdoptionListing: new GetPublicAdoptionListingUseCase(gateway),
  };
};
