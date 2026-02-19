import { Category } from "@/domain/entities";

interface RawCategory {
  id?: string;
  name?: string;
  description?: string;
  active?: boolean;
  created_at?: string;
  updated_at?: string;
  deactivated_at?: string;
}

export function mapCategoryFromApi(raw: Record<string, unknown>): Category {
  const data = raw as RawCategory;

  return {
    id: data.id || "",
    name: data.name || "",
    description: data.description || "",
    active: data.active !== false,
    createdAt: data.created_at || new Date().toISOString(),
    updatedAt: data.updated_at || undefined,
    deactivatedAt: data.deactivated_at || undefined,
  };
}
