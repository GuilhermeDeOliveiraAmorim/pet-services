import { useMemo } from "react";
import { useQuery, type UseQueryOptions } from "@tanstack/react-query";

import type { ListSpeciesOutput } from "@/application/usecases/specie";
import { createSpecieUseCases } from "@/application/factories/specie-usecase-factory";
import { createApiContext } from "@/infra";

const useSpecieUseCases = () => {
  return useMemo(() => {
    const { specieGateway } = createApiContext();
    return createSpecieUseCases(specieGateway);
  }, []);
};

type ListSpeciesOptions = Omit<
  UseQueryOptions<ListSpeciesOutput, Error>,
  "queryKey" | "queryFn"
>;

export const useSpeciesList = (options?: ListSpeciesOptions) => {
  const { listSpeciesUseCase } = useSpecieUseCases();

  return useQuery({
    queryKey: ["species"],
    queryFn: () => listSpeciesUseCase.execute(),
    ...options,
  });
};
