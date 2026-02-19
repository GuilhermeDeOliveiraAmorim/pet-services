import { ProviderGateway } from "../ports/provider-gateway";
import {
  AddProviderUseCase,
  AddProviderPhotoUseCase,
  DeleteProviderUseCase,
  DeleteProviderPhotoUseCase,
  GetProviderUseCase,
  ListProvidersUseCase,
  UpdateProviderUseCase,
} from "../usecases/provider";

export const createProviderCases = (gateway: ProviderGateway) => {
  return {
    addProvider: new AddProviderUseCase(gateway),
    getProvider: new GetProviderUseCase(gateway),
    updateProvider: new UpdateProviderUseCase(gateway),
    deleteProvider: new DeleteProviderUseCase(gateway),
    listProviders: new ListProvidersUseCase(gateway),
    addProviderPhoto: new AddProviderPhotoUseCase(gateway),
    deleteProviderPhoto: new DeleteProviderPhotoUseCase(gateway),
  };
};
