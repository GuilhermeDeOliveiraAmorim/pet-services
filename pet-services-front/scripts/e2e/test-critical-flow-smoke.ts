import {
  createAuthUseCases,
  createPetCases,
  createRequestCases,
  createReviewCases,
  createServiceCases,
  createUserUseCases,
} from "../../src/application";
import type { Request as RequestEntity } from "../../src/domain/entities/request";
import { UserTypes } from "../../src/domain";
import { createApiContext } from "../../src/infra";

type SessionInfo = {
  userId: string;
  userType: string;
  providerId?: string;
};

const apiBaseUrl = process.env.API_URL;

const ownerEmail =
  process.env.TEST_OWNER_EMAIL ?? "owner.seed@petservices.local";
const ownerPassword = process.env.TEST_OWNER_PASSWORD ?? "Owner@123";
const providerEmail =
  process.env.TEST_PROVIDER_EMAIL ?? "provider.seed@petservices.local";
const providerPassword = process.env.TEST_PROVIDER_PASSWORD ?? "Provider@123";

const searchQuery = process.env.TEST_SERVICE_QUERY;
const maxServices = Number(process.env.TEST_SERVICE_PAGE_SIZE ?? "50");
const rating = Number(process.env.TEST_REVIEW_RATING ?? "5");
const providedPetId = process.env.TEST_OWNER_PET_ID;

const ownerContext = createApiContext(apiBaseUrl);
const providerContext = createApiContext(apiBaseUrl);

const ownerAuth = createAuthUseCases(ownerContext.authGateway);
const ownerUsers = createUserUseCases(ownerContext.userGateway);
const ownerPets = createPetCases(ownerContext.petGateway);
const ownerServices = createServiceCases(ownerContext.serviceGateway);
const ownerRequests = createRequestCases(ownerContext.requestGateway);
const ownerReviews = createReviewCases(ownerContext.reviewGateway);

const providerAuth = createAuthUseCases(providerContext.authGateway);
const providerUsers = createUserUseCases(providerContext.userGateway);
const providerRequests = createRequestCases(providerContext.requestGateway);

const assertCondition = (condition: unknown, message: string) => {
  if (!condition) {
    throw new Error(message);
  }
};

const login = async (
  label: string,
  email: string,
  password: string,
  expectedType: string,
  setAccessToken: (token: string) => void,
  authUseCases: ReturnType<typeof createAuthUseCases>,
  getProfile: ReturnType<typeof createUserUseCases>["getProfileUseCase"],
): Promise<SessionInfo> => {
  console.log(`-> Login ${label}`);

  const auth = await authUseCases.loginUseCase.execute({ email, password });
  setAccessToken(auth.accessToken);

  assertCondition(
    auth.user.userType === expectedType,
    `${label}: tipo de usuario inesperado (${auth.user.userType}). Esperado ${expectedType}.`,
  );

  const profile = await getProfile.execute();

  console.log(`${label} autenticado:`, {
    userId: auth.user.id,
    userType: auth.user.userType,
    providerId: profile.providerId,
  });

  return {
    userId: auth.user.id,
    userType: auth.user.userType,
    providerId: profile.providerId,
  };
};

