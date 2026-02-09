import type { AddPetInput, AddPetOutput, PetGateway } from "@/application";
import type { AxiosInstance } from "axios";

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
      pet?: Record<string, any>;
    }>("/pets", payload);

    return {
      message: data.message,
      detail: data.detail,
      pet: data.pet
        ? {
            id: data.pet.id,
            name: data.pet.name,
            specieId:
              data.pet.specie?.id ??
              data.pet.specie_id ??
              data.pet.specieId ??
              "",
            age: data.pet.age,
            weight: data.pet.weight,
            notes: data.pet.notes,
          }
        : undefined,
    };
  }
}
