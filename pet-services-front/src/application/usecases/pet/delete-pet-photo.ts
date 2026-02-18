import type { PetGateway } from "@/application/ports/pet-gateway";

export interface DeletePetPhotoOutput {
  message?: string;
  detail?: string;
}

export class DeletePetPhotoUseCase {
  constructor(private readonly petGateway: PetGateway) {}

  execute(petId: string | number, photoId: string | number): Promise<DeletePetPhotoOutput> {
    return this.petGateway.deletePetPhoto(petId, photoId);
  }
}
