"use client";

import Link from "next/link";
import { useMemo, useState } from "react";
import {
  Badge,
  Box,
  Button,
  Flex,
  Grid,
  Heading,
  HStack,
  NativeSelect,
  SimpleGrid,
  Skeleton,
  Stack,
  Text,
  VStack,
} from "@chakra-ui/react";

import { usePublicAdoptionListings } from "@/application";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

const sexOptions = [
  { value: "", label: "Todos os sexos" },
  { value: "male", label: "Macho" },
  { value: "female", label: "Fêmea" },
];

const sizeOptions = [
  { value: "", label: "Todos os portes" },
  { value: "small", label: "Pequeno" },
  { value: "medium", label: "Médio" },
  { value: "large", label: "Grande" },
];

const ageGroupOptions = [
  { value: "", label: "Todas as idades" },
  { value: "puppy", label: "Filhote" },
  { value: "adult", label: "Adulto" },
  { value: "senior", label: "Sênior" },
];

const sexLabels: Record<string, string> = {
  male: "Macho",
  female: "Fêmea",
};

const sizeLabels: Record<string, string> = {
  small: "Pequeno",
  medium: "Médio",
  large: "Grande",
};

const ageGroupLabels: Record<string, string> = {
  puppy: "Filhote",
  adult: "Adulto",
  senior: "Sênior",
};

const booleanHighlights = (listing: {
  goodWithChildren: boolean;
  goodWithDogs: boolean;
  goodWithCats: boolean;
  neutered: boolean;
  dewormed: boolean;
}) => {
  const labels: string[] = [];

  if (listing.goodWithChildren) labels.push("Bom com crianças");
  if (listing.goodWithDogs) labels.push("Bom com cães");
  if (listing.goodWithCats) labels.push("Bom com gatos");
  if (listing.neutered) labels.push("Castrado");
  if (listing.dewormed) labels.push("Vermifugado");

  return labels;
};

