import type { Breed } from "@/domain";

export interface BreedApi {
  id?: string;
  name?: string;
  species_id?: string;
  speciesId?: string;
  active?: boolean;
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
  deactivated_at?: string;
  deactivatedAt?: string;
}

export const mapBreedFromApi = (breed?: BreedApi | null): Breed => {
  const raw = breed ?? {};

  return {
    id: String(raw.id ?? ""),
    name: raw.name ?? "",
    specieId: String(raw.speciesId ?? raw.species_id ?? ""),
    active: raw.active ?? true,
    createdAt: raw.createdAt ?? raw.created_at ?? new Date().toISOString(),
    updatedAt: raw.updatedAt ?? raw.updated_at,
    deactivatedAt: raw.deactivatedAt ?? raw.deactivated_at,
  };
};
