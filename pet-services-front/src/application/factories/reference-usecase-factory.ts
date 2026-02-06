import type { ReferenceGateway } from "../ports";
import {
  ListCitiesUseCase,
  ListCountriesUseCase,
  ListStatesUseCase,
} from "../usecases/reference";

export const createReferenceUseCases = (gateway: ReferenceGateway) => {
  return {
    listCountriesUseCase: new ListCountriesUseCase(gateway),
    listStatesUseCase: new ListStatesUseCase(gateway),
    listCitiesUseCase: new ListCitiesUseCase(gateway),
  };
};
