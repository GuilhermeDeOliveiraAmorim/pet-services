import { createAuthUseCases, createUserUseCases } from "../src/application";
import { createApiContext } from "../src/infra";

const email = process.env.TEST_EMAIL;
const oldPassword = process.env.TEST_PASSWORD ?? "123QWEasd@";
const newPassword = process.env.NEW_PASSWORD ?? "123QWEasd@";

if (!email) {
  console.error("Defina TEST_EMAIL no ambiente.");
  process.exit(1);
}

const apiBaseUrl = process.env.API_URL;
const { authGateway, userGateway } = createApiContext(apiBaseUrl);
const { loginUseCase } = createAuthUseCases(authGateway);
const { changePasswordUseCase } = createUserUseCases(userGateway);

const run = async () => {
  console.log("→ Login");
  const login = await loginUseCase.execute({ email, password: oldPassword });
  console.log("Login ok:", {
    userId: login.user.id,
    email: login.user.login.email,
  });

  console.log("→ Change password");
  userGateway.setAccessToken(login.accessToken);
  const change = await changePasswordUseCase.execute({
    userId: login.user.id,
    oldPassword,
    newPassword,
  });
  console.log("Change password ok:", change);
};

run().catch((error) => {
  console.error("Erro no teste de change-password:", error?.response?.data ?? error);
  process.exit(1);
});
