import { createAuthUseCases, createUserUseCases } from "../src/application";
import { createApiContext } from "../src/infra";

const apiBaseUrl = process.env.API_URL;
const { authGateway, userGateway } = createApiContext(apiBaseUrl);
const { registerUserUseCase } = createUserUseCases(userGateway);
const { resendVerificationEmailUseCase, verifyEmailUseCase, loginUseCase } =
  createAuthUseCases(authGateway);

const password = "123QWEasd@";
const email = "guilherme.o.a.ufal@gmail.com";

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
  const resend = await resendVerificationEmailUseCase.execute({ email });
  console.log("Resend ok:", {
    verifyToken: resend.verifyToken,
    expiresAt: resend.expiresAt,
  });

  if (!resend.verifyToken) {
    throw new Error("Token de verificação não retornado");
  }

  console.log("→ Verify email");
  const verify = await verifyEmailUseCase.execute({ token: resend.verifyToken });
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
