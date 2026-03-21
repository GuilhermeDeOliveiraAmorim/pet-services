import axios from "axios";
import type {
  ChangeAdoptionListingStatusInput,
  ChangeAdoptionListingStatusOutput,
  AdoptionApplicationItem,
  AdoptionGuardianProfile,
  AdoptionApplicationsPagination,
  CreateAdoptionApplicationInput,
  CreateAdoptionApplicationOutput,
  CreateAdoptionGuardianProfileInput,
  CreateAdoptionGuardianProfileOutput,
  CreateAdoptionListingInput,
  CreateAdoptionListingOutput,
  GetMyAdoptionGuardianProfileInput,
  GetMyAdoptionGuardianProfileOutput,
  GetPublicAdoptionListingInput,
  GetPublicAdoptionListingOutput,
  ListMyAdoptionApplicationsInput,
  ListMyAdoptionApplicationsOutput,
  ListMyAdoptionListingsInput,
  ListMyAdoptionListingsOutput,
  ListAdoptionApplicationsByListingInput,
  ListAdoptionApplicationsByListingOutput,
  ListPublicAdoptionListingsInput,
  ListPublicAdoptionListingsOutput,
  MarkAdoptionListingAsAdoptedInput,
  MarkAdoptionListingAsAdoptedOutput,
  ReviewAdoptionApplicationInput,
  ReviewAdoptionApplicationOutput,
  UpdateAdoptionGuardianProfileInput,
  UpdateAdoptionGuardianProfileOutput,
  UpdateAdoptionListingInput,
  UpdateAdoptionListingOutput,
  WithdrawAdoptionApplicationOutput,
} from "@/application/usecases/adoption";
import type { AdoptionListing, Pet } from "@/domain";
import type { AxiosInstance } from "axios";

