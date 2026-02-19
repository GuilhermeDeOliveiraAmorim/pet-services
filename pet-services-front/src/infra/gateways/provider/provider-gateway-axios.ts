import type {
  AddProviderInput,
  AddProviderOutput,
  GetProviderOutput,
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
      name: input.name,
      email: input.email,
      phone: {
        country_code: input.phone.countryCode,
        area_code: input.phone.areaCode,
        number: input.phone.number,
      },
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

    const { data } = await this.http.post<AddProviderOutput>(
      "/providers",
      payload,
    );

    return data;
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
}
