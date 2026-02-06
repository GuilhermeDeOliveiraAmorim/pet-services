import {
  ResendVerificationEmailUseCase,
  VerifyEmailUseCase,
} from "../src/application";
import { AuthGatewayAxios, createApiClient } from "../src/infra";

const email = process.env.TEST_EMAIL;

if (!email) {
  console.error("Defina TEST_EMAIL no ambiente.");
  process.exit(1);
}

const apiBaseUrl = process.env.API_URL;
const http = createApiClient(apiBaseUrl);
const authGateway = new AuthGatewayAxios(http);

const resendUseCase = new ResendVerificationEmailUseCase(authGateway);
const verifyUseCase = new VerifyEmailUseCase(authGateway);

const run = async () => {
  console.log("→ Resend verification email");
  const resend = await resendUseCase.execute({ email });
  console.log("Resend ok:", {
    verifyToken: resend.verifyToken,
    expiresAt: resend.expiresAt,
  });

  if (!resend.verifyToken) {
    console.error("Token de verificação não retornado.");
    process.exit(1);
  }

  console.log("→ Verify email");
  const verify = await verifyUseCase.execute({ token: resend.verifyToken });
  console.log("Verify ok:", verify);
};

run().catch((error) => {
  console.error("Erro no teste de verificação:", error);
  process.exit(1);
});
