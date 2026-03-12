import type {
  AddProviderInput,
  AddProviderOutput,
  AddProviderPhotoOutput,
  DeleteProviderOutput,
  DeleteProviderPhotoOutput,
  GetProviderOutput,
  ListProvidersOutput,
  ProviderGateway,
  UpdateProviderInput,
  UpdateProviderOutput,
} from "@/application";
import { mapProviderFromApi } from "@/infra/mappers/provider-mapper";
import type { AxiosInstance } from "axios";

export class ProviderGatewayAxios implements ProviderGateway {
  constructor(private readonly http: AxiosInstance) {}

  async addProvider(input: AddProviderInput): Promise<AddProviderOutput> {
    const payload = {
      business_name: input.businessName,
      description: input.description,
      price_range: input.priceRange,
      address: {
        street: input.address.street,
        number: input.address.number,
        neighborhood: input.address.neighborhood,
        city: input.address.city,
        zip_code: input.address.zipCode,
        state: input.address.state,
        country: input.address.country,
        complement: input.address.complement,
        location: {
          latitude: input.address.location.latitude,
          longitude: input.address.location.longitude,
        },
      },
    };

    const { data } = await this.http.post<{
      message?: string;
      detail?: string;
      provider?: unknown;
    }>("/providers", payload);

    return {
      message: data.message,
      detail: data.detail,
      provider: data.provider
        ? mapProviderFromApi(data.provider as Record<string, unknown>)
        : undefined,
    };
  }

  async getProvider(providerId: string | number): Promise<GetProviderOutput> {
    const { data } = await this.http.get<{ provider: unknown }>(
      `/providers/${providerId}`,
    );

    return {
      provider: mapProviderFromApi(data.provider as Record<string, unknown>),
    };
  }

  async updateProvider(
    input: UpdateProviderInput,
  ): Promise<UpdateProviderOutput> {
    const payload: {
      business_name?: string;
      address?: {
        street: string;
        number: string;
        neighborhood: string;
        city: string;
        zip_code: string;
        state: string;
        country: string;
        complement?: string;
        location: {
          latitude: number;
          longitude: number;
        };
      };
      description?: string;
      price_range?: string;
    } = {};

    if (input.businessName) {
      payload.business_name = input.businessName;
    }

    if (input.address) {
      payload.address = {
        street: input.address.street,
        number: input.address.number,
        neighborhood: input.address.neighborhood,
        city: input.address.city,
        zip_code: input.address.zipCode,
        state: input.address.state,
        country: input.address.country,
        complement: input.address.complement,
        location: {
          latitude: input.address.location.latitude,
          longitude: input.address.location.longitude,
        },
      };
    }

    if (input.description) {
      payload.description = input.description;
    }

    if (input.priceRange) {
      payload.price_range = input.priceRange;
    }

    const { data } = await this.http.put<{
      message?: string;
      detail?: string;
      provider?: unknown;
    }>(`/providers/${input.providerId}`, payload);

    return {
      message: data.message,
      detail: data.detail,
      provider: data.provider
        ? mapProviderFromApi(data.provider as Record<string, unknown>)
        : undefined,
    };
  }

  async deleteProvider(
    providerId: string | number,
  ): Promise<DeleteProviderOutput> {
    const { data } = await this.http.delete<DeleteProviderOutput>(
      `/providers/${providerId}`,
    );
    return data;
  }

  async listProviders(): Promise<ListProvidersOutput> {
    const { data } = await this.http.get<{ providers: unknown[] }>(
      "/providers",
    );

    return {
      providers: Array.isArray(data.providers)
        ? data.providers.map((p) =>
            mapProviderFromApi(p as Record<string, unknown>),
          )
        : [],
    };
  }

  async addProviderPhoto(
    _providerId: string | number,
    photo: File,
  ): Promise<AddProviderPhotoOutput> {
    const formData = new FormData();
    formData.append("file", photo);

    const { data } = await this.http.post<{
      message?: string;
      detail?: string;
      photo?: { id: string; url: string };
    }>("/providers/photos", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });

    return {
      message: data.message,
      detail: data.detail,
      photo: data.photo,
    };
  }

  async deleteProviderPhoto(
    providerId: string | number,
    photoId: string | number,
  ): Promise<DeleteProviderPhotoOutput> {
    const { data } = await this.http.delete<{
      message?: string;
      detail?: string;
    }>(`/providers/${providerId}/photos/${photoId}`);

    return {
      message: data.message,
      detail: data.detail,
    };
  }
}
