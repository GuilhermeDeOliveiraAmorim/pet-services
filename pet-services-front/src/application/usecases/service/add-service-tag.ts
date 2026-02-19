import { ServiceGateway } from "@/application/ports";

export interface AddServiceTagOutput {
  message?: string;
  detail?: string;
}

export class AddServiceTagUseCase {
  constructor(private readonly serviceGateway: ServiceGateway) {}

  execute(
    serviceId: string | number,
    tagId: string | number,
  ): Promise<AddServiceTagOutput> {
    return this.serviceGateway.addServiceTag(serviceId, tagId);
  }
}
