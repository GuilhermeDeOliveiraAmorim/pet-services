import { ServiceGateway } from "@/application/ports";

export interface DeleteServiceOutput {
  message?: string;
  detail?: string;
}

export class DeleteServiceUseCase {
  constructor(private readonly serviceGateway: ServiceGateway) {}

  execute(serviceId: string | number): Promise<DeleteServiceOutput> {
    return this.serviceGateway.deleteService(serviceId);
  }
}
