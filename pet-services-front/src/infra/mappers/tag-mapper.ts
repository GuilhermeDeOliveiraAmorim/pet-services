import { Tag } from "@/domain/entities";

interface RawTag {
  id?: string;
  name?: string;
  active?: boolean;
  created_at?: string;
  updated_at?: string;
  deactivated_at?: string;
}

export function mapTagFromApi(raw: Record<string, unknown>): Tag {
  const data = raw as RawTag;

  return {
    id: data.id || "",
    name: data.name || "",
    active: data.active !== false,
    createdAt: data.created_at || new Date().toISOString(),
    updatedAt: data.updated_at || undefined,
    deactivatedAt: data.deactivated_at || undefined,
  };
}
