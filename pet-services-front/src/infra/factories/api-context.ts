import { createApiClient } from "../http";
import type { AdoptionGateway } from "@/application/ports";
import {
  AuthGatewayAxios,
  CategoryGatewayAxios,
  PetGatewayAxios,
  ProviderGatewayAxios,
  ReferenceGatewayAxios,
  RequestGatewayAxios,
  ReviewGatewayAxios,
  ServiceGatewayAxios,
  BreedGatewayAxios,
  SpecieGatewayAxios,
  TagGatewayAxios,
  AdoptionGatewayAxios,
  UserGatewayAxios,
} from "../gateways";

type ApiContext = {
  http: ReturnType<typeof createApiClient>;
  authGateway: AuthGatewayAxios;
  userGateway: UserGatewayAxios;
  referenceGateway: ReferenceGatewayAxios;
  petGateway: PetGatewayAxios;
  specieGateway: SpecieGatewayAxios;
  breedGateway: BreedGatewayAxios;
  providerGateway: ProviderGatewayAxios;
  serviceGateway: ServiceGatewayAxios;
  requestGateway: RequestGatewayAxios;
  reviewGateway: ReviewGatewayAxios;
  categoryGateway: CategoryGatewayAxios;
  tagGateway: TagGatewayAxios;
  adoptionGateway: AdoptionGateway;
};

export const createApiContext = (baseURL?: string): ApiContext => {
  const http = createApiClient(baseURL);

  return {
    http,
    authGateway: new AuthGatewayAxios(http),
    userGateway: new UserGatewayAxios(http),
    referenceGateway: new ReferenceGatewayAxios(http),
    petGateway: new PetGatewayAxios(http),
    specieGateway: new SpecieGatewayAxios(http),
    breedGateway: new BreedGatewayAxios(http),
    providerGateway: new ProviderGatewayAxios(http),
    serviceGateway: new ServiceGatewayAxios(http),
    requestGateway: new RequestGatewayAxios(http),
    reviewGateway: new ReviewGatewayAxios(http),
    categoryGateway: new CategoryGatewayAxios(http),
    tagGateway: new TagGatewayAxios(http),
    adoptionGateway: new AdoptionGatewayAxios(http),
  };
};
