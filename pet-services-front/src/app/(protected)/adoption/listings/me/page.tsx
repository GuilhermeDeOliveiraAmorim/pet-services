"use client";

import { Heading, Text, VStack } from "@chakra-ui/react";

import { useGuardianStatus } from "@/application";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import { ErrorState, LoadingState } from "@/components/ui";

export default function MyAdoptionListingsPage() {
  const { isLoading, isApprovedGuardian } = useGuardianStatus({
    enabled: true,
  });

  return (
    <PageWrapper gap={10}>
      <MainNav />

      <VStack align="stretch" gap={2}>
        <Heading size="2xl">Meus anúncios de adoção</Heading>
        <Text color="gray.600">
          Área exclusiva para guardiões com perfil de adoção aprovado.
        </Text>
      </VStack>

      {isLoading ? (
        <LoadingState message="Verificando seu perfil de guardião..." />
      ) : !isApprovedGuardian ? (
        <ErrorState message="Seu perfil de guardião ainda não está aprovado para acessar esta área." />
      ) : (
        <ErrorState message="Painel de anúncios em implementação. Próximo passo: listar anúncios e candidaturas recebidas." />
      )}
    </PageWrapper>
  );
}
