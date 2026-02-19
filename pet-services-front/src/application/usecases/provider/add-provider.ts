import type { ProviderGateway } from "../../ports/provider-gateway";

export interface AddProviderInput {
  name: string;
  email: string;
  phone: {
    countryCode: string;
    areaCode: string;
    number: string;
  };
  address: {
    street: string;
    number: string;
    neighborhood: string;
    city: string;
    zipCode: string;
    state: string;
    country: string;
    complement?: string;
    location: {
      latitude: number;
      longitude: number;
    };
  };
}

export interface AddProviderOutput {
  id: string;
  message?: string;
  detail?: string;
}

export class AddProviderUseCase {
  constructor(private readonly providerGateway: ProviderGateway) {}

  execute(input: AddProviderInput): Promise<AddProviderOutput> {
    return this.providerGateway.addProvider(input);
  }
}
