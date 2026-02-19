import { ServiceGateway } from "@/application/ports";
import type { Service } from "@/domain";

export interface UpdateServiceInput {
  serviceId: string | number;
  name?: string;
  description?: string;
  price?: number;
  priceMinimum?: number;
  priceMaximum?: number;
  duration?: number;
}

export interface UpdateServiceOutput {
  message?: string;
  detail?: string;
  service?: Service;
}

export class UpdateServiceUseCase {
  constructor(private readonly serviceGateway: ServiceGateway) {}

  execute(input: UpdateServiceInput): Promise<UpdateServiceOutput> {
    return this.serviceGateway.updateService(input);
  }
}
