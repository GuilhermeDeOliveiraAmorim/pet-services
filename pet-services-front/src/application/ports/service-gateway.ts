import type {
  AddServiceInput,
  AddServiceOutput,
} from "../usecases/service/add-service";
import type { GetServiceOutput } from "../usecases/service/get-service";
import type {
  UpdateServiceInput,
  UpdateServiceOutput,
} from "../usecases/service/update-service";
import type { DeleteServiceOutput } from "../usecases/service/delete-service";
import type {
  ListServicesInput,
  ListServicesOutput,
} from "../usecases/service/list-services";
import type {
  SearchServicesInput,
  SearchServicesOutput,
} from "../usecases/service/search-services";
import type { AddServicePhotoOutput } from "../usecases/service/add-service-photo";
import type { DeleteServicePhotoOutput } from "../usecases/service/delete-service-photo";
import type { DeleteServiceCategoryOutput } from "../usecases/service/delete-service-category";
import type { AddServiceCategoryOutput } from "../usecases/service/add-service-category";
import type { AddServiceTagOutput } from "../usecases/service/add-service-tag";

export interface ServiceGateway {
  addService(input: AddServiceInput): Promise<AddServiceOutput>;
  getService(serviceId: string | number): Promise<GetServiceOutput>;
  updateService(input: UpdateServiceInput): Promise<UpdateServiceOutput>;
  deleteService(serviceId: string | number): Promise<DeleteServiceOutput>;
  listServices(input?: ListServicesInput): Promise<ListServicesOutput>;
  searchServices(input?: SearchServicesInput): Promise<SearchServicesOutput>;
  addServicePhoto(
    serviceId: string | number,
    photo: File,
  ): Promise<AddServicePhotoOutput>;
  deleteServicePhoto(
    serviceId: string | number,
    photoId: string | number,
  ): Promise<DeleteServicePhotoOutput>;
  addServiceCategory(
    serviceId: string | number,
    categoryId: string | number,
  ): Promise<AddServiceCategoryOutput>;
  deleteServiceCategory(
    serviceId: string | number,
    categoryId: string | number,
  ): Promise<DeleteServiceCategoryOutput>;
  addServiceTag(
    serviceId: string | number,
    tagId: string | number,
  ): Promise<AddServiceTagOutput>;
}
