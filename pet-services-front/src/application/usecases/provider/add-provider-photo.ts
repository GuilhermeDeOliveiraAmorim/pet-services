import type { ProviderGateway } from "@/application/ports/provider-gateway";

export interface AddProviderPhotoInput {
  providerId: string | number;
  photo: File;
}

export interface AddProviderPhotoOutput {
  message?: string;
  detail?: string;
  photo?: {
    id: string;
    url: string;
  };
}

export class AddProviderPhotoUseCase {
  constructor(private readonly providerGateway: ProviderGateway) {}

  execute(input: AddProviderPhotoInput): Promise<AddProviderPhotoOutput> {
    return this.providerGateway.addProviderPhoto(input.providerId, input.photo);
  }
}
