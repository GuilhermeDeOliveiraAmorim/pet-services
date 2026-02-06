import { createReferenceUseCases } from "../../src/application";
import { createApiContext } from "../../src/infra";

const apiBaseUrl = process.env.API_URL;
const { referenceGateway } = createApiContext(apiBaseUrl);
const { listCountriesUseCase, listStatesUseCase, listCitiesUseCase } =
  createReferenceUseCases(referenceGateway);

const run = async () => {
  console.log("→ List countries");
  const countries = await listCountriesUseCase.execute();
  console.log("Countries:", countries.countries.length);

  console.log("→ List states");
  const states = await listStatesUseCase.execute();
  console.log("States:", states.states.length);

  const firstState = states.states[0];
  if (firstState) {
    console.log("→ List cities (first state)");
    const cities = await listCitiesUseCase.execute({ stateId: firstState.id });
    console.log("Cities:", cities.cities.length);
  } else {
    console.log("Sem estados para listar cidades");
  }
};

run().catch((error) => {
  console.error("Erro no teste de referência:", error?.response?.data ?? error);
  process.exit(1);
});
