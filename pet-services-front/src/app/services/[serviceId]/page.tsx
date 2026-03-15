"use client";

import { FormEvent, useMemo, useState } from "react";
import Link from "next/link";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import {
  Badge,
  Box,
  Button,
  chakra,
  Flex,
  Grid,
  Heading,
  HStack,
  NativeSelect,
  Link as ChakraLink,
  Spinner,
  Text,
  Textarea,
  VStack,
} from "@chakra-ui/react";

import {
  useAuthSession,
  usePetListByOwnerId,
  useProviderGet,
  useRequestAdd,
  useServiceGet,
  useUserProfile,
} from "@/application";
import { MainNav, PageWrapper, ProviderRating } from "@/components/common";
import { UserTypes } from "@/domain";
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
  const { isAuthenticated } = useAuthSession();
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
  const requestProviderId = service?.providerId || provider?.id || "";

  const {
    data: profileData,
    isLoading: isLoadingProfile,
    error: profileError,
  } = useUserProfile({
    enabled: isAuthenticated,
  });

  const { mutateAsync: addRequest, isPending: isSubmittingRequest } =
    useRequestAdd();

  const [selectedPetId, setSelectedPetId] = useState("");
  const [requestNotes, setRequestNotes] = useState("");
  const [requestFeedback, setRequestFeedback] = useState<{
    type: "success" | "error";
    message: string;
  } | null>(null);

  const currentUser = profileData?.user;
  const isOwnerUser = currentUser?.userType === UserTypes.Owner;

  const {
    data: ownerPetsData,
    isLoading: isLoadingOwnerPets,
    error: ownerPetsError,
  } = usePetListByOwnerId(isOwnerUser ? currentUser?.id : undefined, {
    enabled: isAuthenticated && isOwnerUser && Boolean(currentUser?.id),
  });

  const ownerPets = ownerPetsData?.pets ?? [];

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

  const profileErrorMessage = profileError
    ? getApiErrorMessage(
        profileError,
        "Não foi possível carregar o perfil para solicitar o serviço.",
      )
    : "";

  const ownerPetsErrorMessage = ownerPetsError
    ? getApiErrorMessage(
        ownerPetsError,
        "Não foi possível carregar seus pets no momento.",
      )
    : "";

  const handleAddRequest = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setRequestFeedback(null);

    if (!service?.id || !requestProviderId) {
      setRequestFeedback({
        type: "error",
        message:
          "Não foi possível identificar os dados do serviço para abrir a solicitação.",
      });
      return;
    }

    if (!selectedPetId.trim()) {
      setRequestFeedback({
        type: "error",
        message: "Selecione um pet para solicitar o serviço.",
      });
      return;
    }

    const normalizedNotes = requestNotes.trim();

    if (normalizedNotes.length > 500) {
      setRequestFeedback({
        type: "error",
        message: "As observações devem ter no máximo 500 caracteres.",
      });
      return;
    }

    try {
      await addRequest({
        providerId: requestProviderId,
        serviceId: service.id,
        petId: selectedPetId,
        notes: normalizedNotes,
      });

      setRequestNotes("");
      setRequestFeedback({
        type: "success",
        message: "Solicitação enviada com sucesso para o provider.",
      });
    } catch (error) {
      setRequestFeedback({
        type: "error",
        message: getApiErrorMessage(
          error,
          "Não foi possível enviar a solicitação no momento.",
        ),
      });
    }
  };

  return (
    <PageWrapper gap={8}>
      <MainNav />

      <VStack align="stretch" gap={{ base: 4, md: 6 }}>
        <ChakraLink
          as={Link}
          href={returnTo ?? "/"}
          fontSize={{ base: "xs", sm: "sm" }}
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
            borderRadius={{ base: "xl", md: "2xl" }}
            borderWidth="1px"
            borderColor="gray.200"
            bg="white"
            py={20}
            justify="center"
            align="center"
            gap={3}
            px={4}
          >
            <Spinner color="teal.500" size="sm" />
            <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
              Carregando...
            </Text>
          </Flex>
        ) : null}

        {!isLoadingService && serviceErrorMessage ? (
          <Box
            borderRadius={{ base: "xl", md: "2xl" }}
            borderWidth="1px"
            borderColor="red.200"
            bg="red.50"
            p={{ base: 3, md: 5 }}
          >
            <Text fontSize={{ base: "xs", sm: "sm" }} color="red.700">
              {serviceErrorMessage}
            </Text>
          </Box>
        ) : null}

        {!isLoadingService && !serviceErrorMessage && service ? (
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
                Serviço
              </Text>
              <Heading
                as="h1"
                size={{ base: "lg", md: "xl" }}
                mt={2}
                color="gray.900"
              >
                {service.name || "Serviço sem nome"}
              </Heading>
              <Text
                mt={3}
                fontSize={{ base: "xs", sm: "sm" }}
                color="gray.600"
                lineHeight="tall"
              >
                {service.description || "Sem descrição informada."}
              </Text>

              <Grid
                mt={6}
                templateColumns={{
                  base: "1fr",
                  md: "repeat(2, 1fr)",
                  lg: "repeat(3, 1fr)",
                }}
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
                    Preço
                  </Text>
                  <Text
                    mt={1}
                    fontWeight="semibold"
                    color="gray.900"
                    fontSize={{ base: "xs", sm: "sm" }}
                  >
                    {formatMoney(service.price)}
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
                    Faixa
                  </Text>
                  <Text
                    mt={1}
                    fontWeight="semibold"
                    color="gray.900"
                    fontSize={{ base: "xs", sm: "sm" }}
                  >
                    {formatMoney(service.priceMinimum)} -{" "}
                    {formatMoney(service.priceMaximum)}
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
                    Duração
                  </Text>
                  <Text
                    mt={1}
                    fontWeight="semibold"
                    color="gray.900"
                    fontSize={{ base: "xs", sm: "sm" }}
                  >
                    {service.duration > 0
                      ? `${service.duration} min`
                      : "Não informada"}
                  </Text>
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
                Fotos
              </Heading>

              {servicePhotos.length ? (
                <Grid
                  mt={4}
                  templateColumns={{
                    base: "1fr",
                    md: "repeat(2, 1fr)",
                    lg: "repeat(3, 1fr)",
                  }}
                  gap={{ base: 2, md: 3 }}
                >
                  {servicePhotos.map((photo) => (
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
                        alt={`Foto do serviço ${service.name || ""}`}
                        w="full"
                        h={{ base: "160px", md: "200px", lg: "220px" }}
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
                  Este serviço ainda não possui fotos.
                </Text>
              )}
            </Box>

            <Grid
              templateColumns={{ base: "1fr", lg: "1fr 1fr" }}
              gap={{ base: 4, md: 6 }}
            >
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
                        fontSize={{ base: "xs" }}
                      >
                        {category.name}
                      </Badge>
                    ))}
                  </HStack>
                ) : (
                  <Text
                    mt={4}
                    fontSize={{ base: "xs", sm: "sm" }}
                    color="gray.500"
                  >
                    Nenhuma categoria vinculada.
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
                <Heading
                  as="h2"
                  size={{ base: "sm", md: "md" }}
                  color="gray.900"
                  fontSize={{ base: "sm", md: "md" }}
                >
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
                        fontSize={{ base: "xs" }}
                      >
                        #{tag.name}
                      </Badge>
                    ))}
                  </HStack>
                ) : (
                  <Text
                    mt={4}
                    fontSize={{ base: "xs", sm: "sm" }}
                    color="gray.500"
                  >
                    Nenhuma tag vinculada.
                  </Text>
                )}
              </Box>
            </Grid>

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
                Solicitar serviço
              </Heading>

              {!isAuthenticated ? (
                <VStack align="stretch" mt={4} gap={3}>
                  <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
                    Entre com sua conta owner para selecionar um pet e enviar
                    uma solicitação para este serviço.
                  </Text>
                  <ChakraLink
                    as={Link}
                    href="/login"
                    fontSize={{ base: "xs", sm: "sm" }}
                    color="teal.600"
                    fontWeight="semibold"
                    _hover={{ textDecoration: "underline" }}
                  >
                    Entrar para solicitar
                  </ChakraLink>
                </VStack>
              ) : null}

              {isAuthenticated && isLoadingProfile ? (
                <HStack mt={4} gap={2} color="gray.500">
                  <Spinner size="xs" />
                  <Text fontSize={{ base: "xs", sm: "sm" }}>
                    Carregando perfil...
                  </Text>
                </HStack>
              ) : null}

              {isAuthenticated && profileErrorMessage ? (
                <Box
                  mt={4}
                  borderRadius={{ base: "lg", md: "xl" }}
                  borderWidth="1px"
                  borderColor="orange.200"
                  bg="orange.50"
                  p={{ base: 3, md: 4 }}
                >
                  <Text fontSize={{ base: "xs", sm: "sm" }} color="orange.700">
                    {profileErrorMessage}
                  </Text>
                </Box>
              ) : null}

              {isAuthenticated &&
              !isLoadingProfile &&
              isOwnerUser &&
              isLoadingOwnerPets ? (
                <HStack mt={4} gap={2} color="gray.500">
                  <Spinner size="xs" />
                  <Text fontSize={{ base: "xs", sm: "sm" }}>
                    Carregando pets...
                  </Text>
                </HStack>
              ) : null}

              {isAuthenticated &&
              !isLoadingProfile &&
              isOwnerUser &&
              ownerPetsErrorMessage ? (
                <Box
                  mt={4}
                  borderRadius={{ base: "lg", md: "xl" }}
                  borderWidth="1px"
                  borderColor="orange.200"
                  bg="orange.50"
                  p={{ base: 3, md: 4 }}
                >
                  <Text fontSize={{ base: "xs", sm: "sm" }} color="orange.700">
                    {ownerPetsErrorMessage}
                  </Text>
                </Box>
              ) : null}

              {isAuthenticated && !isLoadingProfile && !profileErrorMessage ? (
                <>
                  {!isOwnerUser ? (
                    <Text
                      mt={4}
                      fontSize={{ base: "xs", sm: "sm" }}
                      color="gray.600"
                    >
                      Somente usuários owner podem abrir solicitações de
                      serviço.
                    </Text>
                  ) : null}

                  {isOwnerUser &&
                  !isLoadingOwnerPets &&
                  !ownerPetsErrorMessage &&
                  !ownerPets.length ? (
                    <VStack align="stretch" mt={4} gap={3}>
                      <Text
                        fontSize={{ base: "xs", sm: "sm" }}
                        color="gray.600"
                      >
                        Você ainda não possui pets cadastrados para solicitar
                        este serviço.
                      </Text>
                      <ChakraLink
                        as={Link}
                        href="/owner"
                        fontSize={{ base: "xs", sm: "sm" }}
                        color="teal.600"
                        fontWeight="semibold"
                        _hover={{ textDecoration: "underline" }}
                      >
                        Ir para área do owner e cadastrar pet
                      </ChakraLink>
                    </VStack>
                  ) : null}

                  {isOwnerUser &&
                  !isLoadingOwnerPets &&
                  !ownerPetsErrorMessage &&
                  ownerPets.length ? (
                    <chakra.form onSubmit={handleAddRequest}>
                      <VStack align="stretch" mt={4} gap={{ base: 3, md: 4 }}>
                        <Box>
                          <Text
                            fontSize={{ base: "xs", sm: "sm" }}
                            fontWeight="medium"
                            color="gray.700"
                            mb={2}
                          >
                            Pet
                          </Text>
                          <NativeSelect.Root size="md" w="full">
                            <NativeSelect.Field
                              value={selectedPetId}
                              onChange={(event) =>
                                setSelectedPetId(event.target.value)
                              }
                              h={{ base: "10", md: "11" }}
                              borderRadius={{ base: "lg", md: "xl" }}
                              bg="gray.50"
                              borderColor="gray.200"
                              borderWidth="1px"
                              focusRingColor="teal.200"
                              fontSize={{ base: "sm", md: "base" }}
                            >
                              <option value="">Selecione o pet</option>
                              {ownerPets.map((pet) => (
                                <option key={pet.id} value={pet.id}>
                                  {pet.name}
                                </option>
                              ))}
                            </NativeSelect.Field>
                            <NativeSelect.Indicator color="gray.500" />
                          </NativeSelect.Root>
                        </Box>

                        <Box>
                          <Text
                            fontSize={{ base: "xs", sm: "sm" }}
                            fontWeight="medium"
                            color="gray.700"
                            mb={2}
                          >
                            Observações
                          </Text>
                          <Textarea
                            value={requestNotes}
                            onChange={(event) =>
                              setRequestNotes(event.target.value)
                            }
                            minH={{ base: "20", md: "24" }}
                            borderRadius={{ base: "lg", md: "xl" }}
                            bg="gray.50"
                            borderColor="gray.200"
                            focusRingColor="teal.200"
                            placeholder="Detalhes importantes para o provider"
                            fontSize={{ base: "sm", md: "base" }}
                          />
                          <Text
                            mt={1}
                            fontSize={{ base: "xs" }}
                            color="gray.500"
                          >
                            {requestNotes.length}/500
                          </Text>
                        </Box>

                        <HStack gap={3} flexWrap="wrap">
                          <Button
                            type="submit"
                            disabled={
                              isSubmittingRequest ||
                              !selectedPetId.trim() ||
                              !service?.id ||
                              !requestProviderId
                            }
                            borderRadius="full"
                            bg="teal.500"
                            color="white"
                            _hover={{ bg: "teal.600" }}
                            _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
                            fontSize={{ base: "xs", sm: "sm" }}
                            h={{ base: "9", md: "10" }}
                            px={{ base: 3, md: 4 }}
                          >
                            {isSubmittingRequest
                              ? "Enviando..."
                              : "Solicitar serviço"}
                          </Button>
                        </HStack>

                        {requestFeedback ? (
                          <Text
                            fontSize={{ base: "xs", sm: "sm" }}
                            color={
                              requestFeedback.type === "success"
                                ? "green.600"
                                : "red.600"
                            }
                          >
                            {requestFeedback.message}
                          </Text>
                        ) : null}
                      </VStack>
                    </chakra.form>
                  ) : null}
                </>
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
                  Provider
                </Heading>
                {isLoadingProvider ? (
                  <HStack gap={2} color="gray.500">
                    <Spinner size="xs" />
                    <Text fontSize={{ base: "xs", sm: "sm" }}>
                      Carregando...
                    </Text>
                  </HStack>
                ) : null}
              </Flex>

              {providerErrorMessage ? (
                <Box
                  mt={4}
                  borderRadius={{ base: "lg", md: "xl" }}
                  borderWidth="1px"
                  borderColor="orange.200"
                  bg="orange.50"
                  p={{ base: 3, md: 4 }}
                >
                  <Text fontSize={{ base: "xs", sm: "sm" }} color="orange.700">
                    {providerErrorMessage}
                  </Text>
                </Box>
              ) : null}

              {!isLoadingProvider && provider ? (
                <VStack align="stretch" gap={4} mt={4}>
                  <Box>
                    <Text
                      fontSize={{ base: "base", md: "lg" }}
                      fontWeight="semibold"
                      color="gray.900"
                    >
                      {provider.businessName || "Provider sem nome"}
                    </Text>
                    <Text
                      mt={1}
                      fontSize={{ base: "xs", sm: "sm" }}
                      color="gray.600"
                    >
                      {provider.description || "Sem descrição informada."}
                    </Text>
                    <ProviderRating
                      rating={provider.averageRating}
                      hideWhenZero
                      mt={2}
                      fontSize={{ base: "xs" }}
                    />
                  </Box>

                  <Grid
                    templateColumns={{ base: "1fr", md: "repeat(2, 1fr)" }}
                    gap={{ base: 2, md: 3 }}
                  >
                    <Box>
                      <Text
                        fontSize={{ base: "xs" }}
                        textTransform="uppercase"
                        color="gray.500"
                      >
                        Faixa de preço
                      </Text>
                      <Text
                        mt={1}
                        color="gray.800"
                        fontWeight="medium"
                        fontSize={{ base: "xs", sm: "sm" }}
                      >
                        {provider.priceRange || "Não informada"}
                      </Text>
                    </Box>
                    <Box>
                      <Text
                        fontSize={{ base: "xs" }}
                        textTransform="uppercase"
                        color="gray.500"
                      >
                        Endereço
                      </Text>
                      <Text
                        mt={1}
                        color="gray.800"
                        fontWeight="medium"
                        fontSize={{ base: "xs", sm: "sm" }}
                      >
                        {provider.address.street
                          ? `${provider.address.street}, ${provider.address.number}`
                          : "Não informado"}
                      </Text>
                      {provider.address.city || provider.address.state ? (
                        <Text
                          fontSize={{ base: "xs", sm: "xs" }}
                          color="gray.600"
                        >
                          {provider.address.city}
                          {provider.address.city && provider.address.state
                            ? " - "
                            : ""}
                          {provider.address.state}
                        </Text>
                      ) : null}
                    </Box>
                  </Grid>

                  <ChakraLink
                    as={Link}
                    href={`/providers/${provider.id}`}
                    alignSelf="flex-start"
                    display="inline-flex"
                    alignItems="center"
                    justifyContent="center"
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
                    Ver perfil completo do provider
                  </ChakraLink>

                  {providerPhotos.length ? (
                    <Box>
                      <Text
                        fontSize={{ base: "xs", sm: "sm" }}
                        color="gray.600"
                        mb={2}
                      >
                        Fotos
                      </Text>
                      <HStack gap={{ base: 2, md: 3 }} overflowX="auto" pb={1}>
                        {providerPhotos.map((photo) => (
                          <Box
                            key={photo.id}
                            minW={{ base: "140px", sm: "160px", md: "180px" }}
                            w={{ base: "140px", sm: "160px", md: "180px" }}
                            h={{ base: "100px", sm: "110px", md: "120px" }}
                            borderRadius={{ base: "md", md: "lg" }}
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
