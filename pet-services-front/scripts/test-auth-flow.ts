import { createAuthUseCases, createUserUseCases } from "../src/application";
import { createApiContext } from "../src/infra";

const apiBaseUrl = process.env.API_URL;
const { authGateway, userGateway } = createApiContext(apiBaseUrl);
const { registerUserUseCase } = createUserUseCases(userGateway);
const { resendVerificationEmailUseCase, verifyEmailUseCase, loginUseCase } =
  createAuthUseCases(authGateway);

const password = process.env.TEST_PASSWORD ?? "123QWEasd@";
const generatedEmail = `auth-flow-${Date.now()}@example.com`;
const email = process.env.TEST_EMAIL ?? generatedEmail;
const verifyToken = process.env.TEST_VERIFY_TOKEN;

type ErrorWithResponse = {
  response?: {
    status?: number;
    data?: unknown;
  };
  message?: string;
};

const getErrorInfo = (error: unknown) => {
  if (typeof error === "object" && error !== null) {
    const err = error as ErrorWithResponse;
    return {
      status: err.response?.status,
      data: err.response?.data,
      message: err.message,
    };
  }

  return {
    status: undefined,
    data: undefined,
    message: String(error),
  };
};

const run = async () => {
  console.log("→ Register user");
  await registerUserUseCase.execute({
    name: "Guilherme de Oliveira Amorim",
    userType: "owner",
    login: { email, password },
    phone: {
      countryCode: "55",
      areaCode: "82",
      number: "999767761",
    },
  });
  console.log("Register ok:", email);

  console.log("→ Resend verification email");
  const resend = await resendVerificationEmailUseCase.execute({ email });
  console.log("Resend ok:", {
    message: resend.message,
    detail: resend.detail,
  });

  if (!verifyToken) {
    console.log(
      "→ Verify email pulado: defina TEST_VERIFY_TOKEN para validar /auth/verify-email",
    );
  } else {
    console.log("→ Verify email");
    const verify = await verifyEmailUseCase.execute({
      token: verifyToken,
    });
    console.log("Verify ok:", verify);
  }

  console.log("→ Login");
  try {
    const login = await loginUseCase.execute({ email, password });
    console.log("Login ok:", {
      userId: login.user.id,
      email: login.user.login.email,
      expiresIn: login.expiresIn,
    });
  } catch (error) {
    const { status } = getErrorInfo(error);

    if (status === 401) {
      console.log(
        "Login retornou 401 (esperado neste cenário quando a credencial não corresponde ao usuário existente).",
      );
      return;
    }

    if (status === 403 && !verifyToken) {
      console.log(
        "Login retornou 403 (esperado sem TEST_VERIFY_TOKEN, pois o email pode não estar verificado).",
      );
      return;
    }

    throw error;
  }
};

run().catch((error) => {
  const { status, data, message } = getErrorInfo(error);
  console.error("Erro no fluxo de auth:", { status, data, message });
  process.exit(1);
});
