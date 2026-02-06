import { useMemo } from "react";
import { useQuery, type UseQueryOptions } from "@tanstack/react-query";

import type {
  ListCitiesInput,
  ListCitiesOutput,
  ListCountriesOutput,
  ListStatesOutput,
} from "@/application/usecases/reference";
import { createReferenceUseCases } from "@/application/factories/reference-usecase-factory";
import { createApiContext } from "@/infra";

import { referenceQueryKeys } from "./reference-query-keys";

const useReferenceUseCases = () => {
  return useMemo(() => {
    const { referenceGateway } = createApiContext();
    return createReferenceUseCases(referenceGateway);
  }, []);
};

type CountriesOptions = Omit<
  UseQueryOptions<ListCountriesOutput, Error>,
  "queryKey" | "queryFn"
>;

type StatesOptions = Omit<
  UseQueryOptions<ListStatesOutput, Error>,
  "queryKey" | "queryFn"
>;

type CitiesOptions = Omit<
  UseQueryOptions<ListCitiesOutput, Error>,
  "queryKey" | "queryFn"
>;

export const useReferenceCountries = (options?: CountriesOptions) => {
  const { listCountriesUseCase } = useReferenceUseCases();

  return useQuery({
    queryKey: referenceQueryKeys.countries(),
    queryFn: () => listCountriesUseCase.execute(),
    ...options,
  });
};

export const useReferenceStates = (options?: StatesOptions) => {
  const { listStatesUseCase } = useReferenceUseCases();

  return useQuery({
    queryKey: referenceQueryKeys.states(),
    queryFn: () => listStatesUseCase.execute(),
    ...options,
  });
};

export const useReferenceCities = (
  input: ListCitiesInput = {},
  options?: CitiesOptions,
) => {
  const { listCitiesUseCase } = useReferenceUseCases();

  return useQuery({
    queryKey: referenceQueryKeys.cities(input.stateId),
    queryFn: () => listCitiesUseCase.execute(input),
    ...options,
  });
};
