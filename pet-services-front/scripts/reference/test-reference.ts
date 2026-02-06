import {
  ListCitiesUseCase,
  ListCountriesUseCase,
  ListStatesUseCase,
} from "../../src/application";
import { createApiClient, ReferenceGatewayAxios } from "../../src/infra";

const apiBaseUrl = process.env.API_URL;
const http = createApiClient(apiBaseUrl);
const referenceGateway = new ReferenceGatewayAxios(http);

const listCountries = new ListCountriesUseCase(referenceGateway);
const listStates = new ListStatesUseCase(referenceGateway);
const listCities = new ListCitiesUseCase(referenceGateway);

const run = async () => {
  console.log("→ List countries");
  const countries = await listCountries.execute();
  console.log("Countries:", countries.countries.length);

  console.log("→ List states");
  const states = await listStates.execute();
  console.log("States:", states.states.length);

  const firstState = states.states[0];
  if (firstState) {
    console.log("→ List cities (first state)");
    const cities = await listCities.execute({ stateId: firstState.id });
    console.log("Cities:", cities.cities.length);
  } else {
    console.log("Sem estados para listar cidades");
  }
};

run().catch((error) => {
  console.error("Erro no teste de referência:", error?.response?.data ?? error);
  process.exit(1);
});
