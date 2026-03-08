import { ServiceGateway } from "@/application/ports";

export interface AddServiceInput {
  providerId: string;
  name: string;
  description: string;
  price?: number;
  priceMinimum?: number;
  priceMaximum?: number;
  duration: number;
}

export interface AddServiceOutput {
  id: string;
  message?: string;
  detail?: string;
}

export class AddServiceUseCase {
  constructor(private readonly serviceGateway: ServiceGateway) {}

  execute(input: AddServiceInput): Promise<AddServiceOutput> {
    return this.serviceGateway.addService(input);
  }
}
