"use client";

import Link from "next/link";
import { notFound } from "next/navigation";
import {
  Badge,
  Box,
  Button,
  Grid,
  Heading,
  HStack,
  SimpleGrid,
  Skeleton,
  Stack,
  Text,
  VStack,
} from "@chakra-ui/react";

import { usePublicAdoptionListing } from "@/application";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

const labelMap: Record<string, string> = {
  male: "Macho",
  female: "Fêmea",
  small: "Pequeno",
  medium: "Médio",
  large: "Grande",
  puppy: "Filhote",
  adult: "Adulto",
  senior: "Sênior",
};

const yesNo = (value: boolean) => (value ? "Sim" : "Não");

type AdoptionListingDetailPageProps = {
  params: Promise<{ listingId: string }>;
};

export default function AdoptionListingDetailPage({
  params,
}: AdoptionListingDetailPageProps) {
  const resolved = (params as unknown as { listingId: string })?.listingId
    ? (params as unknown as { listingId: string })
    : null;

  const listingId = resolved?.listingId;
  const { data, isLoading, isError } = usePublicAdoptionListing(listingId);

  if (!listingId) {
    notFound();
  }

  const listing = data?.listing;

  return (
    <PageWrapper gap={10}>
      <MainNav />

      {isLoading ? (
        <Stack gap={6}>
          <Skeleton height="360px" borderRadius="3xl" />
          <Skeleton height="220px" borderRadius="3xl" />
        </Stack>
      ) : isError || !listing ? (
        <Box
          bg="white"
          border="1px solid"
          borderColor="gray.200"
          borderRadius="3xl"
          p={8}
        >
          <Heading size="lg" mb={3}>
            Anúncio não encontrado.
          </Heading>
          <Text color="gray.600" mb={5}>
            O perfil solicitado pode ter sido removido ou não está mais
            disponível publicamente.
          </Text>
          <Button asChild borderRadius="full">
            <Link href="/adoption">Voltar para adoção</Link>
          </Button>
        </Box>
      ) : (
        <Stack gap={8}>
          <Grid templateColumns={{ base: "1fr", lg: "1.2fr 0.8fr" }} gap={6}>
            <Box
              minH={{ base: "280px", md: "420px" }}
              borderRadius="3xl"
              bg={
                listing.pet?.photos?.[0]?.url
                  ? `url(${listing.pet.photos[0].url}) center/cover no-repeat`
                  : "linear-gradient(135deg, #fed7aa 0%, #fef3c7 50%, #bbf7d0 100%)"
              }
              border="1px solid"
              borderColor="gray.200"
            />

            <VStack
              align="stretch"
              gap={5}
              bg="white"
              border="1px solid"
              borderColor="gray.200"
              borderRadius="3xl"
              p={{ base: 6, md: 8 }}
            >
              <Badge colorPalette="orange" w="fit-content" borderRadius="full">
                adoção disponível
              </Badge>
              <Box>
                <Heading size="xl" mb={2}>
                  {listing.title}
                </Heading>
                <Text color="gray.500">
                  {listing.pet?.name || "Pet sem nome"}
                  {listing.pet?.specie?.name
                    ? ` • ${listing.pet.specie.name}`
                    : ""}
                  {listing.pet?.breed ? ` • ${listing.pet.breed}` : ""}
                </Text>
              </Box>

              <HStack wrap="wrap">
                {[listing.sex, listing.size, listing.ageGroup]
                  .filter(Boolean)
                  .map((item) => (
                    <Badge key={item} borderRadius="full">
                      {labelMap[item] ?? item}
                    </Badge>
                  ))}
              </HStack>

              <Text color="gray.700">{listing.description}</Text>

              <Box
                bg="orange.50"
                border="1px solid"
                borderColor="orange.100"
                borderRadius="2xl"
                p={4}
              >
                <Text fontWeight="semibold" mb={2}>
                  Motivo da adoção
                </Text>
                <Text color="gray.700">
                  {listing.adoptionReason || "Não informado."}
                </Text>
              </Box>

              <VStack align="stretch" gap={3}>
                <Button
                  asChild
                  borderRadius="full"
                  bg="gray.900"
                  color="white"
                  _hover={{ bg: "gray.700" }}
                >
                  <Link href="/login">Entrar para me candidatar</Link>
                </Button>
                <Button asChild variant="outline" borderRadius="full">
                  <Link href="/adoption">Ver outros pets</Link>
                </Button>
              </VStack>
            </VStack>
          </Grid>

          <SimpleGrid columns={{ base: 1, md: 2 }} gap={6}>
            <Box
              bg="white"
              border="1px solid"
              borderColor="gray.200"
              borderRadius="3xl"
              p={6}
            >
              <Heading size="md" mb={4}>
                Características
              </Heading>
              <VStack align="stretch" gap={3}>
                <HStack justify="space-between">
                  <Text color="gray.500">Bom com crianças</Text>
                  <Text>{yesNo(listing.goodWithChildren)}</Text>
                </HStack>
                <HStack justify="space-between">
                  <Text color="gray.500">Bom com cães</Text>
                  <Text>{yesNo(listing.goodWithDogs)}</Text>
                </HStack>
                <HStack justify="space-between">
                  <Text color="gray.500">Bom com gatos</Text>
                  <Text>{yesNo(listing.goodWithCats)}</Text>
                </HStack>
                <HStack justify="space-between">
                  <Text color="gray.500">Necessidades especiais</Text>
                  <Text>{yesNo(listing.specialNeeds)}</Text>
                </HStack>
                <HStack justify="space-between">
                  <Text color="gray.500">Castrado</Text>
                  <Text>{yesNo(listing.neutered)}</Text>
                </HStack>
                <HStack justify="space-between">
                  <Text color="gray.500">Vermifugado</Text>
                  <Text>{yesNo(listing.dewormed)}</Text>
                </HStack>
                <HStack justify="space-between">
                  <Text color="gray.500">Triagem residencial</Text>
                  <Text>{yesNo(listing.requiresHouseScreening)}</Text>
                </HStack>
              </VStack>
            </Box>

            <Box
              bg="white"
              border="1px solid"
              borderColor="gray.200"
              borderRadius="3xl"
              p={6}
            >
              <Heading size="md" mb={4}>
                Sobre o responsável
              </Heading>
              <VStack align="stretch" gap={3}>
                <HStack justify="space-between">
                  <Text color="gray.500">Nome</Text>
                  <Text>
                    {listing.guardianProfile?.displayName || "Não informado"}
                  </Text>
                </HStack>
                <HStack justify="space-between">
                  <Text color="gray.500">Contato</Text>
                  <Text>
                    {listing.guardianProfile?.whatsapp ||
                      "Disponível após candidatura"}
                  </Text>
                </HStack>
                <Box>
                  <Text color="gray.500" mb={2}>
                    Apresentação
                  </Text>
                  <Text color="gray.700">
                    {listing.guardianProfile?.about ||
                      "Responsável ainda não adicionou uma descrição pública."}
                  </Text>
                </Box>
              </VStack>
            </Box>
          </SimpleGrid>
        </Stack>
      )}
    </PageWrapper>
  );
}
