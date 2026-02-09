import type { AddPetInput, AddPetOutput, PetGateway } from "@/application";
import type { AxiosInstance } from "axios";

type PetApi = {
  id?: string | number;
  name?: string;
  specie?: { id?: string };
  specie_id?: string;
  specieId?: string;
  age?: number;
  weight?: number;
  notes?: string;
};

export class PetGatewayAxios implements PetGateway {
  constructor(private readonly http: AxiosInstance) {}

  async addPet(input: AddPetInput): Promise<AddPetOutput> {
    const payload = {
      name: input.name,
      specie_id: input.specieId,
      age: input.age,
      weight: input.weight,
      notes: input.notes,
    };

    const { data } = await this.http.post<{
      message?: string;
      detail?: string;
      pet?: PetApi;
    }>("/pets/", payload);

    return {
      message: data.message,
      detail: data.detail,
      pet: data.pet
        ? {
            id: Number(data.pet.id ?? 0),
            name: data.pet.name ?? "",
            specieId:
              data.pet.specie?.id ??
              data.pet.specie_id ??
              data.pet.specieId ??
              "",
            age: data.pet.age ?? 0,
            weight: data.pet.weight ?? 0,
            notes: data.pet.notes ?? "",
          }
        : undefined,
    };
  }
}
