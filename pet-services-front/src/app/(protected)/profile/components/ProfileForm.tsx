"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import Image from "next/image";
import { useQueryClient } from "@tanstack/react-query";
import {
  Box,
  Button,
  Flex,
  Grid,
  HStack,
  Input,
  Portal,
  createListCollection,
  Select,
  Text,
  VStack,
  chakra,
} from "@chakra-ui/react";

import {
  useReferenceCities,
  useReferenceCountries,
  useReferenceStates,
  useUserAddPhoto,
  useUserUpdate,
} from "@/application";
import type { UpdateUserInput } from "@/application";
import type { User } from "@/domain/entities/user";
import { getApiErrorMessage } from "@/lib/api-error";

type ProfileFormProps = {
  user: User;
};

export default function ProfileForm({ user }: ProfileFormProps) {
  const queryClient = useQueryClient();
  const {
    mutateAsync: updateUser,
    isPending: isSaving,
    error: updateError,
    isSuccess: updateSuccess,
  } = useUserUpdate();
  const {
    mutateAsync: addUserPhoto,
    isPending: isUploadingPhoto,
    error: uploadError,
    isSuccess: uploadSuccess,
  } = useUserAddPhoto({
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["user-profile"] }),
  });

  const initialValues = useMemo(
    () => ({
      name: user.name ?? "",
      phoneCountryCode: user.phone?.countryCode ?? "",
      phoneAreaCode: user.phone?.areaCode ?? "",
      phoneNumber: user.phone?.number ?? "",
      street: user.address?.street ?? "",
      addressNumber: user.address?.number ?? "",
      neighborhood: user.address?.neighborhood ?? "",
      city: user.address?.city ?? "",
      zipCode: user.address?.zipCode ?? "",
      state: user.address?.state ?? "",
      country: user.address?.country ?? "",
      complement: user.address?.complement ?? "",
      latitude:
        user.address?.location?.latitude !== undefined
          ? String(user.address.location.latitude)
          : "",
      longitude:
        user.address?.location?.longitude !== undefined
          ? String(user.address.location.longitude)
          : "",
    }),
    [user],
  );

  const [name, setName] = useState(initialValues.name);
  const [email] = useState(user.login?.email ?? "");
  const [phoneCountryCode, setPhoneCountryCode] = useState(
    initialValues.phoneCountryCode,
  );
  const [phoneAreaCode, setPhoneAreaCode] = useState(
    initialValues.phoneAreaCode,
  );
  const [phoneNumber, setPhoneNumber] = useState(initialValues.phoneNumber);

  const [street, setStreet] = useState(initialValues.street);
  const [addressNumber, setAddressNumber] = useState(
    initialValues.addressNumber,
  );
  const [neighborhood, setNeighborhood] = useState(initialValues.neighborhood);
  const [city, setCity] = useState(initialValues.city);
  const [zipCode, setZipCode] = useState(initialValues.zipCode);
  const [state, setState] = useState(initialValues.state);
  const [country, setCountry] = useState(initialValues.country);
  const [complement, setComplement] = useState(initialValues.complement);
  const [latitude, setLatitude] = useState(initialValues.latitude);
  const [longitude, setLongitude] = useState(initialValues.longitude);
  const [isResolvingCoordinates, setIsResolvingCoordinates] = useState(false);
  const [geocodeStatus, setGeocodeStatus] = useState<
    "idle" | "success" | "error"
  >("idle");
  const [geocodeMessage, setGeocodeMessage] = useState("");
  const [selectedStateId, setSelectedStateId] = useState<number | undefined>(
    undefined,
  );
  const [selectedPhoto, setSelectedPhoto] = useState<File | null>(null);
  const photoInputRef = useRef<HTMLInputElement | null>(null);

  const isBrazil = country.trim().toUpperCase() === "BR";

  const { data: countriesData } = useReferenceCountries();
  const { data: statesData } = useReferenceStates();
  const { data: citiesData } = useReferenceCities(
    { stateId: selectedStateId },
    { enabled: isBrazil && selectedStateId !== undefined },
  );

  const countries = useMemo(() => {
    const seen = new Set<string>();
    return (countriesData?.countries ?? []).filter((c) => {
      if (seen.has(c.code)) return false;
      seen.add(c.code);
      return true;
    });
  }, [countriesData]);
  const states = useMemo(() => statesData?.states ?? [], [statesData]);
  const cities = useMemo(() => citiesData?.cities ?? [], [citiesData]);

  const countryCollection = useMemo(
    () =>
      createListCollection({
        items: countries.map((c) => ({
          label: `${c.flag} ${c.name}`,
          value: c.code,
        })),
      }),
    [countries],
  );

  const stateCollection = useMemo(
    () =>
      createListCollection({
        items: states.map((s) => ({ label: s.name, value: s.code })),
      }),
    [states],
  );

  const cityCollection = useMemo(
    () =>
      createListCollection({
        items: cities.map((c) => ({ label: c.name, value: c.name })),
      }),
    [cities],
  );

  useEffect(() => {
    if (!states.length || !initialValues.state) {
      return;
    }
    const match = states.find(
      (s) => s.code.toUpperCase() === initialValues.state.toUpperCase(),
    );
    if (match) {
      setSelectedStateId(match.id);
    }
  }, [states, initialValues.state]);

  const hasChanges = useMemo(() => {
    const normalize = (value: string) => value.trim();

    return (
      normalize(name) !== normalize(initialValues.name) ||
      normalize(phoneCountryCode) !==
        normalize(initialValues.phoneCountryCode) ||
      normalize(phoneAreaCode) !== normalize(initialValues.phoneAreaCode) ||
      normalize(phoneNumber) !== normalize(initialValues.phoneNumber) ||
      normalize(street) !== normalize(initialValues.street) ||
      normalize(addressNumber) !== normalize(initialValues.addressNumber) ||
      normalize(neighborhood) !== normalize(initialValues.neighborhood) ||
      normalize(city) !== normalize(initialValues.city) ||
      normalize(zipCode) !== normalize(initialValues.zipCode) ||
      normalize(state) !== normalize(initialValues.state) ||
      normalize(country) !== normalize(initialValues.country) ||
      normalize(complement) !== normalize(initialValues.complement) ||
      normalize(latitude) !== normalize(initialValues.latitude) ||
      normalize(longitude) !== normalize(initialValues.longitude)
    );
  }, [
    addressNumber,
    city,
    complement,
    country,
    initialValues,
    latitude,
    longitude,
    name,
    neighborhood,
    phoneAreaCode,
    phoneCountryCode,
    phoneNumber,
    state,
    street,
    zipCode,
  ]);

  const feedback = useMemo(() => {
    if (!updateError) {
      return "";
    }

    return getApiErrorMessage(
      updateError,
      "Não foi possível atualizar o perfil.",
    );
  }, [updateError]);

  const photoFeedback = useMemo(() => {
    if (!uploadError) {
      return "";
    }

    return getApiErrorMessage(uploadError, "Não foi possível enviar a foto.");
  }, [uploadError]);

  const currentPhotoUrl = user.photos?.[0]?.url;

  const normalizedLatitude = latitude.replace(",", ".").trim();
  const normalizedLongitude = longitude.replace(",", ".").trim();
  const parsedLatitude = Number(normalizedLatitude);
  const parsedLongitude = Number(normalizedLongitude);
  const hasValidCoordinates =
    Number.isFinite(parsedLatitude) &&
    Number.isFinite(parsedLongitude) &&
    parsedLatitude >= -90 &&
    parsedLatitude <= 90 &&
    parsedLongitude >= -180 &&
    parsedLongitude <= 180;
  const mapsQuery = hasValidCoordinates
    ? `${parsedLatitude},${parsedLongitude}`
    : "";
  const mapsUrl = mapsQuery
    ? `https://www.google.com/maps?q=${encodeURIComponent(mapsQuery)}`
    : "";
  const mapsEmbedUrl = mapsQuery
    ? `https://www.google.com/maps?q=${encodeURIComponent(mapsQuery)}&z=15&output=embed`
    : "";

  const handlePhotoUpload = async () => {
    if (!selectedPhoto) {
      return;
    }

    await addUserPhoto({ file: selectedPhoto });
    setSelectedPhoto(null);
  };

  const fetchCoordinatesByZipCode = async (zipCodeValue: string) => {
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
      throw new Error("Coordenadas inválidas retornadas para o CEP informado.");
    }

    return { latitude: nextLatitude, longitude: nextLongitude };
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
      setLatitude(String(coords.latitude));
      setLongitude(String(coords.longitude));
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

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const payload: UpdateUserInput = {};

    if (name.trim()) {
      payload.name = name.trim();
    }

    if (phoneCountryCode || phoneAreaCode || phoneNumber) {
      payload.phone = {
        countryCode: phoneCountryCode.trim(),
        areaCode: phoneAreaCode.trim(),
        number: phoneNumber.trim(),
      };
    }

    const hasAddressField =
      street ||
      addressNumber ||
      neighborhood ||
      city ||
      zipCode ||
      state ||
      country ||
      complement ||
      latitude ||
      longitude;

    if (hasAddressField) {
      let resolvedLatitude = Number(latitude.replace(",", ".").trim());
      let resolvedLongitude = Number(longitude.replace(",", ".").trim());

      const hasValidManualCoordinates =
        Number.isFinite(resolvedLatitude) &&
        Number.isFinite(resolvedLongitude) &&
        resolvedLatitude >= -90 &&
        resolvedLatitude <= 90 &&
        resolvedLongitude >= -180 &&
        resolvedLongitude <= 180;

      const normalizedZipCode = zipCode.trim();

      if (!hasValidManualCoordinates && normalizedZipCode) {
        setIsResolvingCoordinates(true);
        setGeocodeStatus("idle");
        setGeocodeMessage("");

        try {
          const coords = await fetchCoordinatesByZipCode(normalizedZipCode);
          resolvedLatitude = coords.latitude;
          resolvedLongitude = coords.longitude;
          setLatitude(String(coords.latitude));
          setLongitude(String(coords.longitude));
          setGeocodeStatus("success");
          setGeocodeMessage(
            "Coordenadas preenchidas automaticamente pelo CEP.",
          );
        } catch (error) {
          setGeocodeStatus("error");
          setGeocodeMessage(
            error instanceof Error
              ? error.message
              : "Não foi possível buscar coordenadas pelo CEP.",
          );
          setIsResolvingCoordinates(false);
          return;
        } finally {
          setIsResolvingCoordinates(false);
        }
      }

      payload.address = {
        street: street.trim(),
        number: addressNumber.trim(),
        neighborhood: neighborhood.trim(),
        city: city.trim(),
        zipCode: zipCode.trim(),
        state: state.trim(),
        country: country.trim(),
        complement: complement.trim(),
        location: {
          latitude: Number.isFinite(resolvedLatitude) ? resolvedLatitude : 0,
          longitude: Number.isFinite(resolvedLongitude) ? resolvedLongitude : 0,
        },
      };
    }

    await updateUser(payload);
  };

  return (
    <chakra.form onSubmit={handleSubmit}>
      <VStack align="stretch" gap={4}>
        <Box
          borderRadius={{ base: "xl", md: "2xl" }}
          borderWidth="1px"
          borderColor="gray.200"
          bg="white"
          p={{ base: 3, sm: 4, md: 4 }}
        >
          <Input
            ref={photoInputRef}
            type="file"
            accept="image/*"
            display="none"
            onChange={(event) =>
              setSelectedPhoto(event.target.files?.[0] ?? null)
            }
          />

          <Grid
            gap={{ base: 4, md: 5 }}
            templateColumns={{ base: "1fr", lg: "180px 1fr" }}
            alignItems="start"
          >
            <VStack align={{ base: "center", lg: "start" }} gap={3} minW={0}>
              <Box
                h={{ base: "24", sm: "20" }}
                w={{ base: "24", sm: "20" }}
                overflow="hidden"
                borderRadius={{ base: "xl", md: "2xl" }}
                bg="gray.200"
                flexShrink={0}
              >
                {currentPhotoUrl ? (
                  <Image
                    src={currentPhotoUrl}
                    alt="Foto do usuário"
                    width={96}
                    height={96}
                    style={{
                      width: "100%",
                      height: "100%",
                      objectFit: "cover",
                    }}
                    unoptimized
                  />
                ) : null}
              </Box>

              <Text
                fontSize={{ base: "xs", sm: "sm" }}
                fontWeight="medium"
                color="gray.800"
                textAlign={{ base: "center", lg: "left" }}
              >
                Foto do perfil
              </Text>

              <Button
                type="button"
                borderRadius="full"
                size={{ base: "xs", sm: "sm" }}
                variant="outline"
                borderColor="gray.300"
                color="gray.700"
                onClick={() => photoInputRef.current?.click()}
                fontSize={{ base: "xs", sm: "sm" }}
                w={{ base: "full", lg: "auto" }}
              >
                Escolher foto
              </Button>
              <Button
                type="button"
                onClick={handlePhotoUpload}
                disabled={!selectedPhoto || isUploadingPhoto}
                borderRadius="full"
                size={{ base: "xs", sm: "sm" }}
                bg="green.400"
                color="white"
                _hover={{ bg: "green.500" }}
                _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
                fontSize={{ base: "xs", sm: "sm" }}
                w={{ base: "full", lg: "auto" }}
              >
                {isUploadingPhoto ? "Enviando..." : "Enviar"}
              </Button>

              <Text
                fontSize={{ base: "xs" }}
                color="gray.500"
                maxW={{ base: "full", lg: "160px" }}
                overflow="hidden"
                textOverflow="ellipsis"
                whiteSpace="nowrap"
                textAlign={{ base: "center", lg: "left" }}
              >
                {selectedPhoto
                  ? selectedPhoto.name
                  : "Nenhum arquivo selecionado"}
              </Text>

              {uploadSuccess ? (
                <Text fontSize={{ base: "xs" }} color="green.600">
                  Foto enviada com sucesso.
                </Text>
              ) : null}
              {photoFeedback ? (
                <Text fontSize={{ base: "xs" }} color="red.600">
                  {photoFeedback}
                </Text>
              ) : null}
            </VStack>

            <VStack align="stretch" gap={{ base: 3, md: 4 }} minW={0}>
              <Box minW={0}>
                <Text
                  fontSize={{ base: "xs", sm: "sm" }}
                  fontWeight="medium"
                  color="gray.700"
                  mb={2}
                >
                  Nome
                </Text>
                <Input
                  value={name}
                  onChange={(event) => setName(event.target.value)}
                  h={{ base: "10", md: "11" }}
                  borderRadius={{ base: "lg", md: "xl" }}
                  bg="gray.50"
                  borderColor="gray.200"
                  focusRingColor="teal.200"
                  placeholder="Seu nome"
                  required
                  fontSize={{ base: "sm" }}
                />
              </Box>

              <Box minW={0}>
                <Text
                  fontSize={{ base: "xs", sm: "sm" }}
                  fontWeight="medium"
                  color="gray.700"
                  mb={2}
                >
                  Email
                </Text>
                <Input
                  value={email}
                  readOnly
                  h={{ base: "10", md: "11" }}
                  borderRadius={{ base: "lg", md: "xl" }}
                  bg="gray.100"
                  borderColor="gray.200"
                  color="gray.500"
                  fontSize={{ base: "sm" }}
                />
                <Text mt={1.5} fontSize={{ base: "xs" }} color="gray.400">
                  O email é fixo e não pode ser alterado aqui.
                </Text>
              </Box>

              <Box minW={0}>
                <Text
                  fontSize={{ base: "xs", sm: "sm" }}
                  fontWeight="medium"
                  color="gray.700"
                  mb={2}
                >
                  Telefone
                </Text>
                <Grid
                  gap={{ base: 2, md: 4 }}
                  templateColumns={{
                    base: "1fr",
                    sm: "86px 86px minmax(0, 1fr)",
                  }}
                >
                  <Box minW={0}>
                    <Input
                      value={phoneCountryCode}
                      onChange={(event) =>
                        setPhoneCountryCode(event.target.value)
                      }
                      h={{ base: "10", md: "11" }}
                      borderRadius={{ base: "lg", md: "xl" }}
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="+55"
                      fontSize={{ base: "sm" }}
                    />
                  </Box>

                  <Box minW={0}>
                    <Input
                      value={phoneAreaCode}
                      onChange={(event) => setPhoneAreaCode(event.target.value)}
                      h={{ base: "10", md: "11" }}
                      borderRadius={{ base: "lg", md: "xl" }}
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="82"
                      fontSize={{ base: "sm" }}
                    />
                  </Box>

                  <Box minW={0}>
                    <Input
                      value={phoneNumber}
                      onChange={(event) => setPhoneNumber(event.target.value)}
                      h={{ base: "10", md: "11" }}
                      borderRadius={{ base: "lg", md: "xl" }}
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="999776761"
                      fontSize={{ base: "sm" }}
                    />
                  </Box>
                </Grid>
              </Box>
            </VStack>
          </Grid>
        </Box>

        <Box
          borderRadius={{ base: "xl", md: "2xl" }}
          borderWidth="1px"
          borderColor="teal.200"
          bg="teal.50"
          p={{ base: 3, sm: 4 }}
        >
          <Text
            fontSize={{ base: "xs" }}
            fontWeight="semibold"
            textTransform="uppercase"
            color="teal.500"
          >
            Complete seu perfil
          </Text>
          <Text mt={1} fontSize={{ base: "xs", sm: "sm" }} color="teal.700">
            Preencha o endereço e a localização para ficar visível na
            plataforma.
          </Text>
        </Box>

        <Box
          borderRadius={{ base: "xl", md: "2xl" }}
          borderWidth="1px"
          borderColor="gray.200"
          bg="gray.50"
          p={{ base: 3, sm: 4, md: 4 }}
        >
          <Text
            fontSize={{ base: "sm", md: "sm" }}
            fontWeight="semibold"
            color="gray.900"
          >
            Endereço
          </Text>
          <Grid
            mt={3}
            gap={{ base: 2, md: 4 }}
            templateColumns={{ base: "1fr", sm: "1fr 1fr", md: "1fr 1fr" }}
          >
            <Box minW={0}>
              <Text
                fontSize={{ base: "xs", sm: "sm" }}
                fontWeight="medium"
                color="gray.700"
                mb={2}
              >
                Rua
              </Text>
              <Input
                value={street}
                onChange={(event) => setStreet(event.target.value)}
                h={{ base: "10", md: "11" }}
                borderRadius={{ base: "lg", md: "xl" }}
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Rua"
                fontSize={{ base: "sm" }}
              />
            </Box>

            <Box minW={0}>
              <Text
                fontSize={{ base: "xs", sm: "sm" }}
                fontWeight="medium"
                color="gray.700"
                mb={2}
              >
                Número
              </Text>
              <Input
                value={addressNumber}
                onChange={(event) => setAddressNumber(event.target.value)}
                h={{ base: "10", md: "11" }}
                borderRadius={{ base: "lg", md: "xl" }}
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="123"
                fontSize={{ base: "sm" }}
              />
            </Box>

            <Box minW={0}>
              <Text
                fontSize={{ base: "xs", sm: "sm" }}
                fontWeight="medium"
                color="gray.700"
                mb={2}
              >
                Bairro
              </Text>
              <Input
                value={neighborhood}
                onChange={(event) => setNeighborhood(event.target.value)}
                h={{ base: "10", md: "11" }}
                borderRadius={{ base: "lg", md: "xl" }}
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Centro"
                fontSize={{ base: "sm" }}
              />
            </Box>

            <Box minW={0}>
              <Text
                fontSize={{ base: "xs", sm: "sm" }}
                fontWeight="medium"
                color="gray.700"
                mb={2}
              >
                Cidade
              </Text>
              {isBrazil ? (
                <Select.Root
                  collection={cityCollection}
                  value={city ? [city] : []}
                  onValueChange={({ value }) => setCity(value[0] ?? "")}
                  disabled={selectedStateId === undefined}
                >
                  <Select.HiddenSelect />
                  <Select.Trigger
                    h={{ base: "10", md: "11" }}
                    borderRadius={{ base: "lg", md: "xl" }}
                    bg="gray.50"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                    fontSize={{ base: "sm" }}
                    w="full"
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
                  onChange={(event) => setCity(event.target.value)}
                  h={{ base: "10", md: "11" }}
                  borderRadius={{ base: "lg", md: "xl" }}
                  bg="gray.50"
                  borderColor="gray.200"
                  focusRingColor="teal.200"
                  placeholder="Maceió"
                  fontSize={{ base: "sm" }}
                />
              )}
            </Box>

            <Box minW={0}>
              <Text
                fontSize={{ base: "xs", sm: "sm" }}
                fontWeight="medium"
                color="gray.700"
                mb={2}
              >
                CEP
              </Text>
              <Input
                value={zipCode}
                onChange={(event) => setZipCode(event.target.value)}
                h={{ base: "10", md: "11" }}
                borderRadius={{ base: "lg", md: "xl" }}
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="00000-000"
                fontSize={{ base: "sm" }}
              />
            </Box>

            <Box minW={0}>
              <Text
                fontSize={{ base: "xs", sm: "sm" }}
                fontWeight="medium"
                color="gray.700"
                mb={2}
              >
                Estado
              </Text>
              {isBrazil ? (
                <Select.Root
                  collection={stateCollection}
                  value={state ? [state] : []}
                  onValueChange={({ value }) => {
                    const code = value[0] ?? "";
                    setState(code);
                    const match = states.find((s) => s.code === code);
                    setSelectedStateId(match?.id ?? undefined);
                    setCity("");
                  }}
                >
                  <Select.HiddenSelect />
                  <Select.Trigger
                    h={{ base: "10", md: "11" }}
                    borderRadius={{ base: "lg", md: "xl" }}
                    bg="gray.50"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                    fontSize={{ base: "sm" }}
                    w="full"
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
                  onChange={(event) => setState(event.target.value)}
                  h={{ base: "10", md: "11" }}
                  borderRadius={{ base: "lg", md: "xl" }}
                  bg="gray.50"
                  borderColor="gray.200"
                  focusRingColor="teal.200"
                  placeholder="AL"
                  fontSize={{ base: "sm" }}
                />
              )}
            </Box>

            <Box minW={0}>
              <Text
                fontSize={{ base: "xs", sm: "sm" }}
                fontWeight="medium"
                color="gray.700"
                mb={2}
              >
                País
              </Text>
              <Select.Root
                collection={countryCollection}
                value={country ? [country] : []}
                onValueChange={({ value }) => {
                  const next = value[0] ?? "";
                  setCountry(next);
                  if (next.toUpperCase() !== "BR") {
                    setState("");
                    setCity("");
                    setSelectedStateId(undefined);
                  }
                }}
              >
                <Select.HiddenSelect />
                <Select.Trigger
                  h={{ base: "10", md: "11" }}
                  borderRadius={{ base: "lg", md: "xl" }}
                  bg="gray.50"
                  borderColor="gray.200"
                  focusRingColor="teal.200"
                  fontSize={{ base: "sm" }}
                  w="full"
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
              <Text
                fontSize={{ base: "xs", sm: "sm" }}
                fontWeight="medium"
                color="gray.700"
                mb={2}
              >
                Complemento
              </Text>
              <Input
                value={complement}
                onChange={(event) => setComplement(event.target.value)}
                h={{ base: "10", md: "11" }}
                borderRadius={{ base: "lg", md: "xl" }}
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Apto"
                fontSize={{ base: "sm" }}
              />
            </Box>

            <Box minW={0} gridColumn={{ base: "auto", sm: "1 / -1" }}>
              <Flex gap={{ base: 2, md: 4 }} wrap="wrap" align="end">
                <Box flex="1" minW="120px">
                  <Text
                    fontSize={{ base: "xs", sm: "sm" }}
                    fontWeight="medium"
                    color="gray.700"
                    mb={2}
                  >
                    Latitude
                  </Text>
                  <Input
                    value={latitude}
                    onChange={(event) => setLatitude(event.target.value)}
                    h={{ base: "10", md: "11" }}
                    borderRadius={{ base: "lg", md: "xl" }}
                    bg="gray.50"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                    placeholder="-10.000"
                    fontSize={{ base: "sm" }}
                  />
                </Box>

                <Box flex="1" minW="120px">
                  <Text
                    fontSize={{ base: "xs", sm: "sm" }}
                    fontWeight="medium"
                    color="gray.700"
                    mb={2}
                  >
                    Longitude
                  </Text>
                  <Input
                    value={longitude}
                    onChange={(event) => setLongitude(event.target.value)}
                    h={{ base: "10", md: "11" }}
                    borderRadius={{ base: "lg", md: "xl" }}
                    bg="gray.50"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                    placeholder="-37.000"
                    fontSize={{ base: "sm" }}
                  />
                </Box>

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
                  fontSize={{ base: "xs", sm: "sm" }}
                  h={{ base: "10", md: "11" }}
                  px={{ base: 3, sm: 4 }}
                  flexShrink={0}
                >
                  {isResolvingCoordinates
                    ? "Buscando..."
                    : "Buscar coordenadas pelo CEP"}
                </Button>
              </Flex>
            </Box>

            {geocodeStatus !== "idle" ? (
              <Box minW={0} gridColumn={{ base: "auto", sm: "1 / -1" }}>
                <Text
                  fontSize={{ base: "xs", sm: "sm" }}
                  fontWeight="medium"
                  color="gray.700"
                  mb={2}
                >
                  Localização no mapa
                </Text>

                {geocodeMessage ? (
                  <Box
                    mb={3}
                    borderRadius={{ base: "lg", md: "xl" }}
                    borderWidth="1px"
                    borderColor={
                      geocodeStatus === "success" ? "green.200" : "red.200"
                    }
                    bg={geocodeStatus === "success" ? "green.50" : "red.50"}
                    px={{ base: 3, md: 4 }}
                    py={2}
                  >
                    <Text
                      fontSize={{ base: "xs", sm: "sm" }}
                      color={
                        geocodeStatus === "success" ? "green.700" : "red.600"
                      }
                    >
                      {geocodeMessage}
                    </Text>
                  </Box>
                ) : null}

                {geocodeStatus === "success" && hasValidCoordinates ? (
                  <VStack align="stretch" gap={2}>
                    <Box
                      borderRadius={{ base: "lg", md: "xl" }}
                      borderWidth="1px"
                      borderColor="gray.200"
                      overflow="hidden"
                      bg="gray.100"
                      h={{ base: "220px", md: "260px" }}
                    >
                      <chakra.iframe
                        title="Mapa com as coordenadas do perfil"
                        src={mapsEmbedUrl}
                        w="full"
                        h="full"
                        border={0}
                        loading="lazy"
                        referrerPolicy="no-referrer-when-downgrade"
                      />
                    </Box>

                    <chakra.a
                      href={mapsUrl}
                      target="_blank"
                      rel="noopener noreferrer"
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
                      Abrir no Google Maps
                    </chakra.a>
                  </VStack>
                ) : null}
              </Box>
            ) : null}
          </Grid>
        </Box>

        {feedback ? (
          <Box
            borderRadius={{ base: "xl", md: "2xl" }}
            borderWidth="1px"
            borderColor="red.200"
            bg="red.50"
            px={{ base: 3, md: 4 }}
            py={3}
          >
            <Text fontSize={{ base: "xs", sm: "sm" }} color="red.600">
              {feedback}
            </Text>
          </Box>
        ) : null}

        {updateSuccess ? (
          <Box
            borderRadius={{ base: "xl", md: "2xl" }}
            borderWidth="1px"
            borderColor="green.200"
            bg="green.50"
            px={{ base: 3, md: 4 }}
            py={3}
          >
            <Text fontSize={{ base: "xs", sm: "sm" }} color="green.700">
              Perfil atualizado com sucesso.
            </Text>
          </Box>
        ) : null}

        <Button
          type="submit"
          disabled={isSaving || !hasChanges}
          h={{ base: "10", md: "11" }}
          w="full"
          borderRadius="full"
          bg="green.400"
          color="white"
          _hover={{ bg: "green.500" }}
          _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
          fontSize={{ base: "sm" }}
        >
          {isSaving ? "Salvando..." : "Salvar alterações"}
        </Button>
      </VStack>
    </chakra.form>
  );
}
