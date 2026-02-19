import { PetGateway } from "../ports/pet-gateway";
import {
  AddPetUseCase,
  AddPetPhotoUseCase,
  GetPetUseCase,
  UpdatePetUseCase,
  DeletePetUseCase,
  ListPetsUseCase,
  DeletePetPhotoUseCase,
} from "../usecases/pet";

export const createPetCases = (gateway: PetGateway) => {
  return {
    addPet: new AddPetUseCase(gateway),
    addPetPhoto: new AddPetPhotoUseCase(gateway),
    getPet: new GetPetUseCase(gateway),
    updatePet: new UpdatePetUseCase(gateway),
    deletePet: new DeletePetUseCase(gateway),
    listPets: new ListPetsUseCase(gateway),
    deletePetPhoto: new DeletePetPhotoUseCase(gateway),
  };
};
