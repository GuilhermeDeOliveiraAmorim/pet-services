"use client";

import { useMemo, useRef, useState } from "react";
import Image from "next/image";
import { useQueryClient } from "@tanstack/react-query";
import {
  Box,
  Button,
  Flex,
  Grid,
  HStack,
  Input,
  Text,
  VStack,
  chakra,
} from "@chakra-ui/react";

import { useUserAddPhoto, useUserUpdate } from "@/application";
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
  const [selectedPhoto, setSelectedPhoto] = useState<File | null>(null);
  const photoInputRef = useRef<HTMLInputElement | null>(null);

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

  const handlePhotoUpload = async () => {
    if (!selectedPhoto) {
      return;
    }

    await addUserPhoto({ file: selectedPhoto });
    setSelectedPhoto(null);
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
          latitude: Number(latitude || 0),
          longitude: Number(longitude || 0),
        },
      };
    }

    await updateUser(payload);
  };

  return (
    <chakra.form onSubmit={handleSubmit}>
      <VStack align="stretch" gap={6}>
        <Box
          borderRadius="2xl"
          borderWidth="1px"
          borderColor="gray.200"
          bg="white"
          p={4}
        >
          <Flex wrap="wrap" align="center" gap={4}>
            <Box
              h="20"
              w="20"
              overflow="hidden"
              borderRadius="2xl"
              bg="gray.200"
            >
              {currentPhotoUrl ? (
                <Image
                  src={currentPhotoUrl}
                  alt="Foto do usuário"
                  width={80}
                  height={80}
                  style={{ width: "100%", height: "100%", objectFit: "cover" }}
                  unoptimized
                />
              ) : null}
            </Box>
            <VStack align="start" gap={2} flex={1} minW="260px">
              <Text fontSize="sm" fontWeight="medium" color="gray.800">
                Foto do perfil
              </Text>
              <Input
                ref={photoInputRef}
                type="file"
                accept="image/*"
                display="none"
                onChange={(event) =>
                  setSelectedPhoto(event.target.files?.[0] ?? null)
                }
              />
              <HStack gap={3} flexWrap="wrap">
                <Button
                  type="button"
                  borderRadius="full"
                  size="sm"
                  variant="outline"
                  borderColor="gray.300"
                  color="gray.700"
                  onClick={() => photoInputRef.current?.click()}
                >
                  Escolher foto
                </Button>
                <Text
                  fontSize="xs"
                  color="gray.500"
                  maxW="300px"
                  overflow="hidden"
                  textOverflow="ellipsis"
                  whiteSpace="nowrap"
                >
                  {selectedPhoto
                    ? selectedPhoto.name
                    : "Nenhum arquivo selecionado"}
                </Text>
                <Button
                  type="button"
                  onClick={handlePhotoUpload}
                  disabled={!selectedPhoto || isUploadingPhoto}
                  borderRadius="full"
                  size="sm"
                  bg="green.400"
                  color="white"
                  _hover={{ bg: "green.500" }}
                  _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
                >
                  {isUploadingPhoto ? "Enviando..." : "Enviar foto"}
                </Button>
              </HStack>
              {uploadSuccess ? (
                <Text fontSize="xs" color="green.600">
                  Foto enviada com sucesso.
                </Text>
              ) : null}
              {photoFeedback ? (
                <Text fontSize="xs" color="red.600">
                  {photoFeedback}
                </Text>
              ) : null}
            </VStack>
          </Flex>
        </Box>

        <Grid gap={4} templateColumns={{ base: "1fr", md: "1fr 1fr" }}>
          <Box minW={0}>
            <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
              Nome
            </Text>
            <Input
              value={name}
              onChange={(event) => setName(event.target.value)}
              h="11"
              borderRadius="xl"
              bg="gray.50"
              borderColor="gray.200"
              focusRingColor="teal.200"
              placeholder="Seu nome"
              required
            />
          </Box>

          <Box minW={0}>
            <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
              Email
            </Text>
            <Input
              value={email}
              readOnly
              h="11"
              borderRadius="xl"
              bg="gray.100"
              borderColor="gray.200"
              color="gray.500"
            />
            <Text mt={1.5} fontSize="xs" color="gray.400">
              O email é fixo e não pode ser alterado aqui.
            </Text>
          </Box>
        </Grid>

        <Box
          borderRadius="2xl"
          borderWidth="1px"
          borderColor="gray.200"
          bg="gray.50"
          p={4}
        >
          <Text fontSize="sm" fontWeight="semibold" color="gray.900">
            Telefone
          </Text>
          <Grid
            mt={3}
            gap={4}
            templateColumns={{ base: "1fr", md: "repeat(3, minmax(0, 1fr))" }}
          >
            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                DDI
              </Text>
              <Input
                value={phoneCountryCode}
                onChange={(event) => setPhoneCountryCode(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="55"
              />
            </Box>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                DDD
              </Text>
              <Input
                value={phoneAreaCode}
                onChange={(event) => setPhoneAreaCode(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="82"
              />
            </Box>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Número
              </Text>
              <Input
                value={phoneNumber}
                onChange={(event) => setPhoneNumber(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="999999999"
              />
            </Box>
          </Grid>
        </Box>

        <Box
          borderRadius="2xl"
          borderWidth="1px"
          borderColor="gray.200"
          bg="gray.50"
          p={4}
        >
          <Text fontSize="sm" fontWeight="semibold" color="gray.900">
            Endereço
          </Text>
          <Grid mt={3} gap={4} templateColumns={{ base: "1fr", md: "1fr 1fr" }}>
            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Rua
              </Text>
              <Input
                value={street}
                onChange={(event) => setStreet(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Rua"
              />
            </Box>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Número
              </Text>
              <Input
                value={addressNumber}
                onChange={(event) => setAddressNumber(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="123"
              />
            </Box>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Bairro
              </Text>
              <Input
                value={neighborhood}
                onChange={(event) => setNeighborhood(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Centro"
              />
            </Box>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Cidade
              </Text>
              <Input
                value={city}
                onChange={(event) => setCity(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Maceió"
              />
            </Box>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                CEP
              </Text>
              <Input
                value={zipCode}
                onChange={(event) => setZipCode(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="00000-000"
              />
            </Box>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Estado
              </Text>
              <Input
                value={state}
                onChange={(event) => setState(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="AL"
              />
            </Box>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                País
              </Text>
              <Input
                value={country}
                onChange={(event) => setCountry(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Brasil"
              />
            </Box>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Complemento
              </Text>
              <Input
                value={complement}
                onChange={(event) => setComplement(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Apartamento"
              />
            </Box>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Latitude
              </Text>
              <Input
                value={latitude}
                onChange={(event) => setLatitude(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="-10.000"
              />
            </Box>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Longitude
              </Text>
              <Input
                value={longitude}
                onChange={(event) => setLongitude(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="-37.000"
              />
            </Box>
          </Grid>
        </Box>

        {feedback ? (
          <Box
            borderRadius="2xl"
            borderWidth="1px"
            borderColor="red.200"
            bg="red.50"
            px={4}
            py={3}
          >
            <Text fontSize="sm" color="red.600">
              {feedback}
            </Text>
          </Box>
        ) : null}

        {updateSuccess ? (
          <Box
            borderRadius="2xl"
            borderWidth="1px"
            borderColor="green.200"
            bg="green.50"
            px={4}
            py={3}
          >
            <Text fontSize="sm" color="green.700">
              Perfil atualizado com sucesso.
            </Text>
          </Box>
        ) : null}

        <Button
          type="submit"
          disabled={isSaving || !hasChanges}
          h="11"
          w="full"
          borderRadius="full"
          bg="green.400"
          color="white"
          _hover={{ bg: "green.500" }}
          _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
        >
          {isSaving ? "Salvando..." : "Salvar alterações"}
        </Button>
      </VStack>
    </chakra.form>
  );
}
