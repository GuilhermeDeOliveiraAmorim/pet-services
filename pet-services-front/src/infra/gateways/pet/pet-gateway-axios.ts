import type {
  AddPetInput,
  AddPetOutput,
  AddPetPhotoOutput,
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
  active?: boolean;
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
  deactivated_at?: string | null;
  deactivatedAt?: string | null;
  specie?: {
    id?: string;
    name?: string;
    active?: boolean;
    created_at?: string;
    createdAt?: string;
  };
  species?: {
    id?: string;
    name?: string;
    active?: boolean;
    created_at?: string;
    createdAt?: string;
  };
  species_id?: string;
  speciesId?: string;
  breed?: string;
  age?: number;
  weight?: number;
  notes?: string;
  user_id?: string;
  photos?: Array<{ id?: string | number; url?: string }> | null;
};

const mapPetApiToDomain = (pet: PetApi): Pet => {
  const rawSpecie = pet.specie ?? pet.species;
  const specieName =
    typeof rawSpecie?.name === "string" ? rawSpecie.name.trim() : "";
  const breedName = typeof pet.breed === "string" ? pet.breed.trim() : "";
  const notes = typeof pet.notes === "string" ? pet.notes.trim() : "";

  return {
    id: String(pet.id ?? ""),
    userId: pet.user_id ?? "",
    name: pet.name ?? "",
    specie: {
      id: String(rawSpecie?.id ?? ""),
      name: specieName,
      active: rawSpecie?.active ?? true,
      createdAt:
        rawSpecie?.createdAt ??
        rawSpecie?.created_at ??
        new Date().toISOString(),
      updatedAt: undefined,
      deactivatedAt: undefined,
    },
    breed: breedName,
    age: pet.age ?? 0,
    weight: pet.weight ?? 0,
    notes,
    photos: (pet.photos ?? []).map((p) => ({
      id: String(p.id ?? ""),
      url: p.url ?? "",
      active: true,
      createdAt: new Date().toISOString(),
      updatedAt: undefined,
      deactivatedAt: undefined,
    })),
    active: pet.active ?? true,
    createdAt: pet.createdAt ?? pet.created_at ?? new Date().toISOString(),
    updatedAt: pet.updatedAt ?? pet.updated_at,
    deactivatedAt: pet.deactivatedAt ?? pet.deactivated_at ?? undefined,
  };
};

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
      breed: input.breed,
      age: input.age,
      weight: input.weight,
      notes: input.notes,
    };

    const { data } = await this.http.post<{
      message?: string;
      detail?: string;
      pet?: PetApi;
    }>("/pets", payload);

    return {
      message: data.message,
      detail: data.detail,
      pet: data.pet
        ? {
            id: Number(data.pet.id ?? 0),
            name: data.pet.name ?? "",
            speciesId:
              data.pet.specie?.id ??
              data.pet.species?.id ??
              data.pet.species_id ??
              data.pet.speciesId ??
              "",
            breed: data.pet.breed ?? "",
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
    if (input.speciesId !== undefined) payload.species_id = input.speciesId;
    if (input.breed !== undefined) payload.breed = input.breed;
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
    const { data } = await this.http.get<{ pets: PetApi[] }>("/pets");
    return {
      pets: (data.pets ?? []).map(mapPetApiToDomain),
    };
  }

  async listPetsByOwnerId(ownerId: string): Promise<ListPetsOutput> {
    const { data } = await this.http.get<{ pets: PetApi[] }>(
      `/users/${ownerId}/pets`,
    );

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

  async addPetPhoto(
    petId: string | number,
    photo: File,
  ): Promise<AddPetPhotoOutput> {
    const formData = new FormData();
    formData.append("file", photo);

    const { data } = await this.http.post<{
      message?: string;
      detail?: string;
      photo?: { id: string; url: string };
    }>(`/pets/${petId}/photos`, formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });

    return {
      message: data.message,
      detail: data.detail,
      photo: data.photo,
    };
  }
}
