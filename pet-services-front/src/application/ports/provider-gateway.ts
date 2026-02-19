import type { AddProviderInput, AddProviderOutput } from "../usecases/provider/add-provider";

export interface ProviderGateway {
  addProvider(input: AddProviderInput): Promise<AddProviderOutput>;
  // ...outros métodos futuros
}
