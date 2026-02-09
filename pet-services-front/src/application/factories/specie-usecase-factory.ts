import type { SpecieGateway } from "../ports";
import { ListSpeciesUseCase } from "../usecases/specie";

export const createSpecieUseCases = (gateway: SpecieGateway) => {
  return {
    listSpeciesUseCase: new ListSpeciesUseCase(gateway),
  };
};
