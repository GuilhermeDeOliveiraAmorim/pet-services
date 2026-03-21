import { createApiClient } from "../http";
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

export const createApiContext = (baseURL?: string) => {
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
