import type { AddPetInput, AddPetOutput } from "../usecases/pet/add_pet";

export interface PetGateway {
  addPet(input: AddPetInput): Promise<AddPetOutput>;
}
