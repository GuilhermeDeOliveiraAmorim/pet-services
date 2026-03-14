import { createAuthUseCases, createRequestCases } from "../../src/application";
import { createApiContext } from "../../src/infra";

const email = process.env.TEST_EMAIL ?? "provider.seed@petservices.local";
const password = process.env.TEST_PASSWORD ?? "Provider@123";
const apiBaseUrl = process.env.API_URL;

const { authGateway, requestGateway } = createApiContext(apiBaseUrl);
const { loginUseCase } = createAuthUseCases(authGateway);
const { listRequests } = createRequestCases(requestGateway);

const run = async () => {
  console.log("-> Login provider");

  const login = await loginUseCase.execute({ email, password });
  authGateway.setAccessToken(login.accessToken);

  console.log("Login ok:", {
    userId: login.user.id,
    email: login.user.login.email,
    userType: login.user.userType,
  });

  console.log("-> Listar solicitacoes");
  const output = await listRequests.execute({
    page: 1,
    pageSize: 50,
  });

  const requests = output.requests ?? [];

  const groupedByStatus = requests.reduce<Record<string, number>>((acc, req) => {
    const key = req.status || "unknown";
    acc[key] = (acc[key] ?? 0) + 1;
    return acc;
  }, {});

  console.log("Resumo:", {
    total: requests.length,
    page: output.page,
    pageSize: output.pageSize,
    groupedByStatus,
  });

  if (!requests.length) {
    console.log("Resultado: provider sem solicitacoes no momento.");
    return;
  }

  console.log("Primeiras solicitacoes:");
  requests.slice(0, 5).forEach((request, index) => {
    console.log(`${index + 1}.`, {
      id: request.id,
      status: request.status,
      serviceId: request.serviceId,
      serviceName: request.serviceName,
      userId: request.userId,
      userName: request.userName,
      petId: request.petId,
      petName: request.pet?.name,
      createdAt: request.createdAt,
    });
  });
};

run().catch((error: unknown) => {
  const err = error as {
    message?: string;
    response?: { status?: number; data?: unknown };
  };

  console.error("Erro ao testar solicitacoes do provider:", {
    message: err.message,
    status: err.response?.status,
    data: err.response?.data,
  });

  process.exit(1);
});
