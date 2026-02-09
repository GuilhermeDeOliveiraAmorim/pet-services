import type { Species } from "@/domain";

export interface SpecieApi {
  id?: string;
  active?: boolean;
  created_at?: string;
  createdAt?: string;
  updated_at?: string | null;
  updatedAt?: string | null;
  deactivated_at?: string | null;
  deactivatedAt?: string | null;
  name?: string;
}

export const mapSpecieFromApi = (specie?: SpecieApi | null): Species => {
  const raw = specie ?? {};

  return {
    id: raw.id ?? "",
    active: raw.active ?? false,
    createdAt: raw.created_at ?? raw.createdAt ?? "",
    updatedAt: raw.updated_at ?? raw.updatedAt ?? null,
    deactivatedAt: raw.deactivated_at ?? raw.deactivatedAt ?? null,
    name: raw.name ?? "",
  };
};
