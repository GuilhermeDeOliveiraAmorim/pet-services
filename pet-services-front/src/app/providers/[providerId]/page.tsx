"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import {
  Badge,
  Box,
  chakra,
  Flex,
  Grid,
  Heading,
  HStack,
  Link as ChakraLink,
  Spinner,
  Text,
  VStack,
} from "@chakra-ui/react";

import { useProviderGet, useReviewList, useServiceList } from "@/application";
import { MainNav, PageWrapper, ProviderRating } from "@/components/common";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui";
import { getApiErrorMessage } from "@/lib/api-error";

type PhotoView = {
  id: string;
  url: string;
};

const normalizePhotos = (items?: unknown[]): PhotoView[] => {
  if (!Array.isArray(items)) {
    return [];
  }

  return items
    .map((item, index) => {
      if (!item || typeof item !== "object") {
        return null;
      }

      const maybePhoto = item as { id?: unknown; url?: unknown };
      const url =
        typeof maybePhoto.url === "string" ? maybePhoto.url.trim() : "";

      if (!url) {
        return null;
      }

      return {
        id:
          typeof maybePhoto.id === "string" && maybePhoto.id
            ? maybePhoto.id
            : `${index}-${url}`,
        url,
      };
    })
    .filter((photo): photo is PhotoView => Boolean(photo));
};

const formatMoney = (value: number): string => {
  if (!Number.isFinite(value)) {
    return "R$ 0,00";
  }

  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
    minimumFractionDigits: 2,
  }).format(value);
};

const formatRating = (value: number): string => {
  if (!Number.isFinite(value)) {
    return "0.0";
  }

  return value.toFixed(1);
};

