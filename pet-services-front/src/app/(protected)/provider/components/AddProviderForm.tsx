import { useEffect, useMemo, useState, type FormEvent } from "react";
import {
  Box,
  Button,
  Flex,
  Grid,
  HStack,
  Input,
  Portal,
  Select,
  Text,
  Textarea,
  VStack,
  chakra,
  createListCollection,
} from "@chakra-ui/react";
import {
  useReferenceCities,
  useReferenceCountries,
  useReferenceStates,
} from "@/application";
import {
  BRAZIL_COUNTRY_CODE,
  fetchAddressByZipCode,
  fetchCoordinatesByZipCode,
  isBrazilCountry,
} from "@/lib/address-lookup";

type Feedback = {
  type: "success" | "error";
  message: string;
};

type AddProviderFormProps = {
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
  onPrefillFromProfile: () => void;
  canPrefillFromProfile: boolean;
  isSubmitting: boolean;
  isLoadingProviderContext: boolean;
  maxPriceRangeLength: number;
  feedback: Feedback | null;
  businessName: string;
  onBusinessNameChange: (value: string) => void;
  priceRange: string;
  onPriceRangeChange: (value: string) => void;
  description: string;
  onDescriptionChange: (value: string) => void;
  street: string;
  onStreetChange: (value: string) => void;
  addressNumber: string;
  onAddressNumberChange: (value: string) => void;
  neighborhood: string;
  onNeighborhoodChange: (value: string) => void;
  city: string;
  onCityChange: (value: string) => void;
  zipCode: string;
  onZipCodeChange: (value: string) => void;
  state: string;
  onStateChange: (value: string) => void;
  country: string;
  onCountryChange: (value: string) => void;
  latitude: string;
  onLatitudeChange: (value: string) => void;
  longitude: string;
  onLongitudeChange: (value: string) => void;
  complement: string;
  onComplementChange: (value: string) => void;
};

