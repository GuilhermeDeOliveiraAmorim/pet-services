"use client";

import Link from "next/link";
import { useMemo, useState } from "react";
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
  useAdoptionApplicationWithdraw,
  useMyAdoptionApplications,
} from "@/application";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui";
import { getApiErrorMessage } from "@/lib/api-error";

const statusLabels: Record<string, string> = {
  submitted: "Enviada",
  under_review: "Em análise",
  interview: "Entrevista",
  approved: "Aprovada",
  rejected: "Recusada",
  withdrawn: "Retirada",
};

const statusColorPalette = (status: string) => {
  if (status === "approved") return "green";
  if (status === "rejected") return "red";
  if (status === "withdrawn") return "gray";
  if (status === "interview") return "purple";
  if (status === "under_review") return "blue";
  return "yellow";
};

const formatDate = (value?: string) => {
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

const canWithdraw = (status: string) => {
  return (
    status !== "withdrawn" && status !== "approved" && status !== "rejected"
  );
};

export default function AdoptionApplicationsPage() {
  const [actionError, setActionError] = useState("");
  const [actionSuccess, setActionSuccess] = useState("");
  const [pendingWithdrawId, setPendingWithdrawId] = useState<string | null>(
    null,
  );

  const {
    data,
    isLoading,
    error: listError,
  } = useMyAdoptionApplications({
    page: 1,
    pageSize: 50,
  });

  const { mutateAsync: withdrawApplication, isPending: isWithdrawing } =
    useAdoptionApplicationWithdraw();

  const listErrorMessage = useMemo(() => {
    if (!listError) {
      return "";
    }

    return getApiErrorMessage(
      listError,
      "Não foi possível carregar suas candidaturas de adoção.",
    );
  }, [listError]);

  const applications: AdoptionApplicationItem[] = data?.applications ?? [];

  const handleWithdraw = async (applicationId: string) => {
    setActionError("");
    setActionSuccess("");
    setPendingWithdrawId(applicationId);

    try {
      await withdrawApplication({ applicationId });
      setActionSuccess("Candidatura retirada com sucesso.");
    } catch (error) {
      setActionError(
        getApiErrorMessage(error, "Não foi possível retirar a candidatura."),
      );
    } finally {
      setPendingWithdrawId(null);
    }
  };

  return (
    <PageWrapper gap={10}>
      <MainNav />

      <VStack align="stretch" gap={2}>
        <Heading size="2xl">Minhas candidaturas de adoção</Heading>
        <Text color="gray.600">
          Acompanhe o andamento e retire candidaturas ainda não finalizadas.
        </Text>
      </VStack>

      {actionSuccess ? (
        <Box
          borderRadius="2xl"
          borderWidth="1px"
          borderColor="green.200"
          bg="green.50"
          p={4}
        >
          <Text fontSize="sm" color="green.700">
            {actionSuccess}
          </Text>
        </Box>
      ) : null}

      {actionError ? (
        <Box
          borderRadius="2xl"
          borderWidth="1px"
          borderColor="red.200"
          bg="red.50"
          p={4}
        >
          <Text fontSize="sm" color="red.700">
            {actionError}
          </Text>
        </Box>
      ) : null}

      {isLoading ? (
        <LoadingState message="Carregando candidaturas..." />
      ) : listErrorMessage ? (
        <ErrorState message={listErrorMessage} />
      ) : applications.length === 0 ? (
        <EmptyState message="Você ainda não enviou candidaturas de adoção.">
          <Button asChild mt={4} borderRadius="full">
            <Link href="/adoption">Explorar pets para adoção</Link>
          </Button>
        </EmptyState>
      ) : (
        <VStack align="stretch" gap={4}>
          {applications.map((application) => {
            const withdrawDisabled =
              !canWithdraw(application.status) || isWithdrawing;
            const isCurrentWithdraw = pendingWithdrawId === application.id;

            return (
              <Box
                key={application.id}
                borderRadius="2xl"
                borderWidth="1px"
                borderColor="gray.200"
                bg="white"
                p={{ base: 4, md: 5 }}
              >
                <VStack align="stretch" gap={3}>
                  <Flex justify="space-between" align="start" gap={3}>
                    <VStack align="stretch" gap={1}>
                      <Text fontWeight="semibold">
                        Candidatura #{application.id}
                      </Text>
                      <Text fontSize="sm" color="gray.500">
                        Anúncio: {application.listingId}
                      </Text>
                    </VStack>
                    <Badge
                      colorPalette={statusColorPalette(application.status)}
                      borderRadius="full"
                    >
                      {statusLabels[application.status] ?? application.status}
                    </Badge>
                  </Flex>

                  <Text color="gray.700" fontSize="sm">
                    {application.motivation || "Sem motivação registrada."}
                  </Text>

                  <HStack gap={6} align="start" wrap="wrap">
                    <VStack align="stretch" gap={0}>
                      <Text fontSize="xs" color="gray.500">
                        Telefone
                      </Text>
                      <Text fontSize="sm">
                        {application.contactPhone || "-"}
                      </Text>
                    </VStack>
                    <VStack align="stretch" gap={0}>
                      <Text fontSize="xs" color="gray.500">
                        Criada em
                      </Text>
                      <Text fontSize="sm">
                        {formatDate(application.createdAt)}
                      </Text>
                    </VStack>
                    <VStack align="stretch" gap={0}>
                      <Text fontSize="xs" color="gray.500">
                        Revisada em
                      </Text>
                      <Text fontSize="sm">
                        {formatDate(application.reviewedAt)}
                      </Text>
                    </VStack>
                  </HStack>

                  <HStack justify="space-between" align="center" pt={2}>
                    <Button
                      asChild
                      variant="outline"
                      borderRadius="full"
                      size="sm"
                    >
                      <Link href={`/adoption/${application.listingId}`}>
                        Ver anúncio
                      </Link>
                    </Button>

                    <Button
                      borderRadius="full"
                      size="sm"
                      colorPalette="red"
                      variant="subtle"
                      disabled={withdrawDisabled}
                      onClick={() => handleWithdraw(application.id)}
                    >
                      {isCurrentWithdraw
                        ? "Retirando..."
                        : "Retirar candidatura"}
                    </Button>
                  </HStack>
                </VStack>
              </Box>
            );
          })}
        </VStack>
      )}
    </PageWrapper>
  );
}
