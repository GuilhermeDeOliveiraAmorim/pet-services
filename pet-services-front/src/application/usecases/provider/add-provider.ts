import type { ProviderGateway } from "../../ports/provider-gateway";
import type { Provider } from "@/domain";

export interface AddProviderInput {
  businessName: string;
  description: string;
  priceRange: string;
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
  message?: string;
  detail?: string;
  provider?: Provider;
}

export class AddProviderUseCase {
  constructor(private readonly providerGateway: ProviderGateway) {}

  execute(input: AddProviderInput): Promise<AddProviderOutput> {
    return this.providerGateway.addProvider(input);
  }
}
