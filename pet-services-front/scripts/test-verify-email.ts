import { createAuthUseCases } from "../src/application";
import { createApiContext } from "../src/infra";

const email = process.env.TEST_EMAIL;
const verifyToken = process.env.TEST_VERIFY_TOKEN;

if (!email) {
  console.error("Defina TEST_EMAIL no ambiente.");
  process.exit(1);
}

if (!verifyToken) {
  console.error("Defina TEST_VERIFY_TOKEN no ambiente.");
  process.exit(1);
}

const apiBaseUrl = process.env.API_URL;
const { authGateway } = createApiContext(apiBaseUrl);
const { resendVerificationEmailUseCase, verifyEmailUseCase } =
  createAuthUseCases(authGateway);

const run = async () => {
  console.log("→ Resend verification email");
  const resend = await resendVerificationEmailUseCase.execute({ email });
  console.log("Resend ok:", {
    message: resend.message,
    detail: resend.detail,
  });

  console.log("→ Verify email");
  const verify = await verifyEmailUseCase.execute({
    token: verifyToken,
  });
  console.log("Verify ok:", verify);
};

run().catch((error) => {
  console.error("Erro no teste de verificação:", error);
  process.exit(1);
});
