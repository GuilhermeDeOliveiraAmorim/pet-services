"use client";

import { useMemo } from "react";
import Link from "next/link";
import { useParams, useRouter, useSearchParams } from "next/navigation";
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

import { useProviderGet, useServiceGet } from "@/application";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
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

export default function ServiceDetailsPage() {
  const router = useRouter();
  const params = useParams<{ serviceId: string }>();
  const searchParams = useSearchParams();
  const serviceId = params?.serviceId;
  const returnToParam = searchParams.get("from");
  const returnTo =
    returnToParam &&
    returnToParam.startsWith("/") &&
    !returnToParam.startsWith("//")
      ? returnToParam
      : null;

  const {
    data: serviceData,
    isLoading: isLoadingService,
    error: serviceError,
  } = useServiceGet(serviceId);

  const providerId = serviceData?.service?.providerId;

  const {
    data: providerData,
    isLoading: isLoadingProvider,
    error: providerError,
  } = useProviderGet(providerId, {
    enabled: Boolean(providerId),
  });

  const service = serviceData?.service;
  const provider = providerData?.provider;

  const servicePhotos = useMemo(
    () => normalizePhotos(service?.photos),
    [service?.photos],
  );

  const providerPhotos = useMemo(
    () => normalizePhotos(provider?.photos),
    [provider?.photos],
  );

  const serviceErrorMessage = serviceError
    ? getApiErrorMessage(serviceError, "Não foi possível carregar o serviço.")
    : "";

  const providerErrorMessage = providerError
    ? getApiErrorMessage(providerError, "Não foi possível carregar o provider.")
    : "";

  return (
    <PageWrapper gap={8}>
      <MainNav />

      <VStack align="stretch" gap={6}>
        <ChakraLink
          as={Link}
          href={returnTo ?? "/"}
          fontSize="sm"
          color="teal.600"
          onClick={(event) => {
            event.preventDefault();

            if (returnTo) {
              router.push(returnTo);
              return;
            }

            router.back();
          }}
        >
          ← Voltar para início
        </ChakraLink>

        {isLoadingService ? (
          <Flex
            borderRadius="2xl"
            borderWidth="1px"
            borderColor="gray.200"
            bg="white"
            py={20}
            justify="center"
            align="center"
            gap={3}
          >
            <Spinner color="teal.500" size="sm" />
            <Text color="gray.600">Carregando detalhes do serviço...</Text>
          </Flex>
        ) : null}

        {!isLoadingService && serviceErrorMessage ? (
          <Box
            borderRadius="2xl"
            borderWidth="1px"
            borderColor="red.200"
            bg="red.50"
            p={5}
          >
            <Text fontSize="sm" color="red.700">
              {serviceErrorMessage}
            </Text>
          </Box>
        ) : null}

        {!isLoadingService && !serviceErrorMessage && service ? (
          <VStack align="stretch" gap={6}>
            <Box
              borderRadius="3xl"
              bg="white"
              borderWidth="1px"
              borderColor="gray.200"
              p={{ base: 5, md: 7 }}
            >
              <Text
                fontSize="xs"
                fontWeight="semibold"
                textTransform="uppercase"
                color="teal.600"
              >
                Serviço
              </Text>
              <Heading as="h1" size="xl" mt={2} color="gray.900">
                {service.name || "Serviço sem nome"}
              </Heading>
              <Text mt={3} color="gray.600" lineHeight="tall">
                {service.description || "Sem descrição informada."}
              </Text>

              <Grid
                mt={6}
                templateColumns={{ base: "1fr", md: "repeat(3, 1fr)" }}
                gap={3}
              >
                <Box
                  borderWidth="1px"
                  borderColor="gray.200"
                  borderRadius="xl"
                  p={4}
                >
                  <Text
                    fontSize="xs"
                    color="gray.500"
                    textTransform="uppercase"
                  >
                    Preço
                  </Text>
                  <Text mt={1} fontWeight="semibold" color="gray.900">
                    {formatMoney(service.price)}
                  </Text>
                </Box>

                <Box
                  borderWidth="1px"
                  borderColor="gray.200"
                  borderRadius="xl"
                  p={4}
                >
                  <Text
                    fontSize="xs"
                    color="gray.500"
                    textTransform="uppercase"
                  >
                    Faixa
                  </Text>
                  <Text mt={1} fontWeight="semibold" color="gray.900">
                    {formatMoney(service.priceMinimum)} -{" "}
                    {formatMoney(service.priceMaximum)}
                  </Text>
                </Box>

                <Box
                  borderWidth="1px"
                  borderColor="gray.200"
                  borderRadius="xl"
                  p={4}
                >
                  <Text
                    fontSize="xs"
                    color="gray.500"
                    textTransform="uppercase"
                  >
                    Duração
                  </Text>
                  <Text mt={1} fontWeight="semibold" color="gray.900">
                    {service.duration > 0
                      ? `${service.duration} min`
                      : "Não informada"}
                  </Text>
                </Box>
              </Grid>
            </Box>

            <Box
              borderRadius="3xl"
              bg="white"
              borderWidth="1px"
              borderColor="gray.200"
              p={{ base: 5, md: 7 }}
            >
              <Heading as="h2" size="md" color="gray.900">
                Fotos do serviço
              </Heading>

              {servicePhotos.length ? (
                <Grid
                  mt={4}
                  templateColumns={{
                    base: "1fr",
                    sm: "repeat(2, 1fr)",
                    lg: "repeat(3, 1fr)",
                  }}
                  gap={3}
                >
                  {servicePhotos.map((photo) => (
                    <Box
                      key={photo.id}
                      borderRadius="xl"
                      overflow="hidden"
                      borderWidth="1px"
                      borderColor="gray.200"
                      bg="gray.100"
                    >
                      <chakra.img
                        src={photo.url}
                        alt={`Foto do serviço ${service.name || ""}`}
                        w="full"
                        h="220px"
                        objectFit="cover"
                      />
                    </Box>
                  ))}
                </Grid>
              ) : (
                <Text mt={4} fontSize="sm" color="gray.500">
                  Este serviço ainda não possui fotos.
                </Text>
              )}
            </Box>

            <Grid templateColumns={{ base: "1fr", lg: "1fr 1fr" }} gap={6}>
              <Box
                borderRadius="3xl"
                bg="white"
                borderWidth="1px"
                borderColor="gray.200"
                p={{ base: 5, md: 7 }}
              >
                <Heading as="h2" size="md" color="gray.900">
                  Categorias
                </Heading>

                {service.categories.length ? (
                  <HStack mt={4} wrap="wrap" gap={2}>
                    {service.categories.map((category) => (
                      <Badge
                        key={category.id}
                        borderRadius="full"
                        px={3}
                        py={1}
                        colorPalette="cyan"
                      >
                        {category.name}
                      </Badge>
                    ))}
                  </HStack>
                ) : (
                  <Text mt={4} fontSize="sm" color="gray.500">
                    Nenhuma categoria vinculada.
                  </Text>
                )}
              </Box>

              <Box
                borderRadius="3xl"
                bg="white"
                borderWidth="1px"
                borderColor="gray.200"
                p={{ base: 5, md: 7 }}
              >
                <Heading as="h2" size="md" color="gray.900">
                  Tags
                </Heading>

                {service.tags.length ? (
                  <HStack mt={4} wrap="wrap" gap={2}>
                    {service.tags.map((tag) => (
                      <Badge
                        key={tag.id}
                        borderRadius="full"
                        px={3}
                        py={1}
                        colorPalette="purple"
                      >
                        #{tag.name}
                      </Badge>
                    ))}
                  </HStack>
                ) : (
                  <Text mt={4} fontSize="sm" color="gray.500">
                    Nenhuma tag vinculada.
                  </Text>
                )}
              </Box>
            </Grid>

            <Box
              borderRadius="3xl"
              bg="white"
              borderWidth="1px"
              borderColor="gray.200"
              p={{ base: 5, md: 7 }}
            >
              <Flex align="center" justify="space-between" wrap="wrap" gap={3}>
                <Heading as="h2" size="md" color="gray.900">
                  Provider
                </Heading>
                {isLoadingProvider ? (
                  <HStack gap={2} color="gray.500">
                    <Spinner size="xs" />
                    <Text fontSize="sm">Carregando provider...</Text>
                  </HStack>
                ) : null}
              </Flex>

              {providerErrorMessage ? (
                <Box
                  mt={4}
                  borderRadius="xl"
                  borderWidth="1px"
                  borderColor="orange.200"
                  bg="orange.50"
                  p={4}
                >
                  <Text fontSize="sm" color="orange.700">
                    {providerErrorMessage}
                  </Text>
                </Box>
              ) : null}

              {!isLoadingProvider && provider ? (
                <VStack align="stretch" gap={4} mt={4}>
                  <Box>
                    <Text fontSize="lg" fontWeight="semibold" color="gray.900">
                      {provider.businessName || "Provider sem nome"}
                    </Text>
                    <Text mt={1} color="gray.600">
                      {provider.description || "Sem descrição informada."}
                    </Text>
                  </Box>

                  <Grid
                    templateColumns={{ base: "1fr", md: "repeat(2, 1fr)" }}
                    gap={3}
                  >
                    <Box>
                      <Text
                        fontSize="xs"
                        textTransform="uppercase"
                        color="gray.500"
                      >
                        Faixa de preço
                      </Text>
                      <Text mt={1} color="gray.800" fontWeight="medium">
                        {provider.priceRange || "Não informada"}
                      </Text>
                    </Box>
                    <Box>
                      <Text
                        fontSize="xs"
                        textTransform="uppercase"
                        color="gray.500"
                      >
                        Endereço
                      </Text>
                      <Text mt={1} color="gray.800" fontWeight="medium">
                        {provider.address.street
                          ? `${provider.address.street}, ${provider.address.number}`
                          : "Não informado"}
                      </Text>
                      {provider.address.city || provider.address.state ? (
                        <Text fontSize="sm" color="gray.600">
                          {provider.address.city}
                          {provider.address.city && provider.address.state
                            ? " - "
                            : ""}
                          {provider.address.state}
                        </Text>
                      ) : null}
                    </Box>
                  </Grid>

                  {providerPhotos.length ? (
                    <Box>
                      <Text fontSize="sm" color="gray.600" mb={2}>
                        Fotos do provider
                      </Text>
                      <HStack gap={3} overflowX="auto" pb={1}>
                        {providerPhotos.map((photo) => (
                          <Box
                            key={photo.id}
                            minW="180px"
                            w="180px"
                            h="120px"
                            borderRadius="lg"
                            overflow="hidden"
                            borderWidth="1px"
                            borderColor="gray.200"
                            bg="gray.100"
                            flex="0 0 auto"
                          >
                            <chakra.img
                              src={photo.url}
                              alt={`Foto do provider ${provider.businessName || ""}`}
                              w="full"
                              h="full"
                              objectFit="cover"
                            />
                          </Box>
                        ))}
                      </HStack>
                    </Box>
                  ) : null}
                </VStack>
              ) : null}
            </Box>
          </VStack>
        ) : null}
      </VStack>
    </PageWrapper>
  );
}
