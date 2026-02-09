import type { User } from "@/domain";

export const mapUserFromApi = (user?: Record<string, any> | null): User => {
  const raw = user ?? {};

  return {
    id: raw.id ?? "",
    active: raw.active ?? false,
    createdAt: raw.created_at ?? raw.createdAt ?? "",
    updatedAt: raw.updated_at ?? raw.updatedAt ?? null,
    deactivatedAt: raw.deactivated_at ?? raw.deactivatedAt ?? null,
    name: raw.name ?? "",
    userType: raw.userType ?? raw.user_type ?? "owner",
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
    photos: raw.photos ?? [],
    pets: raw.pets ?? [],
  } as User;
};