export default function AddProviderForm({
  onSubmit,
  onPrefillFromProfile,
  canPrefillFromProfile,
  isSubmitting,
  isLoadingProviderContext,
  maxPriceRangeLength,
  feedback,
  businessName,
  onBusinessNameChange,
  priceRange,
  onPriceRangeChange,
  description,
  onDescriptionChange,
  street,
  onStreetChange,
  addressNumber,
  onAddressNumberChange,
  neighborhood,
  onNeighborhoodChange,
  city,
  onCityChange,
  zipCode,
  onZipCodeChange,
  state,
  onStateChange,
  country,
  onCountryChange,
  latitude,
  onLatitudeChange,
  longitude,
  onLongitudeChange,
  complement,
  onComplementChange,
}: AddProviderFormProps) {
  const [isResolvingCoordinates, setIsResolvingCoordinates] = useState(false);
  const [geocodeStatus, setGeocodeStatus] = useState<
    "idle" | "success" | "error"
  >("idle");
  const [geocodeMessage, setGeocodeMessage] = useState("");
  const [isResolvingAddress, setIsResolvingAddress] = useState(false);
  const [addressLookupStatus, setAddressLookupStatus] = useState<
    "idle" | "success" | "error"
  >("idle");
  const [addressLookupMessage, setAddressLookupMessage] = useState("");
  const [selectedStateId, setSelectedStateId] = useState<number | undefined>(
    undefined,
  );

  const { data: countriesData } = useReferenceCountries();
  const { data: statesData } = useReferenceStates();
  const { data: citiesData } = useReferenceCities(
    { stateId: selectedStateId },
    {
      enabled: isBrazilCountry(country) && selectedStateId !== undefined,
    },
  );

  const countries = useMemo(() => {
    const seen = new Set<string>();
    return (countriesData?.countries ?? []).filter((item) => {
      if (seen.has(item.code)) {
        return false;
      }
      seen.add(item.code);
      return true;
    });
  }, [countriesData]);
  const states = useMemo(() => statesData?.states ?? [], [statesData]);
  const cities = useMemo(() => citiesData?.cities ?? [], [citiesData]);

  const countryCollection = useMemo(
    () =>
      createListCollection({
        items: countries.map((item) => ({
          label: `${item.flag} ${item.name}`,
          value: item.code,
        })),
      }),
    [countries],
  );

  const stateCollection = useMemo(
    () =>
      createListCollection({
        items: states.map((item) => ({
          label: item.name,
          value: item.code,
        })),
      }),
    [states],
  );

  const cityCollection = useMemo(
    () =>
      createListCollection({
        items: cities.map((item) => ({
          label: item.name,
          value: item.name,
        })),
      }),
    [cities],
  );

  useEffect(() => {
    if (isBrazilCountry(country) && country.trim().toUpperCase() !== "BR") {
      onCountryChange(BRAZIL_COUNTRY_CODE);
    }
  }, [country, onCountryChange]);

  useEffect(() => {
    if (!states.length || !state) {
      return;
    }
    const match = states.find(
      (item) => item.code.toUpperCase() === state.toUpperCase(),
    );
    setSelectedStateId(match?.id ?? undefined);
  }, [state, states]);

  const handleResolveAddressFromZipCode = async () => {
    const normalizedZipCode = zipCode.trim();

    if (!normalizedZipCode) {
      setAddressLookupStatus("error");
      setAddressLookupMessage("Informe o CEP para buscar o endereço.");
      return;
    }

    setIsResolvingAddress(true);
    setAddressLookupStatus("idle");
    setAddressLookupMessage("");

    try {
      const nextAddress = await fetchAddressByZipCode(normalizedZipCode);
      onStreetChange(nextAddress.street);
      onNeighborhoodChange(nextAddress.neighborhood);
      onCityChange(nextAddress.city);
      onStateChange(nextAddress.state);
      onCountryChange(nextAddress.country);

      const match = states.find((item) => item.code === nextAddress.state);
      setSelectedStateId(match?.id ?? undefined);

      setAddressLookupStatus("success");
      setAddressLookupMessage("Endereço preenchido automaticamente pelo CEP.");
    } catch (error) {
      setAddressLookupStatus("error");
      setAddressLookupMessage(
        error instanceof Error
          ? error.message
          : "Não foi possível buscar o endereço pelo CEP.",
      );
    } finally {
      setIsResolvingAddress(false);
    }
  };

  const handleResolveCoordinatesFromZipCode = async () => {
    const normalizedZipCode = zipCode.trim();

    if (!normalizedZipCode) {
      setGeocodeStatus("error");
      setGeocodeMessage("Informe o CEP para buscar as coordenadas.");
      return;
    }

    setIsResolvingCoordinates(true);
    setGeocodeStatus("idle");
    setGeocodeMessage("");

    try {
      const coords = await fetchCoordinatesByZipCode(normalizedZipCode);
      onLatitudeChange(String(coords.latitude));
      onLongitudeChange(String(coords.longitude));
      setGeocodeStatus("success");
      setGeocodeMessage("Coordenadas preenchidas automaticamente pelo CEP.");
    } catch (error) {
      setGeocodeStatus("error");
      setGeocodeMessage(
        error instanceof Error
          ? error.message
          : "Não foi possível buscar coordenadas pelo CEP.",
      );
    } finally {
      setIsResolvingCoordinates(false);
    }
  };

  const isBrazil = isBrazilCountry(country);

  return (
    <Box
      borderRadius="3xl"
      bg="white"
      p={{ base: 5, md: 6 }}
      borderWidth="1px"
      borderColor="gray.200"
      shadow="sm"
    >
      <Box mb={4}>
        <Text
          fontSize="xs"
          fontWeight="semibold"
          textTransform="uppercase"
          color="green.500"
        >
          Provider
        </Text>
        <Text mt={2} fontSize="xl" fontWeight="semibold" color="gray.900">
          Cadastrar provider
        </Text>
        <Text mt={2} fontSize="sm" color="gray.600">
          Antes de criar serviços, você precisa cadastrar os dados do seu
          provider.
        </Text>
      </Box>

      <chakra.form onSubmit={onSubmit}>
        <VStack align="stretch" gap={4}>
          <HStack justify="space-between" flexWrap="wrap">
            <Text fontSize="sm" color="gray.600">
              Use seus dados do perfil como base e ajuste se necessário.
            </Text>
            <Button
              type="button"
              size="sm"
              variant="outline"
              borderRadius="full"
              onClick={onPrefillFromProfile}
              disabled={!canPrefillFromProfile || isSubmitting}
            >
              Preencher com dados do perfil
            </Button>
          </HStack>

          <Grid gap={4} templateColumns={{ base: "1fr", md: "1fr 1fr" }}>
            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Nome comercial
              </Text>
              <Input
                value={businessName}
                onChange={(event) => onBusinessNameChange(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Ex: Clínica Pet Saúde"
                required
              />
            </Box>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Faixa de preço
              </Text>
              <Input
                value={priceRange}
                onChange={(event) =>
                  onPriceRangeChange(
                    event.target.value.slice(0, maxPriceRangeLength),
                  )
                }
                maxLength={maxPriceRangeLength}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Ex: 80-180"
                required
              />
              <Text mt={1.5} fontSize="xs" color="gray.500">
                Máximo de {maxPriceRangeLength} caracteres.
              </Text>
            </Box>

            <Box minW={0} gridColumn={{ base: "auto", md: "1 / -1" }}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Descrição
              </Text>
              <Textarea
                value={description}
                onChange={(event) => onDescriptionChange(event.target.value)}
                minH="20"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Descreva o seu negócio e os serviços prestados"
                required
              />
            </Box>
          </Grid>

          <Box
            borderRadius="2xl"
            borderWidth="1px"
            borderColor="gray.200"
            bg="gray.50"
            p={{ base: 3, md: 4 }}
          >
            <Text fontSize="sm" fontWeight="semibold" color="gray.900" mb={3}>
              Endereço
            </Text>

            <Grid gap={4} templateColumns={{ base: "1fr", md: "1fr 1fr 1fr" }}>
              <Box minW={0} gridColumn={{ base: "auto", md: "1 / 2" }}>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  CEP
                </Text>
                <Flex gap={2} direction={{ base: "column", sm: "row" }}>
                  <Input
                    value={zipCode}
                    onChange={(event) => onZipCodeChange(event.target.value)}
                    h="11"
                    borderRadius="xl"
                    bg="white"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                    placeholder="57000-000"
                    required
                    flex={1}
                  />
                  <Button
                    type="button"
                    onClick={handleResolveAddressFromZipCode}
                    disabled={isResolvingAddress || !zipCode.trim()}
                    borderRadius="full"
                    borderWidth="1px"
                    borderColor="gray.300"
                    bg="white"
                    color="gray.700"
                    _hover={{ bg: "gray.50" }}
                    _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
                    h="11"
                    px={4}
                    flexShrink={0}
                  >
                    {isResolvingAddress ? "Buscando..." : "Buscar"}
                  </Button>
                </Flex>
                {addressLookupStatus !== "idle" && addressLookupMessage ? (
                  <Box
                    mt={2}
                    borderRadius="xl"
                    borderWidth="1px"
                    borderColor={
                      addressLookupStatus === "success"
                        ? "green.200"
                        : "red.200"
                    }
                    bg={
                      addressLookupStatus === "success" ? "green.50" : "red.50"
                    }
                    px={3}
                    py={2}
                  >
                    <Text
                      fontSize="xs"
                      color={
                        addressLookupStatus === "success"
                          ? "green.700"
                          : "red.600"
                      }
                    >
                      {addressLookupMessage}
                    </Text>
                  </Box>
                ) : null}
              </Box>

              <Box minW={0} gridColumn={{ base: "auto", md: "2 / 4" }}>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  Rua
                </Text>
                <Input
                  value={street}
                  onChange={(event) => onStreetChange(event.target.value)}
                  h="11"
                  borderRadius="xl"
                  bg="white"
                  borderColor="gray.200"
                  focusRingColor="teal.200"
                  placeholder="Rua"
                  required
                />
              </Box>

              <Box minW={0} gridColumn={{ base: "auto", md: "1 / 2" }}>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  Número
                </Text>
                <Input
                  value={addressNumber}
                  onChange={(event) =>
                    onAddressNumberChange(event.target.value)
                  }
                  h="11"
                  borderRadius="xl"
                  bg="white"
                  borderColor="gray.200"
                  focusRingColor="teal.200"
                  placeholder="123"
                  required
                />
              </Box>

              <Box minW={0} gridColumn={{ base: "auto", md: "2 / 4" }}>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  Complemento (opcional)
                </Text>
                <Input
                  value={complement}
                  onChange={(event) => onComplementChange(event.target.value)}
                  h="11"
                  borderRadius="xl"
                  bg="white"
                  borderColor="gray.200"
                  focusRingColor="teal.200"
                  placeholder="Sala, bloco, referência"
                />
              </Box>

              <Box minW={0}>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  País
                </Text>
                <Select.Root
                  collection={countryCollection}
                  value={country ? [country] : []}
                  onValueChange={({ value }) => {
                    const nextCountry = value[0] ?? "";
                    onCountryChange(nextCountry);

                    if (!isBrazilCountry(nextCountry)) {
                      onStateChange("");
                      onCityChange("");
                      setSelectedStateId(undefined);
                    }
                  }}
                >
                  <Select.HiddenSelect />
                  <Select.Trigger
                    h="11"
                    borderRadius="xl"
                    bg="white"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                  >
                    <Select.ValueText placeholder="Selecione o país" />
                    <Select.Indicator />
                  </Select.Trigger>
                  <Portal>
                    <Select.Positioner>
                      <Select.Content zIndex={1500}>
                        {countryCollection.items.map((item) => (
                          <Select.Item key={item.value} item={item}>
                            <Select.ItemText>{item.label}</Select.ItemText>
                          </Select.Item>
                        ))}
                      </Select.Content>
                    </Select.Positioner>
                  </Portal>
                </Select.Root>
              </Box>

              <Box minW={0}>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  Estado
                </Text>
                {isBrazil ? (
                  <Select.Root
                    collection={stateCollection}
                    value={state ? [state] : []}
                    onValueChange={({ value }) => {
                      const code = value[0] ?? "";
                      onStateChange(code);

                      const match = states.find((item) => item.code === code);
                      setSelectedStateId(match?.id ?? undefined);
                      onCityChange("");
                    }}
                  >
                    <Select.HiddenSelect />
                    <Select.Trigger
                      h="11"
                      borderRadius="xl"
                      bg="white"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                    >
                      <Select.ValueText placeholder="Selecione o estado" />
                      <Select.Indicator />
                    </Select.Trigger>
                    <Portal>
                      <Select.Positioner>
                        <Select.Content zIndex={1500}>
                          {stateCollection.items.map((item) => (
                            <Select.Item key={item.value} item={item}>
                              <Select.ItemText>{item.label}</Select.ItemText>
                            </Select.Item>
                          ))}
                        </Select.Content>
                      </Select.Positioner>
                    </Portal>
                  </Select.Root>
                ) : (
                  <Input
                    value={state}
                    onChange={(event) => onStateChange(event.target.value)}
                    h="11"
                    borderRadius="xl"
                    bg="white"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                    placeholder="AL"
                    required
                  />
                )}
              </Box>

              <Box minW={0}>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  Cidade
                </Text>
                {isBrazil ? (
                  <Select.Root
                    collection={cityCollection}
                    value={city ? [city] : []}
                    onValueChange={({ value }) => onCityChange(value[0] ?? "")}
                    disabled={selectedStateId === undefined}
                  >
                    <Select.HiddenSelect />
                    <Select.Trigger
                      h="11"
                      borderRadius="xl"
                      bg="white"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                    >
                      <Select.ValueText
                        placeholder={
                          selectedStateId !== undefined
                            ? "Selecione a cidade"
                            : "Selecione um estado primeiro"
                        }
                      />
                      <Select.Indicator />
                    </Select.Trigger>
                    <Portal>
                      <Select.Positioner>
                        <Select.Content zIndex={1500}>
                          {cityCollection.items.map((item) => (
                            <Select.Item key={item.value} item={item}>
                              <Select.ItemText>{item.label}</Select.ItemText>
                            </Select.Item>
                          ))}
                        </Select.Content>
                      </Select.Positioner>
                    </Portal>
                  </Select.Root>
                ) : (
                  <Input
                    value={city}
                    onChange={(event) => onCityChange(event.target.value)}
                    h="11"
                    borderRadius="xl"
                    bg="white"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                    placeholder="Cidade"
                    required
                  />
                )}
              </Box>

              <Box minW={0} gridColumn={{ base: "auto", md: "1 / -1" }}>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  Bairro
                </Text>
                <Input
                  value={neighborhood}
                  onChange={(event) => onNeighborhoodChange(event.target.value)}
                  h="11"
                  borderRadius="xl"
                  bg="white"
                  borderColor="gray.200"
                  focusRingColor="teal.200"
                  placeholder="Bairro"
                  required
                />
              </Box>

              <Box minW={0} gridColumn={{ base: "auto", md: "1 / -1" }}>
                <Text
                  fontSize="sm"
                  fontWeight="semibold"
                  color="gray.900"
                  mb={2}
                >
                  Coordenadas
                </Text>
                <Grid
                  gap={4}
                  templateColumns={{ base: "1fr", md: "1fr 1fr auto" }}
                >
                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      Latitude
                    </Text>
                    <Input
                      type="number"
                      inputMode="decimal"
                      value={latitude}
                      onChange={(event) => onLatitudeChange(event.target.value)}
                      h="11"
                      borderRadius="xl"
                      bg="white"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="-9.6498"
                      step="0.000001"
                      required
                    />
                  </Box>

                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      Longitude
                    </Text>
                    <Input
                      type="number"
                      inputMode="decimal"
                      value={longitude}
                      onChange={(event) =>
                        onLongitudeChange(event.target.value)
                      }
                      h="11"
                      borderRadius="xl"
                      bg="white"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="-35.7089"
                      step="0.000001"
                      required
                    />
                  </Box>

                  <Flex align="end">
                    <Button
                      type="button"
                      onClick={handleResolveCoordinatesFromZipCode}
                      disabled={isResolvingCoordinates || !zipCode.trim()}
                      borderRadius="full"
                      borderWidth="1px"
                      borderColor="gray.300"
                      bg="white"
                      color="gray.700"
                      _hover={{ bg: "gray.50" }}
                      _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
                      h="11"
                      px={4}
                      w={{ base: "full", md: "auto" }}
                    >
                      {isResolvingCoordinates
                        ? "Buscando..."
                        : "Atualizar pelo CEP"}
                    </Button>
                  </Flex>
                </Grid>

                {geocodeStatus !== "idle" && geocodeMessage ? (
                  <Box
                    mt={2}
                    borderRadius="xl"
                    borderWidth="1px"
                    borderColor={
                      geocodeStatus === "success" ? "green.200" : "red.200"
                    }
                    bg={geocodeStatus === "success" ? "green.50" : "red.50"}
                    px={3}
                    py={2}
                  >
                    <Text
                      fontSize="xs"
                      color={
                        geocodeStatus === "success" ? "green.700" : "red.600"
                      }
                    >
                      {geocodeMessage}
                    </Text>
                  </Box>
                ) : null}
              </Box>
            </Grid>
          </Box>

          <HStack gap={3} flexWrap="wrap">
            <Button
              type="submit"
              borderRadius="full"
              bg="green.400"
              color="white"
              _hover={{ bg: "green.500" }}
              _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
              disabled={isSubmitting || isLoadingProviderContext}
            >
              {isSubmitting ? "Cadastrando provider..." : "Cadastrar provider"}
            </Button>
          </HStack>

          {feedback ? (
            <Text
              fontSize="xs"
              color={feedback.type === "error" ? "red.600" : "green.600"}
            >
              {feedback.message}
            </Text>
          ) : null}
        </VStack>
      </chakra.form>
    </Box>
  );
}
