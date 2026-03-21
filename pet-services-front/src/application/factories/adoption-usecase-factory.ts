import type { AdoptionGateway } from "../ports/adoption-gateway";
import {
  CreateAdoptionApplicationUseCase,
  GetPublicAdoptionListingUseCase,
  ListPublicAdoptionListingsUseCase,
} from "../usecases/adoption";

export const createAdoptionCases = (gateway: AdoptionGateway) => {
  return {
    listPublicAdoptionListings: new ListPublicAdoptionListingsUseCase(gateway),
    getPublicAdoptionListing: new GetPublicAdoptionListingUseCase(gateway),
    createAdoptionApplication: new CreateAdoptionApplicationUseCase(gateway),
  };
};
