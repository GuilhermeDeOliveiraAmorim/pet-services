"use client";

import Link from "next/link";
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
  useAdoptionApplicationsByListing,
  useGuardianStatus,
  useMyAdoptionListings,
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
  const { data, isLoading, error } = useAdoptionApplicationsByListing({
    listingId,
    page: 1,
    pageSize: 5,
  });

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

  if (applications.length === 0) {
    return (
      <EmptyState message="Nenhuma candidatura recebida para este anúncio até o momento." />
    );
  }

  return (
    <VStack align="stretch" gap={3}>
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
          </VStack>
        </Box>
      ))}
    </VStack>
  );
}

export default function MyAdoptionListingsPage() {
  const { isLoading: isGuardianLoading, isApprovedGuardian } =
    useGuardianStatus({
      enabled: true,
    });

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

  const listingsErrorMessage = listingsError
    ? getApiErrorMessage(
        listingsError,
        "Não foi possível carregar seus anúncios de adoção.",
      )
    : "";

  return (
    <PageWrapper gap={10}>
      <MainNav />

      <VStack align="stretch" gap={2}>
        <Heading size="2xl">Meus anúncios de adoção</Heading>
        <Text color="gray.600">
          Área exclusiva para guardiões com perfil de adoção aprovado.
        </Text>
      </VStack>

      {isGuardianLoading ? (
        <LoadingState message="Verificando seu perfil de guardião..." />
      ) : !isApprovedGuardian ? (
        <ErrorState message="Seu perfil de guardião ainda não está aprovado para acessar esta área." />
      ) : isListingsLoading ? (
        <LoadingState message="Carregando seus anúncios de adoção..." />
      ) : listingsErrorMessage ? (
        <ErrorState message={listingsErrorMessage} />
      ) : listings.length === 0 ? (
        <EmptyState message="Você ainda não possui anúncios de adoção cadastrados." />
      ) : (
        <VStack align="stretch" gap={5}>
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
