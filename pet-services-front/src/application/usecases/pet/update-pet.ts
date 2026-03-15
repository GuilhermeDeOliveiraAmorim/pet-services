import type { PetGateway } from "@/application/ports/pet-gateway";

export interface UpdatePetInput {
  petId: string | number;
  name?: string;
  speciesId?: string;
  breed?: string;
  age?: number;
  weight?: number;
  notes?: string;
}

export interface UpdatePetOutput {
  message?: string;
  detail?: string;
}

export class UpdatePetUseCase {
  constructor(private readonly petGateway: PetGateway) {}

  execute(input: UpdatePetInput): Promise<UpdatePetOutput> {
    return this.petGateway.updatePet(input);
  }
}
