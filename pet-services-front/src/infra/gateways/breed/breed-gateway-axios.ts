import type { BreedGateway } from "@/application/ports/breed-gateway";
import type { ListBreedsOutput } from "@/application/usecases/breed";
import type { AxiosInstance } from "axios";
import { mapBreedFromApi, type BreedApi } from "@/infra/mappers/breed-mapper";

export class BreedGatewayAxios implements BreedGateway {
  constructor(private readonly http: AxiosInstance) {}

  async listBreedsBySpecies(speciesId: string): Promise<ListBreedsOutput> {
    const { data } = await this.http.get<{ breeds: BreedApi[] }>(
      `/util/species/${speciesId}/breeds`,
    );

    return {
      breeds: data.breeds?.map(mapBreedFromApi) ?? [],
    };
  }
}
