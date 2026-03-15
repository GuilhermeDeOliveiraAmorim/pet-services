import type { Service } from "@/domain";

interface RawCategory {
  id?: string;
  name?: string;
  description?: string;
}

interface RawTag {
  id?: string;
  name?: string;
}

interface RawService {
  id?: string;
  active?: boolean;
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
  deactivated_at?: string;
  deactivatedAt?: string;
  providerId?: string;
  provider_id?: string;
  averageRating?: number;
  average_rating?: number;
  reviewCount?: number;
  review_count?: number;
  distanceKm?: number;
  distance_km?: number;
  name?: string;
  description?: string;
  price?: number;
  priceMinimum?: number;
  price_minimum?: number;
  priceMaximum?: number;
  price_maximum?: number;
  duration?: number;
  photos?: unknown[];
  categories?: RawCategory[];
  tags?: RawTag[];
}

export const mapServiceFromApi = (service?: RawService | null): Service => {
  const raw = service ?? {};

  return {
    id: raw.id ?? "",
    active: raw.active ?? false,
    createdAt: raw.created_at ?? raw.createdAt ?? "",
    updatedAt: raw.updated_at ?? raw.updatedAt ?? null,
    deactivatedAt: raw.deactivated_at ?? raw.deactivatedAt ?? null,
    providerId: raw.providerId ?? raw.provider_id ?? "",
    averageRating: raw.averageRating ?? raw.average_rating ?? 0,
    reviewCount: raw.reviewCount ?? raw.review_count ?? 0,
    distanceKm: raw.distanceKm ?? raw.distance_km,
    name: raw.name ?? "",
    description: raw.description ?? "",
    price: raw.price ?? 0,
    priceMinimum: raw.priceMinimum ?? raw.price_minimum ?? 0,
    priceMaximum: raw.priceMaximum ?? raw.price_maximum ?? 0,
    duration: raw.duration ?? 0,
    photos: Array.isArray(raw.photos) ? raw.photos : [],
    categories: Array.isArray(raw.categories) ? raw.categories : [],
    tags: Array.isArray(raw.tags) ? raw.tags : [],
  } as Service;
};