export default function AdoptionCatalogPage() {
  const [sex, setSex] = useState("");
  const [size, setSize] = useState("");
  const [ageGroup, setAgeGroup] = useState("");

  const filters = useMemo(
    () => ({
      page: 1,
      pageSize: 12,
      sex: sex || undefined,
      size: size || undefined,
      ageGroup: ageGroup || undefined,
    }),
    [ageGroup, sex, size],
  );

  const { data, isLoading, isError } = usePublicAdoptionListings(filters);
  const listings = data?.listings ?? [];

  return (
    <PageWrapper gap={12}>
      <MainNav />

      <Box
        borderRadius="3xl"
        px={{ base: 6, md: 10 }}
        py={{ base: 8, md: 12 }}
        bg="linear-gradient(135deg, rgba(255,245,235,1) 0%, rgba(236,253,245,1) 55%, rgba(239,246,255,1) 100%)"
        border="1px solid"
        borderColor="orange.100"
        boxShadow="0 20px 60px rgba(15, 23, 42, 0.08)"
      >
        <Stack gap={6}>
          <VStack align="start" gap={3}>
            <Badge colorPalette="orange" borderRadius="full" px={3} py={1}>
              adoção responsável
            </Badge>
            <Heading size={{ base: "2xl", md: "3xl" }} maxW="3xl">
              Encontre um novo companheiro com informações claras e processo de
              adoção organizado.
            </Heading>
            <Text maxW="2xl" color="gray.600" fontSize="lg">
              Explore animais disponíveis, entenda o perfil de cada pet e avance
              para a candidatura quando fizer sentido para sua família.
            </Text>
          </VStack>

          <SimpleGrid columns={{ base: 1, md: 3 }} gap={4}>
            <Box
              bg="whiteAlpha.900"
              borderRadius="2xl"
              p={4}
              border="1px solid"
              borderColor="whiteAlpha.700"
            >
              <Text fontSize="sm" color="gray.500">
                Anúncios visíveis
              </Text>
              <Text fontSize="3xl" fontWeight="bold">
                {data?.pagination.totalRecords ?? 0}
              </Text>
            </Box>
            <Box
              bg="whiteAlpha.900"
              borderRadius="2xl"
              p={4}
              border="1px solid"
              borderColor="whiteAlpha.700"
            >
              <Text fontSize="sm" color="gray.500">
                Página atual
              </Text>
              <Text fontSize="3xl" fontWeight="bold">
                {data?.pagination.currentPage ?? 1}
              </Text>
            </Box>
            <Box
              bg="whiteAlpha.900"
              borderRadius="2xl"
              p={4}
              border="1px solid"
              borderColor="whiteAlpha.700"
            >
              <Text fontSize="sm" color="gray.500">
                Páginas
              </Text>
              <Text fontSize="3xl" fontWeight="bold">
                {data?.pagination.totalPages ?? 1}
              </Text>
            </Box>
          </SimpleGrid>
        </Stack>
      </Box>

      <Box
        bg="white"
        borderRadius="2xl"
        border="1px solid"
        borderColor="gray.200"
        p={{ base: 5, md: 6 }}
      >
        <Flex
          direction={{ base: "column", md: "row" }}
          gap={4}
          align={{ md: "end" }}
        >
          <Box flex="1">
            <Text fontSize="sm" color="gray.500" mb={2}>
              Sexo
            </Text>
            <NativeSelect.Root>
              <NativeSelect.Field
                value={sex}
                onChange={(event) => setSex(event.target.value)}
              >
                {sexOptions.map((option) => (
                  <option key={option.value || "all-sex"} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </NativeSelect.Field>
              <NativeSelect.Indicator />
            </NativeSelect.Root>
          </Box>
          <Box flex="1">
            <Text fontSize="sm" color="gray.500" mb={2}>
              Porte
            </Text>
            <NativeSelect.Root>
              <NativeSelect.Field
                value={size}
                onChange={(event) => setSize(event.target.value)}
              >
                {sizeOptions.map((option) => (
                  <option key={option.value || "all-size"} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </NativeSelect.Field>
              <NativeSelect.Indicator />
            </NativeSelect.Root>
          </Box>
          <Box flex="1">
            <Text fontSize="sm" color="gray.500" mb={2}>
              Faixa etária
            </Text>
            <NativeSelect.Root>
              <NativeSelect.Field
                value={ageGroup}
                onChange={(event) => setAgeGroup(event.target.value)}
              >
                {ageGroupOptions.map((option) => (
                  <option key={option.value || "all-age"} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </NativeSelect.Field>
              <NativeSelect.Indicator />
            </NativeSelect.Root>
          </Box>
          <Button
            variant="outline"
            borderRadius="full"
            onClick={() => {
              setSex("");
              setSize("");
              setAgeGroup("");
            }}
          >
            Limpar filtros
          </Button>
        </Flex>
      </Box>

      {isLoading ? (
        <Grid
          templateColumns={{
            base: "1fr",
            md: "repeat(2, 1fr)",
            xl: "repeat(3, 1fr)",
          }}
          gap={6}
        >
          {Array.from({ length: 6 }).map((_, index) => (
            <Skeleton key={index} height="320px" borderRadius="2xl" />
          ))}
        </Grid>
      ) : isError ? (
        <Box
          bg="red.50"
          border="1px solid"
          borderColor="red.200"
          borderRadius="2xl"
          p={6}
        >
          <Heading size="md" mb={2}>
            Não foi possível carregar os anúncios.
          </Heading>
          <Text color="gray.600">
            Verifique a API de adoção e tente novamente.
          </Text>
        </Box>
      ) : listings.length === 0 ? (
        <Box
          bg="white"
          border="1px solid"
          borderColor="gray.200"
          borderRadius="2xl"
          p={8}
        >
          <Heading size="md" mb={2}>
            Nenhum pet encontrado com esse recorte.
          </Heading>
          <Text color="gray.600">Ajuste os filtros para ampliar a busca.</Text>
        </Box>
      ) : (
        <Grid
          templateColumns={{
            base: "1fr",
            md: "repeat(2, 1fr)",
            xl: "repeat(3, 1fr)",
          }}
          gap={6}
        >
          {listings.map((listing) => {
            const highlights = booleanHighlights(listing);
            const coverPhoto = listing.pet?.photos?.[0]?.url;

            return (
              <Box
                key={listing.id}
                bg="white"
                border="1px solid"
                borderColor="gray.200"
                borderRadius="2xl"
                overflow="hidden"
                boxShadow="sm"
              >
                <Box
                  h="220px"
                  bg={
                    coverPhoto
                      ? `url(${coverPhoto}) center/cover no-repeat`
                      : "linear-gradient(135deg, #fef3c7 0%, #fde68a 50%, #fca5a5 100%)"
                  }
                />
                <VStack align="stretch" gap={4} p={5}>
                  <VStack align="stretch" gap={2}>
                    <HStack justify="space-between" align="start">
                      <Box>
                        <Heading size="md">{listing.title}</Heading>
                        <Text color="gray.500">
                          {listing.pet?.name || "Pet sem nome"}
                          {listing.pet?.specie?.name
                            ? ` • ${listing.pet.specie.name}`
                            : ""}
                        </Text>
                      </Box>
                      <Badge colorPalette="green" borderRadius="full">
                        {listing.status || "published"}
                      </Badge>
                    </HStack>
                    <Text
                      color="gray.600"
                      css={{
                        display: "-webkit-box",
                        WebkitLineClamp: 3,
                        WebkitBoxOrient: "vertical",
                        overflow: "hidden",
                      }}
                    >
                      {listing.description}
                    </Text>
                  </VStack>

                  <HStack wrap="wrap">
                    {listing.sex ? (
                      <Badge borderRadius="full">
                        {sexLabels[listing.sex] ?? listing.sex}
                      </Badge>
                    ) : null}
                    {listing.size ? (
                      <Badge borderRadius="full">
                        {sizeLabels[listing.size] ?? listing.size}
                      </Badge>
                    ) : null}
                    {listing.ageGroup ? (
                      <Badge borderRadius="full">
                        {ageGroupLabels[listing.ageGroup] ?? listing.ageGroup}
                      </Badge>
                    ) : null}
                  </HStack>

                  <HStack wrap="wrap">
                    {highlights.slice(0, 3).map((label) => (
                      <Badge
                        key={label}
                        colorPalette="orange"
                        borderRadius="full"
                      >
                        {label}
                      </Badge>
                    ))}
                  </HStack>

                  <Text fontSize="sm" color="gray.500">
                    {listing.requiresHouseScreening
                      ? "Processo inclui triagem residencial."
                      : "Processo simplificado de candidatura."}
                  </Text>

                  <Button
                    asChild
                    borderRadius="full"
                    bg="gray.900"
                    color="white"
                    _hover={{ bg: "gray.700" }}
                  >
                    <Link href={`/adoption/${listing.id}`}>
                      Ver perfil do pet
                    </Link>
                  </Button>
                </VStack>
              </Box>
            );
          })}
        </Grid>
      )}
    </PageWrapper>
  );
}
