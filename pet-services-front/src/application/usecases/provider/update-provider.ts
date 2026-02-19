import type { ProviderGateway } from "@/application/ports/provider-gateway";
import type { Provider } from "@/domain";

export interface UpdateProviderInput {
  providerId: string | number;
  businessName?: string;
  address?: {
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
  description?: string;
  priceRange?: string;
}

export interface UpdateProviderOutput {
  message?: string;
  detail?: string;
  provider?: Provider;
}

export class UpdateProviderUseCase {
  constructor(private readonly providerGateway: ProviderGateway) {}

  execute(input: UpdateProviderInput): Promise<UpdateProviderOutput> {
    return this.providerGateway.updateProvider(input);
  }
}
