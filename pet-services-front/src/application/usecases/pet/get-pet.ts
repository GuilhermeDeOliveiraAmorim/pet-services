import type { Pet } from "@/domain";
import type { PetGateway } from "@/application/ports/pet-gateway";

export interface GetPetOutput {
  pet: Pet;
}

export class GetPetUseCase {
  constructor(private readonly petGateway: PetGateway) {}

  execute(petId: string | number): Promise<GetPetOutput> {
    return this.petGateway.getPet(petId);
  }
}
