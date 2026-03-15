import type { BreedGateway } from "@/application/ports/breed-gateway";
import type { Breed } from "@/domain";

export interface ListBreedsOutput {
  breeds: Breed[];
}

export class ListBreedsUseCase {
  constructor(private readonly breedGateway: BreedGateway) {}

  execute(speciesId: string): Promise<ListBreedsOutput> {
    return this.breedGateway.listBreedsBySpecies(speciesId);
  }
}