export default function ProviderDetailsPage() {
  const params = useParams<{ providerId: string }>();
  const providerId = params?.providerId;

  const {
    data: providerData,
    isLoading: isLoadingProvider,
    error: providerError,
  } = useProviderGet(providerId, {
    enabled: Boolean(providerId),
  });

  const {
    data: servicesData,
    isLoading: isLoadingServices,
    error: servicesError,
  } = useServiceList({
    input: providerId ? { providerId } : undefined,
    enabled: Boolean(providerId),
  });

  const {
    data: reviewsData,
    isLoading: isLoadingReviews,
    error: reviewsError,
  } = useReviewList(
    providerId
      ? {
          providerId,
          page: 1,
          pageSize: 10,
        }
      : undefined,
    {
      enabled: Boolean(providerId),
    },
  );

  const provider = providerData?.provider;
  const services = servicesData?.services ?? [];
  const reviews = reviewsData?.reviews ?? [];
  const providerPhotos = normalizePhotos(provider?.photos);

  const providerErrorMessage = providerError
    ? getApiErrorMessage(providerError, "Não foi possível carregar o provider.")
    : "";

  const servicesErrorMessage = servicesError
    ? getApiErrorMessage(
        servicesError,
        "Não foi possível carregar os serviços publicados.",
      )
    : "";

  const reviewsErrorMessage = reviewsError
    ? getApiErrorMessage(
        reviewsError,
        "Não foi possível carregar as avaliações deste provider.",
      )
    : "";

  const hasAddress =
    Boolean(provider?.address.street) ||
    Boolean(provider?.address.city) ||
    Boolean(provider?.address.state);

  return (
    <PageWrapper gap={8}>
      <MainNav />

      <VStack align="stretch" gap={{ base: 4, md: 6 }}>
        <ChakraLink
          as={Link}
          href="/"
          fontSize={{ base: "xs", sm: "sm" }}
          color="teal.600"
        >
          ← Voltar para início
        </ChakraLink>

        {isLoadingProvider ? (
          <LoadingState message="Carregando dados do provider..." />
        ) : null}

        {!isLoadingProvider && providerErrorMessage ? (
          <ErrorState message={providerErrorMessage} />
        ) : null}

        {!isLoadingProvider && !providerErrorMessage && provider ? (
          <VStack align="stretch" gap={{ base: 4, md: 6 }}>
            <Box
              borderRadius={{ base: "2xl", md: "3xl" }}
              bg="white"
              borderWidth="1px"
              borderColor="gray.200"
              p={{ base: 4, sm: 5, md: 7 }}
            >
              <Text
                fontSize={{ base: "xs" }}
                fontWeight="semibold"
                textTransform="uppercase"
                color="teal.600"
              >
                Provider
              </Text>
              <Heading
                as="h1"
                size={{ base: "lg", md: "xl" }}
                mt={2}
                color="gray.900"
              >
                {provider.businessName || "Provider sem nome"}
              </Heading>
              <Text
                mt={3}
                fontSize={{ base: "xs", sm: "sm" }}
                color="gray.600"
                lineHeight="tall"
              >
                {provider.description || "Sem descrição informada."}
              </Text>

              <HStack mt={4} gap={2} wrap="wrap">
                {provider.priceRange ? (
                  <Badge borderRadius="full" px={3} py={1} colorPalette="teal">
                    Faixa: {provider.priceRange}
                  </Badge>
                ) : null}
                <ProviderRating
                  rating={provider.averageRating}
                  totalReviews={reviewsData?.total}
                  showCount
                />
              </HStack>

              <Grid
                mt={6}
                templateColumns={{ base: "1fr", md: "repeat(2, 1fr)" }}
                gap={{ base: 2, md: 3 }}
              >
                <Box
                  borderWidth="1px"
                  borderColor="gray.200"
                  borderRadius={{ base: "lg", md: "xl" }}
                  p={{ base: 3, md: 4 }}
                >
                  <Text
                    fontSize={{ base: "xs" }}
                    color="gray.500"
                    textTransform="uppercase"
                  >
                    Serviços publicados
                  </Text>
                  <Text
                    mt={1}
                    fontWeight="semibold"
                    color="gray.900"
                    fontSize={{ base: "xs", sm: "sm" }}
                  >
                    {services.length}
                  </Text>
                </Box>

                <Box
                  borderWidth="1px"
                  borderColor="gray.200"
                  borderRadius={{ base: "lg", md: "xl" }}
                  p={{ base: 3, md: 4 }}
                >
                  <Text
                    fontSize={{ base: "xs" }}
                    color="gray.500"
                    textTransform="uppercase"
                  >
                    Endereço
                  </Text>
                  <Text
                    mt={1}
                    fontWeight="semibold"
                    color="gray.900"
                    fontSize={{ base: "xs", sm: "sm" }}
                  >
                    {hasAddress
                      ? `${provider.address.street || ""}${provider.address.street ? ", " : ""}${provider.address.number || ""}`.trim()
                      : "Não informado"}
                  </Text>
                  {provider.address.city || provider.address.state ? (
                    <Text mt={1} fontSize={{ base: "xs" }} color="gray.600">
                      {provider.address.city}
                      {provider.address.city && provider.address.state
                        ? " - "
                        : ""}
                      {provider.address.state}
                    </Text>
                  ) : null}
                </Box>
              </Grid>
            </Box>

            <Box
              borderRadius={{ base: "2xl", md: "3xl" }}
              bg="white"
              borderWidth="1px"
              borderColor="gray.200"
              p={{ base: 4, sm: 5, md: 7 }}
            >
              <Heading
                as="h2"
                size={{ base: "sm", md: "md" }}
                color="gray.900"
                fontSize={{ base: "sm", md: "md" }}
              >
                Fotos do provider
              </Heading>

              {providerPhotos.length ? (
                <Grid
                  mt={4}
                  templateColumns={{
                    base: "1fr",
                    md: "repeat(2, 1fr)",
                    lg: "repeat(3, 1fr)",
                  }}
                  gap={{ base: 2, md: 3 }}
                >
                  {providerPhotos.map((photo) => (
                    <Box
                      key={photo.id}
                      borderRadius={{ base: "lg", md: "xl" }}
                      overflow="hidden"
                      borderWidth="1px"
                      borderColor="gray.200"
                      bg="gray.100"
                    >
                      <chakra.img
                        src={photo.url}
                        alt={`Foto do provider ${provider.businessName || ""}`}
                        w="full"
                        h={{ base: "170px", md: "210px", lg: "220px" }}
                        objectFit="cover"
                      />
                    </Box>
                  ))}
                </Grid>
              ) : (
                <Text
                  mt={4}
                  fontSize={{ base: "xs", sm: "sm" }}
                  color="gray.500"
                >
                  Este provider ainda não possui fotos.
                </Text>
              )}
            </Box>

            <Box
              borderRadius={{ base: "2xl", md: "3xl" }}
              bg="white"
              borderWidth="1px"
              borderColor="gray.200"
              p={{ base: 4, sm: 5, md: 7 }}
            >
              <Flex align="center" justify="space-between" wrap="wrap" gap={3}>
                <Heading
                  as="h2"
                  size={{ base: "sm", md: "md" }}
                  color="gray.900"
                  fontSize={{ base: "sm", md: "md" }}
                >
                  Avaliações
                </Heading>
                {isLoadingReviews ? (
                  <HStack gap={2} color="gray.500">
                    <Spinner size="xs" />
                    <Text fontSize={{ base: "xs", sm: "sm" }}>
                      Carregando...
                    </Text>
                  </HStack>
                ) : null}
              </Flex>

              {reviewsErrorMessage ? (
                <Box
                  mt={4}
                  borderRadius={{ base: "lg", md: "xl" }}
                  borderWidth="1px"
                  borderColor="orange.200"
                  bg="orange.50"
                  p={{ base: 3, md: 4 }}
                >
                  <Text fontSize={{ base: "xs", sm: "sm" }} color="orange.700">
                    {reviewsErrorMessage}
                  </Text>
                </Box>
              ) : null}

              {!isLoadingReviews && !reviewsErrorMessage ? (
                reviews.length ? (
                  <VStack mt={4} align="stretch" gap={3}>
                    {reviews.map((review) => (
                      <Box
                        key={review.id}
                        borderWidth="1px"
                        borderColor="gray.200"
                        borderRadius={{ base: "xl", md: "2xl" }}
                        p={{ base: 3, md: 4 }}
                      >
                        <Flex
                          justify="space-between"
                          align={{ base: "flex-start", sm: "center" }}
                          direction={{ base: "column", sm: "row" }}
                          gap={2}
                        >
                          <Text
                            fontSize={{ base: "xs" }}
                            color="gray.500"
                            textTransform="uppercase"
                          >
                            Avaliação
                          </Text>
                          <Badge
                            borderRadius="full"
                            px={3}
                            py={1}
                            colorPalette="teal"
                            variant="subtle"
                            fontSize={{ base: "xs" }}
                          >
                            Nota {formatRating(review.rating)}
                          </Badge>
                        </Flex>

                        <Text
                          mt={2}
                          fontSize={{ base: "xs", sm: "sm" }}
                          color="gray.700"
                        >
                          {review.comment || "Sem comentário."}
                        </Text>

                        <Text mt={2} fontSize={{ base: "xs" }} color="gray.500">
                          Publicado em{" "}
                          {new Date(review.createdAt).toLocaleDateString(
                            "pt-BR",
                          )}
                        </Text>
                      </Box>
                    ))}
                  </VStack>
                ) : (
                  <EmptyState message="Este provider ainda não possui avaliações." />
                )
              ) : null}
            </Box>

            <Box
              borderRadius={{ base: "2xl", md: "3xl" }}
              bg="white"
              borderWidth="1px"
              borderColor="gray.200"
              p={{ base: 4, sm: 5, md: 7 }}
            >
              <Flex align="center" justify="space-between" wrap="wrap" gap={3}>
                <Heading
                  as="h2"
                  size={{ base: "sm", md: "md" }}
                  color="gray.900"
                  fontSize={{ base: "sm", md: "md" }}
                >
                  Serviços publicados
                </Heading>
                {isLoadingServices ? (
                  <HStack gap={2} color="gray.500">
                    <Spinner size="xs" />
                    <Text fontSize={{ base: "xs", sm: "sm" }}>
                      Carregando...
                    </Text>
                  </HStack>
                ) : null}
              </Flex>

              {servicesErrorMessage ? (
                <Box
                  mt={4}
                  borderRadius={{ base: "lg", md: "xl" }}
                  borderWidth="1px"
                  borderColor="orange.200"
                  bg="orange.50"
                  p={{ base: 3, md: 4 }}
                >
                  <Text fontSize={{ base: "xs", sm: "sm" }} color="orange.700">
                    {servicesErrorMessage}
                  </Text>
                </Box>
              ) : null}

              {!isLoadingServices && !servicesErrorMessage ? (
                services.length ? (
                  <VStack mt={4} align="stretch" gap={3}>
                    {services.map((service) => (
                      <Box
                        key={service.id}
                        borderWidth="1px"
                        borderColor="gray.200"
                        borderRadius={{ base: "xl", md: "2xl" }}
                        p={{ base: 3, md: 4 }}
                      >
                        <Flex
                          direction={{ base: "column", md: "row" }}
                          justify="space-between"
                          align={{ base: "start", md: "center" }}
                          gap={3}
                        >
                          <Box>
                            <Text
                              fontSize={{ base: "sm", md: "md" }}
                              fontWeight="semibold"
                              color="gray.900"
                            >
                              {service.name}
                            </Text>
                            <Text
                              mt={1}
                              fontSize={{ base: "xs", sm: "sm" }}
                              color="gray.600"
                            >
                              {service.description ||
                                "Sem descrição informada."}
                            </Text>
                            <Text
                              mt={2}
                              fontSize={{ base: "xs" }}
                              color="gray.500"
                            >
                              {formatMoney(service.price)} · {service.duration}{" "}
                              min
                            </Text>
                          </Box>

                          <ChakraLink
                            as={Link}
                            href={`/services/${service.id}?from=/providers/${provider.id}`}
                            borderRadius="full"
                            borderWidth="1px"
                            borderColor="gray.300"
                            bg="white"
                            color="gray.700"
                            _hover={{ bg: "gray.50", textDecoration: "none" }}
                            h={{ base: "8", sm: "9" }}
                            px={{ base: 3, sm: 4 }}
                            fontSize={{ base: "xs", sm: "sm" }}
                          >
                            Ver detalhes
                          </ChakraLink>
                        </Flex>
                      </Box>
                    ))}
                  </VStack>
                ) : (
                  <EmptyState message="Este provider ainda não possui serviços publicados." />
                )
              ) : null}
            </Box>
          </VStack>
        ) : null}
      </VStack>
    </PageWrapper>
  );
}
