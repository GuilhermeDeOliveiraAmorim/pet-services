import type {
  AddProviderInput,
  AddProviderOutput,
} from "../usecases/provider/add-provider";
import type { GetProviderOutput } from "../usecases/provider/get-provider";
import type {
  UpdateProviderInput,
  UpdateProviderOutput,
} from "../usecases/provider/update-provider";

export interface ProviderGateway {
  addProvider(input: AddProviderInput): Promise<AddProviderOutput>;
  getProvider(providerId: string | number): Promise<GetProviderOutput>;
  updateProvider(input: UpdateProviderInput): Promise<UpdateProviderOutput>;
  // ...outros métodos futuros
}
