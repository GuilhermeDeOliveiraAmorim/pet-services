import { Category } from "@/domain/entities";

export interface ListCategoriesInput {
  name?: string;
  page?: number;
  pageSize?: number;
}

export interface ListCategoriesOutput {
  categories: Category[];
  total: number;
  page: number;
  pageSize: number;
}

export interface ListCategoriesUseCase {
  execute(input?: ListCategoriesInput): Promise<ListCategoriesOutput>;
}
