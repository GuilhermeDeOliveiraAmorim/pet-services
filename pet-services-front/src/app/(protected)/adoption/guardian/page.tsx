"use client";

import Link from "next/link";
import { useMemo, useState } from "react";
import {
  Badge,
  Box,
  Button,
  Heading,
  HStack,
  Text,
  VStack,
} from "@chakra-ui/react";

import {
  useAdoptionGuardianProfileCreate,
  useAdoptionGuardianProfileUpdate,
  useMyAdoptionGuardianProfile,
} from "@/application";
import GuardianProfileForm from "@/components/adoption/GuardianProfileForm";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import { ErrorState, LoadingState } from "@/components/ui";
import { getApiErrorMessage } from "@/lib/api-error";

const approvalLabels: Record<string, string> = {
  pending: "Em análise",
  approved: "Aprovado",
  rejected: "Rejeitado",
};

const approvalPalette = (status?: string) => {
  if (status === "approved") return "green";
  if (status === "rejected") return "red";
  return "yellow";
};

export default function AdoptionGuardianPage() {
  const [feedback, setFeedback] = useState<{
    type: "success" | "error";
    message: string;
  } | null>(null);
  const { data, isLoading, error } = useMyAdoptionGuardianProfile();
  const { mutateAsync: createProfile, isPending: isCreating } =
    useAdoptionGuardianProfileCreate();
  const { mutateAsync: updateProfile, isPending: isUpdating } =
    useAdoptionGuardianProfileUpdate();

  const profile = data?.profile;

  const helperText = useMemo(() => {
    if (!profile) {
      return "Cadastre seu perfil para publicar pets para adoção e acompanhar candidaturas.";
    }
    if (profile.approvalStatus === "approved") {
      return "Seu perfil está apto para publicar anúncios e gerenciar o processo de adoção.";
    }
    if (profile.approvalStatus === "rejected") {
      return "Seu perfil foi rejeitado anteriormente. Atualize os dados e envie novamente para análise.";
    }
    return "Seu perfil está em análise. Você ainda pode ajustar os dados enquanto aguarda revisão.";
  }, [profile]);

  const handleSubmit = async (values: {
    displayName: string;
    guardianType: string;
    document: string;
    phone: string;
    whatsapp: string;
    about: string;
    cityId: string;
    stateId: string;
  }) => {
    setFeedback(null);

    try {
      if (profile) {
        const result = await updateProfile(values);
        setFeedback({
          type: "success",
          message: result.message ?? "Perfil atualizado com sucesso.",
        });
        return;
      }

      const result = await createProfile(values);
      setFeedback({
        type: "success",
        message:
          result.detail ?? result.message ?? "Perfil criado com sucesso.",
      });
    } catch (submitError) {
      setFeedback({
        type: "error",
        message: getApiErrorMessage(
          submitError,
          "Não foi possível salvar o perfil de guardião.",
        ),
      });
    }
  };

  return (
    <PageWrapper gap={10}>
      <MainNav />

      <VStack align="stretch" gap={2}>
        <Heading size="2xl">Perfil de guardião</Heading>
        <Text color="gray.600">{helperText}</Text>
      </VStack>

      {isLoading ? (
        <LoadingState message="Carregando perfil de guardião..." />
      ) : error ? (
        <ErrorState
          message={getApiErrorMessage(
            error,
            "Não foi possível carregar seu perfil de guardião.",
          )}
        />
      ) : (
        <VStack align="stretch" gap={6}>
          <Box
            borderRadius="2xl"
            borderWidth="1px"
            borderColor="gray.200"
            bg="white"
            p={{ base: 5, md: 6 }}
          >
            <VStack align="stretch" gap={4}>
              <HStack justify="space-between" align="start" wrap="wrap">
                <VStack align="stretch" gap={1}>
                  <Heading size="md">
                    {profile ? "Dados do perfil" : "Criar perfil de guardião"}
                  </Heading>
                  <Text fontSize="sm" color="gray.500">
                    O perfil é usado para identificar publicações de adoção e
                    liberar acesso à área de gestão.
                  </Text>
                </VStack>

                {profile ? (
                  <Badge
                    colorPalette={approvalPalette(profile.approvalStatus)}
                    borderRadius="full"
                  >
                    {approvalLabels[profile.approvalStatus] ??
                      profile.approvalStatus}
                  </Badge>
                ) : null}
              </HStack>

              {feedback ? (
                <Box
                  borderRadius="xl"
                  borderWidth="1px"
                  borderColor={
                    feedback.type === "success" ? "green.200" : "red.200"
                  }
                  bg={feedback.type === "success" ? "green.50" : "red.50"}
                  p={4}
                >
                  <Text
                    fontSize="sm"
                    color={
                      feedback.type === "success" ? "green.700" : "red.700"
                    }
                  >
                    {feedback.message}
                  </Text>
                </Box>
              ) : null}

              <GuardianProfileForm
                initialValues={profile ?? undefined}
                submitLabel={profile ? "Atualizar perfil" : "Criar perfil"}
                isSubmitting={isCreating || isUpdating}
                onSubmit={handleSubmit}
              />
            </VStack>
          </Box>

          {profile?.approvalStatus === "approved" ? (
            <HStack wrap="wrap" gap={3}>
              <Button asChild borderRadius="full">
                <Link href="/adoption/listings/new">Criar anúncio</Link>
              </Button>
              <Button asChild borderRadius="full" variant="outline">
                <Link href="/adoption/listings/me">Gerenciar anúncios</Link>
              </Button>
            </HStack>
          ) : null}
        </VStack>
      )}
    </PageWrapper>
  );
}
