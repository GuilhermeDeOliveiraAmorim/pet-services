import type {
  AddProviderInput,
  AddProviderOutput,
  ProviderGateway,
} from "@/application";
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
}
