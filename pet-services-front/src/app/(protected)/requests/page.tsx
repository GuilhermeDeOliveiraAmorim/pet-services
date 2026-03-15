"use client";

import { useMemo, useState } from "react";
import {
  Badge,
  Box,
  Button,
  Flex,
  Grid,
  Heading,
  HStack,
  Spinner,
  Textarea,
  Text,
  VStack,
} from "@chakra-ui/react";

import {
  useRequestAccept,
  useRequestComplete,
  useRequestList,
  useRequestReject,
  useReviewCreate,
  useReviewList,
  useUserProfile,
} from "@/application";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import { UserTypes } from "@/domain";
import { getApiErrorMessage } from "@/lib/api-error";

type RequestStatusFilter =
  | "all"
  | "pending"
  | "accepted"
  | "rejected"
  | "completed";

const STATUS_OPTIONS: { value: RequestStatusFilter; label: string }[] = [
  { value: "all", label: "Todos" },
  { value: "pending", label: "Pendentes" },
  { value: "accepted", label: "Aceitos" },
  { value: "rejected", label: "Recusados" },
  { value: "completed", label: "Concluídos" },
];

const STATUS_LABELS: Record<Exclude<RequestStatusFilter, "all">, string> = {
  pending: "Pendente",
  accepted: "Aceita",
  rejected: "Recusada",
  completed: "Concluída",
};

const STATUS_COLOR: Record<
  Exclude<RequestStatusFilter, "all">,
  "yellow" | "green" | "red" | "blue"
> = {
  pending: "yellow",
  accepted: "green",
  rejected: "red",
  completed: "blue",
};

const normalizeStatus = (
  status: string,
): Exclude<RequestStatusFilter, "all"> => {
  if (
    status === "pending" ||
    status === "accepted" ||
    status === "rejected" ||
    status === "completed"
  ) {
    return status;
  }

  return "pending";
};

const formatDateTime = (value?: string): string => {
  if (!value) {
    return "-";
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return "-";
  }

  return new Intl.DateTimeFormat("pt-BR", {
    dateStyle: "short",
    timeStyle: "short",
  }).format(date);
};

