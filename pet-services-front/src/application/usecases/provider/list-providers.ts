import type { ProviderGateway } from "@/application/ports/provider-gateway";
import type { Provider } from "@/domain";

export interface ListProvidersOutput {
  providers: Provider[];
}

export class ListProvidersUseCase {
  constructor(private readonly providerGateway: ProviderGateway) {}

  execute(): Promise<ListProvidersOutput> {
    return this.providerGateway.listProviders();
  }
}
