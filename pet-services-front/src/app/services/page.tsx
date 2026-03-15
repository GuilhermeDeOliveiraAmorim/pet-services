"use client";

import { Suspense, useCallback, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import {
  Badge,
  Box,
  Button,
  chakra,
  Flex,
  Grid,
  Heading,
  HStack,
  Input,
  NativeSelect,
  Spinner,
  Text,
  VStack,
} from "@chakra-ui/react";

import {
  useCategoryList,
  useServiceList,
  useServiceSearch,
  useTagList,
} from "@/application";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import { Link as ChakraLink } from "@chakra-ui/react";

const PAGE_SIZE = 12;

type LocationFeedback = {
  type: "success" | "error";
  message: string;
};

const formatMoney = (value: number): string => {
  if (!Number.isFinite(value) || value <= 0) return "";
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
    minimumFractionDigits: 2,
  }).format(value);
};

const priceLabel = (price: number, min: number, max: number): string => {
  if (price > 0) return formatMoney(price);
  if (min > 0 && max > 0) return `${formatMoney(min)} – ${formatMoney(max)}`;
  if (min > 0) return `A partir de ${formatMoney(min)}`;
  if (max > 0) return `Até ${formatMoney(max)}`;
  return "Consulte";
};

const getMapZoom = (radiusKm: number): number => {
  if (radiusKm <= 0) return 14;
  if (radiusKm <= 3) return 15;
  if (radiusKm <= 10) return 13;
  if (radiusKm <= 25) return 11;
  return 10;
};

function ServicesCatalogPageContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [zipCode, setZipCode] = useState(searchParams.get("zip_code") ?? "");
  const [qInput, setQInput] = useState(searchParams.get("q") ?? "");
  const [priceMinInput, setPriceMinInput] = useState(
    searchParams.get("price_min") ?? "",
  );
  const [priceMaxInput, setPriceMaxInput] = useState(
    searchParams.get("price_max") ?? "",
  );
  const [latitudeInput, setLatitudeInput] = useState(
    searchParams.get("latitude") ?? "",
  );
  const [longitudeInput, setLongitudeInput] = useState(
    searchParams.get("longitude") ?? "",
  );
  const [radiusKmInput, setRadiusKmInput] = useState(
    searchParams.get("radius_km") ?? "",
  );
  const [locationFeedback, setLocationFeedback] =
    useState<LocationFeedback | null>(null);
  const [isResolvingZipCode, setIsResolvingZipCode] = useState(false);

  const q = searchParams.get("q") ?? "";
  const categoryId = searchParams.get("category_id") ?? "";
  const tagId = searchParams.get("tag_id") ?? "";
  const latitude = searchParams.get("latitude") ?? "";
  const longitude = searchParams.get("longitude") ?? "";
  const radiusKm = Number(searchParams.get("radius_km") ?? 0);
  const priceMin = Number(searchParams.get("price_min") ?? 0);
  const priceMax = Number(searchParams.get("price_max") ?? 0);
  const page = Math.max(1, Number(searchParams.get("page") ?? 1));

  const { data: categoriesData } = useCategoryList();
  const { data: tagsData } = useTagList();

  const categories = categoriesData?.categories ?? [];
  const tags = tagsData?.tags ?? [];

  const hasTextQuery = q.trim().length > 0;
  const normalizedLatitude = latitude.replace(",", ".").trim();
  const normalizedLongitude = longitude.replace(",", ".").trim();
  const hasLatitude = normalizedLatitude.length > 0;
  const hasLongitude = normalizedLongitude.length > 0;
  const parsedLatitude = Number(normalizedLatitude);
  const parsedLongitude = Number(normalizedLongitude);
  const hasCoordinateInput = hasLatitude || hasLongitude;
  const hasValidCoordinates =
    hasLatitude &&
    hasLongitude &&
    Number.isFinite(parsedLatitude) &&
    Number.isFinite(parsedLongitude) &&
    parsedLatitude >= -90 &&
    parsedLatitude <= 90 &&
    parsedLongitude >= -180 &&
    parsedLongitude <= 180;
  const shouldUseSearch = hasTextQuery || hasValidCoordinates;
  const mapsQuery = hasValidCoordinates
    ? `${parsedLatitude},${parsedLongitude}`
    : "";
  const mapsUrl = mapsQuery
    ? `https://www.google.com/maps?q=${encodeURIComponent(mapsQuery)}`
    : "";
  const mapsEmbedUrl = mapsQuery
    ? `https://www.google.com/maps?q=${encodeURIComponent(mapsQuery)}&z=${getMapZoom(radiusKm)}&output=embed`
    : "";

  const listInput = useMemo(
    () => ({
      categoryId: categoryId || undefined,
      tagId: tagId || undefined,
      priceMin: priceMin > 0 ? priceMin : undefined,
      priceMax: priceMax > 0 ? priceMax : undefined,
      page,
      pageSize: PAGE_SIZE,
    }),
    [categoryId, tagId, priceMin, priceMax, page],
  );

  const searchInput = useMemo(
    () => ({
      query: q.trim(),
      categoryId: categoryId || undefined,
      tagId: tagId || undefined,
      latitude: hasValidCoordinates ? parsedLatitude : undefined,
      longitude: hasValidCoordinates ? parsedLongitude : undefined,
      radiusKm: hasValidCoordinates && radiusKm > 0 ? radiusKm : undefined,
      priceMin: priceMin > 0 ? priceMin : undefined,
      priceMax: priceMax > 0 ? priceMax : undefined,
      page,
      pageSize: PAGE_SIZE,
    }),
    [
      q,
      categoryId,
      tagId,
      hasValidCoordinates,
      parsedLatitude,
      parsedLongitude,
      radiusKm,
      priceMin,
      priceMax,
      page,
    ],
  );

  const listResult = useServiceList({
    input: listInput,
    enabled: !shouldUseSearch,
  });

  const searchResult = useServiceSearch({
    input: searchInput,
    enabled: shouldUseSearch,
  });

  const { data, isLoading, isError } = shouldUseSearch
    ? searchResult
    : listResult;

  const services = data?.services ?? [];
  const total = data?.total ?? 0;
  const totalPages = total > 0 ? Math.ceil(total / PAGE_SIZE) : 0;

  useEffect(() => {
    setZipCode(searchParams.get("zip_code") ?? "");
    setQInput(searchParams.get("q") ?? "");
    setPriceMinInput(searchParams.get("price_min") ?? "");
    setPriceMaxInput(searchParams.get("price_max") ?? "");
    setLatitudeInput(searchParams.get("latitude") ?? "");
    setLongitudeInput(searchParams.get("longitude") ?? "");
    setRadiusKmInput(searchParams.get("radius_km") ?? "");
  }, [searchParams]);

  useEffect(() => {
    if (isLoading || isError || totalPages === 0 || page <= totalPages) {
      return;
    }

    const params = new URLSearchParams(searchParams.toString());
    params.set("page", String(totalPages));
    router.replace(`/services?${params.toString()}`);
  }, [isError, isLoading, page, router, searchParams, totalPages]);

  const setParam = useCallback(
    (updates: Record<string, string>) => {
      const params = new URLSearchParams(searchParams.toString());
      Object.entries(updates).forEach(([key, value]) => {
        if (value) {
          params.set(key, value);
        } else {
          params.delete(key);
        }
      });
      params.delete("page");
      router.push(`/services?${params.toString()}`);
    },
    [router, searchParams],
  );

  const setPage = useCallback(
    (newPage: number) => {
      const params = new URLSearchParams(searchParams.toString());
      params.set("page", String(newPage));
      router.push(`/services?${params.toString()}`);
    },
    [router, searchParams],
  );

  const handleSearch = useCallback(
    (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault();
      setParam({ q: qInput });
    },
    [qInput, setParam],
  );

  const fetchCoordinatesByZipCode = useCallback(
    async (zipCodeValue: string) => {
      const response = await fetch(
        `/api/geocode?address=${encodeURIComponent(zipCodeValue)}`,
      );
      const data = await response.json();

      if (!response.ok) {
        const message =
          typeof data?.error === "string"
            ? data.error
            : "Não foi possível buscar coordenadas pelo CEP.";
        throw new Error(message);
      }

      const nextLatitude = Number(data?.latitude);
      const nextLongitude = Number(data?.longitude);

      const hasValidCoordinates =
        Number.isFinite(nextLatitude) &&
        Number.isFinite(nextLongitude) &&
        nextLatitude >= -90 &&
        nextLatitude <= 90 &&
        nextLongitude >= -180 &&
        nextLongitude <= 180;

      if (!hasValidCoordinates) {
        throw new Error(
          "Coordenadas inválidas retornadas para o CEP informado.",
        );
      }

      return { latitude: nextLatitude, longitude: nextLongitude };
    },
    [],
  );

  const handleSearchByZipCode = useCallback(async () => {
    const normalizedZipCode = zipCode.trim();

    if (!normalizedZipCode) {
      setLocationFeedback({
        type: "error",
        message: "Informe o CEP para buscar as coordenadas.",
      });
      return;
    }

    setIsResolvingZipCode(true);
    setLocationFeedback(null);

    try {
      const coords = await fetchCoordinatesByZipCode(normalizedZipCode);
      setParam({
        zip_code: normalizedZipCode,
        latitude: String(coords.latitude),
        longitude: String(coords.longitude),
        radius_km: radiusKm > 0 ? String(radiusKm) : "10",
      });
      setLocationFeedback({
        type: "success",
        message: "Coordenadas preenchidas automaticamente pelo CEP.",
      });
    } catch (error) {
      setLocationFeedback({
        type: "error",
        message:
          error instanceof Error
            ? error.message
            : "Não foi possível buscar coordenadas pelo CEP.",
      });
    } finally {
      setIsResolvingZipCode(false);
    }
  }, [fetchCoordinatesByZipCode, radiusKm, setParam, zipCode]);

  const handleClearLocationForm = useCallback(() => {
    setZipCode("");
    setLocationFeedback(null);
    setParam({
      zip_code: "",
      latitude: "",
      longitude: "",
      radius_km: "",
    });
  }, [setParam]);

  return (
    <PageWrapper gap={8}>
      <MainNav />

      <VStack align="stretch" gap={{ base: 4, md: 6, lg: 8 }}>
        {/* Cabeçalho */}
        <Box>
          <Text
            fontSize={{ base: "xs" }}
            fontWeight="semibold"
            textTransform="uppercase"
            color="teal.600"
          >
            Catálogo
          </Text>
          <Heading
            as="h1"
            size={{ base: "lg", md: "xl" }}
            mt={1}
            color="gray.900"
          >
            Encontre Serviços
          </Heading>
          <Text mt={1} fontSize={{ base: "xs", sm: "sm" }} color="gray.500">
            Encontre o cuidado ideal para o seu pet
          </Text>
        </Box>

        {/* Busca + filtros */}
        <Box
          borderRadius={{ base: "2xl", md: "3xl" }}
          bg="white"
          borderWidth="1px"
          borderColor="gray.200"
          p={{ base: 4, sm: 5, md: 7 }}
        >
          <form onSubmit={handleSearch}>
            <VStack gap={{ base: 3, md: 4 }} align="stretch" mb={4}>
              <HStack gap={2}>
                <Input
                  name="q"
                  value={qInput}
                  onChange={(e) => setQInput(e.target.value)}
                  placeholder="Buscar por nome..."
                  size={{ base: "sm", md: "md" }}
                  flex={1}
                  borderRadius={{ base: "lg", md: "xl" }}
                  fontSize={{ base: "sm" }}
                />
                <Button
                  type="submit"
                  colorPalette="teal"
                  borderRadius={{ base: "lg", md: "xl" }}
                  h={{ base: "9", md: "10" }}
                  fontSize={{ base: "sm" }}
                >
                  Buscar
                </Button>
              </HStack>
            </VStack>
          </form>

          <Flex
            gap={{ base: 2, md: 3 }}
            wrap="wrap"
            align="center"
            justify={{ base: "stretch", md: "flex-start" }}
          >
            <NativeSelect.Root
              size={{ base: "sm", md: "md" }}
              minW={{ base: "100%", sm: "180px" }}
            >
              <NativeSelect.Field
                value={categoryId}
                onChange={(e) => setParam({ category_id: e.target.value })}
                borderRadius={{ base: "lg", md: "lg" }}
                fontSize={{ base: "sm" }}
              >
                <option value="">Categorias</option>
                {categories.map((cat) => (
                  <option key={String(cat.id)} value={String(cat.id)}>
                    {cat.name}
                  </option>
                ))}
              </NativeSelect.Field>
              <NativeSelect.Indicator />
            </NativeSelect.Root>

            <NativeSelect.Root
              size={{ base: "sm", md: "md" }}
              minW={{ base: "100%", sm: "160px" }}
            >
              <NativeSelect.Field
                value={tagId}
                onChange={(e) => setParam({ tag_id: e.target.value })}
                borderRadius={{ base: "lg", md: "lg" }}
                fontSize={{ base: "sm" }}
              >
                <option value="">Tags</option>
                {tags.map((tag) => (
                  <option key={String(tag.id)} value={String(tag.id)}>
                    {tag.name}
                  </option>
                ))}
              </NativeSelect.Field>
              <NativeSelect.Indicator />
            </NativeSelect.Root>

            <Input
              type="number"
              placeholder="Mín"
              size={{ base: "sm", md: "md" }}
              flex={{ base: 1, sm: "auto" }}
              minW={{ base: "auto", sm: "120px" }}
              borderRadius={{ base: "lg", md: "lg" }}
              value={priceMinInput}
              onChange={(e) => setPriceMinInput(e.target.value)}
              onBlur={(e) => setParam({ price_min: e.target.value })}
              fontSize={{ base: "sm" }}
            />

            <Input
              type="number"
              placeholder="Máx"
              size={{ base: "sm", md: "md" }}
              flex={{ base: 1, sm: "auto" }}
              minW={{ base: "auto", sm: "120px" }}
              borderRadius={{ base: "lg", md: "lg" }}
              value={priceMaxInput}
              onChange={(e) => setPriceMaxInput(e.target.value)}
              onBlur={(e) => setParam({ price_max: e.target.value })}
              fontSize={{ base: "sm" }}
            />

            {(q || categoryId || tagId || priceMin > 0 || priceMax > 0) && (
              <Button
                size={{ base: "sm", md: "md" }}
                variant="ghost"
                colorPalette="gray"
                borderRadius={{ base: "lg", md: "lg" }}
                onClick={() => router.push("/services")}
                fontSize={{ base: "xs", sm: "sm" }}
                flex={{ base: 1, sm: "auto" }}
              >
                Limpar
              </Button>
            )}
          </Flex>

          <Box
            mt={{ base: 5, md: 6 }}
            pt={{ base: 5, md: 6 }}
            borderTopWidth="1px"
            borderColor="gray.100"
          >
            <Flex
              gap={{ base: 3, md: 4 }}
              align={{ base: "stretch", md: "center" }}
              justify="space-between"
              direction={{ base: "column", md: "row" }}
            >
              <Box>
                <Text
                  fontSize={{ base: "xs" }}
                  fontWeight="semibold"
                  textTransform="uppercase"
                  color="teal.600"
                >
                  Mapa
                </Text>
                <Text
                  mt={1}
                  fontSize={{ base: "sm", md: "md" }}
                  fontWeight="medium"
                  color="gray.900"
                >
                  Buscar serviços próximos
                </Text>
                <Text
                  mt={1}
                  fontSize={{ base: "xs", sm: "sm" }}
                  color="gray.500"
                >
                  Informe um CEP para centralizar a busca ou preencha as
                  coordenadas manualmente.
                </Text>
              </Box>

              <HStack gap={2} align="stretch" flexWrap="wrap">
                <Input
                  type="text"
                  inputMode="numeric"
                  placeholder="Buscar por CEP"
                  size={{ base: "sm", md: "md" }}
                  borderRadius={{ base: "lg", md: "xl" }}
                  value={zipCode}
                  onChange={(e) => setZipCode(e.target.value)}
                  maxW={{ base: "100%", md: "180px" }}
                  fontSize={{ base: "sm" }}
                />

                <Button
                  type="button"
                  variant="outline"
                  colorPalette="teal"
                  borderRadius={{ base: "lg", md: "xl" }}
                  onClick={handleSearchByZipCode}
                  loading={isResolvingZipCode}
                  fontSize={{ base: "sm" }}
                >
                  Buscar por CEP
                </Button>

                {(zipCode || latitude || longitude || radiusKm > 0) && (
                  <Button
                    type="button"
                    variant="ghost"
                    colorPalette="gray"
                    borderRadius={{ base: "lg", md: "xl" }}
                    onClick={handleClearLocationForm}
                    fontSize={{ base: "sm" }}
                  >
                    Limpar
                  </Button>
                )}
              </HStack>
            </Flex>

            <Grid
              mt={4}
              templateColumns={{
                base: "1fr",
                sm: "repeat(2, 1fr)",
                lg: "repeat(3, 1fr)",
              }}
              gap={{ base: 3, md: 4 }}
            >
              <Input
                type="text"
                placeholder="Latitude"
                size={{ base: "sm", md: "md" }}
                borderRadius={{ base: "lg", md: "xl" }}
                value={latitudeInput}
                onChange={(e) => setLatitudeInput(e.target.value)}
                onBlur={(e) => setParam({ latitude: e.target.value })}
                fontSize={{ base: "sm" }}
              />

              <Input
                type="text"
                placeholder="Longitude"
                size={{ base: "sm", md: "md" }}
                borderRadius={{ base: "lg", md: "xl" }}
                value={longitudeInput}
                onChange={(e) => setLongitudeInput(e.target.value)}
                onBlur={(e) => setParam({ longitude: e.target.value })}
                fontSize={{ base: "sm" }}
              />

              <Input
                type="number"
                min={1}
                placeholder="Raio (km)"
                size={{ base: "sm", md: "md" }}
                borderRadius={{ base: "lg", md: "xl" }}
                value={radiusKmInput}
                onChange={(e) => setRadiusKmInput(e.target.value)}
                onBlur={(e) => setParam({ radius_km: e.target.value })}
                fontSize={{ base: "sm" }}
              />
            </Grid>

            {locationFeedback ? (
              <Box
                mt={4}
                borderRadius={{ base: "lg", md: "xl" }}
                borderWidth="1px"
                borderColor={
                  locationFeedback.type === "error" ? "red.200" : "green.200"
                }
                bg={locationFeedback.type === "error" ? "red.50" : "green.50"}
                px={{ base: 3, md: 4 }}
                py={3}
              >
                <Text
                  fontSize={{ base: "xs", sm: "sm" }}
                  color={
                    locationFeedback.type === "error" ? "red.700" : "green.700"
                  }
                >
                  {locationFeedback.message}
                </Text>
              </Box>
            ) : null}

            {hasValidCoordinates ? (
              <VStack mt={4} align="stretch" gap={3}>
                <Box
                  borderRadius={{ base: "xl", md: "2xl" }}
                  borderWidth="1px"
                  borderColor="gray.200"
                  overflow="hidden"
                  bg="gray.100"
                  h={{ base: "220px", md: "280px" }}
                >
                  <chakra.iframe
                    title="Mapa da busca por serviços"
                    src={mapsEmbedUrl}
                    w="full"
                    h="full"
                    border={0}
                    loading="lazy"
                    referrerPolicy="no-referrer-when-downgrade"
                  />
                </Box>

                <Flex
                  gap={3}
                  direction={{ base: "column", md: "row" }}
                  align={{ base: "stretch", md: "center" }}
                  justify="space-between"
                >
                  <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
                    {radiusKm > 0
                      ? `Buscando serviços em até ${radiusKm} km deste ponto.`
                      : "Buscando serviços próximos deste ponto."}
                  </Text>

                  <ChakraLink
                    href={mapsUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    alignSelf={{ base: "flex-start", md: "auto" }}
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
                    Abrir no Google Maps
                  </ChakraLink>
                </Flex>
              </VStack>
            ) : hasCoordinateInput ? (
              <Box
                mt={4}
                borderRadius={{ base: "lg", md: "xl" }}
                borderWidth="1px"
                borderColor="orange.200"
                bg="orange.50"
                px={{ base: 3, md: 4 }}
                py={3}
              >
                <Text fontSize={{ base: "xs", sm: "sm" }} color="orange.700">
                  Informe latitude entre -90 e 90 e longitude entre -180 e 180
                  para visualizar o mapa e filtrar por proximidade.
                </Text>
              </Box>
            ) : (
              <Text mt={4} fontSize={{ base: "xs", sm: "sm" }} color="gray.500">
                Defina um ponto no mapa para encontrar serviços próximos.
              </Text>
            )}
          </Box>
        </Box>

        {/* Estado de carregamento */}
        {isLoading && (
          <Flex
            borderRadius={{ base: "xl", md: "2xl" }}
            borderWidth="1px"
            borderColor="gray.200"
            bg="white"
            py={20}
            justify="center"
            align="center"
            gap={3}
          >
            <Spinner color="teal.500" size="sm" />
            <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
              Carregando...
            </Text>
          </Flex>
        )}

        {/* Estado de erro */}
        {isError && !isLoading && (
          <Box
            borderRadius={{ base: "xl", md: "2xl" }}
            borderWidth="1px"
            borderColor="red.200"
            bg="red.50"
            p={{ base: 3, md: 5 }}
          >
            <Text fontSize={{ base: "xs", sm: "sm" }} color="red.700">
              Não foi possível carregar os serviços. Tente novamente.
            </Text>
          </Box>
        )}

        {/* Lista vazia */}
        {!isLoading && !isError && services.length === 0 && total === 0 && (
          <Flex
            borderRadius={{ base: "xl", md: "2xl" }}
            borderWidth="1px"
            borderColor="gray.200"
            bg="white"
            py={20}
            justify="center"
            align="center"
            direction="column"
            gap={1}
            px={4}
          >
            <Text
              fontSize={{ base: "sm", md: "md" }}
              fontWeight="medium"
              color="gray.700"
            >
              Nenhum serviço encontrado
            </Text>
            <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.500">
              Tente ajustar os filtros de busca
            </Text>
          </Flex>
        )}

        {!isLoading && !isError && services.length === 0 && total > 0 && (
          <Box
            borderRadius={{ base: "xl", md: "2xl" }}
            borderWidth="1px"
            borderColor="orange.200"
            bg="orange.50"
            p={{ base: 3, md: 5 }}
          >
            <Text fontSize={{ base: "xs", sm: "sm" }} color="orange.700">
              A página atual não possui itens. Use a paginação abaixo para
              navegar pelos resultados.
            </Text>
          </Box>
        )}

        {/* Grid de cards */}
        {!isLoading && !isError && services.length > 0 && (
          <>
            <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.500">
              {total} serviço{total !== 1 ? "s" : ""} encontrado
              {total !== 1 ? "s" : ""}
            </Text>

            <Grid
              templateColumns={{
                base: "1fr",
                sm: "repeat(2, 1fr)",
                md: "repeat(2, 1fr)",
                lg: "repeat(3, 1fr)",
              }}
              gap={{ base: 3, md: 4 }}
            >
              {services.map((service) => (
                <ChakraLink
                  key={String(service.id)}
                  as={Link}
                  href={`/services/${service.id}?from=/services`}
                  _hover={{ textDecoration: "none" }}
                >
                  <Box
                    borderWidth="1px"
                    borderColor="gray.200"
                    borderRadius={{ base: "2xl", md: "3xl" }}
                    bg="white"
                    p={{ base: 3, sm: 4, md: 6 }}
                    h="full"
                    _hover={{ shadow: "sm", borderColor: "teal.300" }}
                    transition="all 0.15s"
                  >
                    <VStack align="start" gap={3} h="full">
                      <Box flex={1}>
                        <Text
                          fontSize={{ base: "xs" }}
                          fontWeight="semibold"
                          textTransform="uppercase"
                          color="teal.600"
                          mb={1}
                        >
                          Serviço
                        </Text>
                        <Heading
                          as="h3"
                          size={{ base: "sm", md: "sm" }}
                          color="gray.900"
                          lineClamp={2}
                          fontSize={{ base: "sm", md: "md" }}
                        >
                          {service.name}
                        </Heading>
                        <Text
                          mt={2}
                          fontSize={{ base: "xs", sm: "sm" }}
                          color="gray.600"
                          lineClamp={3}
                          lineHeight="tall"
                        >
                          {service.description}
                        </Text>
                      </Box>

                      <Box
                        borderWidth="1px"
                        borderColor="gray.200"
                        borderRadius={{ base: "lg", md: "xl" }}
                        px={3}
                        py={2}
                        w="full"
                      >
                        <Text
                          fontSize={{ base: "xs" }}
                          color="gray.500"
                          textTransform="uppercase"
                        >
                          Preço
                        </Text>
                        <Text
                          mt={0.5}
                          fontWeight="semibold"
                          color="gray.900"
                          fontSize={{ base: "xs", sm: "sm" }}
                        >
                          {priceLabel(
                            service.price,
                            service.priceMinimum,
                            service.priceMaximum,
                          )}
                        </Text>
                      </Box>

                      {service.categories?.length > 0 && (
                        <HStack gap={1} wrap="wrap">
                          {service.categories.map((cat) => (
                            <Badge
                              key={String(cat.id)}
                              colorPalette="cyan"
                              borderRadius="full"
                              px={2}
                              py={0.5}
                              fontSize={{ base: "xs" }}
                            >
                              {cat.name}
                            </Badge>
                          ))}
                        </HStack>
                      )}

                      {service.tags?.length > 0 && (
                        <HStack gap={1} wrap="wrap">
                          {service.tags.map((tag) => (
                            <Badge
                              key={String(tag.id)}
                              colorPalette="purple"
                              borderRadius="full"
                              px={2}
                              py={0.5}
                              fontSize={{ base: "xs" }}
                            >
                              #{tag.name}
                            </Badge>
                          ))}
                        </HStack>
                      )}
                    </VStack>
                  </Box>
                </ChakraLink>
              ))}
            </Grid>
          </>
        )}

        {/* Paginação */}
        {!isLoading && !isError && totalPages > 1 && (
          <Flex
            justify="center"
            gap={{ base: 1, md: 2 }}
            wrap="wrap"
            mt={4}
            flexShrink={0}
          >
            <Button
              size={{ base: "xs", md: "sm" }}
              variant="outline"
              borderRadius="full"
              disabled={page <= 1}
              onClick={() => setPage(page - 1)}
              fontSize={{ base: "xs", md: "sm" }}
            >
              ← Ant
            </Button>

            {Array.from({ length: totalPages }, (_, i) => i + 1)
              .filter(
                (p) => p === 1 || p === totalPages || Math.abs(p - page) <= 2,
              )
              .reduce<(number | "...")[]>((acc, p, i, arr) => {
                if (i > 0 && p - (arr[i - 1] as number) > 1) {
                  acc.push("...");
                }
                acc.push(p);
                return acc;
              }, [])
              .map((item, i) =>
                item === "..." ? (
                  <Text
                    key={`ellipsis-${i}`}
                    alignSelf="center"
                    px={1}
                    color="gray.400"
                    fontSize={{ base: "xs", md: "sm" }}
                  >
                    …
                  </Text>
                ) : (
                  <Button
                    key={item}
                    size={{ base: "xs", md: "sm" }}
                    borderRadius="full"
                    variant={item === page ? "solid" : "outline"}
                    colorPalette={item === page ? "teal" : "gray"}
                    onClick={() => setPage(item as number)}
                    fontSize={{ base: "xs", md: "sm" }}
                  >
                    {item}
                  </Button>
                ),
              )}

            <Button
              size={{ base: "xs", md: "sm" }}
              variant="outline"
              borderRadius="full"
              disabled={page >= totalPages}
              onClick={() => setPage(page + 1)}
              fontSize={{ base: "xs", md: "sm" }}
            >
              Prox →
            </Button>
          </Flex>
        )}
      </VStack>
    </PageWrapper>
  );
}

export default function ServicesCatalogPage() {
  return (
    <Suspense
      fallback={
        <PageWrapper gap={8}>
          <MainNav />
          <Flex
            borderRadius={{ base: "xl", md: "2xl" }}
            borderWidth="1px"
            borderColor="gray.200"
            bg="white"
            py={20}
            justify="center"
            align="center"
            gap={3}
          >
            <Spinner color="teal.500" size="sm" />
            <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
              Carregando...
            </Text>
          </Flex>
        </PageWrapper>
      }
    >
      <ServicesCatalogPageContent />
    </Suspense>
  );
}
