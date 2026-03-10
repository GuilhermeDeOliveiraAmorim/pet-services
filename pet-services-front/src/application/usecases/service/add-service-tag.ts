import { ServiceGateway } from "@/application/ports";
import type { Tag } from "@/domain";

export interface AddServiceTagOutput {
  message?: string;
  detail?: string;
  tag?: Tag;
}

export interface AddServiceTagPayload {
  tagId?: string | number;
  tagName?: string;
}

export class AddServiceTagUseCase {
  constructor(private readonly serviceGateway: ServiceGateway) {}

  execute(
    serviceId: string | number,
    payload: AddServiceTagPayload,
  ): Promise<AddServiceTagOutput> {
    return this.serviceGateway.addServiceTag(serviceId, payload);
  }
}
