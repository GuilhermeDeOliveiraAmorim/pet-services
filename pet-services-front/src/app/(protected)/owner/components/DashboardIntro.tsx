import Link from "next/link";
import { Box, Grid, HStack, Text, VStack } from "@chakra-ui/react";

type DashboardIntroProps = {
  userName?: string;
  profileComplete: boolean;
  hasPets: boolean;
  petsCount: number;
  petsWithPhotosCount: number;
  totalPetPhotosCount: number;
  activeRequestsCount: number;
  pendingRequestsCount: number;
  completedRequestsCount: number;
  isLoadingRequests: boolean;
  requestsErrorMessage?: string;
  latestPetName?: string;
  latestRequestLabel?: string;
};

type SummaryMetricCardProps = {
  label: string;
  value: string;
  helper: string;
};

function SummaryMetricCard({ label, value, helper }: SummaryMetricCardProps) {
  return (
    <Box
      borderRadius={{ base: "xl", md: "2xl" }}
      borderWidth="1px"
      borderColor="green.100"
      bg="whiteAlpha.900"
      p={{ base: 4, md: 5 }}
    >
      <Text
        fontSize="xs"
        fontWeight="semibold"
        textTransform="uppercase"
        color="green.700"
      >
        {label}
      </Text>
      <Text
        mt={2}
        fontSize={{ base: "2xl", md: "3xl" }}
        fontWeight="semibold"
        color="gray.900"
      >
        {value}
      </Text>
      <Text mt={2} fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
        {helper}
      </Text>
    </Box>
  );
}

const formatMetricValue = (
  value: number,
  isLoading: boolean,
  hasError: boolean,
): string => {
  if (isLoading) {
    return "...";
  }

  if (hasError) {
    return "--";
  }

  return String(value);
};

