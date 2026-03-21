import type { AdoptionGateway } from "../ports/adoption-gateway";
import {
  CreateAdoptionApplicationUseCase,
  GetMyAdoptionGuardianProfileUseCase,
  GetPublicAdoptionListingUseCase,
  ListAdoptionApplicationsByListingUseCase,
  ListMyAdoptionApplicationsUseCase,
  ListMyAdoptionListingsUseCase,
  ListPublicAdoptionListingsUseCase,
  WithdrawAdoptionApplicationUseCase,
} from "../usecases/adoption";

export const createAdoptionCases = (gateway: AdoptionGateway) => {
  return {
    listPublicAdoptionListings: new ListPublicAdoptionListingsUseCase(gateway),
    getPublicAdoptionListing: new GetPublicAdoptionListingUseCase(gateway),
    getMyAdoptionGuardianProfile: new GetMyAdoptionGuardianProfileUseCase(
      gateway,
    ),
    createAdoptionApplication: new CreateAdoptionApplicationUseCase(gateway),
    listMyAdoptionApplications: new ListMyAdoptionApplicationsUseCase(gateway),
    listMyAdoptionListings: new ListMyAdoptionListingsUseCase(gateway),
    listAdoptionApplicationsByListing:
      new ListAdoptionApplicationsByListingUseCase(gateway),
    withdrawAdoptionApplication: new WithdrawAdoptionApplicationUseCase(
      gateway,
    ),
  };
};
