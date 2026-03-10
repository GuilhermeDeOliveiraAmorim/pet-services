import { ServiceGateway } from "@/application/ports";

export interface DeleteServiceCategoryOutput {
  message?: string;
  detail?: string;
}

export class DeleteServiceCategoryUseCase {
  constructor(private readonly serviceGateway: ServiceGateway) {}

  execute(
    serviceId: string | number,
    categoryId: string | number,
  ): Promise<DeleteServiceCategoryOutput> {
    return this.serviceGateway.deleteServiceCategory(serviceId, categoryId);
  }
}
