import type { Pet } from "@/domain";
import type { PetGateway } from "@/application/ports/pet-gateway";

export interface ListPetsOutput {
  pets: Pet[];
}

export class ListPetsUseCase {
  constructor(private readonly petGateway: PetGateway) {}

  execute(): Promise<ListPetsOutput> {
    return this.petGateway.listPets();
  }
}
