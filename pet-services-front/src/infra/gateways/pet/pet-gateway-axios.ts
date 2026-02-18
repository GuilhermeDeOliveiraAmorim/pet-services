import type {
  AddPetInput,
  AddPetOutput,
  DeletePetOutput,
  DeletePetPhotoOutput,
  GetPetOutput,
  ListPetsOutput,
  PetGateway,
  UpdatePetInput,
  UpdatePetOutput,
} from "@/application";
import type { Pet } from "@/domain";
import type { AxiosInstance } from "axios";

type PetApi = {
  id?: string | number;
  name?: string;
  species?: {
    id?: string;
    name?: string;
    active?: boolean;
    created_at?: string;
    createdAt?: string;
  };
  species_id?: string;
  speciesId?: string;
  age?: number;
  weight?: number;
  notes?: string;
  user_id?: string;
  photos?: Array<{ id?: string | number; url?: string }>;
};

const mapPetApiToDomain = (pet: PetApi): Pet => ({
  id: String(pet.id ?? ""),
  userId: pet.user_id ?? "",
  name: pet.name ?? "",
  specie: {
    id: String(pet.species?.id ?? ""),
    name: pet.species?.name ?? "",
    active: pet.species?.active ?? true,
    createdAt:
      pet.species?.createdAt ??
      pet.species?.created_at ??
      new Date().toISOString(),
    updatedAt: undefined,
    deactivatedAt: undefined,
  },
  age: pet.age ?? 0,
  weight: pet.weight ?? 0,
  notes: pet.notes ?? "",
  photos: (pet.photos ?? []).map((p) => ({
    id: String(p.id ?? ""),
    url: p.url ?? "",
    active: true,
    createdAt: new Date().toISOString(),
    updatedAt: undefined,
    deactivatedAt: undefined,
  })),
  active: true,
  createdAt: new Date().toISOString(),
  updatedAt: undefined,
  deactivatedAt: undefined,
});

export class PetGatewayAxios implements PetGateway {
  constructor(private readonly http: AxiosInstance) {}

  setAccessToken(token?: string) {
    if (token) {
      this.http.defaults.headers.common.Authorization = `Bearer ${token}`;
      return;
    }

    delete this.http.defaults.headers.common.Authorization;
  }

  async addPet(input: AddPetInput): Promise<AddPetOutput> {
    const payload = {
      name: input.name,
      species_id: input.speciesId,
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
            speciesId:
              data.pet.species?.id ??
              data.pet.species_id ??
              data.pet.speciesId ??
              "",
            age: data.pet.age ?? 0,
            weight: data.pet.weight ?? 0,
            notes: data.pet.notes ?? "",
          }
        : undefined,
    };
  }

  async getPet(petId: string | number): Promise<GetPetOutput> {
    const { data } = await this.http.get<PetApi>(`/pets/${petId}`);
    return { pet: mapPetApiToDomain(data) };
  }

  async updatePet(input: UpdatePetInput): Promise<UpdatePetOutput> {
    const payload: Record<string, unknown> = {};
    if (input.name !== undefined) payload.name = input.name;
    if (input.age !== undefined) payload.age = input.age;
    if (input.weight !== undefined) payload.weight = input.weight;
    if (input.notes !== undefined) payload.notes = input.notes;

    const { data } = await this.http.put<{
      message?: string;
      detail?: string;
    }>(`/pets/${input.petId}`, payload);

    return {
      message: data.message,
      detail: data.detail,
    };
  }

  async deletePet(petId: string | number): Promise<DeletePetOutput> {
    const { data } = await this.http.delete<{
      message?: string;
      detail?: string;
    }>(`/pets/${petId}`);

    return {
      message: data.message,
      detail: data.detail,
    };
  }

  async listPets(): Promise<ListPetsOutput> {
    const { data } = await this.http.get<{ pets: PetApi[] }>("/pets/");
    return {
      pets: (data.pets ?? []).map(mapPetApiToDomain),
    };
  }

  async deletePetPhoto(
    petId: string | number,
    photoId: string | number,
  ): Promise<DeletePetPhotoOutput> {
    const { data } = await this.http.delete<{
      message?: string;
      detail?: string;
    }>(`/pets/${petId}/photos/${photoId}`);

    return {
      message: data.message,
      detail: data.detail,
    };
  }
}
