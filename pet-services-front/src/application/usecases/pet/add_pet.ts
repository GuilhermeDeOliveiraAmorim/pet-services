import { PetGateway } from "@/application/ports/pet-gateway";

export interface AddPetInput {
  name: string;
  speciesId: string;
  breed?: string;
  age: number;
  weight: number;
  notes: string;
}

export interface AddPetOutput {
  message?: string;
  detail?: string;
  pet?: {
    id: number;
    name: string;
    speciesId: string;
    breed?: string;
    age: number;
    weight: number;
    notes: string;
  };
}

export class AddPetUseCase {
  constructor(private readonly petGateway: PetGateway) {}

  execute(input: AddPetInput): Promise<AddPetOutput> {
    return this.petGateway.addPet(input);
  }
}
