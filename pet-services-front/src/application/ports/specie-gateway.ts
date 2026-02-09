import type { ListSpeciesOutput } from "@/application/usecases/specie/list-species";

export interface SpecieGateway {
  listSpecies(): Promise<ListSpeciesOutput>;
}