type AdoptionGuardianProfileApi = {
  id?: string;
  user_id?: string;
  userId?: string;
  display_name?: string;
  displayName?: string;
  guardian_type?: string;
  guardianType?: string;
  document?: string;
  phone?: string;
  about?: string;
  whatsapp?: string;
  city_id?: string;
  cityId?: string;
  state_id?: string;
  stateId?: string;
  approval_status?: string;
  approvalStatus?: string;
  approved_by?: string;
  approvedBy?: string;
  approved_at?: string;
  approvedAt?: string;
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
  vaccinated?: boolean;
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

type AdoptionApplicationApi = {
  id?: string;
  listing_id?: string;
  listingId?: string;
  status?: string;
  motivation?: string;
  contact_phone?: string;
  contactPhone?: string;
  reviewed_at?: string;
  reviewedAt?: string;
  created_at?: string;
  createdAt?: string;
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
    vaccinated: listing.vaccinated ?? false,
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

export class AdoptionGatewayAxios {
  constructor(private readonly http: AxiosInstance) {}

  async createGuardianProfile(
    input: CreateAdoptionGuardianProfileInput,
  ): Promise<CreateAdoptionGuardianProfileOutput> {
    const { data } = await this.http.post<{
      message?: string;
      detail?: string;
      profile?: AdoptionGuardianProfileApi;
    }>("/adoption/guardian-profile", {
      display_name: input.displayName,
      guardian_type: input.guardianType,
      document: input.document,
      phone: input.phone,
      whatsapp: input.whatsapp,
      about: input.about,
      city_id: input.cityId,
      state_id: input.stateId,
    });

    const profileApi = data.profile;

    return {
      message: data.message,
      detail: data.detail,
      profile: profileApi
        ? {
            id: profileApi.id ?? "",
            userId: profileApi.userId ?? profileApi.user_id ?? "",
            displayName:
              profileApi.displayName ?? profileApi.display_name ?? "",
            guardianType:
              profileApi.guardianType ?? profileApi.guardian_type ?? "",
            document: profileApi.document ?? "",
            phone: profileApi.phone ?? "",
            whatsapp: profileApi.whatsapp ?? "",
            about: profileApi.about ?? "",
            cityId: profileApi.cityId ?? profileApi.city_id ?? "",
            stateId: profileApi.stateId ?? profileApi.state_id ?? "",
            approvalStatus:
              profileApi.approvalStatus ?? profileApi.approval_status ?? "",
            approvedBy:
              profileApi.approvedBy ?? profileApi.approved_by ?? undefined,
            approvedAt:
              profileApi.approvedAt ?? profileApi.approved_at ?? undefined,
          }
        : undefined,
    };
  }

  async updateGuardianProfile(
    input: UpdateAdoptionGuardianProfileInput,
  ): Promise<UpdateAdoptionGuardianProfileOutput> {
    const { data } = await this.http.put<{
      message?: string;
      profile?: AdoptionGuardianProfileApi;
    }>("/adoption/guardian-profile/me", {
      display_name: input.displayName,
      guardian_type: input.guardianType,
      document: input.document,
      phone: input.phone,
      whatsapp: input.whatsapp,
      about: input.about,
      city_id: input.cityId,
      state_id: input.stateId,
    });

    const profileApi = data.profile;

    return {
      message: data.message,
      profile: profileApi
        ? {
            id: profileApi.id ?? "",
            userId: profileApi.userId ?? profileApi.user_id ?? "",
            displayName:
              profileApi.displayName ?? profileApi.display_name ?? "",
            guardianType:
              profileApi.guardianType ?? profileApi.guardian_type ?? "",
            document: profileApi.document ?? "",
            phone: profileApi.phone ?? "",
            whatsapp: profileApi.whatsapp ?? "",
            about: profileApi.about ?? "",
            cityId: profileApi.cityId ?? profileApi.city_id ?? "",
            stateId: profileApi.stateId ?? profileApi.state_id ?? "",
            approvalStatus:
              profileApi.approvalStatus ?? profileApi.approval_status ?? "",
            approvedBy:
              profileApi.approvedBy ?? profileApi.approved_by ?? undefined,
            approvedAt:
              profileApi.approvedAt ?? profileApi.approved_at ?? undefined,
          }
        : undefined,
    };
  }

  async listMyListings(
    input?: ListMyAdoptionListingsInput,
  ): Promise<ListMyAdoptionListingsOutput> {
    const params: Record<string, number> = {};

    if (input?.page) params.page = input.page;
    if (input?.pageSize) params.page_size = input.pageSize;

    const { data } = await this.http.get<{
      listings?: AdoptionListingApi[];
      pagination?: PaginationApi;
    }>("/adoption/listings/me", { params });

    const paginationApi = data.pagination ?? {};

    return {
      listings: (data.listings ?? []).map(mapListingApiToDomain),
      pagination: {
        currentPage:
          paginationApi.currentPage ?? paginationApi.current_page ?? 1,
        perPage:
          paginationApi.perPage ??
          paginationApi.per_page ??
          input?.pageSize ??
          10,
        totalPages: paginationApi.totalPages ?? paginationApi.total_pages ?? 1,
        totalRecords:
          paginationApi.totalRecords ??
          paginationApi.total_records ??
          data.listings?.length ??
          0,
      },
    };
  }

  async createListing(
    input: CreateAdoptionListingInput,
  ): Promise<CreateAdoptionListingOutput> {
    const { data } = await this.http.post<{
      message?: string;
      listing?: AdoptionListingApi;
    }>("/adoption/listings", {
      pet_id: input.petId,
      title: input.title,
      description: input.description,
      adoption_reason: input.adoptionReason,
      sex: input.sex,
      size: input.size,
      age_group: input.ageGroup,
      city_id: input.cityId,
      state_id: input.stateId,
      latitude: input.latitude,
      longitude: input.longitude,
      vaccinated: input.vaccinated ?? false,
      neutered: input.neutered ?? false,
      dewormed: input.dewormed ?? false,
      special_needs: input.specialNeeds ?? false,
      good_with_children: input.goodWithChildren ?? false,
      good_with_dogs: input.goodWithDogs ?? false,
      good_with_cats: input.goodWithCats ?? false,
      requires_house_screening: input.requiresHouseScreening ?? false,
    });

    return {
      message: data.message,
      listing: data.listing ? mapListingApiToDomain(data.listing) : undefined,
    };
  }

  async updateListing(
    input: UpdateAdoptionListingInput,
  ): Promise<UpdateAdoptionListingOutput> {
    const { data } = await this.http.put<{
      message?: string;
      listing?: AdoptionListingApi;
    }>(`/adoption/listings/${input.listingId}`, {
      title: input.title,
      description: input.description,
      adoption_reason: input.adoptionReason,
      sex: input.sex,
      size: input.size,
      age_group: input.ageGroup,
      city_id: input.cityId,
      state_id: input.stateId,
      latitude: input.latitude,
      longitude: input.longitude,
      vaccinated: input.vaccinated,
      neutered: input.neutered,
      dewormed: input.dewormed,
      special_needs: input.specialNeeds,
      good_with_children: input.goodWithChildren,
      good_with_dogs: input.goodWithDogs,
      good_with_cats: input.goodWithCats,
      requires_house_screening: input.requiresHouseScreening,
    });

    return {
      message: data.message,
      listing: data.listing ? mapListingApiToDomain(data.listing) : undefined,
    };
  }

  async changeListingStatus(
    input: ChangeAdoptionListingStatusInput,
  ): Promise<ChangeAdoptionListingStatusOutput> {
    const { data } = await this.http.patch<{
      message?: string;
      listing?: AdoptionListingApi;
    }>(`/adoption/listings/${input.listingId}/${input.action}`);

    return {
      message: data.message,
      listing: data.listing ? mapListingApiToDomain(data.listing) : undefined,
    };
  }

  async getMyGuardianProfile(
    input?: GetMyAdoptionGuardianProfileInput,
  ): Promise<GetMyAdoptionGuardianProfileOutput | null> {
    void input;
    try {
      const { data } = await this.http.get<{
        profile?: AdoptionGuardianProfileApi;
      }>("/adoption/guardian-profile/me");

      const profileApi = data.profile;
      if (!profileApi) {
        return null;
      }

      const profile: AdoptionGuardianProfile = {
        id: profileApi.id ?? "",
        userId: profileApi.userId ?? profileApi.user_id ?? "",
        displayName: profileApi.displayName ?? profileApi.display_name ?? "",
        guardianType: profileApi.guardianType ?? profileApi.guardian_type ?? "",
        document: profileApi.document ?? "",
        phone: profileApi.phone ?? "",
        whatsapp: profileApi.whatsapp ?? "",
        about: profileApi.about ?? "",
        cityId: profileApi.cityId ?? profileApi.city_id ?? "",
        stateId: profileApi.stateId ?? profileApi.state_id ?? "",
        approvalStatus:
          profileApi.approvalStatus ?? profileApi.approval_status ?? "",
        approvedBy:
          profileApi.approvedBy ?? profileApi.approved_by ?? undefined,
        approvedAt:
          profileApi.approvedAt ?? profileApi.approved_at ?? undefined,
      };

      return { profile };
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const status = error.response?.status;
        if (status === 403 || status === 404) {
          return null;
        }
      }
      throw error;
    }
  }

  async listMyApplications(
    input?: ListMyAdoptionApplicationsInput,
  ): Promise<ListMyAdoptionApplicationsOutput> {
    const params: Record<string, number> = {};

    if (input?.page) params.page = input.page;
    if (input?.pageSize) params.page_size = input.pageSize;

    const { data } = await this.http.get<{
      applications?: AdoptionApplicationApi[];
      pagination?: PaginationApi;
    }>("/adoption/applications/me", { params });

    const applications: AdoptionApplicationItem[] = (
      data.applications ?? []
    ).map((item) => ({
      id: item.id ?? "",
      listingId: item.listingId ?? item.listing_id ?? "",
      status: item.status ?? "",
      motivation: item.motivation ?? "",
      contactPhone: item.contactPhone ?? item.contact_phone ?? "",
      reviewedAt: item.reviewedAt ?? item.reviewed_at ?? undefined,
      createdAt: item.createdAt ?? item.created_at ?? "",
    }));

    const paginationApi = data.pagination ?? {};
    const pagination: AdoptionApplicationsPagination = {
      currentPage: paginationApi.currentPage ?? paginationApi.current_page ?? 1,
      perPage:
        paginationApi.perPage ??
        paginationApi.per_page ??
        input?.pageSize ??
        10,
      totalPages: paginationApi.totalPages ?? paginationApi.total_pages ?? 1,
      totalRecords:
        paginationApi.totalRecords ??
        paginationApi.total_records ??
        applications.length,
    };

    return {
      applications,
      pagination,
    };
  }

  async listApplicationsByListing(
    input: ListAdoptionApplicationsByListingInput,
  ): Promise<ListAdoptionApplicationsByListingOutput> {
    const params: Record<string, number> = {};

    if (input.page) params.page = input.page;
    if (input.pageSize) params.page_size = input.pageSize;

    const { data } = await this.http.get<{
      applications?: AdoptionApplicationApi[];
      pagination?: PaginationApi;
    }>(`/adoption/listings/${input.listingId}/applications`, { params });

    const applications: AdoptionApplicationItem[] = (
      data.applications ?? []
    ).map((item) => ({
      id: item.id ?? "",
      listingId: item.listingId ?? item.listing_id ?? "",
      status: item.status ?? "",
      motivation: item.motivation ?? "",
      contactPhone: item.contactPhone ?? item.contact_phone ?? "",
      reviewedAt: item.reviewedAt ?? item.reviewed_at ?? undefined,
      createdAt: item.createdAt ?? item.created_at ?? "",
    }));

    const paginationApi = data.pagination ?? {};

    return {
      applications,
      pagination: {
        currentPage:
          paginationApi.currentPage ?? paginationApi.current_page ?? 1,
        perPage:
          paginationApi.perPage ??
          paginationApi.per_page ??
          input.pageSize ??
          10,
        totalPages: paginationApi.totalPages ?? paginationApi.total_pages ?? 1,
        totalRecords:
          paginationApi.totalRecords ??
          paginationApi.total_records ??
          applications.length,
      },
    };
  }

  async reviewApplication(
    input: ReviewAdoptionApplicationInput,
  ): Promise<ReviewAdoptionApplicationOutput> {
    const { data } = await this.http.post<{ id?: string; status?: string }>(
      `/adoption/applications/${input.applicationId}/review`,
      {
        action: input.action,
        notes_internal: input.notesInternal,
      },
    );

    return {
      id: data.id ?? input.applicationId,
      status: data.status ?? input.action,
    };
  }

  async markListingAsAdopted(
    input: MarkAdoptionListingAsAdoptedInput,
  ): Promise<MarkAdoptionListingAsAdoptedOutput> {
    const { data } = await this.http.post<{ id?: string; status?: string }>(
      `/adoption/listings/${input.listingId}/mark-adopted`,
    );

    return {
      id: data.id ?? input.listingId,
      status: data.status ?? "adopted",
    };
  }

  async withdrawApplication(
    applicationId: string,
  ): Promise<WithdrawAdoptionApplicationOutput> {
    const { data } = await this.http.post<{ id?: string; status?: string }>(
      `/adoption/applications/${applicationId}/withdraw`,
    );

    return {
      id: data.id ?? applicationId,
      status: data.status ?? "withdrawn",
    };
  }

  async createApplication(
    input: CreateAdoptionApplicationInput,
  ): Promise<CreateAdoptionApplicationOutput> {
    const payload = {
      listing_id: input.listingId,
      motivation: input.motivation,
      housing_type: input.housingType,
      pet_experience: input.petExperience,
      contact_phone: input.contactPhone,
      family_members: input.familyMembers,
      agrees_home_visit: input.agreesHomeVisit,
      has_other_pets: input.hasOtherPets,
    };

    const { data } = await this.http.post<{
      id?: string;
      listing_id?: string;
      applicant_user_id?: string;
      status?: string;
    }>("/adoption/applications", payload);

    return {
      id: data.id ?? "",
      listingId: data.listing_id ?? input.listingId,
      applicantUserId: data.applicant_user_id ?? "",
      status: data.status ?? "pending",
    };
  }

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
    input: GetPublicAdoptionListingInput,
  ): Promise<GetPublicAdoptionListingOutput> {
    const { data } = await this.http.get<{ listing?: AdoptionListingApi }>(
      `/adoption/listings/${input.listingId}`,
    );

    return {
      listing: mapListingApiToDomain(data.listing ?? {}),
    };
  }
}
