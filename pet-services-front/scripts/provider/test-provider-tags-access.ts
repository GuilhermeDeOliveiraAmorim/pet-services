import {
  createAuthUseCases,
  createTagCases,
  createUserUseCases,
} from "../../src/application";
import { UserTypes } from "../../src/domain";
import { createApiContext } from "../../src/infra";

const email = process.env.TEST_PROVIDER_EMAIL ?? "provider@bol.com";
const password = process.env.TEST_PROVIDER_PASSWORD ?? "123QWEasd@";
const apiBaseUrl =
  process.env.API_URL ?? "https://pet-services-j7x3.onrender.com";

const { authGateway, userGateway, tagGateway, http } =
  createApiContext(apiBaseUrl);
const { loginUseCase } = createAuthUseCases(authGateway);
const { getProfileUseCase } = createUserUseCases(userGateway);
const { listTags } = createTagCases(tagGateway);

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

  console.log("→ Testar acesso às tags (use case)");
  const tagsResult = await listTags.execute({ page: 1, pageSize: 5 });
  console.log("Tags (use case) ok:", {
    total: tagsResult.total,
    returned: tagsResult.tags.length,
    sample: tagsResult.tags.slice(0, 3).map((item) => item.name),
  });

  console.log("→ Testar endpoint bruto /tags");
  const rawResponse = await http.get<{
    tags?: Array<{ id?: string; name?: string }>;
    total?: number;
  }>("/tags", {
    params: { page: 1, page_size: 5 },
  });

  console.log("Endpoint bruto ok:", {
    status: rawResponse.status,
    total: rawResponse.data.total,
    returned: rawResponse.data.tags?.length ?? 0,
  });

  console.log("\n✔ Teste concluído com sucesso.");
};

run().catch((error: any) => {
  const status = error?.response?.status;
  const data = error?.response?.data;
  const message = error?.message;

  console.error("\n✖ Erro no teste de tags para provider:", {
    status,
    message,
    data,
  });

  if (status === 401 || status === 403) {
    console.log("\nℹ Possíveis causas:");
    console.log("  - Ambiente remoto com autorização incorreta para /tags");
    console.log("  - Token inválido/expirado");
    console.log("  - Usuário sem role provider no backend");
  }

  process.exit(1);
});
