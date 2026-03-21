import type { AdoptionGateway } from "../ports/adoption-gateway";
import {
  CreateAdoptionApplicationUseCase,
  GetMyAdoptionGuardianProfileUseCase,
  GetPublicAdoptionListingUseCase,
  ListMyAdoptionApplicationsUseCase,
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
    withdrawAdoptionApplication: new WithdrawAdoptionApplicationUseCase(
      gateway,
    ),
  };
};
