import {
  createAuthUseCases,
  createCategoryCases,
  createUserUseCases,
} from "../../src/application";
import { UserTypes } from "../../src/domain";
import { createApiContext } from "../../src/infra";

const email = process.env.TEST_PROVIDER_EMAIL ?? "provider@bol.com";
const password = process.env.TEST_PROVIDER_PASSWORD ?? "123QWEasd@";
const apiBaseUrl = process.env.API_URL ?? "http://localhost:8080";

const { authGateway, userGateway, categoryGateway, http } =
  createApiContext(apiBaseUrl);
const { loginUseCase } = createAuthUseCases(authGateway);
const { getProfileUseCase } = createUserUseCases(userGateway);
const { listCategories } = createCategoryCases(categoryGateway);

const run = async () => {
  console.log("→ API:", apiBaseUrl);
  console.log("→ Login provider");

  const login = await loginUseCase.execute({ email, password });
  authGateway.setAccessToken(login.accessToken);

  console.log("Login ok:", {
    userId: login.user.id,
    email: login.user.login.email,
    userType: login.user.userType,
  });

  console.log("→ Validar perfil autenticado");
  const profile = await getProfileUseCase.execute();

  if (profile.user.userType !== UserTypes.Provider) {
    throw new Error(
      `Usuário autenticado não é provider (userType=${profile.user.userType}).`,
    );
  }

  console.log("Perfil ok (provider):", {
    userId: profile.user.id,
    userType: profile.user.userType,
  });

  console.log("→ Testar acesso às categorias (use case)");
  const categoriesResult = await listCategories.execute({
    page: 1,
    pageSize: 5,
  });
  console.log("Categorias (use case) ok:", {
    total: categoriesResult.total,
    returned: categoriesResult.categories.length,
    sample: categoriesResult.categories.slice(0, 3).map((item) => item.name),
  });

  console.log("→ Testar endpoint bruto /util/categories");
  const rawResponse = await http.get<{
    categories?: Array<{ id?: string; name?: string }>;
    total?: number;
  }>("/util/categories", {
    params: { page: 1, page_size: 5 },
  });

  console.log("Endpoint bruto ok:", {
    status: rawResponse.status,
    total: rawResponse.data.total,
    returned: rawResponse.data.categories?.length ?? 0,
  });

  console.log("\n✔ Teste concluído com sucesso.");
};

run().catch((error: any) => {
  const status = error?.response?.status;
  const data = error?.response?.data;
  const message = error?.message;

  console.error("\n✖ Erro no teste de categorias para provider:", {
    status,
    message,
    data,
  });

  if (status === 401 || status === 403) {
    console.log("\nℹ Possíveis causas:");
    console.log(
      "  - Ambiente remoto com autorização incorreta para /util/categories",
    );
    console.log("  - Token inválido/expirado");
    console.log("  - Usuário sem role provider no backend");
  }

  process.exit(1);
});
