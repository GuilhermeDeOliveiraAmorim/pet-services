import type {
  AdoptionGateway,
  GetPublicAdoptionListingOutput,
  ListPublicAdoptionListingsInput,
  ListPublicAdoptionListingsOutput,
} from "@/application";
import type { AdoptionListing, Pet } from "@/domain";
import type { AxiosInstance } from "axios";

type AdoptionGuardianProfileApi = {
  id?: string;
  display_name?: string;
  displayName?: string;
  about?: string;
  whatsapp?: string;
  city_id?: string;
  cityId?: string;
  state_id?: string;
  stateId?: string;
};

type PetApi = {
  id?: string | number;
  user_id?: string;
  userId?: string;
  name?: string;
  breed?: string;
  age?: number;
  weight?: number;
  notes?: string;
  active?: boolean;
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
  deactivated_at?: string | null;
  deactivatedAt?: string | null;
  specie?: {
    id?: string;
    name?: string;
    active?: boolean;
    created_at?: string;
    createdAt?: string;
  };
  species?: {
    id?: string;
    name?: string;
    active?: boolean;
    created_at?: string;
    createdAt?: string;
  };
  photos?: Array<{ id?: string | number; url?: string }> | null;
};

type AdoptionListingApi = {
  id?: string;
  active?: boolean;
  created_at?: string;
  createdAt?: string;
  updated_at?: string | null;
  updatedAt?: string | null;
  deactivated_at?: string | null;
  deactivatedAt?: string | null;
  title?: string;
  description?: string;
  adoption_reason?: string;
  adoptionReason?: string;
  status?: string;
  sex?: string;
  size?: string;
  age_group?: string;
  ageGroup?: string;
  state_id?: string;
  stateId?: string;
  city_id?: string;
  cityId?: string;
  latitude?: number;
  longitude?: number;
  good_with_children?: boolean;
  goodWithChildren?: boolean;
  good_with_dogs?: boolean;
  goodWithDogs?: boolean;
  good_with_cats?: boolean;
  goodWithCats?: boolean;
  special_needs?: boolean;
  specialNeeds?: boolean;
  neutered?: boolean;
  dewormed?: boolean;
  requires_house_screening?: boolean;
  requiresHouseScreening?: boolean;
  pet_id?: string;
  petId?: string;
  guardian_profile_id?: string;
  guardianProfileId?: string;
  published_at?: string | null;
  publishedAt?: string | null;
  adopted_at?: string | null;
  adoptedAt?: string | null;
  pet?: PetApi;
  guardian_profile?: AdoptionGuardianProfileApi;
  guardianProfile?: AdoptionGuardianProfileApi;
};

type PaginationApi = {
  current_page?: number;
  currentPage?: number;
  per_page?: number;
  perPage?: number;
  total_pages?: number;
  totalPages?: number;
  total_records?: number;
  totalRecords?: number;
};

const mapPetApiToDomain = (pet?: PetApi): Pet | undefined => {
  if (!pet) {
    return undefined;
  }

  const specie = pet.specie ?? pet.species;

  return {
    id: String(pet.id ?? ""),
    userId: String(pet.userId ?? pet.user_id ?? ""),
    name: pet.name ?? "",
    specie: {
      id: String(specie?.id ?? ""),
      name: specie?.name ?? "",
      active: specie?.active ?? true,
      createdAt:
        specie?.createdAt ?? specie?.created_at ?? new Date().toISOString(),
      updatedAt: undefined,
      deactivatedAt: undefined,
    },
    breed: pet.breed ?? "",
    age: pet.age ?? 0,
    weight: pet.weight ?? 0,
    notes: pet.notes ?? "",
    photos: (pet.photos ?? []).map((photo) => ({
      id: String(photo.id ?? ""),
      url: photo.url ?? "",
      active: true,
      createdAt: new Date().toISOString(),
      updatedAt: undefined,
      deactivatedAt: undefined,
    })),
    active: pet.active ?? true,
    createdAt: pet.createdAt ?? pet.created_at ?? new Date().toISOString(),
    updatedAt: pet.updatedAt ?? pet.updated_at ?? undefined,
    deactivatedAt: pet.deactivatedAt ?? pet.deactivated_at ?? undefined,
  };
};