export default function RequestsPage() {
  const [statusFilter, setStatusFilter] = useState<RequestStatusFilter>("all");
  const [actionSuccessMessage, setActionSuccessMessage] = useState("");
  const [actionErrorMessage, setActionErrorMessage] = useState("");
  const [rejectingRequestId, setRejectingRequestId] = useState<string | null>(
    null,
  );
  const [actionRequestId, setActionRequestId] = useState<string | null>(null);
  const [rejectReason, setRejectReason] = useState("");
  const [reviewingRequestId, setReviewingRequestId] = useState<string | null>(
    null,
  );
  const [reviewRating, setReviewRating] = useState<number>(5);
  const [reviewComment, setReviewComment] = useState("");
  const {
    data: profileData,
    isLoading: isLoadingProfile,
    error: profileError,
  } = useUserProfile();

  const { mutateAsync: acceptRequest, isPending: isAccepting } =
    useRequestAccept();
  const { mutateAsync: rejectRequest, isPending: isRejecting } =
    useRequestReject();
  const { mutateAsync: completeRequest, isPending: isCompleting } =
    useRequestComplete();
  const { mutateAsync: createReview, isPending: isCreatingReview } =
    useReviewCreate();

  const user = profileData?.user;
  const canListRequests =
    user?.userType === UserTypes.Owner || user?.userType === UserTypes.Provider;

  const listInput = useMemo(
    () => ({
      status: statusFilter === "all" ? undefined : statusFilter,
      page: 1,
      pageSize: 50,
    }),
    [statusFilter],
  );

  const {
    data: requestsData,
    isLoading: isLoadingRequests,
    error: requestsError,
  } = useRequestList(listInput, {
    enabled: !isLoadingProfile && canListRequests,
  });

  const requests = requestsData?.requests ?? [];

  const { data: ownerReviewsData } = useReviewList(
    {
      userId: user?.id,
      page: 1,
      pageSize: 100,
    },
    {
      enabled: !isLoadingProfile && user?.userType === UserTypes.Owner,
    },
  );

  const reviewedProviderIds = useMemo(() => {
    const ids = new Set<string>();
    for (const review of ownerReviewsData?.reviews ?? []) {
      if (review.providerId) {
        ids.add(review.providerId);
      }
    }
    return ids;
  }, [ownerReviewsData?.reviews]);
  const requestsErrorMessage = requestsError
    ? getApiErrorMessage(
        requestsError,
        "Não foi possível carregar suas solicitações.",
      )
    : "";

  const profileErrorMessage = profileError
    ? getApiErrorMessage(
        profileError,
        "Não foi possível carregar o perfil do usuário.",
      )
    : "";

  const isProviderUser = user?.userType === UserTypes.Provider;
  const isOwnerUser = user?.userType === UserTypes.Owner;

  const clearActionMessages = () => {
    setActionSuccessMessage("");
    setActionErrorMessage("");
  };

  const handleAcceptRequest = async (requestId: string) => {
    clearActionMessages();
    setActionRequestId(requestId);

    try {
      await acceptRequest({ id: requestId });
      setActionSuccessMessage("Solicitação aceita com sucesso.");
    } catch (error) {
      setActionErrorMessage(
        getApiErrorMessage(error, "Não foi possível aceitar a solicitação."),
      );
    } finally {
      setActionRequestId(null);
    }
  };

  const handleCompleteRequest = async (requestId: string) => {
    clearActionMessages();
    setActionRequestId(requestId);

    try {
      await completeRequest({ id: requestId });
      setActionSuccessMessage("Solicitação concluída com sucesso.");
    } catch (error) {
      setActionErrorMessage(
        getApiErrorMessage(error, "Não foi possível concluir a solicitação."),
      );
    } finally {
      setActionRequestId(null);
    }
  };

  const handleOpenReject = (requestId: string) => {
    clearActionMessages();
    setRejectingRequestId(requestId);
    setRejectReason("");
  };

  const handleCancelReject = () => {
    setRejectingRequestId(null);
    setRejectReason("");
  };

  const handleRejectRequest = async () => {
    if (!rejectingRequestId) {
      return;
    }

    const reason = rejectReason.trim();
    if (!reason) {
      setActionErrorMessage("Informe o motivo da recusa.");
      return;
    }

    clearActionMessages();
    setActionRequestId(rejectingRequestId);

    try {
      await rejectRequest({ id: rejectingRequestId, rejectReason: reason });
      setActionSuccessMessage("Solicitação recusada com sucesso.");
      setRejectingRequestId(null);
      setRejectReason("");
    } catch (error) {
      setActionErrorMessage(
        getApiErrorMessage(error, "Não foi possível recusar a solicitação."),
      );
    } finally {
      setActionRequestId(null);
    }
  };

  const handleOpenReview = (requestId: string) => {
    clearActionMessages();
    setReviewingRequestId(requestId);
    setReviewRating(5);
    setReviewComment("");
  };

  const handleCancelReview = () => {
    setReviewingRequestId(null);
    setReviewRating(5);
    setReviewComment("");
  };

  const handleCreateReview = async (providerId: string, requestId: string) => {
    clearActionMessages();
    setActionRequestId(requestId);

    const comment = reviewComment.trim();
    if (!comment) {
      setActionErrorMessage("Escreva um comentário para enviar a avaliação.");
      setActionRequestId(null);
      return;
    }

    try {
      await createReview({
        providerId,
        rating: reviewRating,
        comment,
      });

      setActionSuccessMessage("Avaliação enviada com sucesso.");
      setReviewingRequestId(null);
      setReviewRating(5);
      setReviewComment("");
    } catch (error) {
      setActionErrorMessage(
        getApiErrorMessage(error, "Não foi possível enviar a avaliação."),
      );
    } finally {
      setActionRequestId(null);
    }
  };

  return (
    <PageWrapper gap={8}>
      <MainNav />

      <Box
        borderRadius={{ base: "2xl", md: "3xl" }}
        bg="white"
        borderWidth="1px"
        borderColor="gray.200"
        p={{ base: 4, sm: 5, md: 7 }}
      >
        <Flex justify="space-between" align="start" gap={4} wrap="wrap">
          <Box>
            <Text
              fontSize={{ base: "xs" }}
              fontWeight="semibold"
              textTransform="uppercase"
              color="teal.600"
            >
              Solicitações
            </Text>
            <Heading
              mt={2}
              as="h1"
              size={{ base: "lg", md: "xl" }}
              color="gray.900"
            >
              Minhas solicitações
            </Heading>
            <Text mt={2} fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
              Acompanhe o andamento das solicitações de serviço.
            </Text>
          </Box>

          <Badge
            borderRadius="full"
            px={3}
            py={1}
            colorPalette="teal"
            variant="subtle"
            fontSize={{ base: "xs" }}
          >
            {requests.length} registros
          </Badge>
        </Flex>

        <HStack mt={5} gap={2} wrap="wrap">
          {STATUS_OPTIONS.map((option) => (
            <Button
              key={option.value}
              type="button"
              size="sm"
              borderRadius="full"
              variant={statusFilter === option.value ? "solid" : "outline"}
              colorPalette={statusFilter === option.value ? "teal" : "gray"}
              onClick={() => setStatusFilter(option.value)}
            >
              {option.label}
            </Button>
          ))}
        </HStack>
      </Box>

      {isLoadingProfile ? (
        <Flex
          borderRadius={{ base: "2xl", md: "3xl" }}
          borderWidth="1px"
          borderColor="gray.200"
          bg="white"
          py={12}
          justify="center"
          align="center"
          gap={3}
          px={4}
        >
          <Spinner color="teal.500" size="sm" />
          <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
            Carregando perfil...
          </Text>
        </Flex>
      ) : null}

      {!isLoadingProfile && !canListRequests ? (
        <Box
          borderRadius={{ base: "2xl", md: "3xl" }}
          borderWidth="1px"
          borderColor={profileErrorMessage ? "red.200" : "orange.200"}
          bg={profileErrorMessage ? "red.50" : "orange.50"}
          p={{ base: 4, md: 6 }}
        >
          <Text
            fontSize={{ base: "xs", sm: "sm" }}
            color={profileErrorMessage ? "red.700" : "orange.700"}
          >
            {profileErrorMessage ||
              "Apenas usuários owner e provider podem visualizar solicitações."}
          </Text>
        </Box>
      ) : null}

      {!isLoadingProfile && canListRequests ? (
        <>
          {actionSuccessMessage ? (
            <Box
              borderRadius={{ base: "2xl", md: "3xl" }}
              borderWidth="1px"
              borderColor="green.200"
              bg="green.50"
              p={{ base: 4, md: 6 }}
            >
              <Text fontSize={{ base: "xs", sm: "sm" }} color="green.700">
                {actionSuccessMessage}
              </Text>
            </Box>
          ) : null}

          {actionErrorMessage ? (
            <Box
              borderRadius={{ base: "2xl", md: "3xl" }}
              borderWidth="1px"
              borderColor="red.200"
              bg="red.50"
              p={{ base: 4, md: 6 }}
            >
              <Text fontSize={{ base: "xs", sm: "sm" }} color="red.700">
                {actionErrorMessage}
              </Text>
            </Box>
          ) : null}

          {isLoadingRequests ? (
            <Flex
              borderRadius={{ base: "2xl", md: "3xl" }}
              borderWidth="1px"
              borderColor="gray.200"
              bg="white"
              py={12}
              justify="center"
              align="center"
              gap={3}
              px={4}
            >
              <Spinner color="teal.500" size="sm" />
              <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
                Carregando solicitações...
              </Text>
            </Flex>
          ) : null}

          {!isLoadingRequests && requestsErrorMessage ? (
            <Box
              borderRadius={{ base: "2xl", md: "3xl" }}
              borderWidth="1px"
              borderColor="red.200"
              bg="red.50"
              p={{ base: 4, md: 6 }}
            >
              <Text fontSize={{ base: "xs", sm: "sm" }} color="red.700">
                {requestsErrorMessage}
              </Text>
            </Box>
          ) : null}

          {!isLoadingRequests && !requestsErrorMessage && !requests.length ? (
            <Box
              borderRadius={{ base: "2xl", md: "3xl" }}
              borderWidth="1px"
              borderColor="gray.200"
              bg="white"
              p={{ base: 4, md: 6 }}
            >
              <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
                Nenhuma solicitação encontrada para o filtro selecionado.
              </Text>
            </Box>
          ) : null}

          {!isLoadingRequests && !requestsErrorMessage && requests.length ? (
            <VStack align="stretch" gap={4}>
              {requests.map((request) => {
                const status = normalizeStatus(request.status);
                const isLoadingCurrentAction =
                  actionRequestId === request.id &&
                  (isAccepting || isRejecting || isCompleting);
                const canAccept =
                  isProviderUser && request.status === "pending";
                const canReject =
                  isProviderUser && request.status === "pending";
                const canComplete =
                  isProviderUser && request.status === "accepted";
                const isRejectFormOpen = rejectingRequestId === request.id;
                const canReview =
                  isOwnerUser &&
                  request.status === "completed" &&
                  Boolean(request.providerId);
                const isAlreadyReviewed =
                  !!request.providerId &&
                  reviewedProviderIds.has(request.providerId);
                const isReviewFormOpen = reviewingRequestId === request.id;

                return (
                  <Box
                    key={request.id}
                    borderRadius={{ base: "2xl", md: "3xl" }}
                    borderWidth="1px"
                    borderColor="gray.200"
                    bg="white"
                    p={{ base: 4, md: 6 }}
                  >
                    <Flex
                      justify="space-between"
                      align="start"
                      gap={3}
                      wrap="wrap"
                    >
                      <Box>
                        <Text
                          fontSize={{ base: "xs" }}
                          textTransform="uppercase"
                          color="gray.500"
                        >
                          Solicitação
                        </Text>
                        <Text
                          mt={1}
                          fontSize={{ base: "xs", sm: "sm" }}
                          color="gray.700"
                        >
                          #{request.id}
                        </Text>
                      </Box>

                      <Badge
                        borderRadius="full"
                        px={3}
                        py={1}
                        colorPalette={STATUS_COLOR[status]}
                        variant="subtle"
                        fontSize={{ base: "xs" }}
                      >
                        {STATUS_LABELS[status]}
                      </Badge>
                    </Flex>

                    <Grid
                      mt={4}
                      templateColumns={{ base: "1fr", md: "1fr 1fr" }}
                      gap={{ base: 3, md: 4 }}
                    >
                      <Box>
                        <Text
                          fontSize={{ base: "xs" }}
                          color="gray.500"
                          textTransform="uppercase"
                        >
                          Serviço
                        </Text>
                        <Text
                          mt={1}
                          fontSize={{ base: "xs", sm: "sm" }}
                          color="gray.800"
                          fontWeight="medium"
                        >
                          {request.serviceName ||
                            request.serviceId ||
                            "Não informado"}
                        </Text>
                      </Box>

                      <Box>
                        <Text
                          fontSize={{ base: "xs" }}
                          color="gray.500"
                          textTransform="uppercase"
                        >
                          Pet
                        </Text>
                        <Text
                          mt={1}
                          fontSize={{ base: "xs", sm: "sm" }}
                          color="gray.800"
                          fontWeight="medium"
                        >
                          {request.pet?.name ||
                            request.petId ||
                            "Não informado"}
                        </Text>
                      </Box>

                      <Box>
                        <Text
                          fontSize={{ base: "xs" }}
                          color="gray.500"
                          textTransform="uppercase"
                        >
                          Owner
                        </Text>
                        <Text
                          mt={1}
                          fontSize={{ base: "xs", sm: "sm" }}
                          color="gray.800"
                          fontWeight="medium"
                        >
                          {request.userName ||
                            request.userId ||
                            "Não informado"}
                        </Text>
                      </Box>

                      <Box>
                        <Text
                          fontSize={{ base: "xs" }}
                          color="gray.500"
                          textTransform="uppercase"
                        >
                          Provider
                        </Text>
                        <Text
                          mt={1}
                          fontSize={{ base: "xs", sm: "sm" }}
                          color="gray.800"
                          fontWeight="medium"
                        >
                          {request.businessName ||
                            request.providerId ||
                            "Não informado"}
                        </Text>
                      </Box>
                    </Grid>

                    <Box mt={4}>
                      <Text
                        fontSize={{ base: "xs" }}
                        color="gray.500"
                        textTransform="uppercase"
                      >
                        Observações
                      </Text>
                      <Text
                        mt={1}
                        fontSize={{ base: "xs", sm: "sm" }}
                        color="gray.700"
                      >
                        {request.notes || "Sem observações"}
                      </Text>
                    </Box>

                    {request.status === "rejected" && request.rejectReason ? (
                      <Box
                        mt={4}
                        borderRadius={{ base: "lg", md: "xl" }}
                        borderWidth="1px"
                        borderColor="red.200"
                        bg="red.50"
                        p={3}
                      >
                        <Text
                          fontSize={{ base: "xs", sm: "sm" }}
                          color="red.700"
                        >
                          Motivo da recusa: {request.rejectReason}
                        </Text>
                      </Box>
                    ) : null}

                    <HStack mt={4} gap={4} wrap="wrap">
                      <Text fontSize={{ base: "xs" }} color="gray.500">
                        Criada em: {formatDateTime(request.createdAt)}
                      </Text>
                      <Text fontSize={{ base: "xs" }} color="gray.500">
                        Atualizada em:{" "}
                        {formatDateTime(request.updatedAt || undefined)}
                      </Text>
                    </HStack>

                    {(canAccept || canReject || canComplete) && (
                      <Box mt={4}>
                        <Text
                          fontSize={{ base: "xs" }}
                          color="gray.500"
                          textTransform="uppercase"
                        >
                          Ações do provider
                        </Text>

                        <HStack mt={2} gap={2} wrap="wrap">
                          {canAccept ? (
                            <Button
                              size="sm"
                              colorPalette="green"
                              onClick={() =>
                                void handleAcceptRequest(request.id)
                              }
                              disabled={isLoadingCurrentAction}
                            >
                              {isAccepting && actionRequestId === request.id
                                ? "Aceitando..."
                                : "Aceitar"}
                            </Button>
                          ) : null}

                          {canReject ? (
                            <Button
                              size="sm"
                              variant="outline"
                              colorPalette="red"
                              onClick={() => handleOpenReject(request.id)}
                              disabled={isLoadingCurrentAction}
                            >
                              Recusar
                            </Button>
                          ) : null}

                          {canComplete ? (
                            <Button
                              size="sm"
                              colorPalette="blue"
                              onClick={() =>
                                void handleCompleteRequest(request.id)
                              }
                              disabled={isLoadingCurrentAction}
                            >
                              {isCompleting && actionRequestId === request.id
                                ? "Concluindo..."
                                : "Concluir"}
                            </Button>
                          ) : null}
                        </HStack>

                        {isRejectFormOpen ? (
                          <Box
                            mt={3}
                            borderRadius={{ base: "lg", md: "xl" }}
                            borderWidth="1px"
                            borderColor="red.200"
                            bg="red.50"
                            p={3}
                          >
                            <Text
                              fontSize={{ base: "xs", sm: "sm" }}
                              color="red.700"
                              mb={2}
                            >
                              Informe o motivo da recusa
                            </Text>

                            <Textarea
                              value={rejectReason}
                              onChange={(event) =>
                                setRejectReason(event.target.value)
                              }
                              placeholder="Descreva por que a solicitação foi recusada"
                              maxLength={500}
                              bg="white"
                              size="sm"
                            />

                            <Flex mt={2} justify="space-between" align="center">
                              <Text fontSize={{ base: "xs" }} color="gray.500">
                                {rejectReason.trim().length}/500
                              </Text>

                              <HStack>
                                <Button
                                  size="sm"
                                  variant="ghost"
                                  onClick={handleCancelReject}
                                  disabled={isLoadingCurrentAction}
                                >
                                  Cancelar
                                </Button>
                                <Button
                                  size="sm"
                                  colorPalette="red"
                                  onClick={() => void handleRejectRequest()}
                                  disabled={isLoadingCurrentAction}
                                >
                                  {isRejecting && actionRequestId === request.id
                                    ? "Recusando..."
                                    : "Confirmar recusa"}
                                </Button>
                              </HStack>
                            </Flex>
                          </Box>
                        ) : null}
                      </Box>
                    )}

                    {canReview && (
                      <Box mt={4}>
                        <Text
                          fontSize={{ base: "xs" }}
                          color="gray.500"
                          textTransform="uppercase"
                        >
                          Ações do owner
                        </Text>

                        {isAlreadyReviewed ? (
                          <Badge
                            mt={2}
                            borderRadius="full"
                            px={3}
                            py={1}
                            colorPalette="green"
                            variant="subtle"
                            fontSize={{ base: "xs" }}
                          >
                            Você já avaliou este provider
                          </Badge>
                        ) : (
                          <HStack mt={2} gap={2} wrap="wrap">
                            <Button
                              size="sm"
                              colorPalette="teal"
                              onClick={() => handleOpenReview(request.id)}
                              disabled={
                                isCreatingReview ||
                                actionRequestId === request.id
                              }
                            >
                              Avaliar provider
                            </Button>
                          </HStack>
                        )}

                        {isReviewFormOpen && !isAlreadyReviewed ? (
                          <Box
                            mt={3}
                            borderRadius={{ base: "lg", md: "xl" }}
                            borderWidth="1px"
                            borderColor="teal.200"
                            bg="teal.50"
                            p={3}
                          >
                            <Text
                              fontSize={{ base: "xs", sm: "sm" }}
                              color="teal.700"
                              mb={2}
                            >
                              Avalie o atendimento do provider
                            </Text>

                            <HStack mb={2} gap={2} wrap="wrap">
                              {[1, 2, 3, 4, 5].map((value) => (
                                <Button
                                  key={value}
                                  size="sm"
                                  variant={
                                    reviewRating === value ? "solid" : "outline"
                                  }
                                  colorPalette={
                                    reviewRating === value ? "teal" : "gray"
                                  }
                                  onClick={() => setReviewRating(value)}
                                  disabled={
                                    isCreatingReview ||
                                    (actionRequestId === request.id &&
                                      isCreatingReview)
                                  }
                                >
                                  {value}
                                </Button>
                              ))}
                            </HStack>

                            <Textarea
                              value={reviewComment}
                              onChange={(event) =>
                                setReviewComment(event.target.value)
                              }
                              placeholder="Conte como foi sua experiência"
                              maxLength={500}
                              bg="white"
                              size="sm"
                            />

                            <Flex mt={2} justify="space-between" align="center">
                              <Text fontSize={{ base: "xs" }} color="gray.500">
                                {reviewComment.trim().length}/500
                              </Text>

                              <HStack>
                                <Button
                                  size="sm"
                                  variant="ghost"
                                  onClick={handleCancelReview}
                                  disabled={
                                    isCreatingReview ||
                                    (actionRequestId === request.id &&
                                      isCreatingReview)
                                  }
                                >
                                  Cancelar
                                </Button>
                                <Button
                                  size="sm"
                                  colorPalette="teal"
                                  onClick={() =>
                                    void handleCreateReview(
                                      request.providerId,
                                      request.id,
                                    )
                                  }
                                  disabled={
                                    !request.providerId ||
                                    isCreatingReview ||
                                    (actionRequestId === request.id &&
                                      isCreatingReview)
                                  }
                                >
                                  {isCreatingReview &&
                                  actionRequestId === request.id
                                    ? "Enviando..."
                                    : "Enviar avaliação"}
                                </Button>
                              </HStack>
                            </Flex>
                          </Box>
                        ) : null}
                      </Box>
                    )}
                  </Box>
                );
              })}
            </VStack>
          ) : null}
        </>
      ) : null}
    </PageWrapper>
  );
}
