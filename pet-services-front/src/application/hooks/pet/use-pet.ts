import { useMemo } from "react";
import {
  useMutation,
  useQuery,
  type UseMutationOptions,
  type UseQueryOptions,
} from "@tanstack/react-query";

import type {
  AddPetInput,
  AddPetOutput,
  DeletePetOutput,
  DeletePetPhotoOutput,
  GetPetOutput,
  ListPetsOutput,
  UpdatePetInput,
  UpdatePetOutput,
} from "@/application";
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

type GetPetOptions = Omit<
  UseQueryOptions<GetPetOutput, Error>,
  "queryKey" | "queryFn"
>;

type UpdatePetOptions = Omit<
  UseMutationOptions<UpdatePetOutput, Error, UpdatePetInput>,
  "mutationFn"
>;

type DeletePetOptions = Omit<
  UseMutationOptions<DeletePetOutput, Error, string | number>,
  "mutationFn"
>;

type ListPetsOptions = Omit<
  UseQueryOptions<ListPetsOutput, Error>,
  "queryKey" | "queryFn"
>;

type DeletePetPhotoOptions = Omit<
  UseMutationOptions<
    DeletePetPhotoOutput,
    Error,
    { petId: string | number; photoId: string | number }
  >,
  "mutationFn"
>;

export const usePetAdd = (options?: AddPetOptions) => {
  const { addPet } = usePetUseCases();

  return useMutation({
    mutationFn: (input) => addPet.execute(input),
    ...options,
  });
};

export const usePetGet = (
  petId?: string | number,
  options?: GetPetOptions,
) => {
  const { getPet } = usePetUseCases();

  return useQuery({
    queryKey: ["pet", petId],
    queryFn: () => getPet.execute(petId!),
    enabled: Boolean(petId),
    ...options,
  });
};

export const usePetUpdate = (options?: UpdatePetOptions) => {
  const { updatePet } = usePetUseCases();

  return useMutation({
    mutationFn: (input) => updatePet.execute(input),
    ...options,
  });
};

export const usePetDelete = (options?: DeletePetOptions) => {
  const { deletePet } = usePetUseCases();

  return useMutation({
    mutationFn: (petId) => deletePet.execute(petId),
    ...options,
  });
};

export const usePetList = (options?: ListPetsOptions) => {
  const { listPets } = usePetUseCases();

  return useQuery({
    queryKey: ["pets"],
    queryFn: () => listPets.execute(),
    ...options,
  });
};

export const usePetDeletePhoto = (options?: DeletePetPhotoOptions) => {
  const { deletePetPhoto } = usePetUseCases();

  return useMutation({
    mutationFn: (input) =>
      deletePetPhoto.execute(input.petId, input.photoId),
    ...options,
  });
};
