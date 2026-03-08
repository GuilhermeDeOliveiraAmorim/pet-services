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

  if (profile.user.userType !== UserTypes.Owner) {
    throw new Error(
      `Esperado userType=owner, recebido userType=${profile.user.userType}.`,
    );
  }

  if (profile.providerId) {
    throw new Error(
      "Usuário owner não deveria receber provider_id no profile.",
    );
  }

  console.log(
    "✔ Teste de owner/profile sem provider_id concluído com sucesso.",
  );
};

run().catch((error) => {
  const status = error?.response?.status;
  const data = error?.response?.data;
  const message = error?.message;
  console.error("Erro no teste de get-profile-owner:", {
    status,
    data,
    message,
  });
  process.exit(1);
});
