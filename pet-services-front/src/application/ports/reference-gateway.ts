import type {
  ListCitiesInput,
  ListCitiesOutput,
  ListCountriesOutput,
  ListStatesOutput,
} from "../usecases/reference";

export interface ReferenceGateway {
  listCountries(): Promise<ListCountriesOutput>;
  listStates(): Promise<ListStatesOutput>;
  listCities(input: ListCitiesInput): Promise<ListCitiesOutput>;
}
