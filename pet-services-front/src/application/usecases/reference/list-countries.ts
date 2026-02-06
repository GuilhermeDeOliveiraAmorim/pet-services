import type { Country } from "@/domain";
import type { ReferenceGateway } from "@/application/ports";

export interface ListCountriesOutput {
  countries: Country[];
}

export class ListCountriesUseCase {
  constructor(private readonly gateway: ReferenceGateway) {}

  execute(): Promise<ListCountriesOutput> {
    return this.gateway.listCountries();
  }
}
