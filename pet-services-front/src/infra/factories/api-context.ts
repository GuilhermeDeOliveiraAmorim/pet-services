import { createApiClient } from "../http";
import {
  AuthGatewayAxios,
  ReferenceGatewayAxios,
  UserGatewayAxios,
} from "../gateways";

export const createApiContext = (baseURL?: string) => {
  const http = createApiClient(baseURL);

  return {
    http,
    authGateway: new AuthGatewayAxios(http),
    userGateway: new UserGatewayAxios(http),
    referenceGateway: new ReferenceGatewayAxios(http),
  };
};
