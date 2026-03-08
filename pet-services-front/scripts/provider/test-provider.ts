import {
  createAuthUseCases,
  createProviderCases,
  createUserUseCases,
} from "../../src/application";
import { UserTypes } from "../../src/domain";
import { createApiContext } from "../../src/infra";

const email = process.env.TEST_EMAIL;
const password = process.env.TEST_PASSWORD ?? "123QWEasd@";

if (!email) {
  console.error("Defina TEST_EMAIL no ambiente.");
  process.exit(1);
}

const apiBaseUrl = process.env.API_URL;
const { authGateway, providerGateway, userGateway } =
  createApiContext(apiBaseUrl);
const { loginUseCase } = createAuthUseCases(authGateway);
const { getProfileUseCase } = createUserUseCases(userGateway);
const { addProvider: addProviderUseCase } =
  createProviderCases(providerGateway);

const providerPayload = {
  businessName:
    process.env.TEST_PROVIDER_BUSINESS_NAME ?? `Provider ${Date.now()}`,
  description:
    process.env.TEST_PROVIDER_DESCRIPTION ??
    "Prestador de serviços pet para testes automatizados",
  priceRange: process.env.TEST_PROVIDER_PRICE_RANGE ?? "80-180",
  address: {
    street: process.env.TEST_PROVIDER_STREET ?? "Rua de Teste",
    number: process.env.TEST_PROVIDER_NUMBER ?? "123",
    neighborhood: process.env.TEST_PROVIDER_NEIGHBORHOOD ?? "Centro",
    city: process.env.TEST_PROVIDER_CITY ?? "Maceió",
    zipCode: process.env.TEST_PROVIDER_ZIP_CODE ?? "57000000",
    state: process.env.TEST_PROVIDER_STATE ?? "AL",
    country: process.env.TEST_PROVIDER_COUNTRY ?? "Brasil",
    complement: process.env.TEST_PROVIDER_COMPLEMENT ?? "",
    location: {
      latitude: Number(process.env.TEST_PROVIDER_LATITUDE ?? "-9.6498"),
      longitude: Number(process.env.TEST_PROVIDER_LONGITUDE ?? "-35.7089"),
    },
  },
};

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

  if (profile.user.userType !== UserTypes.Provider) {
    throw new Error(
      `Usuário autenticado não é provider (userType=${profile.user.userType}).`,
    );
  }

  console.log("Profile ok:", {
    userId: profile.user.id,
    userType: profile.user.userType,
  });

  if (providerPayload.priceRange.length > 10) {
    throw new Error(
      `TEST_PROVIDER_PRICE_RANGE excede 10 caracteres (${providerPayload.priceRange.length}).`,
    );
  }

  console.log("→ Add provider (validação de unicidade)");
  try {
    const addResult = await addProviderUseCase.execute(providerPayload);
    console.log("✔ Provider criado com sucesso:", {
      message: addResult.message,
      detail: addResult.detail,
      providerId: addResult.provider?.id,
      businessName: addResult.provider?.businessName,
    });
  } catch (error: unknown) {
    const responseError = error as {
      response?: {
        status?: number;
        data?: {
          errors?: Array<{ detail?: string }>;
          detail?: string;
        };
      };
    };

    const status = responseError.response?.status;
    const data = responseError.response?.data;
    const detail = data?.errors?.[0]?.detail || data?.detail;

    if (status === 409) {
      console.log(
        "✔ Regra atendida: usuário já possui provider (API bloqueou novo cadastro).",
      );
      console.log("Detalhe:", detail ?? "Conflito retornado pela API.");
      return;
    }

    throw responseError;
  }
};

run().catch((error) => {
  const status = error?.response?.status;
  const data = error?.response?.data;
  const message = error?.message;
  console.error("Erro no teste de provider:", { status, data, message });

  if (status === 403) {
    console.log("\nℹ Possíveis causas:");
    console.log("  - Usuário sem tipo provider");
    console.log("  - Cadastro do usuário ainda incompleto para criar provider");
  }

  process.exit(1);
});
