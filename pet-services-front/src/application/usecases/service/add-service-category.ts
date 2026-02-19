import { ServiceGateway } from "@/application/ports";

export interface AddServiceCategoryOutput {
  message?: string;
  detail?: string;
}

export class AddServiceCategoryUseCase {
  constructor(private readonly serviceGateway: ServiceGateway) {}

  execute(
    serviceId: string | number,
    categoryId: string | number,
  ): Promise<AddServiceCategoryOutput> {
    return this.serviceGateway.addServiceCategory(serviceId, categoryId);
  }
}
