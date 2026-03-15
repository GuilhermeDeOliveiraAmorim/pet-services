import {
  Box,
  Button,
  Grid,
  HStack,
  Spinner,
  Text,
  VStack,
  chakra,
} from "@chakra-ui/react";
import { useRef, type ChangeEvent } from "react";
import type { Pet } from "@/domain";

type PetListCardProps = {
  pets: Pet[];
  isLoading: boolean;
  errorMessage: string;
  isUpdatingPet: boolean;
  deletingPetId: string | null;
  addingPhotoPetId: string | null;
  isUploadingPetPhoto: boolean;
  onEditPet: (pet: Pet) => void;
  onDeletePet: (petId: string) => void;
  onAddPhotoToPet: (petId: string, file: File) => Promise<void>;
};

export default function PetListCard({
  pets,
  isLoading,
  errorMessage,
  isUpdatingPet,
  deletingPetId,
  addingPhotoPetId,
  isUploadingPetPhoto,
  onEditPet,
  onDeletePet,
  onAddPhotoToPet,
}: PetListCardProps) {
  const petPhotoInputRef = useRef<HTMLInputElement | null>(null);

  const sortedPets = [...pets].sort((a, b) => {
    const aTime = a.createdAt ? new Date(a.createdAt).getTime() : 0;
    const bTime = b.createdAt ? new Date(b.createdAt).getTime() : 0;
    return bTime - aTime;
  });

  const formatWeight = (weight: number) =>
    Number.isFinite(weight)
      ? new Intl.NumberFormat("pt-BR", {
          minimumFractionDigits: 0,
          maximumFractionDigits: 1,
        }).format(weight)
      : "0";

  const formatDate = (value?: string) => {
    if (!value) {
      return "Data indisponível";
    }

    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return "Data indisponível";
    }

    return new Intl.DateTimeFormat("pt-BR", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
    }).format(date);
  };

  const handleSelectPetPhoto = (petId: string) => {
    const input = petPhotoInputRef.current;
    if (!input) {
      return;
    }

    input.dataset.petId = petId;
    input.click();
  };

  const handlePetPhotoChange = async (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0] ?? null;
    const petId = event.target.dataset.petId;

    if (!file || !petId) {
      event.target.value = "";
      return;
    }

    await onAddPhotoToPet(petId, file);
    event.target.value = "";
  };

  return (
    <Box
      borderRadius={{ base: "2xl", md: "3xl" }}
      bg="white"
      p={{ base: 4, sm: 5, md: 7 }}
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
        <Text
          mt={2}
          fontSize={{ base: "lg", md: "xl" }}
          fontWeight="semibold"
          color="gray.900"
        >
          Seus pets
        </Text>
        <Text mt={2} fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
          Lista de pets vinculados ao seu perfil.
        </Text>
      </Box>

      {isLoading ? (
        <HStack gap={2} color="gray.500" py={2}>
          <Spinner size="xs" />
          <Text fontSize={{ base: "xs", sm: "sm" }}>Carregando pets...</Text>
        </HStack>
      ) : null}

      <input
        ref={petPhotoInputRef}
        type="file"
        accept="image/*"
        style={{ display: "none" }}
        onChange={(event) => {
          void handlePetPhotoChange(event);
        }}
      />

      {!isLoading && errorMessage ? (
        <Text fontSize={{ base: "xs", sm: "sm" }} color="red.600">
          {errorMessage}
        </Text>
      ) : null}

      {!isLoading && !errorMessage && !pets.length ? (
        <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
          Você ainda não possui pets cadastrados.
        </Text>
      ) : null}

      {!isLoading && !errorMessage && pets.length ? (
        <Grid
          gap={{ base: 3, md: 4 }}
          templateColumns={{ base: "1fr", md: "1fr 1fr" }}
        >
          {sortedPets.map((pet) => (
            <Box
              key={pet.id}
              borderWidth="1px"
              borderColor="gray.200"
              borderRadius={{ base: "lg", md: "xl" }}
              p={{ base: 3, md: 4 }}
              bg="gray.50"
            >
              <HStack align="start" gap={{ base: 3, md: 4 }}>
                <Box
                  flexShrink={0}
                  w={{ base: "58px", md: "72px" }}
                  h={{ base: "58px", md: "72px" }}
                  borderRadius={{ base: "lg", md: "xl" }}
                  overflow="hidden"
                  borderWidth="1px"
                  borderColor="gray.200"
                  bg="white"
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                >
                  {pet.photos?.[0]?.url ? (
                    <chakra.img
                      src={pet.photos[0].url}
                      alt={`Foto de ${pet.name || "pet"}`}
                      w="full"
                      h="full"
                      objectFit="cover"
                    />
                  ) : (
                    <Text fontSize={{ base: "lg", md: "xl" }} color="gray.400">
                      {"🐾"}
                    </Text>
                  )}
                </Box>

                <VStack align="stretch" gap={1} flex="1">
                  <Text fontWeight="semibold" color="gray.900" fontSize="sm">
                    {pet.name || "Sem nome"}
                  </Text>
                  <Text fontSize="xs" color="gray.500">
                    Cadastrado em: {formatDate(pet.createdAt)}
                  </Text>
                  <Text fontSize="xs" color="gray.600">
                    Espécie: {pet.specie?.name?.trim() || "Não informada"}
                  </Text>
                  <Text fontSize="xs" color="gray.600">
                    Raça: {pet.breed?.trim() || "Sem raça definida"}
                  </Text>
                  <Text fontSize="xs" color="gray.600">
                    Idade: {pet.age || 0} {pet.age === 1 ? "ano" : "anos"}
                  </Text>
                  <Text fontSize="xs" color="gray.600">
                    Peso: {formatWeight(pet.weight || 0)} kg
                  </Text>
                  <Text fontSize="xs" color="gray.600">
                    Observações: {pet.notes?.trim() || "Nenhuma observação"}
                  </Text>

                  <HStack mt={2} gap={2}>
                    <Button
                      size="xs"
                      borderRadius="full"
                      variant="subtle"
                      onClick={() => handleSelectPetPhoto(pet.id)}
                      disabled={
                        isUpdatingPet ||
                        Boolean(deletingPetId) ||
                        isUploadingPetPhoto
                      }
                    >
                      {addingPhotoPetId === pet.id
                        ? "Enviando foto..."
                        : "Adicionar foto"}
                    </Button>
                    <Button
                      size="xs"
                      variant="outline"
                      borderRadius="full"
                      onClick={() => onEditPet(pet)}
                      disabled={
                        isUpdatingPet ||
                        Boolean(deletingPetId) ||
                        isUploadingPetPhoto
                      }
                    >
                      Editar
                    </Button>
                    <Button
                      size="xs"
                      borderRadius="full"
                      bg="red.500"
                      color="white"
                      _hover={{ bg: "red.600" }}
                      _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
                      onClick={() => onDeletePet(pet.id)}
                      disabled={
                        isUpdatingPet ||
                        Boolean(deletingPetId) ||
                        isUploadingPetPhoto
                      }
                    >
                      {deletingPetId === pet.id ? "Excluindo..." : "Excluir"}
                    </Button>
                  </HStack>
                </VStack>
              </HStack>
            </Box>
          ))}
        </Grid>
      ) : null}
    </Box>
  );
}
