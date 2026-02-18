import type { AddPetInput, AddPetOutput } from "../usecases/pet/add_pet";
import type { GetPetOutput } from "../usecases/pet/get-pet";
import type { UpdatePetInput, UpdatePetOutput } from "../usecases/pet/update-pet";
import type { DeletePetOutput } from "../usecases/pet/delete-pet";
import type { ListPetsOutput } from "../usecases/pet/list-pets";
import type { DeletePetPhotoOutput } from "../usecases/pet/delete-pet-photo";

export interface PetGateway {
  addPet(input: AddPetInput): Promise<AddPetOutput>;
  getPet(petId: string | number): Promise<GetPetOutput>;
  updatePet(input: UpdatePetInput): Promise<UpdatePetOutput>;
  deletePet(petId: string | number): Promise<DeletePetOutput>;
  listPets(): Promise<ListPetsOutput>;
  deletePetPhoto(petId: string | number, photoId: string | number): Promise<DeletePetPhotoOutput>;
}
