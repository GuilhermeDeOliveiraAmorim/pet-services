import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { ListCategoriesInput } from "@/application/usecases/category";
import { createCategoryCases } from "@/application/factories/category-usecase-factory";
import { createApiContext } from "@/infra";
import { CATEGORY_KEYS } from "./category-query-keys";

const useCategoryUseCases = () => {
  return useMemo(() => {
    const { categoryGateway } = createApiContext();
    return createCategoryCases(categoryGateway);
  }, []);
};

export function useCategoryList(input?: ListCategoriesInput) {
  const { listCategories } = useCategoryUseCases();

  return useQuery({
    queryKey: CATEGORY_KEYS.list(input),
    queryFn: () => listCategories.execute(input),
  });
}
