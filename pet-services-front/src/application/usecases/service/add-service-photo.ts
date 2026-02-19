import { ServiceGateway } from "@/application/ports";

export interface AddServicePhotoInput {
  serviceId: string | number;
  photo: File;
}

export interface AddServicePhotoOutput {
  message?: string;
  detail?: string;
  photo?: {
    id: string;
    url: string;
  };
}

export class AddServicePhotoUseCase {
  constructor(private readonly serviceGateway: ServiceGateway) {}

  execute(input: AddServicePhotoInput): Promise<AddServicePhotoOutput> {
    return this.serviceGateway.addServicePhoto(input.serviceId, input.photo);
  }
}
