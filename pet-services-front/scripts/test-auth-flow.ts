import {
  LoginUserUseCase,
  RegisterUserUseCase,
  ResendVerificationEmailUseCase,
  VerifyEmailUseCase,
} from "../src/application";
import {
  AuthGatewayAxios,
  createApiClient,
  UserGatewayAxios,
} from "../src/infra";

const apiBaseUrl = process.env.API_URL;
const http = createApiClient(apiBaseUrl);

const authGateway = new AuthGatewayAxios(http);
const userGateway = new UserGatewayAxios(http);

const registerUseCase = new RegisterUserUseCase(userGateway);
const resendUseCase = new ResendVerificationEmailUseCase(authGateway);
const verifyUseCase = new VerifyEmailUseCase(authGateway);
const loginUseCase = new LoginUserUseCase(authGateway);

const password = "123QWEasd@";
const email = "guilherme.o.a.ufal@gmail.com";

const run = async () => {
  console.log("→ Register user");
  await registerUseCase.execute({
    name: "Guilherme de Oliveira Amorim",
    userType: "owner",
    login: { email, password },
    phone: {
      countryCode: "55",
      areaCode: "82",
      number: "999767761",
    },
    address: {
      street: "Rua Rafael Pereira Rodrigues",
      number: "28",
      neighborhood: "Grageru",
      city: "Aracaju",
      zipCode: "49027015",
      state: "SE",
      country: "Brasil",
      complement: "Condomínio Verdes Mares, Bloco T, Apartamento 402",
      location: {
        latitude: -10.941262807413592,
        longitude: -37.06499679433729,
      },
    },
  });
  console.log("Register ok:", email);

  console.log("→ Resend verification email");
  const resend = await resendUseCase.execute({ email });
  console.log("Resend ok:", {
    verifyToken: resend.verifyToken,
    expiresAt: resend.expiresAt,
  });

  if (!resend.verifyToken) {
    throw new Error("Token de verificação não retornado");
  }

  console.log("→ Verify email");
  const verify = await verifyUseCase.execute({ token: resend.verifyToken });
  console.log("Verify ok:", verify);

  console.log("→ Login");
  const login = await loginUseCase.execute({ email, password });
  console.log("Login ok:", {
    userId: login.user.id,
    email: login.user.login.email,
    expiresIn: login.expiresIn,
  });
};

run().catch((error) => {
  const status = error?.response?.status;
  const data = error?.response?.data;
  const message = error?.message;
  console.error("Erro no fluxo de auth:", { status, data, message });
  process.exit(1);
});
