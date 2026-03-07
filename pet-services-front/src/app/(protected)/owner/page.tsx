"use client";

import { useMemo, useRef, useState } from "react";
import Image from "next/image";
import { useQueryClient } from "@tanstack/react-query";
import { isAxiosError } from "axios";
import {
  Box,
  Button,
  Flex,
  Grid,
  HStack,
  Input,
  NativeSelect,
  Text,
  Textarea,
  VStack,
  chakra,
} from "@chakra-ui/react";

import {
  type ProblemDetailsResponse,
  usePetAdd,
  useSpeciesList,
  useUserAddPhoto,
  useUserProfile,
} from "@/application";
import ChangePasswordCard from "@/components/account/ChangePasswordCard";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

export default function OwnerDashboardPage() {
  const queryClient = useQueryClient();
  const { data, isLoading } = useUserProfile();
  const user = data?.user;
  const [selectedPhoto, setSelectedPhoto] = useState<File | null>(null);
  const {
    mutateAsync: addUserPhoto,
    isPending: isUploadingPhoto,
    error: uploadError,
    isSuccess: uploadSuccess,
  } = useUserAddPhoto({
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["user-profile"] }),
  });
  const {
    mutateAsync: addPet,
    isPending: isAddingPet,
    error: addPetError,
    isSuccess: addPetSuccess,
  } = usePetAdd({
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user-profile"] });
      setPetName("");
      setPetSpeciesId("");
      setPetAge("");
      setPetWeight("");
      setPetNotes("");
    },
  });

  const [petName, setPetName] = useState("");
  const [petSpeciesId, setPetSpeciesId] = useState("");
  const [petAge, setPetAge] = useState("");
  const [petWeight, setPetWeight] = useState("");
  const [petNotes, setPetNotes] = useState("");
  const photoInputRef = useRef<HTMLInputElement | null>(null);
  const {
    data: speciesData,
    isLoading: isLoadingSpecies,
    error: speciesError,
  } = useSpeciesList();

  const specieOptions = useMemo(
    () =>
      speciesData?.species?.map((specie) => ({
        value: specie.id,
        label: specie.name,
      })) ?? [],
    [speciesData?.species],
  );

  const photoFeedback = useMemo(() => {
    if (!uploadError) {
      return "";
    }

    if (isAxiosError<ProblemDetailsResponse>(uploadError)) {
      const problem = uploadError.response?.data?.errors?.[0];
      return (
        problem?.detail || problem?.title || "Não foi possível enviar a foto."
      );
    }

    return "Não foi possível enviar a foto.";
  }, [uploadError]);

  const petFeedback = useMemo(() => {
    if (!addPetError) {
      return "";
    }

    if (isAxiosError<ProblemDetailsResponse>(addPetError)) {
      const problem = addPetError.response?.data?.errors?.[0];
      return (
        problem?.detail || problem?.title || "Não foi possível cadastrar o pet."
      );
    }

    return "Não foi possível cadastrar o pet.";
  }, [addPetError]);

  const isPetFormValid =
    Boolean(petName.trim()) &&
    Boolean(petSpeciesId.trim()) &&
    Number(petAge) > 0 &&
    Number(petWeight) > 0;

  const handlePhotoUpload = async () => {
    if (!selectedPhoto) {
      return;
    }

    await addUserPhoto({ file: selectedPhoto });
    setSelectedPhoto(null);
  };

  const handleAddPet = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    await addPet({
      name: petName.trim(),
      speciesId: petSpeciesId.trim(),
      age: Number(petAge),
      weight: Number(petWeight),
      notes: petNotes.trim(),
    });
  };

  return (
    <PageWrapper gap={10}>
      <MainNav />

      <VStack align="stretch" gap={6}>
        <Box>
          <Text
            fontSize="xs"
            fontWeight="semibold"
            textTransform="uppercase"
            color="green.500"
          >
            Dashboard
          </Text>
          <Text mt={2} fontSize="2xl" fontWeight="semibold" color="gray.900">
            Olá, tutor
          </Text>
          <Text mt={2} fontSize="sm" color="gray.600">
            Este é o seu painel inicial. Aqui vão aparecer seus pets e
            agendamentos.
          </Text>
        </Box>

        <Box
          borderRadius="3xl"
          borderWidth="1px"
          borderStyle="dashed"
          borderColor="gray.300"
          bg="white"
          px={6}
          py={16}
          textAlign="center"
        >
          <Flex
            mx="auto"
            h="12"
            w="12"
            align="center"
            justify="center"
            borderRadius="2xl"
            bg="green.50"
            color="green.500"
            fontSize="lg"
            fontWeight="semibold"
          >
            🐾
          </Flex>
          <Text mt={4} fontSize="lg" fontWeight="semibold" color="gray.900">
            Sem dados ainda
          </Text>
          <Text mt={2} fontSize="sm" color="gray.600">
            Quando você cadastrar seu primeiro pet ou agendar um serviço, as
            informações vão aparecer aqui.
          </Text>
        </Box>
      </VStack>

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
            Pets
          </Text>
          <Text mt={2} fontSize="xl" fontWeight="semibold" color="gray.900">
            Cadastrar pet
          </Text>
          <Text mt={2} fontSize="sm" color="gray.600">
            Adicione um novo pet ao seu perfil.
          </Text>
        </Box>

        <chakra.form onSubmit={handleAddPet}>
          <VStack align="stretch" gap={4}>
            <Grid gap={4} templateColumns={{ base: "1fr", md: "1fr 1fr" }}>
              <Box minW={0}>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  Nome
                </Text>
                <Input
                  type="text"
                  value={petName}
                  onChange={(event) => setPetName(event.target.value)}
                  h="11"
                  borderRadius="xl"
                  bg="gray.50"
                  borderColor="gray.200"
                  focusRingColor="teal.200"
                  placeholder="Ex: Thor"
                  required
                />
              </Box>

              <Box minW={0}>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  Espécie
                </Text>
                <NativeSelect.Root
                  size="md"
                  w="full"
                  minW={0}
                  disabled={isLoadingSpecies || Boolean(speciesError)}
                >
                  <NativeSelect.Field
                    name="specie"
                    value={petSpeciesId}
                    onChange={(event) => setPetSpeciesId(event.target.value)}
                    h="11"
                    borderRadius="xl"
                    bg="gray.50"
                    borderColor="gray.200"
                    borderWidth="1px"
                    focusRingColor="teal.200"
                    required
                  >
                    <option value="">
                      {isLoadingSpecies
                        ? "Carregando espécies..."
                        : "Selecione a espécie"}
                    </option>
                    {specieOptions.map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                  </NativeSelect.Field>
                  <NativeSelect.Indicator color="gray.500" />
                </NativeSelect.Root>
              </Box>

              <Box minW={0}>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  Idade
                </Text>
                <Input
                  type="number"
                  inputMode="numeric"
                  value={petAge}
                  onChange={(event) => setPetAge(event.target.value)}
                  h="11"
                  borderRadius="xl"
                  bg="gray.50"
                  borderColor="gray.200"
                  focusRingColor="teal.200"
                  placeholder="Ex: 3"
                  min={0}
                  required
                />
              </Box>

              <Box minW={0}>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  Peso (kg)
                </Text>
                <Input
                  type="number"
                  inputMode="decimal"
                  value={petWeight}
                  onChange={(event) => setPetWeight(event.target.value)}
                  h="11"
                  borderRadius="xl"
                  bg="gray.50"
                  borderColor="gray.200"
                  focusRingColor="teal.200"
                  placeholder="Ex: 12.5"
                  min={0}
                  step="0.1"
                  required
                />
              </Box>
            </Grid>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Observações
              </Text>
              <Textarea
                value={petNotes}
                onChange={(event) => setPetNotes(event.target.value)}
                minH="24"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Comportamento, cuidados especiais, etc."
              />
            </Box>

            <HStack gap={3} flexWrap="wrap">
              <Button
                type="submit"
                disabled={!isPetFormValid || isAddingPet}
                borderRadius="full"
                bg="green.400"
                color="white"
                _hover={{ bg: "green.500" }}
                _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
              >
                {isAddingPet ? "Cadastrando..." : "Cadastrar pet"}
              </Button>
              {addPetSuccess ? (
                <Text fontSize="xs" color="green.600">
                  Pet cadastrado com sucesso.
                </Text>
              ) : null}
            </HStack>

            {petFeedback ? (
              <Text fontSize="xs" color="red.600">
                {petFeedback}
              </Text>
            ) : null}
            {speciesError ? (
              <Text fontSize="xs" color="red.600">
                Não foi possível carregar as espécies.
              </Text>
            ) : null}
          </VStack>
        </chakra.form>
      </Box>

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
            Fotos
          </Text>
          <Text mt={2} fontSize="xl" fontWeight="semibold" color="gray.900">
            Fotos do perfil
          </Text>
          <Text mt={2} fontSize="sm" color="gray.600">
            Envie uma foto para aparecer no seu perfil.
          </Text>
        </Box>

        {isLoading ? (
          <Text fontSize="sm" color="gray.500">
            Carregando fotos...
          </Text>
        ) : (
          <VStack align="stretch" gap={4}>
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
                type="file"
                type="button"
                borderRadius="full"
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

            <Grid
              gap={3}
              templateColumns={{
                base: "1fr 1fr",
                sm: "repeat(3, 1fr)",
                lg: "repeat(4, 1fr)",
              }}
            >
              {user?.photos?.length ? (
                user.photos.map((photo) => (
                  <Box
                    key={photo.id}
                    h="28"
                    w="full"
                    overflow="hidden"
                    borderRadius="2xl"
                    bg="gray.100"
                  >
                    <Image
                      src={photo.url}
                      alt="Foto do usuário"
                      width={112}
                      height={112}
                      style={{
                        width: "100%",
                        height: "100%",
                        objectFit: "cover",
                      }}
                      unoptimized
                    />
                  </Box>
                ))
              ) : (
                <Text fontSize="sm" color="gray.500">
                  Nenhuma foto enviada ainda.
                </Text>
              )}
            </Grid>
          </VStack>
        )}
      </Box>

      <ChangePasswordCard />
    </PageWrapper>
  );
}
