import { ServiceGateway } from "@/application/ports";

export interface DeleteServiceTagOutput {
  message?: string;
  detail?: string;
}

export class DeleteServiceTagUseCase {
  constructor(private readonly serviceGateway: ServiceGateway) {}

  execute(
    serviceId: string | number,
    tagId: string | number,
  ): Promise<DeleteServiceTagOutput> {
    return this.serviceGateway.deleteServiceTag(serviceId, tagId);
  }
}
