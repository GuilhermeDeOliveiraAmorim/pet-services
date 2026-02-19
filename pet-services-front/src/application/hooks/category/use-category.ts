import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { ListCategoriesInput } from "@/application/usecases/category";
import { createCategoryCases } from "@/application/factories/category-usecase-factory";
import { createApiContext } from "@/infra";

const useCategoryUseCases = () => {
  return useMemo(() => {
    const { categoryGateway } = createApiContext();
    return createCategoryCases(categoryGateway);
  }, []);
};

const CATEGORY_KEYS = {
  all: ["categories"] as const,
  lists: () => [...CATEGORY_KEYS.all, "list"] as const,
  list: (filters?: ListCategoriesInput) =>
    [...CATEGORY_KEYS.lists(), filters] as const,
};

export function useCategoryList(input?: ListCategoriesInput) {
  const { listCategories } = useCategoryUseCases();

  return useQuery({
    queryKey: CATEGORY_KEYS.list(input),
    queryFn: () => listCategories.execute(input),
  });
}
