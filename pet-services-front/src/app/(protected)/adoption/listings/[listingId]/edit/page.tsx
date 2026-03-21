"use client";

import Link from "next/link";
import { use, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { Box, Button, Heading, Text, VStack } from "@chakra-ui/react";

import {
  useAdoptionListingUpdate,
  useGuardianStatus,
  useMyAdoptionListings,
} from "@/application";
import AdoptionListingForm from "@/components/adoption/AdoptionListingForm";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import { ErrorState, LoadingState } from "@/components/ui";
import { getApiErrorMessage } from "@/lib/api-error";

type EditAdoptionListingPageProps = {
  params: Promise<{ listingId: string }>;
};

export default function EditAdoptionListingPage({
  params,
}: EditAdoptionListingPageProps) {
  const router = useRouter();
  const [feedback, setFeedback] = useState("");
  const resolvedParams = use(params);
  const listingId = resolvedParams.listingId;
  const { isApprovedGuardian, isLoading: isGuardianLoading } =
    useGuardianStatus({
      enabled: true,
    });
  const { data, isLoading, error } = useMyAdoptionListings(
    { page: 1, pageSize: 100 },
    { enabled: isApprovedGuardian },
  );
  const { mutateAsync: updateListing, isPending } = useAdoptionListingUpdate();

  const listing = useMemo(
    () => data?.listings.find((item) => item.id === listingId),
    [data?.listings, listingId],
  );

  const canEdit = listing?.status === "draft" || listing?.status === "paused";

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
      await updateListing({
        listingId,
        title: values.title,
        description: values.description,
        adoptionReason: values.adoptionReason,
        sex: values.sex,
        size: values.size,
        ageGroup: values.ageGroup,
        cityId: values.cityId,
        stateId: values.stateId,
        latitude: values.latitude,
        longitude: values.longitude,
        vaccinated: values.vaccinated,
        neutered: values.neutered,
        dewormed: values.dewormed,
        specialNeeds: values.specialNeeds,
        goodWithChildren: values.goodWithChildren,
        goodWithDogs: values.goodWithDogs,
        goodWithCats: values.goodWithCats,
        requiresHouseScreening: values.requiresHouseScreening,
      });
      router.push("/adoption/listings/me");
    } catch (submitError) {
      setFeedback(
        getApiErrorMessage(
          submitError,
          "Não foi possível atualizar o anúncio de adoção.",
        ),
      );
    }
  };

  return (
    <PageWrapper gap={10}>
      <MainNav />

      <VStack align="stretch" gap={2}>
        <Heading size="2xl">Editar anúncio</Heading>
        <Text color="gray.600">
          Ajuste os dados do anúncio antes de publicar ou enquanto ele estiver
          pausado.
        </Text>
      </VStack>

      {isGuardianLoading || isLoading ? (
        <LoadingState message="Carregando anúncio..." />
      ) : !isApprovedGuardian ? (
        <ErrorState message="Seu perfil de guardião precisa estar aprovado para editar anúncios.">
          <Button asChild mt={4} borderRadius="full">
            <Link href="/adoption/guardian">Ir para perfil de guardião</Link>
          </Button>
        </ErrorState>
      ) : error ? (
        <ErrorState
          message={getApiErrorMessage(
            error,
            "Não foi possível carregar seus anúncios de adoção.",
          )}
        />
      ) : !listing ? (
        <ErrorState message="Anúncio não encontrado na sua área de gestão." />
      ) : !canEdit ? (
        <ErrorState message="Somente anúncios em rascunho ou pausados podem ser editados." />
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
              listing={listing}
              submitLabel="Salvar alterações"
              isSubmitting={isPending}
              onSubmit={handleSubmit}
            />
          </VStack>
        </Box>
      )}
    </PageWrapper>
  );
}
