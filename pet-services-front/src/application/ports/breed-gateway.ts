import type { ListBreedsOutput } from "@/application/usecases/breed/list-breeds";

export interface BreedGateway {
  listBreedsBySpecies(speciesId: string): Promise<ListBreedsOutput>;
}
