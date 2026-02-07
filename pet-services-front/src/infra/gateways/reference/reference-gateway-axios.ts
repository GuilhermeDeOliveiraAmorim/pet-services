import type { ReferenceGateway } from "@/application";
import type { City, Country, State } from "@/domain";
import type { AxiosInstance } from "axios";

export class ReferenceGatewayAxios implements ReferenceGateway {
  constructor(private readonly http: AxiosInstance) {}

  async listCountries() {
    const { data } = await this.http.get<{
      countries: Array<{
        name: string;
        dial_code?: string;
        dialCode?: string;
        code: string;
        flag: string;
      }>;
    }>("/reference/countries");

    const countries: Country[] = data.countries.map((country) => ({
      name: country.name,
      dialCode: country.dialCode ?? country.dial_code ?? "",
      code: country.code,
      flag: country.flag,
    }));

    return { countries };
  }

  async listStates() {
    const { data } = await this.http.get<{ states: State[] }>(
      "/reference/states",
    );
    return { states: data.states };
  }

  async listCities(input: { stateId?: number }) {
    const params = input.stateId ? { state_id: input.stateId } : undefined;
    const { data } = await this.http.get<{ cities: City[] }>(
      "/reference/cities",
      { params },
    );
    return { cities: data.cities };
  }
}
