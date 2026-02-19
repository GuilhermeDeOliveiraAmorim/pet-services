import { Review } from "@/domain/entities";

interface RawReview {
  id?: string;
  user_id?: string;
  provider_id?: string;
  rating?: number;
  comment?: string;
  active?: boolean;
  created_at?: string;
  updated_at?: string;
  deactivated_at?: string;
}

export function mapReviewFromApi(raw: Record<string, unknown>): Review {
  const data = raw as RawReview;

  return {
    id: data.id || "",
    userId: data.user_id || "",
    providerId: data.provider_id || "",
    rating: data.rating || 0,
    comment: data.comment || "",
    active: data.active !== false,
    createdAt: data.created_at || new Date().toISOString(),
    updatedAt: data.updated_at || undefined,
    deactivatedAt: data.deactivated_at || undefined,
  };
}
