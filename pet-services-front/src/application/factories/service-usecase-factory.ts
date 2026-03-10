import { ServiceGateway } from "../ports/service-gateway";
import {
  AddServiceUseCase,
  AddServiceCategoryUseCase,
  AddServicePhotoUseCase,
  AddServiceTagUseCase,
  DeleteServiceUseCase,
  DeleteServiceCategoryUseCase,
  DeleteServicePhotoUseCase,
  GetServiceUseCase,
  ListServicesUseCase,
  SearchServicesUseCase,
  UpdateServiceUseCase,
} from "../usecases/service";

export const createServiceCases = (gateway: ServiceGateway) => {
  return {
    addService: new AddServiceUseCase(gateway),
    getService: new GetServiceUseCase(gateway),
    updateService: new UpdateServiceUseCase(gateway),
    deleteService: new DeleteServiceUseCase(gateway),
    listServices: new ListServicesUseCase(gateway),
    searchServices: new SearchServicesUseCase(gateway),
    addServicePhoto: new AddServicePhotoUseCase(gateway),
    deleteServicePhoto: new DeleteServicePhotoUseCase(gateway),
    addServiceCategory: new AddServiceCategoryUseCase(gateway),
    deleteServiceCategory: new DeleteServiceCategoryUseCase(gateway),
    addServiceTag: new AddServiceTagUseCase(gateway),
  };
};
