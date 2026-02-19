import type {
  AddProviderInput,
  AddProviderOutput,
} from "../usecases/provider/add-provider";
import type { GetProviderOutput } from "../usecases/provider/get-provider";
import type {
  UpdateProviderInput,
  UpdateProviderOutput,
} from "../usecases/provider/update-provider";
import type { DeleteProviderOutput } from "../usecases/provider/delete-provider";
import type { ListProvidersOutput } from "../usecases/provider/list-providers";
import type { AddProviderPhotoOutput } from "../usecases/provider/add-provider-photo";
import type { DeleteProviderPhotoOutput } from "../usecases/provider/delete-provider-photo";

export interface ProviderGateway {
  addProvider(input: AddProviderInput): Promise<AddProviderOutput>;
  getProvider(providerId: string | number): Promise<GetProviderOutput>;
  updateProvider(input: UpdateProviderInput): Promise<UpdateProviderOutput>;
  deleteProvider(providerId: string | number): Promise<DeleteProviderOutput>;
  listProviders(): Promise<ListProvidersOutput>;
  addProviderPhoto(
    providerId: string | number,
    photo: File,
  ): Promise<AddProviderPhotoOutput>;
  deleteProviderPhoto(
    providerId: string | number,
    photoId: string | number,
  ): Promise<DeleteProviderPhotoOutput>;
}
