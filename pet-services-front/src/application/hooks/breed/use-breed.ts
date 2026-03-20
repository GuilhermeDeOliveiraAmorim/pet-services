import { useMemo } from "react";
import { useQuery, type UseQueryOptions } from "@tanstack/react-query";

import type { ListBreedsOutput } from "@/application/usecases/breed";
import { createBreedUseCases } from "@/application/factories/breed-usecase-factory";
import { BREED_KEYS } from "@/application/hooks/breed/breed-query-keys";
import { createApiContext } from "@/infra";

const useBreedUseCases = () => {
  return useMemo(() => {
    const { breedGateway } = createApiContext();
    return createBreedUseCases(breedGateway);
  }, []);
};

type ListBreedsOptions = Omit<
  UseQueryOptions<ListBreedsOutput, Error>,
  "queryKey" | "queryFn"
>;

export const useBreedsList = (
  speciesId?: string,
  options?: ListBreedsOptions,
) => {
  const { listBreedsUseCase } = useBreedUseCases();

  return useQuery({
    queryKey: BREED_KEYS.list(speciesId),
    queryFn: () => listBreedsUseCase.execute(speciesId ?? ""),
    enabled: Boolean(speciesId),
    ...options,
  });
};
