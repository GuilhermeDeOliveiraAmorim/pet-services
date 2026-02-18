import { createAuthUseCases } from "../../src/application";
import { createPetCases } from "../../src/application/factories/pet-usecase-factory";
import { createApiContext } from "../../src/infra";

const email = process.env.TEST_EMAIL;
const password = process.env.TEST_PASSWORD ?? "123QWEasd@";

if (!email) {
  console.error("Defina TEST_EMAIL no ambiente.");
  process.exit(1);
}

const apiBaseUrl = process.env.API_URL;
const { authGateway, petGateway } = createApiContext(apiBaseUrl);
const { loginUseCase } = createAuthUseCases(authGateway);
const {
  listPets: listPetsUseCase,
  getPet: getPetUseCase,
  addPet: addPetUseCase,
  updatePet: updatePetUseCase,
  deletePet: deletePetUseCase,
} = createPetCases(petGateway);

const run = async () => {
  console.log("→ Login");
  const login = await loginUseCase.execute({ email, password });
  console.log("Login ok:", {
    userId: login.user.id,
    email: login.user.login.email,
  });

  petGateway.setAccessToken(login.accessToken);

  console.log("→ List pets");
  const listResult = await listPetsUseCase.execute();
  console.log("List ok:", { count: listResult.pets.length });

  const existingPetId = listResult.pets[0]?.id;

  if (existingPetId) {
    console.log("→ Get pet");
    const getResult = await getPetUseCase.execute(existingPetId);
    console.log("Get ok:", {
      id: getResult.pet.id,
      name: getResult.pet.name,
      age: getResult.pet.age,
    });

    console.log("→ Update pet");
    const updateResult = await updatePetUseCase.execute({
      petId: existingPetId,
      age: getResult.pet.age + 1,
      notes: `Updated at ${Date.now()}`,
    });
    console.log(
      "Update ok:",
      updateResult.message ?? updateResult.detail ?? "ok",
    );
  } else {
    console.log("ℹ No pets found. Skipping get/update tests.");
  }

  console.log("→ Add pet");
  const addResult = await addPetUseCase.execute({
    name: `Pet Test ${Date.now()}`,
    speciesId: "01KG7BG1XN0BQ4KHKPHY5V5ZEW", // Cachorro
    age: 2,
    weight: 5.5,
    notes: "Test pet created at " + new Date().toISOString(),
  });
  console.log("Add ok:", {
    id: addResult.pet?.id,
    name: addResult.pet?.name,
    message: addResult.message,
  });

  if (addResult.pet?.id) {
    console.log("→ Delete pet");
    const deleteResult = await deletePetUseCase.execute(addResult.pet.id);
    console.log(
      "Delete ok:",
      deleteResult.message ?? deleteResult.detail ?? "ok",
    );
  }
};

run().catch((error) => {
  const status = error?.response?.status;
  const data = error?.response?.data;
  const message = error?.message;
  console.error("Erro no teste de pets:", { status, data, message });

  if (status === 404) {
    console.log("\nℹ Verificar se:");
    console.log("  - API está rodando em", apiBaseUrl);
    console.log("  - Usuário tem tipo 'owner'");
    console.log("  - Perfil do usuário está completo");
  }

  process.exit(1);
});
