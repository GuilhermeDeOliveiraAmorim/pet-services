import type { SpecieGateway } from "@/application/ports/specie-gateway";
import type { Specie } from "@/domain";

export interface ListSpeciesOutput {
  species: Specie[];
}

export class ListSpeciesUseCase {
  constructor(private readonly specieGateway: SpecieGateway) {}

  execute(): Promise<ListSpeciesOutput> {
    return this.specieGateway.listSpecies();
  }
}
