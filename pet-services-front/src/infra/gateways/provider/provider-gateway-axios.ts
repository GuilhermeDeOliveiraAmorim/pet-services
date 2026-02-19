import type {
  AddProviderInput,
  AddProviderOutput,
  GetProviderOutput,
  ProviderGateway,
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
}
