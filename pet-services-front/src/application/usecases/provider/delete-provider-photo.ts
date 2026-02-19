import type { ProviderGateway } from "@/application/ports/provider-gateway";

export interface DeleteProviderPhotoOutput {
  message?: string;
  detail?: string;
}

export class DeleteProviderPhotoUseCase {
  constructor(private readonly providerGateway: ProviderGateway) {}

  execute(
    providerId: string | number,
    photoId: string | number,
  ): Promise<DeleteProviderPhotoOutput> {
    return this.providerGateway.deleteProviderPhoto(providerId, photoId);
  }
}
