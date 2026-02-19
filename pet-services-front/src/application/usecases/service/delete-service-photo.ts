import { ServiceGateway } from "@/application/ports";

export interface DeleteServicePhotoOutput {
  message?: string;
  detail?: string;
}

export class DeleteServicePhotoUseCase {
  constructor(private readonly serviceGateway: ServiceGateway) {}

  execute(
    serviceId: string | number,
    photoId: string | number,
  ): Promise<DeleteServicePhotoOutput> {
    return this.serviceGateway.deleteServicePhoto(serviceId, photoId);
  }
}
