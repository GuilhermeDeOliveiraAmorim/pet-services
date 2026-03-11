import type { User } from "@/domain";

interface RawLogin {
  email?: string;
  password?: string;
}
interface RawPhone {
  countryCode?: string;
  areaCode?: string;
  number?: string;
  country_code?: string;
  area_code?: string;
}
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
interface RawUser {
  id?: string;
  active?: boolean;
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
  deactivated_at?: string;
  deactivatedAt?: string;
  name?: string;
  userType?: string;
  user_type?: string;
  login?: RawLogin;
  phone?: RawPhone;
  address?: RawAddress;
  emailVerified?: boolean;
  email_verified?: boolean;
  profile_complete?: boolean;
  photos?: unknown[];
  pets?: unknown[];
}

export const mapUserFromApi = (user?: RawUser | null): User => {
  const raw = user ?? {};

  return {
    id: raw.id ?? "",
    active: raw.active ?? false,
    createdAt: raw.created_at ?? raw.createdAt ?? "",
    updatedAt: raw.updated_at ?? raw.updatedAt ?? null,
    deactivatedAt: raw.deactivated_at ?? raw.deactivatedAt ?? null,
    name: raw.name ?? "",
    userType: (raw.userType ?? raw.user_type ?? "owner") as User["userType"],
    login: {
      email: raw.login?.email ?? "",
      password: raw.login?.password ?? "",
    },
    phone: {
      countryCode: raw.phone?.countryCode ?? raw.phone?.country_code ?? "",
      areaCode: raw.phone?.areaCode ?? raw.phone?.area_code ?? "",
      number: raw.phone?.number ?? "",
    },
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
    emailVerified: raw.emailVerified ?? raw.email_verified ?? false,
    profileComplete: raw.profile_complete ?? false,
    photos: Array.isArray(raw.photos) ? raw.photos : [],
    pets: Array.isArray(raw.pets) ? raw.pets : [],
  } as User;
};
