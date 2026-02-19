import type { Provider } from "@/domain";

interface RawLocation {
  latitude?: number;
  longitude?: number;
  lat?: number;
  lng?: number;
}

interface RawAddress {
  street?: string;
  number?: string;
  neighborhood?: string;
  city?: string;
  zipCode?: string;
  zip_code?: string;
  state?: string;
  country?: string;
  complement?: string;
  location?: RawLocation;
}

interface RawProvider {
  id?: string;
  active?: boolean;
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
  deactivated_at?: string;
  deactivatedAt?: string;
  userId?: string;
  user_id?: string;
  businessName?: string;
  business_name?: string;
  address?: RawAddress;
  description?: string;
  priceRange?: string;
  price_range?: string;
  averageRating?: number;
  average_rating?: number;
  photos?: unknown[];
}

export const mapProviderFromApi = (
  provider?: RawProvider | null,
): Provider => {
  const raw = provider ?? {};

  return {
    id: raw.id ?? "",
    active: raw.active ?? false,
    createdAt: raw.created_at ?? raw.createdAt ?? "",
    updatedAt: raw.updated_at ?? raw.updatedAt ?? null,
    deactivatedAt: raw.deactivated_at ?? raw.deactivatedAt ?? null,
    userId: raw.userId ?? raw.user_id ?? "",
    businessName: raw.businessName ?? raw.business_name ?? "",
    address: {
      street: raw.address?.street ?? "",
      number: raw.address?.number ?? "",
      neighborhood: raw.address?.neighborhood ?? "",
      city: raw.address?.city ?? "",
      zipCode: raw.address?.zipCode ?? raw.address?.zip_code ?? "",
      state: raw.address?.state ?? "",
      country: raw.address?.country ?? "",
      complement: raw.address?.complement ?? "",
      location: {
        latitude:
          raw.address?.location?.latitude ?? raw.address?.location?.lat ?? 0,
        longitude:
          raw.address?.location?.longitude ?? raw.address?.location?.lng ?? 0,
      },
    },
    description: raw.description ?? "",
    priceRange: raw.priceRange ?? raw.price_range ?? "",
    averageRating: raw.averageRating ?? raw.average_rating ?? 0,
    photos: Array.isArray(raw.photos) ? raw.photos : [],
  } as Provider;
};
