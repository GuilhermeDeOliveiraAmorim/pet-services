import type { BreedGateway } from "../ports";
import { ListBreedsUseCase } from "../usecases/breed";

export const createBreedUseCases = (gateway: BreedGateway) => {
  return {
    listBreedsUseCase: new ListBreedsUseCase(gateway),
  };
};
