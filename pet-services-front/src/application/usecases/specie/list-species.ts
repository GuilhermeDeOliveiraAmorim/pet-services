import type { SpecieGateway } from "@/application/ports/specie-gateway";
import type { Species } from "@/domain";

export interface ListSpeciesOutput {
  species: Species[];
}

export class ListSpeciesUseCase {
  constructor(private readonly specieGateway: SpecieGateway) {}

  execute(): Promise<ListSpeciesOutput> {
    return this.specieGateway.listSpecies();
  }
}
