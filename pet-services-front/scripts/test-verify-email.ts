import { createAuthUseCases } from "../src/application";
import { createApiContext } from "../src/infra";

const email = process.env.TEST_EMAIL;

if (!email) {
  console.error("Defina TEST_EMAIL no ambiente.");
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
    verifyToken: resend.verifyToken,
    expiresAt: resend.expiresAt,
  });

  if (!resend.verifyToken) {
    console.error("Token de verificação não retornado.");
    process.exit(1);
  }

  console.log("→ Verify email");
  const verify = await verifyEmailUseCase.execute({
    token: resend.verifyToken,
  });
  console.log("Verify ok:", verify);
};

run().catch((error) => {
  console.error("Erro no teste de verificação:", error);
  process.exit(1);
});
