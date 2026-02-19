import { Category } from "@/domain/entities";
import {
  ListCategoriesInput,
  ListCategoriesOutput,
} from "@/application/usecases/category";

export interface CategoryGateway {
  listCategories(input?: ListCategoriesInput): Promise<ListCategoriesOutput>;
}