export default function DashboardIntro({
  userName,
  profileComplete,
  hasPets,
  petsCount,
  petsWithPhotosCount,
  totalPetPhotosCount,
  activeRequestsCount,
  pendingRequestsCount,
  completedRequestsCount,
  isLoadingRequests,
  requestsErrorMessage,
  latestPetName,
  latestRequestLabel,
}: DashboardIntroProps) {
  const hasRequestsError = Boolean(requestsErrorMessage);
  const greetingName = userName?.trim() || "tutor";
  const headline = hasPets
    ? pendingRequestsCount > 0
      ? `Você tem ${pendingRequestsCount} solicitação${pendingRequestsCount > 1 ? "ões" : ""} aguardando retorno.`
      : activeRequestsCount > 0
        ? `Há ${activeRequestsCount} atendimento${activeRequestsCount > 1 ? "s" : ""} em andamento.`
        : "Seu painel está pronto para acompanhar pets e atendimentos."
    : "Cadastre seu primeiro pet para começar a solicitar atendimentos.";
  const profileStatusLabel = profileComplete
    ? "Perfil completo"
    : "Complete seu perfil";
  const requestsHelper = hasRequestsError
    ? "Resumo indisponível no momento"
    : pendingRequestsCount > 0
      ? "Acompanhe respostas dos providers"
      : "Sem pendências aguardando retorno";

  return (
    <VStack align="stretch" gap={{ base: 4, md: 6 }}>
      <Box>
        <Text
          fontSize="xs"
          fontWeight="semibold"
          textTransform="uppercase"
          color="green.500"
        >
          Dashboard
        </Text>
        <Text
          mt={2}
          fontSize={{ base: "xl", md: "2xl" }}
          fontWeight="semibold"
          color="gray.900"
        >
          Olá, {greetingName}
        </Text>
        <Text mt={2} fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
          Este é o seu painel inicial. Aqui vão aparecer seus pets e
          agendamentos.
        </Text>
      </Box>

      <Box
        borderRadius={{ base: "2xl", md: "3xl" }}
        borderWidth="1px"
        borderColor="green.100"
        bg="linear-gradient(135deg, #f0fdf4 0%, #ffffff 55%, #ecfeff 100%)"
        px={{ base: 4, sm: 6, md: 7 }}
        py={{ base: 5, md: 7 }}
      >
        <Grid
          gap={{ base: 4, md: 6 }}
          templateColumns={{ base: "1fr", xl: "1.3fr 0.9fr" }}
          alignItems="start"
        >
          <VStack align="stretch" gap={{ base: 4, md: 5 }}>
            <Box>
              <Text
                fontSize="xs"
                fontWeight="semibold"
                textTransform="uppercase"
                color="green.700"
              >
                Resumo da área do owner
              </Text>
              <Text
                mt={2}
                fontSize={{ base: "lg", md: "2xl" }}
                fontWeight="semibold"
                color="gray.900"
              >
                {headline}
              </Text>
              <Text
                mt={2}
                fontSize={{ base: "xs", sm: "sm" }}
                color="gray.600"
                lineHeight="tall"
              >
                {hasPets
                  ? `Você já tem ${petsCount} pet${petsCount > 1 ? "s" : ""} cadastrado${petsCount > 1 ? "s" : ""}. Use este resumo para acompanhar fotos, solicitações e próximos passos.`
                  : "Cadastre seu primeiro pet para liberar solicitações, histórico de atendimentos e acompanhamento completo pelo painel."}
              </Text>
            </Box>

            <Grid
              gap={{ base: 3, md: 4 }}
              templateColumns={{ base: "1fr 1fr", xl: "repeat(4, 1fr)" }}
            >
              <SummaryMetricCard
                label="Pets"
                value={String(petsCount)}
                helper={
                  latestPetName
                    ? `Último cadastro: ${latestPetName}`
                    : "Nenhum pet cadastrado ainda"
                }
              />
              <SummaryMetricCard
                label="Pets com fotos"
                value={String(petsWithPhotosCount)}
                helper={`${totalPetPhotosCount} foto${totalPetPhotosCount === 1 ? "" : "s"} no total`}
              />
              <SummaryMetricCard
                label="Solicitações ativas"
                value={formatMetricValue(
                  activeRequestsCount,
                  isLoadingRequests,
                  hasRequestsError,
                )}
                helper={requestsHelper}
              />
              <SummaryMetricCard
                label="Concluídas"
                value={formatMetricValue(
                  completedRequestsCount,
                  isLoadingRequests,
                  hasRequestsError,
                )}
                helper={
                  latestRequestLabel
                    ? `Última movimentação: ${latestRequestLabel}`
                    : "Seu histórico aparecerá aqui"
                }
              />
            </Grid>

            {requestsErrorMessage ? (
              <Text fontSize={{ base: "xs", sm: "sm" }} color="orange.700">
                {requestsErrorMessage}
              </Text>
            ) : null}
          </VStack>

          <Box
            borderRadius={{ base: "xl", md: "2xl" }}
            borderWidth="1px"
            borderColor="green.100"
            bg="white"
            p={{ base: 4, md: 5 }}
            shadow="sm"
          >
            <Text
              fontSize="xs"
              fontWeight="semibold"
              textTransform="uppercase"
              color="green.700"
            >
              Próximos passos
            </Text>
            <VStack align="stretch" gap={3} mt={4}>
              <Box>
                <Text fontSize="sm" fontWeight="semibold" color="gray.900">
                  {profileStatusLabel}
                </Text>
                <Text
                  mt={1}
                  fontSize={{ base: "xs", sm: "sm" }}
                  color="gray.600"
                >
                  {profileComplete
                    ? "Seu cadastro está apto para pedir serviços e acompanhar retornos dos providers."
                    : "Revise seus dados para melhorar a experiência na hora de solicitar atendimentos."}
                </Text>
              </Box>

              <Box>
                <Text fontSize="sm" fontWeight="semibold" color="gray.900">
                  Organize seus pets
                </Text>
                <Text
                  mt={1}
                  fontSize={{ base: "xs", sm: "sm" }}
                  color="gray.600"
                >
                  {petsWithPhotosCount < petsCount
                    ? "Alguns pets ainda estão sem foto. Completar esse cadastro ajuda na identificação durante o atendimento."
                    : "Seus pets já têm mídia cadastrada. Mantenha peso, idade e observações sempre atualizados."}
                </Text>
              </Box>

              <HStack gap={2} wrap="wrap" pt={2}>
                <Link href="/requests">
                  <Box
                    as="span"
                    display="inline-block"
                    px={4}
                    py={2}
                    borderRadius="full"
                    bg="green.500"
                    color="white"
                    fontSize="sm"
                    fontWeight="semibold"
                    transition="background-color 0.2s ease"
                    _hover={{ bg: "green.600" }}
                  >
                    Ver solicitações
                  </Box>
                </Link>
                <Link href="/services">
                  <Box
                    as="span"
                    display="inline-block"
                    px={4}
                    py={2}
                    borderRadius="full"
                    borderWidth="1px"
                    borderColor="gray.200"
                    color="gray.700"
                    fontSize="sm"
                    fontWeight="semibold"
                    bg="white"
                    transition="background-color 0.2s ease"
                    _hover={{ bg: "gray.50" }}
                  >
                    Buscar serviços
                  </Box>
                </Link>
                <Link href="/profile">
                  <Box
                    as="span"
                    display="inline-block"
                    px={4}
                    py={2}
                    borderRadius="full"
                    borderWidth="1px"
                    borderColor="gray.200"
                    color="gray.700"
                    fontSize="sm"
                    fontWeight="semibold"
                    bg="white"
                    transition="background-color 0.2s ease"
                    _hover={{ bg: "gray.50" }}
                  >
                    Revisar perfil
                  </Box>
                </Link>
              </HStack>
            </VStack>
          </Box>
        </Grid>
      </Box>
    </VStack>
  );
}
