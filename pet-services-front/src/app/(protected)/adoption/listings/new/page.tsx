"use client";

import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { Box, Button, Heading, Text, VStack } from "@chakra-ui/react";

import { useAdoptionListingCreate, useGuardianStatus } from "@/application";
import AdoptionListingForm from "@/components/adoption/AdoptionListingForm";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import { ErrorState, LoadingState } from "@/components/ui";
import { getApiErrorMessage } from "@/lib/api-error";

export default function NewAdoptionListingPage() {
  const router = useRouter();
  const [feedback, setFeedback] = useState<string>("");
  const { isLoading: isGuardianLoading, isApprovedGuardian } =
    useGuardianStatus({
      enabled: true,
    });
  const { mutateAsync: createListing, isPending } = useAdoptionListingCreate();

  const handleSubmit = async (values: {
    petId: string;
    title: string;
    description: string;
    adoptionReason: string;
    sex: string;
    size: string;
    ageGroup: string;
    cityId: string;
    stateId: string;
    latitude?: number;
    longitude?: number;
    vaccinated: boolean;
    neutered: boolean;
    dewormed: boolean;
    specialNeeds: boolean;
    goodWithChildren: boolean;
    goodWithDogs: boolean;
    goodWithCats: boolean;
    requiresHouseScreening: boolean;
  }) => {
    setFeedback("");

    try {
      await createListing(values);
      router.push("/adoption/listings/me");
    } catch (error) {
      setFeedback(
        getApiErrorMessage(
          error,
          "Não foi possível criar o anúncio de adoção.",
        ),
      );
    }
  };

  return (
    <PageWrapper gap={10}>
      <MainNav />

      <VStack align="stretch" gap={2}>
        <Heading size="2xl">Novo anúncio de adoção</Heading>
        <Text color="gray.600">
          Cadastre um anúncio em rascunho e publique quando todas as informações
          estiverem revisadas.
        </Text>
      </VStack>

      {isGuardianLoading ? (
        <LoadingState message="Verificando seu perfil de guardião..." />
      ) : !isApprovedGuardian ? (
        <ErrorState message="Seu perfil de guardião precisa estar aprovado para criar anúncios.">
          <Button asChild mt={4} borderRadius="full">
            <Link href="/adoption/guardian">Ir para perfil de guardião</Link>
          </Button>
        </ErrorState>
      ) : (
        <Box
          borderRadius="2xl"
          borderWidth="1px"
          borderColor="gray.200"
          bg="white"
          p={{ base: 5, md: 6 }}
        >
          <VStack align="stretch" gap={4}>
            {feedback ? (
              <Box
                borderRadius="xl"
                borderWidth="1px"
                borderColor="red.200"
                bg="red.50"
                p={4}
              >
                <Text fontSize="sm" color="red.700">
                  {feedback}
                </Text>
              </Box>
            ) : null}

            <AdoptionListingForm
              submitLabel="Criar anúncio"
              isSubmitting={isPending}
              onSubmit={handleSubmit}
            />
          </VStack>
        </Box>
      )}
    </PageWrapper>
  );
}
