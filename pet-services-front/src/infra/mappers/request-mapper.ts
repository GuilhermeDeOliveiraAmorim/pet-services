import { Request } from "@/domain/entities";

interface RawRequest {
  id?: string;
  user_id?: string;
  provider_id?: string;
  service_id?: string;
  pet_id?: string;
  pet?: Record<string, unknown>;
  notes?: string;
  status?: string;
  reject_reason?: string;
  active?: boolean;
  created_at?: string;
  updated_at?: string;
  deactivated_at?: string;
}

export function mapRequestFromApi(raw: Record<string, unknown>): Request {
  const data = raw as RawRequest;

  return {
    id: data.id || "",
    userId: data.user_id || "",
    providerId: data.provider_id || "",
    serviceId: data.service_id || "",
    petId: data.pet_id || "",
    pet: data.pet as Request["pet"],
    notes: data.notes || "",
    status: (data.status as Request["status"]) || "pending",
    rejectReason: data.reject_reason,
    active: data.active !== false,
    createdAt: data.created_at || new Date().toISOString(),
    updatedAt: data.updated_at || undefined,
    deactivatedAt: data.deactivated_at || undefined,
  };
}
