import type { ProviderGateway } from "@/application/ports/provider-gateway";
import type { Provider } from "@/domain";

export interface GetProviderOutput {
  provider: Provider;
}

export class GetProviderUseCase {
  constructor(private readonly providerGateway: ProviderGateway) {}

  execute(providerId: string | number): Promise<GetProviderOutput> {
    return this.providerGateway.getProvider(providerId);
  }
}
