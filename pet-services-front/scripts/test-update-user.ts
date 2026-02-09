import { createAuthUseCases, createUserUseCases } from "../src/application";
import { createApiContext } from "../src/infra";

const email = process.env.TEST_EMAIL;
const password = process.env.TEST_PASSWORD ?? "123QWEasd@";

if (!email) {
  console.error("Defina TEST_EMAIL no ambiente.");
  process.exit(1);
}

const apiBaseUrl = process.env.API_URL;
const { authGateway, userGateway } = createApiContext(apiBaseUrl);
const { loginUseCase } = createAuthUseCases(authGateway);
const { updateUserUseCase, getProfileUseCase } =
  createUserUseCases(userGateway);

const run = async () => {
  console.log("→ Login");
  const login = await loginUseCase.execute({ email, password });
  console.log("Login ok:", {
    userId: login.user.id,
    email: login.user.login.email,
  });

  userGateway.setAccessToken(login.accessToken);

  console.log("→ Update user");
  const update = await updateUserUseCase.execute({
    name: `Teste ${Date.now()}`,
    phone: {
      countryCode: "55",
      areaCode: "82",
      number: "999888777",
    },
    address: {
      street: "Rua de Teste",
      number: "123",
      neighborhood: "Centro",
      city: "Maceió",
      zipCode: "57000000",
      state: "AL",
      country: "Brasil",
      complement: "Bloco B",
      location: {
        latitude: -9.66599,
        longitude: -35.735,
      },
    },
  });
  console.log("Update ok:", update.message ?? update.detail ?? "ok");

  console.log("→ Get profile");
  const profile = await getProfileUseCase.execute();
  console.log("Profile ok:", {
    name: profile.user.name,
    phone: profile.user.phone,
    address: profile.user.address,
  });
};

run().catch((error) => {
  const status = error?.response?.status;
  const data = error?.response?.data;
  const message = error?.message;
  console.error("Erro no teste de update-user:", { status, data, message });
  process.exit(1);
});
