import type { City } from "@/domain";
import type { ReferenceGateway } from "@/application/ports";

export interface ListCitiesInput {
  stateId?: number;
}

export interface ListCitiesOutput {
  cities: City[];
}

export class ListCitiesUseCase {
  constructor(private readonly gateway: ReferenceGateway) {}

  execute(input: ListCitiesInput): Promise<ListCitiesOutput> {
    return this.gateway.listCities(input);
  }
}
