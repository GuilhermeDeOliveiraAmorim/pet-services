import { useMemo } from "react";
import { useMutation, type UseMutationOptions } from "@tanstack/react-query";

import type { AddPetInput, AddPetOutput } from "@/application";
import { createPetCases } from "@/application/factories/pet-usecase-factory";
import { createApiContext } from "@/infra";

const usePetUseCases = () => {
  return useMemo(() => {
    const { petGateway } = createApiContext();
    return createPetCases(petGateway);
  }, []);
};

type AddPetOptions = Omit<
  UseMutationOptions<AddPetOutput, Error, AddPetInput>,
  "mutationFn"
>;

export const usePetAdd = (options?: AddPetOptions) => {
  const { addPet } = usePetUseCases();

  return useMutation({
    mutationFn: (input) => addPet.execute(input),
    ...options,
  });
};
