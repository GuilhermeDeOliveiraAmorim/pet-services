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
  Skeleton,
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

type SortBy = "relevance" | "price_asc" | "price_desc" | "duration";

const formatDuration = (min: number): string => {
  if (min <= 0) return "";
  if (min < 60) return `${min}min`;
  const h = Math.floor(min / 60);
  const m = min % 60;
  return m > 0 ? `${h}h${m}min` : `${h}h`;
};

const effectivePrice = (s: {
  price: number;
  priceMinimum: number;
  priceMaximum: number;
}): number => {
  if (s.price > 0) return s.price;
  if (s.priceMinimum > 0) return s.priceMinimum;
  if (s.priceMaximum > 0) return s.priceMaximum;
  return 0;
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
  const [radiusKmInput, setRadiusKmInput] = useState(
    searchParams.get("radius_km") ?? "",
  );
  const [locationFeedback, setLocationFeedback] =
    useState<LocationFeedback | null>(null);
  const [isResolvingZipCode, setIsResolvingZipCode] = useState(false);
  const [sortBy, setSortBy] = useState<SortBy>("relevance");
  const [showAdvanced, setShowAdvanced] = useState(false);

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

  const services = useMemo(() => data?.services ?? [], [data]);
  const total = data?.total ?? 0;
  const totalPages = total > 0 ? Math.ceil(total / PAGE_SIZE) : 0;

  const sortedServices = useMemo(() => {
    const arr = [...services];
    if (sortBy === "price_asc")
      arr.sort((a, b) => effectivePrice(a) - effectivePrice(b));
    else if (sortBy === "price_desc")
      arr.sort((a, b) => effectivePrice(b) - effectivePrice(a));
    else if (sortBy === "duration") arr.sort((a, b) => b.duration - a.duration);
    return arr;
  }, [services, sortBy]);

  const hasAnyFilter = Boolean(
    q ||
    categoryId ||
    tagId ||
    priceMin > 0 ||
    priceMax > 0 ||
    hasValidCoordinates,
  );

  useEffect(() => {
    setZipCode(searchParams.get("zip_code") ?? "");
    setQInput(searchParams.get("q") ?? "");
    setPriceMinInput(searchParams.get("price_min") ?? "");
    setPriceMaxInput(searchParams.get("price_max") ?? "");
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

  const handleUseMyLocation = useCallback(() => {
    if (!navigator.geolocation) {
      setLocationFeedback({
        type: "error",
        message: "Geolocalização não disponível no seu navegador.",
      });
      return;
    }
    setIsResolvingZipCode(true);
    setLocationFeedback(null);
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        setParam({
          latitude: String(pos.coords.latitude.toFixed(6)),
          longitude: String(pos.coords.longitude.toFixed(6)),
          radius_km: radiusKm > 0 ? String(radiusKm) : "10",
        });
        setLocationFeedback({
          type: "success",
          message: "Localização obtida. Buscando serviços próximos...",
        });
        setIsResolvingZipCode(false);
      },
      () => {
        setLocationFeedback({
          type: "error",
          message: "Não foi possível obter sua localização. Tente pelo CEP.",
        });
        setIsResolvingZipCode(false);
      },
    );
  }, [radiusKm, setParam]);

  const clearAllFilters = useCallback(() => {
    router.push("/services");
    setSortBy("relevance");
  }, [router]);

  useEffect(() => {
    if (tagId || priceMin > 0 || priceMax > 0 || hasCoordinateInput) {
      setShowAdvanced(true);
    }
  }, [tagId, priceMin, priceMax, hasCoordinateInput]);

  // ─────────────────────────────────────────────────────────────────
  return (
    <Box minH="100vh" bg="gray.50" color="gray.900">
      {/* Navbar */}
      <Box bg="white" borderBottomWidth="1px" borderColor="gray.200">
        <Box maxW="7xl" mx="auto" px={{ base: 4, lg: 8 }}>
          <MainNav />
        </Box>
      </Box>

      {/* ─── HERO ──────────────────────────────────────────── */}
      <Box bg="teal.600" px={{ base: 4, md: 8 }} py={{ base: 12, md: 20 }}>
        <VStack gap={{ base: 5, md: 6 }} maxW="2xl" mx="auto">
          <VStack gap={2} textAlign="center">
            <Heading
              as="h1"
              color="white"
              size={{ base: "xl", md: "2xl" }}
              fontWeight="extrabold"
            >
              Encontre serviços para o seu pet
            </Heading>
            <Text color="teal.100" fontSize={{ base: "sm", md: "md" }}>
              Busque por serviço, profissional, bairro ou cidade
            </Text>
          </VStack>

          <Box w="full">
            <form onSubmit={handleSearch}>
              <HStack gap={2}>
                <Box position="relative" flex={1}>
                  <Text
                    position="absolute"
                    left={4}
                    top="50%"
                    transform="translateY(-50%)"
                    color="gray.400"
                    zIndex={1}
                    pointerEvents="none"
                    fontSize={{ base: "md", md: "lg" }}
                  >
                    🔍
                  </Text>
                  <Input
                    name="q"
                    value={qInput}
                    onChange={(e) => setQInput(e.target.value)}
                    placeholder="Buscar serviços, profissionais, bairros..."
                    bg="white"
                    size={{ base: "md", md: "lg" }}
                    pl={{ base: "42px", md: "48px" }}
                    borderRadius="xl"
                    border="none"
                    _placeholder={{ color: "gray.400" }}
                    fontSize={{ base: "sm", md: "md" }}
                  />
                </Box>
                <Button
                  type="submit"
                  size={{ base: "md", md: "lg" }}
                  borderRadius="xl"
                  bg="white"
                  color="teal.700"
                  fontWeight="bold"
                  _hover={{ bg: "teal.50" }}
                  flexShrink={0}
                  fontSize={{ base: "sm", md: "md" }}
                >
                  Buscar
                </Button>
              </HStack>
            </form>
          </Box>
        </VStack>
      </Box>

      {/* ─── CATEGORY CHIPS BAR ────────────────────────────── */}
      <Box
        bg="white"
        borderBottomWidth="1px"
        borderColor="gray.200"
        position="sticky"
        top={0}
        zIndex={10}
      >
        <Flex
          gap={2}
          align="center"
          wrap="wrap"
          maxW="7xl"
          mx="auto"
          px={{ base: 4, lg: 8 }}
          py={3}
        >
          {categories.length > 0 && (
            <Text
              fontSize="xs"
              fontWeight="semibold"
              color="gray.500"
              mr={1}
              flexShrink={0}
            >
              Categoria:
            </Text>
          )}
          {categories.map((cat) => (
            <Button
              key={String(cat.id)}
              size="xs"
              borderRadius="full"
              variant={categoryId === String(cat.id) ? "solid" : "outline"}
              colorPalette={categoryId === String(cat.id) ? "teal" : "gray"}
              onClick={() =>
                setParam({
                  category_id:
                    categoryId === String(cat.id) ? "" : String(cat.id),
                })
              }
              transition="all 0.15s"
            >
              {cat.name}
            </Button>
          ))}

          <Button
            size="xs"
            borderRadius="full"
            variant={showAdvanced ? "solid" : "outline"}
            colorPalette={
              tagId || priceMin > 0 || priceMax > 0 || hasCoordinateInput
                ? "teal"
                : "gray"
            }
            onClick={() => setShowAdvanced((v) => !v)}
            ml={categories.length > 0 ? "auto" : undefined}
            flexShrink={0}
          >
            {showAdvanced ? "▲ Menos filtros" : "▼ Mais filtros"}
            {tagId || priceMin > 0 || priceMax > 0 || hasCoordinateInput
              ? " ●"
              : ""}
          </Button>

          {hasAnyFilter && (
            <Button
              size="xs"
              borderRadius="full"
              variant="ghost"
              colorPalette="red"
              onClick={clearAllFilters}
              flexShrink={0}
            >
              ✕ Limpar
            </Button>
          )}
        </Flex>
      </Box>

      {/* ─── ADVANCED FILTERS ──────────────────────────────── */}
      {showAdvanced && (
        <Box
          bg="gray.50"
          borderBottomWidth="1px"
          borderColor="gray.200"
          px={{ base: 4, lg: 8 }}
          py={{ base: 4, md: 5 }}
        >
          <VStack gap={4} maxW="7xl" mx="auto" align="stretch">
            {tags.length > 0 && (
              <Flex gap={2} align="center" wrap="wrap">
                <Text
                  fontSize="xs"
                  fontWeight="semibold"
                  color="gray.500"
                  flexShrink={0}
                  mr={1}
                >
                  Tags:
                </Text>
                {tags.map((tag) => (
                  <Button
                    key={String(tag.id)}
                    size="xs"
                    borderRadius="full"
                    variant={tagId === String(tag.id) ? "solid" : "outline"}
                    colorPalette={tagId === String(tag.id) ? "purple" : "gray"}
                    onClick={() =>
                      setParam({
                        tag_id: tagId === String(tag.id) ? "" : String(tag.id),
                      })
                    }
                  >
                    #{tag.name}
                  </Button>
                ))}
              </Flex>
            )}

            <Flex gap={3} align="center" wrap="wrap">
              <Text
                fontSize="xs"
                fontWeight="semibold"
                color="gray.500"
                flexShrink={0}
              >
                Preço:
              </Text>
              <HStack gap={2}>
                <Box position="relative">
                  <Text
                    position="absolute"
                    left={2.5}
                    top="50%"
                    transform="translateY(-50%)"
                    fontSize="xs"
                    color="gray.400"
                    pointerEvents="none"
                    zIndex={1}
                  >
                    R$
                  </Text>
                  <Input
                    type="number"
                    placeholder="Mínimo"
                    size="sm"
                    w="28"
                    pl={7}
                    borderRadius="lg"
                    value={priceMinInput}
                    onChange={(e) => setPriceMinInput(e.target.value)}
                    onBlur={(e) => setParam({ price_min: e.target.value })}
                    fontSize="sm"
                  />
                </Box>
                <Text fontSize="sm" color="gray.400">
                  –
                </Text>
                <Box position="relative">
                  <Text
                    position="absolute"
                    left={2.5}
                    top="50%"
                    transform="translateY(-50%)"
                    fontSize="xs"
                    color="gray.400"
                    pointerEvents="none"
                    zIndex={1}
                  >
                    R$
                  </Text>
                  <Input
                    type="number"
                    placeholder="Máximo"
                    size="sm"
                    w="28"
                    pl={7}
                    borderRadius="lg"
                    value={priceMaxInput}
                    onChange={(e) => setPriceMaxInput(e.target.value)}
                    onBlur={(e) => setParam({ price_max: e.target.value })}
                    fontSize="sm"
                  />
                </Box>
              </HStack>
            </Flex>

            {/* Localização */}
            <Box
              borderRadius="xl"
              borderWidth="1px"
              borderColor="gray.200"
              bg="white"
              p={4}
            >
              <Text
                fontSize="xs"
                fontWeight="semibold"
                color="gray.500"
                textTransform="uppercase"
                mb={3}
              >
                📍 Buscar por localização
              </Text>

              <Flex
                gap={2}
                direction={{ base: "column", md: "row" }}
                align={{ base: "stretch", md: "flex-start" }}
                wrap="wrap"
              >
                <Input
                  type="text"
                  inputMode="numeric"
                  placeholder="CEP (ex: 01310-100)"
                  size="sm"
                  borderRadius="lg"
                  value={zipCode}
                  onChange={(e) => setZipCode(e.target.value)}
                  maxW={{ base: "full", md: "180px" }}
                  fontSize="sm"
                />
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  colorPalette="teal"
                  borderRadius="lg"
                  onClick={handleSearchByZipCode}
                  loading={isResolvingZipCode}
                  flexShrink={0}
                >
                  Buscar por CEP
                </Button>
                <Button
                  type="button"
                  size="sm"
                  variant="ghost"
                  colorPalette="teal"
                  borderRadius="lg"
                  onClick={handleUseMyLocation}
                  disabled={isResolvingZipCode}
                  flexShrink={0}
                >
                  📍 Usar minha localização
                </Button>

                {hasCoordinateInput && (
                  <Input
                    type="number"
                    min={1}
                    placeholder="Raio (km)"
                    size="sm"
                    borderRadius="lg"
                    value={radiusKmInput}
                    onChange={(e) => setRadiusKmInput(e.target.value)}
                    onBlur={(e) => setParam({ radius_km: e.target.value })}
                    maxW="28"
                    fontSize="sm"
                  />
                )}

                {(zipCode || latitude || longitude || radiusKm > 0) && (
                  <Button
                    type="button"
                    size="sm"
                    variant="ghost"
                    colorPalette="gray"
                    borderRadius="lg"
                    onClick={handleClearLocationForm}
                    flexShrink={0}
                  >
                    Limpar local
                  </Button>
                )}
              </Flex>

              {locationFeedback && (
                <Box
                  mt={3}
                  borderRadius="lg"
                  borderWidth="1px"
                  borderColor={
                    locationFeedback.type === "error" ? "red.200" : "green.200"
                  }
                  bg={locationFeedback.type === "error" ? "red.50" : "green.50"}
                  px={3}
                  py={2}
                >
                  <Text
                    fontSize="xs"
                    color={
                      locationFeedback.type === "error"
                        ? "red.700"
                        : "green.700"
                    }
                  >
                    {locationFeedback.message}
                  </Text>
                </Box>
              )}

              {hasValidCoordinates && (
                <Box mt={4}>
                  <Box
                    borderRadius="xl"
                    overflow="hidden"
                    h={{ base: "200px", md: "260px" }}
                    bg="gray.100"
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
                    mt={2}
                    justify="space-between"
                    align="center"
                    gap={3}
                    flexWrap="wrap"
                  >
                    <Text fontSize="xs" color="gray.500">
                      {radiusKm > 0
                        ? `Buscando em até ${radiusKm} km deste ponto.`
                        : "Buscando próximo deste ponto."}
                    </Text>
                    <ChakraLink
                      href={mapsUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      display="inline-flex"
                      alignItems="center"
                      borderRadius="full"
                      borderWidth="1px"
                      borderColor="gray.300"
                      bg="white"
                      color="gray.600"
                      _hover={{ bg: "gray.50", textDecoration: "none" }}
                      h="7"
                      px={3}
                      fontSize="xs"
                    >
                      Abrir no Maps ↗
                    </ChakraLink>
                  </Flex>
                </Box>
              )}

              {hasCoordinateInput && !hasValidCoordinates && (
                <Box
                  mt={3}
                  borderRadius="lg"
                  borderWidth="1px"
                  borderColor="orange.200"
                  bg="orange.50"
                  px={3}
                  py={2}
                >
                  <Text fontSize="xs" color="orange.700">
                    Coordenadas inválidas. Use o CEP ou a geolocalização.
                  </Text>
                </Box>
              )}
            </Box>
          </VStack>
        </Box>
      )}

      {/* ─── MAIN RESULTS ──────────────────────────────────── */}
      <Box maxW="7xl" mx="auto" px={{ base: 4, lg: 8 }} py={{ base: 6, md: 8 }}>
        {/* Sort + count bar */}
        {!isLoading && !isError && (total > 0 || hasAnyFilter) && (
          <Flex
            justify="space-between"
            align="center"
            mb={6}
            gap={3}
            flexWrap="wrap"
          >
            <Text
              fontSize={{ base: "sm", md: "md" }}
              color="gray.700"
              fontWeight="medium"
            >
              {total > 0 ? (
                <>
                  <Text as="span" fontWeight="bold" color="teal.700">
                    {total}
                  </Text>{" "}
                  serviço{total !== 1 ? "s" : ""} encontrado
                  {total !== 1 ? "s" : ""}
                  {q ? (
                    <>
                      {" "}
                      para{" "}
                      <Text as="span" fontStyle="italic">
                        &ldquo;{q}&rdquo;
                      </Text>
                    </>
                  ) : null}
                </>
              ) : (
                "Nenhum resultado para os filtros aplicados"
              )}
            </Text>

            <HStack gap={2} flexShrink={0}>
              <Text fontSize="xs" color="gray.500">
                Ordenar:
              </Text>
              <NativeSelect.Root size="sm" minW="160px">
                <NativeSelect.Field
                  value={sortBy}
                  onChange={(e) => setSortBy(e.target.value as SortBy)}
                  borderRadius="lg"
                  fontSize="sm"
                >
                  <option value="relevance">Relevância</option>
                  <option value="price_asc">Menor preço</option>
                  <option value="price_desc">Maior preço</option>
                  <option value="duration">Maior duração</option>
                </NativeSelect.Field>
                <NativeSelect.Indicator />
              </NativeSelect.Root>
            </HStack>
          </Flex>
        )}

        {/* Loading skeleton */}
        {isLoading && (
          <Grid
            templateColumns={{
              base: "1fr",
              sm: "repeat(2, 1fr)",
              lg: "repeat(3, 1fr)",
            }}
            gap={{ base: 4, md: 5 }}
          >
            {Array.from({ length: 6 }).map((_, i) => (
              <Skeleton key={i} borderRadius="2xl" h="320px" />
            ))}
          </Grid>
        )}

        {/* Error state */}
        {isError && !isLoading && (
          <Box
            borderRadius="xl"
            borderWidth="1px"
            borderColor="red.200"
            bg="red.50"
            p={{ base: 4, md: 6 }}
          >
            <Text fontWeight="semibold" color="red.700" mb={1}>
              Erro ao carregar serviços
            </Text>
            <Text fontSize="sm" color="red.600">
              Não foi possível carregar os serviços. Verifique sua conexão e
              tente novamente.
            </Text>
          </Box>
        )}

        {/* Empty state */}
        {!isLoading && !isError && services.length === 0 && total === 0 && (
          <Flex
            direction="column"
            align="center"
            justify="center"
            py={{ base: 16, md: 24 }}
            gap={5}
            textAlign="center"
          >
            <Text fontSize="5xl">🐾</Text>
            <VStack gap={1}>
              <Text
                fontSize={{ base: "md", md: "lg" }}
                fontWeight="semibold"
                color="gray.800"
              >
                Nenhum serviço encontrado
              </Text>
              <Text fontSize="sm" color="gray.500">
                Não encontramos resultados para os filtros aplicados.
              </Text>
            </VStack>
            <Box
              bg="white"
              borderRadius="xl"
              borderWidth="1px"
              borderColor="gray.200"
              px={{ base: 5, md: 6 }}
              py={4}
              textAlign="left"
            >
              <Text fontSize="sm" color="gray.600" fontWeight="medium" mb={2}>
                Tente:
              </Text>
              <VStack align="start" gap={1}>
                <Text fontSize="sm" color="gray.600">
                  • Aumentar o raio de busca
                </Text>
                <Text fontSize="sm" color="gray.600">
                  • Remover alguns filtros
                </Text>
                <Text fontSize="sm" color="gray.600">
                  • Procurar outra categoria
                </Text>
                <Text fontSize="sm" color="gray.600">
                  • Usar termos mais genéricos na busca
                </Text>
              </VStack>
            </Box>
            {hasAnyFilter && (
              <Button
                variant="outline"
                colorPalette="teal"
                borderRadius="xl"
                onClick={clearAllFilters}
              >
                Limpar todos os filtros
              </Button>
            )}
          </Flex>
        )}

        {/* Page out of range */}
        {!isLoading && !isError && services.length === 0 && total > 0 && (
          <Box
            borderRadius="xl"
            borderWidth="1px"
            borderColor="orange.200"
            bg="orange.50"
            p={{ base: 3, md: 5 }}
            mb={5}
          >
            <Text fontSize="sm" color="orange.700">
              A página atual não possui itens. Use a paginação abaixo para
              navegar pelos resultados.
            </Text>
          </Box>
        )}

        {/* ─── Cards grid ────────────────────────────────── */}
        {!isLoading && !isError && sortedServices.length > 0 && (
          <Grid
            templateColumns={{
              base: "1fr",
              sm: "repeat(2, 1fr)",
              lg: "repeat(3, 1fr)",
            }}
            gap={{ base: 4, md: 5 }}
          >
            {sortedServices.map((service) => {
              const coverPhoto =
                Array.isArray(service.photos) && service.photos.length > 0
                  ? ((service.photos[0] as { url?: string })?.url ?? null)
                  : null;
              const durationLabel = formatDuration(service.duration);
              const priceTxt = priceLabel(
                service.price,
                service.priceMinimum,
                service.priceMaximum,
              );

              return (
                <Link
                  key={String(service.id)}
                  href={`/services/${service.id}?from=/services`}
                  style={{ textDecoration: "none", display: "flex" }}
                >
                  <Box
                    borderWidth="1px"
                    borderColor="gray.200"
                    borderRadius="2xl"
                    bg="white"
                    overflow="hidden"
                    w="full"
                    display="flex"
                    flexDirection="column"
                    _hover={{
                      shadow: "md",
                      borderColor: "teal.300",
                      transform: "translateY(-2px)",
                    }}
                    transition="all 0.2s"
                  >
                    {/* Cover photo */}
                    <Box
                      position="relative"
                      h="160px"
                      bg={coverPhoto ? "gray.100" : "teal.50"}
                      flexShrink={0}
                      overflow="hidden"
                    >
                      {coverPhoto ? (
                        <chakra.img
                          src={coverPhoto}
                          alt={service.name}
                          w="full"
                          h="full"
                          objectFit="cover"
                        />
                      ) : (
                        <Flex h="full" align="center" justify="center">
                          <Text fontSize="4xl" opacity={0.25}>
                            🐾
                          </Text>
                        </Flex>
                      )}

                      {durationLabel && (
                        <Box
                          position="absolute"
                          bottom={2}
                          right={2}
                          style={{ background: "rgba(0,0,0,0.62)" }}
                          borderRadius="full"
                          px={2.5}
                          py={0.5}
                        >
                          <Text fontSize="xs" color="white" fontWeight="medium">
                            ⏱ {durationLabel}
                          </Text>
                        </Box>
                      )}
                    </Box>

                    {/* Card body */}
                    <Flex direction="column" flex={1} p={4} gap={3}>
                      {service.categories?.length > 0 && (
                        <HStack gap={1} wrap="wrap">
                          {service.categories.slice(0, 3).map((cat) => (
                            <Badge
                              key={String(cat.id)}
                              colorPalette="cyan"
                              borderRadius="full"
                              px={2}
                              py={0.5}
                              fontSize="xs"
                            >
                              {cat.name}
                            </Badge>
                          ))}
                        </HStack>
                      )}

                      <Box flex={1}>
                        <Heading
                          as="h3"
                          fontSize={{ base: "sm", md: "md" }}
                          fontWeight="semibold"
                          color="gray.900"
                          lineClamp={2}
                        >
                          {service.name}
                        </Heading>
                        {service.description && (
                          <Text
                            mt={1}
                            fontSize="xs"
                            color="gray.500"
                            lineClamp={2}
                            lineHeight="tall"
                          >
                            {service.description}
                          </Text>
                        )}
                      </Box>

                      {service.tags?.length > 0 && (
                        <HStack gap={1} wrap="wrap">
                          {service.tags.slice(0, 3).map((tag) => (
                            <Badge
                              key={String(tag.id)}
                              colorPalette="purple"
                              variant="subtle"
                              borderRadius="full"
                              px={2}
                              py={0.5}
                              fontSize="xs"
                            >
                              #{tag.name}
                            </Badge>
                          ))}
                        </HStack>
                      )}

                      <Flex
                        w="full"
                        justify="space-between"
                        align="center"
                        pt={3}
                        borderTopWidth="1px"
                        borderColor="gray.100"
                        mt="auto"
                        gap={2}
                      >
                        <VStack align="start" gap={0}>
                          <Text
                            fontSize="xs"
                            color="gray.400"
                            textTransform="uppercase"
                            fontWeight="medium"
                          >
                            Preço
                          </Text>
                          <Text
                            fontWeight="bold"
                            fontSize={{ base: "sm", md: "md" }}
                            color="teal.700"
                          >
                            {priceTxt || "Consulte"}
                          </Text>
                        </VStack>
                        <Box
                          bg="teal.600"
                          color="white"
                          borderRadius="lg"
                          px={3}
                          py={1.5}
                          fontSize="xs"
                          fontWeight="semibold"
                          flexShrink={0}
                          _hover={{ bg: "teal.700" }}
                          transition="background 0.15s"
                        >
                          Ver detalhes →
                        </Box>
                      </Flex>
                    </Flex>
                  </Box>
                </Link>
              );
            })}
          </Grid>
        )}

        {/* ─── Pagination ─────────────────────────────────── */}
        {!isLoading && !isError && totalPages > 1 && (
          <Flex
            justify="center"
            gap={{ base: 1, md: 2 }}
            wrap="wrap"
            mt={8}
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
              ← Anterior
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
              Próximo →
            </Button>
          </Flex>
        )}
      </Box>
    </Box>
  );
}

export default function ServicesCatalogPage() {
  return (
    <Suspense
      fallback={
        <Box minH="100vh" bg="gray.50" color="gray.900">
          <Box bg="white" borderBottomWidth="1px" borderColor="gray.200">
            <Box maxW="7xl" mx="auto" px={{ base: 4, lg: 8 }}>
              <MainNav />
            </Box>
          </Box>
          <Flex py={20} justify="center" align="center" gap={3}>
            <Spinner color="teal.500" size="sm" />
            <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
              Carregando...
            </Text>
          </Flex>
        </Box>
      }
    >
      <ServicesCatalogPageContent />
    </Suspense>
  );
}
