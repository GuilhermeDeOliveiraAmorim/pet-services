import type {
  AddProviderInput,
  AddProviderOutput,
} from "../usecases/provider/add-provider";
import type { GetProviderOutput } from "../usecases/provider/get-provider";

export interface ProviderGateway {
  addProvider(input: AddProviderInput): Promise<AddProviderOutput>;
  getProvider(providerId: string | number): Promise<GetProviderOutput>;
  // ...outros métodos futuros
}
