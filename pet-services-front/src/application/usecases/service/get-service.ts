import { ServiceGateway } from "@/application/ports";
import type { Service } from "@/domain";

export interface GetServiceOutput {
  service: Service;
}

export class GetServiceUseCase {
  constructor(private readonly serviceGateway: ServiceGateway) {}

  execute(serviceId: string | number): Promise<GetServiceOutput> {
    return this.serviceGateway.getService(serviceId);
  }
}
