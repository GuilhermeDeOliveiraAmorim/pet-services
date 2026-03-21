"use client";

import Link from "next/link";
import { useState } from "react";
import {
  Badge,
  Box,
  Button,
  Flex,
  Heading,
  HStack,
  Text,
  VStack,
} from "@chakra-ui/react";

import {
  type AdoptionApplicationItem,
  useAdoptionApplicationReview,
  useAdoptionApplicationsByListing,
  useGuardianStatus,
  useAdoptionListingMarkAsAdopted,
  useAdoptionListingStatusChange,
  useMyAdoptionListings,
  useMyAdoptionGuardianProfile,
} from "@/application";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui";
import type { AdoptionListing } from "@/domain";
import { getApiErrorMessage } from "@/lib/api-error";

const listingStatusLabels: Record<string, string> = {
  draft: "Rascunho",
  published: "Publicado",
  paused: "Pausado",
  adopted: "Adotado",
  archived: "Arquivado",
};

const guardianApprovalLabels: Record<string, string> = {
  pending: "Em análise",
  approved: "Aprovado",
  rejected: "Rejeitado",
};

const applicationStatusLabels: Record<string, string> = {
  submitted: "Enviada",
  under_review: "Em análise",
  interview: "Entrevista",
  approved: "Aprovada",
  rejected: "Recusada",
  withdrawn: "Retirada",
};

const listingStatusColorPalette = (status: string) => {
  if (status === "published") return "green";
  if (status === "adopted") return "teal";
  if (status === "paused") return "orange";
  if (status === "archived") return "gray";
  return "yellow";
};

const applicationStatusColorPalette = (status: string) => {
  if (status === "approved") return "green";
  if (status === "rejected") return "red";
  if (status === "withdrawn") return "gray";
  if (status === "interview") return "purple";
  if (status === "under_review") return "blue";
  return "yellow";
};

const formatDate = (value?: string | null) => {
  if (!value) return "-";

  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) {
    return "-";
  }

  return new Intl.DateTimeFormat("pt-BR", {
    dateStyle: "short",
    timeStyle: "short",
  }).format(parsed);
};

function ListingApplicationsPreview({ listingId }: { listingId: string }) {
  const [feedback, setFeedback] = useState<{
    type: "success" | "error";
    message: string;
  } | null>(null);
  const { data, isLoading, error } = useAdoptionApplicationsByListing({
    listingId,
    page: 1,
    pageSize: 5,
  });
  const { mutateAsync: reviewApplication, isPending } =
    useAdoptionApplicationReview();

  if (isLoading) {
    return <LoadingState message="Carregando candidaturas recebidas..." />;
  }

  if (error) {
    return (
      <ErrorState
        message={getApiErrorMessage(
          error,
          "Não foi possível carregar as candidaturas recebidas deste anúncio.",
        )}
      />
    );
  }

  const applications: AdoptionApplicationItem[] = data?.applications ?? [];

  const handleReview = async (
    applicationId: string,
    action: "under_review" | "interview" | "approve" | "reject",
  ) => {
    setFeedback(null);

    try {
      await reviewApplication({ applicationId, action });
      setFeedback({
        type: "success",
        message: "Status da candidatura atualizado com sucesso.",
      });
    } catch (reviewError) {
      setFeedback({
        type: "error",
        message: getApiErrorMessage(
          reviewError,
          "Não foi possível atualizar a candidatura.",
        ),
      });
    }
  };

  if (applications.length === 0) {
    return (
      <EmptyState message="Nenhuma candidatura recebida para este anúncio até o momento." />
    );
  }

  return (
    <VStack align="stretch" gap={3}>
      {feedback ? (
        <Box
          borderRadius="xl"
          borderWidth="1px"
          borderColor={feedback.type === "success" ? "green.200" : "red.200"}
          bg={feedback.type === "success" ? "green.50" : "red.50"}
          p={3}
        >
          <Text
            fontSize="sm"
            color={feedback.type === "success" ? "green.700" : "red.700"}
          >
            {feedback.message}
          </Text>
        </Box>
      ) : null}

      {applications.map((application) => (
        <Box
          key={application.id}
          borderRadius="xl"
          borderWidth="1px"
          borderColor="gray.200"
          bg="white"
          p={4}
        >
          <VStack align="stretch" gap={2}>
            <Flex justify="space-between" align="start" gap={3}>
              <VStack align="stretch" gap={1}>
                <Text fontWeight="semibold">Candidatura #{application.id}</Text>
                <Text fontSize="sm" color="gray.500">
                  Enviada em {formatDate(application.createdAt)}
                </Text>
              </VStack>
              <Badge
                colorPalette={applicationStatusColorPalette(application.status)}
                borderRadius="full"
              >
                {applicationStatusLabels[application.status] ??
                  application.status}
              </Badge>
            </Flex>

            <Text color="gray.700" fontSize="sm">
              {application.motivation || "Sem motivação registrada."}
            </Text>

            <Text fontSize="sm" color="gray.600">
              Contato: {application.contactPhone || "-"}
            </Text>

            <HStack gap={2} wrap="wrap" pt={1}>
              <Button
                size="xs"
                borderRadius="full"
                variant="outline"
                disabled={isPending}
                onClick={() => handleReview(application.id, "under_review")}
              >
                Em análise
              </Button>
              <Button
                size="xs"
                borderRadius="full"
                variant="outline"
                disabled={isPending}
                onClick={() => handleReview(application.id, "interview")}
              >
                Entrevista
              </Button>
              <Button
                size="xs"
                borderRadius="full"
                colorPalette="green"
                variant="subtle"
                disabled={isPending}
                onClick={() => handleReview(application.id, "approve")}
              >
                Aprovar
              </Button>
              <Button
                size="xs"
                borderRadius="full"
                colorPalette="red"
                variant="subtle"
                disabled={isPending}
                onClick={() => handleReview(application.id, "reject")}
              >
                Rejeitar
              </Button>
            </HStack>
          </VStack>
        </Box>
      ))}
    </VStack>
  );
}

