import { isAxiosError } from "axios";
import { createAuthUseCases } from "../src/application";
import { createApiContext } from "../src/infra";

const email = process.env.TEST_EMAIL;

if (!email) {
  console.error("Defina TEST_EMAIL no ambiente.");
  process.exit(1);
}

const apiBaseUrl = process.env.API_URL ?? "http://localhost:8080";
const shouldVerify =
  String(process.env.DO_VERIFY_EMAIL ?? "false").toLowerCase() === "true";

const { authGateway } = createApiContext(apiBaseUrl);
const { resendVerificationEmailUseCase, verifyEmailUseCase } =
  createAuthUseCases(authGateway);

const run = async () => {
  console.log("=== Diagnóstico de envio de e-mail ===");
  console.log(`API_URL: ${apiBaseUrl}`);
  console.log(`TEST_EMAIL: ${email}`);
  console.log(`DO_VERIFY_EMAIL: ${shouldVerify}`);

  console.log("\n→ 1) Reenviando e-mail de verificação...");
  const resend = await resendVerificationEmailUseCase.execute({ email });

  console.log("Resposta do backend:");
  console.log({
    message: resend.message,
    detail: resend.detail,
    verifyTokenReturned: Boolean(resend.verifyToken),
    expiresAt: resend.expiresAt,
  });

  if (!resend.verifyToken) {
    console.error(
      "\nO backend não retornou verify_token. Sem token não dá para validar a sequência.",
    );
    process.exit(1);
  }

  if (!shouldVerify) {
    console.log(
      "\n✅ Envio solicitado com sucesso. Agora valide a caixa de entrada do destinatário.",
    );
    console.log(
      "Se estiver usando Mailpit local, acesse: http://localhost:8025",
    );
    console.log(
      "Se estiver usando Brevo, confira logs/estatísticas de envio no painel.",
    );
    return;
  }

  console.log("\n→ 2) Validando token no endpoint /auth/verify-email...");
  const verify = await verifyEmailUseCase.execute({
    token: resend.verifyToken,
  });
  console.log("Resposta da verificação:", verify);
  console.log("\n✅ Fluxo de reenvio + verificação executado com sucesso.");
};

run().catch((error) => {
  console.error("\n❌ Falha no diagnóstico de e-mail.");

  if (isAxiosError(error)) {
    console.error("HTTP status:", error.response?.status);
    console.error("HTTP data:", error.response?.data);
    console.error("Mensagem:", error.message);
  } else {
    console.error(error);
  }

  process.exit(1);
});