const mapListingApiToDomain = (
  listing: AdoptionListingApi,
): AdoptionListing => {
  const guardian = listing.guardianProfile ?? listing.guardian_profile;

  return {
    id: listing.id ?? "",
    active: listing.active ?? true,
    createdAt:
      listing.createdAt ?? listing.created_at ?? new Date().toISOString(),
    updatedAt: listing.updatedAt ?? listing.updated_at ?? undefined,
    deactivatedAt: listing.deactivatedAt ?? listing.deactivated_at ?? undefined,
    title: listing.title ?? "",
    description: listing.description ?? "",
    adoptionReason: listing.adoptionReason ?? listing.adoption_reason ?? "",
    status: listing.status ?? "",
    sex: listing.sex ?? "",
    size: listing.size ?? "",
    ageGroup: listing.ageGroup ?? listing.age_group ?? "",
    stateId: listing.stateId ?? listing.state_id ?? "",
    cityId: listing.cityId ?? listing.city_id ?? "",
    latitude: listing.latitude,
    longitude: listing.longitude,
    goodWithChildren:
      listing.goodWithChildren ?? listing.good_with_children ?? false,
    goodWithDogs: listing.goodWithDogs ?? listing.good_with_dogs ?? false,
    goodWithCats: listing.goodWithCats ?? listing.good_with_cats ?? false,
    specialNeeds: listing.specialNeeds ?? listing.special_needs ?? false,
    neutered: listing.neutered ?? false,
    dewormed: listing.dewormed ?? false,
    requiresHouseScreening:
      listing.requiresHouseScreening ??
      listing.requires_house_screening ??
      false,
    petId: listing.petId ?? listing.pet_id ?? "",
    guardianProfileId:
      listing.guardianProfileId ?? listing.guardian_profile_id ?? "",
    publishedAt: listing.publishedAt ?? listing.published_at ?? undefined,
    adoptedAt: listing.adoptedAt ?? listing.adopted_at ?? undefined,
    pet: mapPetApiToDomain(listing.pet),
    guardianProfile: guardian
      ? {
          id: guardian.id ?? "",
          displayName: guardian.displayName ?? guardian.display_name ?? "",
          about: guardian.about ?? "",
          whatsapp: guardian.whatsapp ?? "",
          cityId: guardian.cityId ?? guardian.city_id ?? "",
          stateId: guardian.stateId ?? guardian.state_id ?? "",
        }
      : undefined,
  };
};

export class AdoptionGatewayAxios implements AdoptionGateway {
  constructor(private readonly http: AxiosInstance) {}

  async listPublicListings(
    input?: ListPublicAdoptionListingsInput,
  ): Promise<ListPublicAdoptionListingsOutput> {
    const params: Record<string, string | number> = {};

    if (input?.page) params.page = input.page;
    if (input?.pageSize) params.page_size = input.pageSize;
    if (input?.sex) params.sex = input.sex;
    if (input?.size) params.size = input.size;
    if (input?.ageGroup) params.age_group = input.ageGroup;
    if (input?.cityId) params.city_id = input.cityId;
    if (input?.stateId) params.state_id = input.stateId;

    const { data } = await this.http.get<{
      listings?: AdoptionListingApi[];
      pagination?: PaginationApi;
    }>("/adoption/listings", { params });

    const pagination = data.pagination ?? {};

    return {
      listings: (data.listings ?? []).map(mapListingApiToDomain),
      pagination: {
        currentPage: pagination.currentPage ?? pagination.current_page ?? 1,
        perPage:
          pagination.perPage ?? pagination.per_page ?? input?.pageSize ?? 12,
        totalPages: pagination.totalPages ?? pagination.total_pages ?? 1,
        totalRecords:
          pagination.totalRecords ??
          pagination.total_records ??
          data.listings?.length ??
          0,
      },
    };
  }

  async getPublicListing(
    listingId: string | number,
  ): Promise<GetPublicAdoptionListingOutput> {
    const { data } = await this.http.get<{ listing?: AdoptionListingApi }>(
      `/adoption/listings/${listingId}`,
    );

    return {
      listing: mapListingApiToDomain(data.listing ?? {}),
    };
  }
}
