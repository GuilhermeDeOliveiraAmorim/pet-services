import { CategoryGateway } from "@/application/ports/category-gateway";
import {
  ListCategoriesUseCase,
  ListCategoriesInput,
  ListCategoriesOutput,
} from "@/application/usecases/category";

class ListCategoriesUseCaseImpl implements ListCategoriesUseCase {
  constructor(private categoryGateway: CategoryGateway) {}

  async execute(input?: ListCategoriesInput): Promise<ListCategoriesOutput> {
    return this.categoryGateway.listCategories(input);
  }
}

export function createCategoryCases(categoryGateway: CategoryGateway) {
  return {
    listCategories: new ListCategoriesUseCaseImpl(categoryGateway),
  };
}
