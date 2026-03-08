import { createAuthUseCases, createUserUseCases } from "../../src/application";
import { UserTypes } from "../../src/domain";
import { createApiContext } from "../../src/infra";

const email = process.env.TEST_EMAIL;
const password = process.env.TEST_PASSWORD ?? "123QWEasd@";

if (!email) {
  console.error("Defina TEST_EMAIL no ambiente.");
  process.exit(1);
}

const apiBaseUrl = process.env.API_URL;
const { authGateway, userGateway } = createApiContext(apiBaseUrl);
const { loginUseCase } = createAuthUseCases(authGateway);
const { getProfileUseCase } = createUserUseCases(userGateway);

const run = async () => {
  console.log("→ Login");
  const login = await loginUseCase.execute({ email, password });
  console.log("Login ok:", {
    userId: login.user.id,
    email: login.user.login.email,
    userType: login.user.userType,
  });

  authGateway.setAccessToken(login.accessToken);

  console.log("→ Get profile");
  const profile = await getProfileUseCase.execute();
  console.log("Profile response:", {
    userId: profile.user.id,
    userType: profile.user.userType,
    providerId: profile.providerId,
  });

  if (profile.user.userType === UserTypes.Provider && !profile.providerId) {
    throw new Error(
      "Usuário provider sem provider_id no GET /users/profile.",
    );
  }

  console.log("✔ Teste de profile/provider_id concluído com sucesso.");
};

run().catch((error) => {
  const status = error?.response?.status;
  const data = error?.response?.data;
  const message = error?.message;
  console.error("Erro no teste de get-profile-provider-id:", {
    status,
    data,
    message,
  });
  process.exit(1);
});