export default function MyAdoptionListingsPage() {
  const [feedback, setFeedback] = useState<{
    type: "success" | "error";
    message: string;
  } | null>(null);
  const {
    data: guardianStatusData,
    isLoading: isGuardianLoading,
    isApprovedGuardian,
  } = useGuardianStatus({
    enabled: true,
  });
  const { data: guardianProfileData } = useMyAdoptionGuardianProfile({
    enabled: true,
  });
  const { mutateAsync: changeListingStatus, isPending: isUpdatingStatus } =
    useAdoptionListingStatusChange();
  const { mutateAsync: markAsAdopted, isPending: isMarkingAsAdopted } =
    useAdoptionListingMarkAsAdopted();

  const {
    data: listingsData,
    isLoading: isListingsLoading,
    error: listingsError,
  } = useMyAdoptionListings(
    {
      page: 1,
      pageSize: 20,
    },
    {
      enabled: isApprovedGuardian,
    },
  );

  const listings: AdoptionListing[] = listingsData?.listings ?? [];
  const guardianProfile =
    guardianProfileData?.profile ?? guardianStatusData?.profile;

  const listingsErrorMessage = listingsError
    ? getApiErrorMessage(
        listingsError,
        "Não foi possível carregar seus anúncios de adoção.",
      )
    : "";

  const handleStatusAction = async (
    listingId: string,
    action: "publish" | "pause" | "archive",
  ) => {
    setFeedback(null);

    try {
      const result = await changeListingStatus({ listingId, action });
      setFeedback({
        type: "success",
        message: result.message ?? "Status do anúncio atualizado com sucesso.",
      });
    } catch (actionError) {
      setFeedback({
        type: "error",
        message: getApiErrorMessage(
          actionError,
          "Não foi possível atualizar o status do anúncio.",
        ),
      });
    }
  };

  const handleMarkAsAdopted = async (listingId: string) => {
    setFeedback(null);

    try {
      await markAsAdopted({ listingId });
      setFeedback({
        type: "success",
        message: "Anúncio marcado como adotado com sucesso.",
      });
    } catch (actionError) {
      setFeedback({
        type: "error",
        message: getApiErrorMessage(
          actionError,
          "Não foi possível marcar o anúncio como adotado.",
        ),
      });
    }
  };

  return (
    <PageWrapper gap={10}>
      <MainNav />

      <VStack align="stretch" gap={2}>
        <Heading size="2xl">Meus anúncios de adoção</Heading>
        <Text color="gray.600">
          Área exclusiva para guardiões com perfil de adoção aprovado.
        </Text>
      </VStack>

      {feedback ? (
        <Box
          borderRadius="2xl"
          borderWidth="1px"
          borderColor={feedback.type === "success" ? "green.200" : "red.200"}
          bg={feedback.type === "success" ? "green.50" : "red.50"}
          p={4}
        >
          <Text
            fontSize="sm"
            color={feedback.type === "success" ? "green.700" : "red.700"}
          >
            {feedback.message}
          </Text>
        </Box>
      ) : null}

      {isGuardianLoading ? (
        <LoadingState message="Verificando seu perfil de guardião..." />
      ) : !isApprovedGuardian ? (
        <ErrorState
          message={
            guardianProfile
              ? `Seu perfil de guardião está com status ${guardianApprovalLabels[guardianProfile.approvalStatus] ?? guardianProfile.approvalStatus}.`
              : "Você ainda não possui um perfil de guardião cadastrado."
          }
        >
          <Button asChild mt={4} borderRadius="full">
            <Link href="/adoption/guardian">
              {guardianProfile
                ? "Revisar perfil de guardião"
                : "Criar perfil de guardião"}
            </Link>
          </Button>
        </ErrorState>
      ) : isListingsLoading ? (
        <LoadingState message="Carregando seus anúncios de adoção..." />
      ) : listingsErrorMessage ? (
        <ErrorState message={listingsErrorMessage} />
      ) : listings.length === 0 ? (
        <EmptyState message="Você ainda não possui anúncios de adoção cadastrados.">
          <Button asChild mt={4} borderRadius="full">
            <Link href="/adoption/listings/new">Criar primeiro anúncio</Link>
          </Button>
        </EmptyState>
      ) : (
        <VStack align="stretch" gap={5}>
          <HStack justify="space-between" wrap="wrap">
            <Text color="gray.600">
              Gerencie status, revisão de candidaturas e atualização dos seus
              anúncios.
            </Text>
            <Button asChild borderRadius="full">
              <Link href="/adoption/listings/new">Novo anúncio</Link>
            </Button>
          </HStack>

          {listings.map((listing) => (
            <Box
              key={listing.id}
              borderRadius="2xl"
              borderWidth="1px"
              borderColor="gray.200"
              bg="white"
              p={{ base: 4, md: 5 }}
            >
              <VStack align="stretch" gap={4}>
                <Flex justify="space-between" align="start" gap={3}>
                  <VStack align="stretch" gap={1}>
                    <Heading size="md">{listing.title || "Sem título"}</Heading>
                    <Text fontSize="sm" color="gray.500">
                      Pet: {listing.pet?.name || "-"} | Status: {listing.status}
                    </Text>
                  </VStack>

                  <Badge
                    colorPalette={listingStatusColorPalette(listing.status)}
                    borderRadius="full"
                  >
                    {listingStatusLabels[listing.status] ?? listing.status}
                  </Badge>
                </Flex>

                <HStack gap={6} align="start" wrap="wrap">
                  <VStack align="stretch" gap={0}>
                    <Text fontSize="xs" color="gray.500">
                      Local
                    </Text>
                    <Text fontSize="sm">
                      {listing.cityId || "-"}/{listing.stateId || "-"}
                    </Text>
                  </VStack>

                  <VStack align="stretch" gap={0}>
                    <Text fontSize="xs" color="gray.500">
                      Publicado em
                    </Text>
                    <Text fontSize="sm">{formatDate(listing.publishedAt)}</Text>
                  </VStack>

                  <VStack align="stretch" gap={0}>
                    <Text fontSize="xs" color="gray.500">
                      Criado em
                    </Text>
                    <Text fontSize="sm">{formatDate(listing.createdAt)}</Text>
                  </VStack>
                </HStack>

                <HStack justify="space-between" align="center" pt={1}>
                  <HStack wrap="wrap" gap={2}>
                    <Button
                      asChild
                      variant="outline"
                      borderRadius="full"
                      size="sm"
                    >
                      <Link href={`/adoption/${listing.id}`}>
                        Ver anúncio público
                      </Link>
                    </Button>

                    {listing.status === "draft" ||
                    listing.status === "paused" ? (
                      <Button
                        asChild
                        variant="outline"
                        borderRadius="full"
                        size="sm"
                      >
                        <Link href={`/adoption/listings/${listing.id}/edit`}>
                          Editar
                        </Link>
                      </Button>
                    ) : null}

                    {listing.status === "draft" ||
                    listing.status === "paused" ? (
                      <Button
                        borderRadius="full"
                        size="sm"
                        disabled={isUpdatingStatus}
                        onClick={() =>
                          handleStatusAction(listing.id, "publish")
                        }
                      >
                        Publicar
                      </Button>
                    ) : null}

                    {listing.status === "published" ? (
                      <Button
                        borderRadius="full"
                        size="sm"
                        variant="outline"
                        disabled={isUpdatingStatus}
                        onClick={() => handleStatusAction(listing.id, "pause")}
                      >
                        Pausar
                      </Button>
                    ) : null}

                    {listing.status !== "archived" &&
                    listing.status !== "adopted" ? (
                      <Button
                        borderRadius="full"
                        size="sm"
                        variant="outline"
                        disabled={isUpdatingStatus}
                        onClick={() =>
                          handleStatusAction(listing.id, "archive")
                        }
                      >
                        Arquivar
                      </Button>
                    ) : null}

                    {listing.status === "published" ? (
                      <Button
                        borderRadius="full"
                        size="sm"
                        colorPalette="teal"
                        variant="subtle"
                        disabled={isMarkingAsAdopted}
                        onClick={() => handleMarkAsAdopted(listing.id)}
                      >
                        Marcar como adotado
                      </Button>
                    ) : null}
                  </HStack>
                </HStack>

                <VStack align="stretch" gap={3} pt={2}>
                  <Heading size="sm">Candidaturas recebidas</Heading>
                  <ListingApplicationsPreview listingId={listing.id} />
                </VStack>
              </VStack>
            </Box>
          ))}
        </VStack>
      )}
    </PageWrapper>
  );
}
