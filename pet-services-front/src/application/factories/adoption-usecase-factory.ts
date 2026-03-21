import type { AdoptionGateway } from "../ports/adoption-gateway";
import {
  ChangeAdoptionListingStatusUseCase,
  CreateAdoptionApplicationUseCase,
  CreateAdoptionGuardianProfileUseCase,
  CreateAdoptionListingUseCase,
  GetMyAdoptionGuardianProfileUseCase,
  GetPublicAdoptionListingUseCase,
  ListAdoptionApplicationsByListingUseCase,
  ListMyAdoptionApplicationsUseCase,
  ListMyAdoptionListingsUseCase,
  ListPublicAdoptionListingsUseCase,
  MarkAdoptionListingAsAdoptedUseCase,
  ReviewAdoptionApplicationUseCase,
  UpdateAdoptionGuardianProfileUseCase,
  UpdateAdoptionListingUseCase,
  WithdrawAdoptionApplicationUseCase,
} from "../usecases/adoption";

export const createAdoptionCases = (gateway: AdoptionGateway) => {
  return {
    createAdoptionGuardianProfile: new CreateAdoptionGuardianProfileUseCase(
      gateway,
    ),
    updateAdoptionGuardianProfile: new UpdateAdoptionGuardianProfileUseCase(
      gateway,
    ),
    listPublicAdoptionListings: new ListPublicAdoptionListingsUseCase(gateway),
    getPublicAdoptionListing: new GetPublicAdoptionListingUseCase(gateway),
    getMyAdoptionGuardianProfile: new GetMyAdoptionGuardianProfileUseCase(
      gateway,
    ),
    createAdoptionApplication: new CreateAdoptionApplicationUseCase(gateway),
    createAdoptionListing: new CreateAdoptionListingUseCase(gateway),
    updateAdoptionListing: new UpdateAdoptionListingUseCase(gateway),
    changeAdoptionListingStatus: new ChangeAdoptionListingStatusUseCase(
      gateway,
    ),
    reviewAdoptionApplication: new ReviewAdoptionApplicationUseCase(gateway),
    markAdoptionListingAsAdopted: new MarkAdoptionListingAsAdoptedUseCase(
      gateway,
    ),
    listMyAdoptionApplications: new ListMyAdoptionApplicationsUseCase(gateway),
    listMyAdoptionListings: new ListMyAdoptionListingsUseCase(gateway),
    listAdoptionApplicationsByListing:
      new ListAdoptionApplicationsByListingUseCase(gateway),
    withdrawAdoptionApplication: new WithdrawAdoptionApplicationUseCase(
      gateway,
    ),
  };
};
