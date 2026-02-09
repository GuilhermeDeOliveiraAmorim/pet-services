import { PetGateway } from "../ports/pet-gateway";
import { AddPetUseCase } from "../usecases/pet";

export const createPetCases = (gateway: PetGateway) => {
  return {
    addPet: new AddPetUseCase(gateway),
  };
};
