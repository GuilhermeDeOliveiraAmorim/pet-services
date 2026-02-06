import { createAuthUseCases } from "../src/application";
import { createApiContext } from "../src/infra";

const email = process.env.TEST_EMAIL;
const password = process.env.TEST_PASSWORD ?? "123QWEasd@";

if (!email) {
  console.error("Defina TEST_EMAIL no ambiente.");
  process.exit(1);
}

const apiBaseUrl = process.env.API_URL;
const { authGateway } = createApiContext(apiBaseUrl);
const { loginUseCase, refreshTokenUseCase, logoutUseCase } =
  createAuthUseCases(authGateway);

const run = async () => {
  console.log("→ Login");
  const login = await loginUseCase.execute({ email, password });
  console.log("Login ok:", {
    userId: login.user.id,
    email: login.user.login.email,
    expiresIn: login.expiresIn,
  });

  console.log("→ Refresh token");
  const refresh = await refreshTokenUseCase.execute({
    refreshToken: login.refreshToken,
  });
  console.log("Refresh ok:", { expiresIn: refresh.expiresIn });

  console.log("→ Logout (revoke all)");
  authGateway.setAccessToken(refresh.accessToken);
  const logout = await logoutUseCase.execute({
    userId: login.user.id,
    revokeAll: true,
  });
  console.log("Logout ok:", logout);
};

run().catch((error) => {
  console.error("Erro no teste de auth:", error);
  process.exit(1);
});
