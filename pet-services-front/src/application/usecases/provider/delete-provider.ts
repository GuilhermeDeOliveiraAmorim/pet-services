import type { ProviderGateway } from "@/application/ports/provider-gateway";

export interface DeleteProviderOutput {
  message?: string;
  detail?: string;
}

export class DeleteProviderUseCase {
  constructor(private readonly providerGateway: ProviderGateway) {}

  execute(providerId: string | number): Promise<DeleteProviderOutput> {
    return this.providerGateway.deleteProvider(providerId);
  }
}
