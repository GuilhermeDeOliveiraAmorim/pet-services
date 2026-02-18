import type { PetGateway } from "@/application/ports/pet-gateway";

export interface AddPetPhotoInput {
  petId: string | number;
  photo: File;
}

export interface AddPetPhotoOutput {
  message?: string;
  detail?: string;
  photo?: {
    id: string;
    url: string;
  };
}

export class AddPetPhotoUseCase {
  constructor(private readonly petGateway: PetGateway) {}

  execute(input: AddPetPhotoInput): Promise<AddPetPhotoOutput> {
    return this.petGateway.addPetPhoto(input.petId, input.photo);
  }
}