const run = async () => {
  if (rating < 1 || rating > 5) {
    throw new Error("TEST_REVIEW_RATING deve estar entre 1 e 5.");
  }

  if (!Number.isFinite(maxServices) || maxServices < 1) {
    throw new Error(
      "TEST_SERVICE_PAGE_SIZE deve ser um numero inteiro positivo.",
    );
  }

  const ownerSession = await login(
    "owner",
    ownerEmail,
    ownerPassword,
    UserTypes.Owner,
    ownerContext.authGateway.setAccessToken.bind(ownerContext.authGateway),
    ownerAuth,
    ownerUsers.getProfileUseCase,
  );

  const providerSession = await login(
    "provider",
    providerEmail,
    providerPassword,
    UserTypes.Provider,
    providerContext.authGateway.setAccessToken.bind(
      providerContext.authGateway,
    ),
    providerAuth,
    providerUsers.getProfileUseCase,
  );

  assertCondition(
    providerSession.providerId,
    "Provider autenticado sem providerId no profile.",
  );

  console.log("-> Buscar servicos (search)");
  const searchedServices = await ownerServices.searchServices.execute({
    query: searchQuery,
    page: 1,
    pageSize: maxServices,
  });

  console.log("Busca concluida:", {
    total: searchedServices.total,
    retornados: searchedServices.services.length,
    query: searchQuery ?? "<vazio>",
  });

  let targetService = searchedServices.services.find(
    (service) => service.providerId === providerSession.providerId,
  );

  if (!targetService) {
    console.log(
      "Nenhum servico do provider apareceu no search; fazendo fallback para listServices(providerId).",
    );

    const providerServices = await ownerServices.listServices.execute({
      providerId: providerSession.providerId,
      page: 1,
      pageSize: maxServices,
    });

    targetService = providerServices.services[0];
  }

  assertCondition(
    targetService,
    "Nao foi encontrado servico elegivel para o provider informado.",
  );

  console.log("Servico escolhido:", {
    serviceId: targetService.id,
    providerId: targetService.providerId,
    name: targetService.name,
  });

  let petId = providedPetId;

  if (!petId) {
    console.log("-> Buscar pet do owner");
    const petsResult = await ownerPets.listPets.execute();
    petId = petsResult.pets[0]?.id;
  }

  assertCondition(
    petId,
    "Nenhum pet encontrado para o owner. Configure TEST_OWNER_PET_ID ou crie um pet para o usuario.",
  );

  console.log("-> Solicitar servico");
  const requestNotes = `Smoke E2E ${new Date().toISOString()}`;
  let createdRequest: RequestEntity | null = null;

  try {
    const createdRequestResponse = await ownerRequests.addRequest.execute({
      providerId: targetService.providerId,
      serviceId: targetService.id,
      petId,
      notes: requestNotes,
    });

    createdRequest = createdRequestResponse;

    if (!createdRequest.id) {
      const ownerPendingRequests = await ownerRequests.listRequests.execute({
        userId: ownerSession.userId,
        providerId: targetService.providerId,
        serviceId: targetService.id,
        status: "pending",
        page: 1,
        pageSize: 100,
      });

      const matchedCreatedRequest = ownerPendingRequests.requests.find(
        (request) => request.notes === requestNotes,
      );

      if (matchedCreatedRequest) {
        createdRequest = matchedCreatedRequest;
      }
    }
  } catch (error: unknown) {
    const err = error as { response?: { status?: number } };

    if (err.response?.status !== 409) {
      throw error;
    }

    console.log(
      "Conflito 409 ao solicitar; reutilizando solicitacao existente pendente/aceita para seguir o smoke.",
    );

    const existingRequests = await ownerRequests.listRequests.execute({
      userId: ownerSession.userId,
      providerId: targetService.providerId,
      page: 1,
      pageSize: 100,
    });

    const reusableRequest = existingRequests.requests.find(
      (request) =>
        request.serviceId === targetService.id &&
        (request.status === "pending" || request.status === "accepted"),
    );

    assertCondition(
      reusableRequest,
      "Recebido 409 na solicitacao e nao foi encontrada solicitacao reutilizavel.",
    );

    createdRequest = reusableRequest ?? null;
  }

  if (!createdRequest) {
    throw new Error("Solicitacao criada sem payload valido.");
  }

  assertCondition(createdRequest.id, "Solicitacao criada sem id.");
  assertCondition(
    createdRequest.status === "pending" || createdRequest.status === "accepted",
    `Status apos criacao/reuso deveria ser pending ou accepted, recebido ${createdRequest.status}.`,
  );

  console.log("Solicitacao pronta para fluxo:", {
    requestId: createdRequest.id,
    status: createdRequest.status,
  });

  let requestAfterAccept: RequestEntity = createdRequest;

  if (requestAfterAccept.status === "pending") {
    console.log("-> Aceitar solicitacao (provider)");
    requestAfterAccept = await providerRequests.acceptRequest.execute({
      id: requestAfterAccept.id,
    });

    assertCondition(
      requestAfterAccept.status === "accepted",
      `Status apos aceite deveria ser accepted, recebido ${requestAfterAccept.status}.`,
    );

    console.log("Solicitacao aceita:", {
      requestId: requestAfterAccept.id,
      status: requestAfterAccept.status,
    });
  } else {
    console.log("Solicitacao ja estava aceita. Pulando passo de aceite.");
  }

  console.log("-> Concluir solicitacao (provider)");
  const completedRequest = await providerRequests.completeRequest.execute({
    id: requestAfterAccept.id,
  });

  assertCondition(
    completedRequest.status === "completed",
    `Status apos conclusao deveria ser completed, recebido ${completedRequest.status}.`,
  );

  console.log("Solicitacao concluida:", {
    requestId: completedRequest.id,
    status: completedRequest.status,
  });

  console.log("-> Avaliar provider (owner)");
  const createdReview = await ownerReviews.createReview.execute({
    providerId: targetService.providerId,
    rating,
    comment: `Avaliacao smoke ${new Date().toISOString()}`,
  });

  assertCondition(createdReview.id, "Review criada sem id.");

  const ownerReviewList = await ownerReviews.listReviews.execute({
    userId: ownerSession.userId,
    page: 1,
    pageSize: 100,
  });

  const reviewWasListed = ownerReviewList.reviews.some(
    (review) => review.id === createdReview.id,
  );

  assertCondition(
    reviewWasListed,
    "Review criada nao encontrada na listagem do owner.",
  );

  console.log("Review criada:", {
    reviewId: createdReview.id,
    rating: createdReview.rating,
    providerId: createdReview.providerId,
  });

  console.log(
    "\nOK: Smoke E2E concluido (buscar -> solicitar -> aceitar/concluir -> avaliar).",
  );
};

run().catch((error: unknown) => {
  const err = error as {
    message?: string;
    response?: { status?: number; data?: unknown };
  };

  console.error("Erro no smoke E2E de fluxo critico:", {
    message: err.message,
    status: err.response?.status,
    data: err.response?.data,
  });

  process.exit(1);
});
