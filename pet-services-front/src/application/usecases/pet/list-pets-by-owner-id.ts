import type { Pet } from "@/domain";
import type { PetGateway } from "@/application/ports/pet-gateway";

export interface ListPetsByOwnerIdOutput {
  pets: Pet[];
}

export class ListPetsByOwnerIdUseCase {
  constructor(private readonly petGateway: PetGateway) {}

  execute(ownerId: string): Promise<ListPetsByOwnerIdOutput> {
    return this.petGateway.listPetsByOwnerId(ownerId);
  }
}
