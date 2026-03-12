"use client";

import { Suspense, useCallback, useMemo } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import {
  Badge,
  Box,
  Button,
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

function ServicesCatalogPageContent() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const q = searchParams.get("q") ?? "";
  const categoryId = searchParams.get("category_id") ?? "";
  const tagId = searchParams.get("tag_id") ?? "";
  const priceMin = Number(searchParams.get("price_min") ?? 0);
  const priceMax = Number(searchParams.get("price_max") ?? 0);
  const page = Math.max(1, Number(searchParams.get("page") ?? 1));

  const { data: categoriesData } = useCategoryList();
  const { data: tagsData } = useTagList();

  const categories = categoriesData?.categories ?? [];
  const tags = tagsData?.tags ?? [];

  const hasTextQuery = q.trim().length > 0;

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
      priceMin: priceMin > 0 ? priceMin : undefined,
      priceMax: priceMax > 0 ? priceMax : undefined,
      page,
      pageSize: PAGE_SIZE,
    }),
    [q, categoryId, tagId, priceMin, priceMax, page],
  );

  const listResult = useServiceList({
    input: listInput,
    enabled: !hasTextQuery,
  });

  const searchResult = useServiceSearch({
    input: searchInput,
    enabled: hasTextQuery,
  });

  const { data, isLoading, isError } = hasTextQuery ? searchResult : listResult;

  const services = data?.services ?? [];
  const total = data?.total ?? 0;
  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE));

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
      const fd = new FormData(e.currentTarget);
      setParam({ q: String(fd.get("q") ?? "") });
    },
    [setParam],
  );

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
                  defaultValue={q}
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
              defaultValue={priceMin > 0 ? String(priceMin) : ""}
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
              defaultValue={priceMax > 0 ? String(priceMax) : ""}
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
        {!isLoading && !isError && services.length === 0 && (
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

            {/* Paginação */}
            {totalPages > 1 && (
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
                    (p) =>
                      p === 1 || p === totalPages || Math.abs(p - page) <= 2,
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
          </>
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
