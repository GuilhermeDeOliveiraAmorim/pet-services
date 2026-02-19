import type { SpecieGateway } from "@/application/ports/specie-gateway";
import type { ListSpeciesOutput } from "@/application/usecases/specie";
import type { AxiosInstance } from "axios";
import {
  mapSpecieFromApi,
  type SpecieApi,
} from "@/infra/mappers/specie-mapper";

export class SpecieGatewayAxios implements SpecieGateway {
  constructor(private readonly http: AxiosInstance) {}

  async listSpecies(): Promise<ListSpeciesOutput> {
    const { data } = await this.http.get<{ species: SpecieApi[] }>("/species");

    return {
      species: data.species?.map(mapSpecieFromApi) ?? [],
    };
  }
}
