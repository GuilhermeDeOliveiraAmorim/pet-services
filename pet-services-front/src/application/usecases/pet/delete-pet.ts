import type { PetGateway } from "@/application/ports/pet-gateway";

export interface DeletePetOutput {
  message?: string;
  detail?: string;
}

export class DeletePetUseCase {
  constructor(private readonly petGateway: PetGateway) {}

  execute(petId: string | number): Promise<DeletePetOutput> {
    return this.petGateway.deletePet(petId);
  }
}
