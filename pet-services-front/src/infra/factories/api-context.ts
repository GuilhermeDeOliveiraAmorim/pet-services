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
  SpecieGatewayAxios,
  TagGatewayAxios,
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
    providerGateway: new ProviderGatewayAxios(http),
    serviceGateway: new ServiceGatewayAxios(http),
    requestGateway: new RequestGatewayAxios(http),
    reviewGateway: new ReviewGatewayAxios(http),
    categoryGateway: new CategoryGatewayAxios(http),
    tagGateway: new TagGatewayAxios(http),
  };
};
