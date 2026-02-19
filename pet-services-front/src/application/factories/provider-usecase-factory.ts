import { ProviderGateway } from "../ports/provider-gateway";
import {
  AddProviderUseCase,
  GetProviderUseCase,
  UpdateProviderUseCase,
} from "../usecases/provider";

export const createProviderCases = (gateway: ProviderGateway) => {
  return {
    addProvider: new AddProviderUseCase(gateway),
    getProvider: new GetProviderUseCase(gateway),
    updateProvider: new UpdateProviderUseCase(gateway),
  };
};
