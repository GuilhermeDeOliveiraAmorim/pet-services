"use client";

import { useMemo, useRef, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import {
  Button,
  Dialog,
  Input,
  NativeSelect,
  Portal,
  Text,
  Textarea,
  VStack,
} from "@chakra-ui/react";

import {
  useBreedsList,
  usePetAdd,
  usePetAddPhoto,
  usePetDelete,
  usePetListByOwnerId,
  usePetUpdate,
  useSpeciesList,
  useUserAddPhoto,
  useUserProfile,
} from "@/application";
import type { Pet } from "@/domain";
import { UserTypes } from "@/domain";
import ChangePasswordCard from "@/components/account/ChangePasswordCard";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import { getApiErrorMessage } from "@/lib/api-error";

import DashboardIntro from "./components/DashboardIntro";
import PetFormCard from "./components/PetFormCard";
import PetListCard from "./components/PetListCard";
import PhotoUploadCard from "./components/PhotoUploadCard";

type PetActionFeedback = {
  type: "success" | "error";
  message: string;
};

export default function OwnerDashboardPage() {
  const queryClient = useQueryClient();
  const { data, isLoading } = useUserProfile();
  const user = data?.user;
  const isOwnerUser = user?.userType === UserTypes.Owner;

  const [selectedPhoto, setSelectedPhoto] = useState<File | null>(null);
  const [petName, setPetName] = useState("");
  const [petSpeciesId, setPetSpeciesId] = useState("");
  const [petBreed, setPetBreed] = useState("");
  const [petAge, setPetAge] = useState("");
  const [petWeight, setPetWeight] = useState("");
  const [petNotes, setPetNotes] = useState("");
  const [editingPet, setEditingPet] = useState<Pet | null>(null);
  const [editPetName, setEditPetName] = useState("");
  const [editPetSpeciesId, setEditPetSpeciesId] = useState("");
  const [editPetBreed, setEditPetBreed] = useState("");
  const [editPetAge, setEditPetAge] = useState("");
  const [editPetWeight, setEditPetWeight] = useState("");
  const [editPetNotes, setEditPetNotes] = useState("");
  const [confirmDeletePetId, setConfirmDeletePetId] = useState<string | null>(
    null,
  );
  const [deletingPetId, setDeletingPetId] = useState<string | null>(null);
  const [addingPhotoPetId, setAddingPhotoPetId] = useState<string | null>(null);
  const [petActionFeedback, setPetActionFeedback] =
    useState<PetActionFeedback | null>(null);

  const photoInputRef = useRef<HTMLInputElement | null>(null);

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
      queryClient.invalidateQueries({ queryKey: ["pets", "owner"] });
      setPetName("");
      setPetSpeciesId("");
      setPetBreed("");
      setPetAge("");
      setPetWeight("");
      setPetNotes("");
    },
  });

  const { mutateAsync: updatePet, isPending: isUpdatingPet } = usePetUpdate({
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["pets", "owner"] });
      setEditingPet(null);
      setPetActionFeedback({
        type: "success",
        message: "Pet atualizado com sucesso.",
      });
    },
  });

  const { mutateAsync: deletePet, isPending: isDeletingPet } = usePetDelete({
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["pets", "owner"] });
      setPetActionFeedback({
        type: "success",
        message: "Pet excluído com sucesso.",
      });
    },
  });

  const { mutateAsync: addPetPhoto, isPending: isUploadingPetPhoto } =
    usePetAddPhoto({
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: ["pets", "owner"] });
      },
    });

  const {
    data: ownerPetsData,
    isLoading: isLoadingOwnerPets,
    error: ownerPetsError,
  } = usePetListByOwnerId(isOwnerUser ? user?.id : undefined, {
    enabled: isOwnerUser && Boolean(user?.id),
  });

  const ownerPets = ownerPetsData?.pets ?? [];

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

  const {
    data: breedsData,
    isLoading: isLoadingBreeds,
    error: breedsError,
  } = useBreedsList(petSpeciesId || undefined);

  const {
    data: editBreedsData,
    isLoading: isLoadingEditBreeds,
    error: editBreedsError,
  } = useBreedsList(editPetSpeciesId || undefined, {
    enabled: Boolean(editPetSpeciesId),
  });

  const breedOptions = useMemo(() => {
    return (
      breedsData?.breeds?.map((breed) => ({
        value: breed.name,
        label: breed.name,
      })) ?? []
    );
  }, [breedsData?.breeds]);

  const showBreedField = breedOptions.length > 0;

  const editBreedOptions = useMemo(() => {
    return (
      editBreedsData?.breeds?.map((breed) => ({
        value: breed.name,
        label: breed.name,
      })) ?? []
    );
  }, [editBreedsData?.breeds]);

  const showEditBreedField = editBreedOptions.length > 0;

  const photoFeedback = useMemo(() => {
    if (!uploadError) {
      return "";
    }

    return getApiErrorMessage(uploadError, "Não foi possível enviar a foto.");
  }, [uploadError]);

  const petFeedback = useMemo(() => {
    if (!addPetError) {
      return "";
    }

    return getApiErrorMessage(addPetError, "Não foi possível cadastrar o pet.");
  }, [addPetError]);

  const ownerPetsFeedback = useMemo(() => {
    if (!ownerPetsError) {
      return "";
    }

    return getApiErrorMessage(
      ownerPetsError,
      "Não foi possível carregar os pets.",
    );
  }, [ownerPetsError]);

  const petUpdateFeedback = useMemo(() => {
    if (!petActionFeedback) {
      return "";
    }

    return petActionFeedback.message;
  }, [petActionFeedback]);

  const isPetFormValid =
    Boolean(petName.trim()) &&
    Boolean(petSpeciesId.trim()) &&
    (!showBreedField || Boolean(petBreed.trim())) &&
    Number(petAge) > 0 &&
    Number(petWeight) > 0;

  const handleSpeciesChange = (value: string) => {
    setPetSpeciesId(value);
    setPetBreed("");
  };

  const handleOpenEditPet = (pet: Pet) => {
    setPetActionFeedback(null);
    setEditingPet(pet);
    setEditPetName(pet.name ?? "");
    setEditPetSpeciesId(pet.specie?.id ?? "");
    setEditPetBreed(pet.breed ?? "");
    setEditPetAge(String(pet.age ?? 0));
    setEditPetWeight(String(pet.weight ?? 0));
    setEditPetNotes(pet.notes ?? "");
  };

  const handleCloseEditPet = () => {
    if (isUpdatingPet) {
      return;
    }
    setEditingPet(null);
  };

  const handleEditSpeciesChange = (value: string) => {
    setEditPetSpeciesId(value);
    setEditPetBreed("");
  };

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
      breed: petBreed.trim() || undefined,
      age: Number(petAge),
      weight: Number(petWeight),
      notes: petNotes.trim(),
    });
  };

  const handleEditPet = async () => {
    if (!editingPet) {
      return;
    }

    setPetActionFeedback(null);

    try {
      await updatePet({
        petId: editingPet.id,
        name: editPetName.trim(),
        speciesId: editPetSpeciesId.trim() || undefined,
        breed: editPetBreed.trim() || undefined,
        age: Number(editPetAge),
        weight: Number(editPetWeight),
        notes: editPetNotes.trim(),
      });
    } catch (error) {
      setPetActionFeedback({
        type: "error",
        message: getApiErrorMessage(error, "Não foi possível atualizar o pet."),
      });
    }
  };

  const handleRequestDeletePet = (petId: string) => {
    setPetActionFeedback(null);
    setConfirmDeletePetId(petId);
  };

  const handleCancelDeletePet = () => {
    if (isDeletingPet) {
      return;
    }
    setConfirmDeletePetId(null);
  };

  const handleConfirmDeletePet = async () => {
    if (!confirmDeletePetId) {
      return;
    }

    setPetActionFeedback(null);
    setDeletingPetId(confirmDeletePetId);

    try {
      await deletePet(confirmDeletePetId);
      setConfirmDeletePetId(null);
    } catch (error) {
      setPetActionFeedback({
        type: "error",
        message: getApiErrorMessage(error, "Não foi possível excluir o pet."),
      });
    } finally {
      setDeletingPetId(null);
    }
  };

  const handleAddPhotosToPet = async (petId: string, files: File[]) => {
    setPetActionFeedback(null);
    setAddingPhotoPetId(petId);

    try {
      for (const file of files) {
        await addPetPhoto({ petId, photo: file });
      }

      setPetActionFeedback({
        type: "success",
        message:
          files.length === 1
            ? "Foto do pet adicionada com sucesso."
            : `${files.length} fotos do pet adicionadas com sucesso.`,
      });
    } catch (error) {
      setPetActionFeedback({
        type: "error",
        message: getApiErrorMessage(
          error,
          "Não foi possível enviar as fotos do pet.",
        ),
      });
    } finally {
      setAddingPhotoPetId(null);
    }
  };

  return (
    <PageWrapper gap={10}>
      <MainNav />

      <DashboardIntro
        hasPets={ownerPets.length > 0}
        petsCount={ownerPets.length}
      />

      <PetListCard
        pets={ownerPets}
        isLoading={isLoadingOwnerPets}
        errorMessage={ownerPetsFeedback}
        isUpdatingPet={isUpdatingPet}
        deletingPetId={deletingPetId}
        addingPhotoPetId={addingPhotoPetId}
        isUploadingPetPhoto={isUploadingPetPhoto}
        onEditPet={handleOpenEditPet}
        onDeletePet={handleRequestDeletePet}
        onAddPhotosToPet={handleAddPhotosToPet}
      />

      {petUpdateFeedback ? (
        <Text
          fontSize={{ base: "xs", sm: "sm" }}
          color={
            petActionFeedback?.type === "success" ? "green.600" : "red.600"
          }
        >
          {petUpdateFeedback}
        </Text>
      ) : null}

      <PetFormCard
        petName={petName}
        petSpeciesId={petSpeciesId}
        petBreed={petBreed}
        petAge={petAge}
        petWeight={petWeight}
        petNotes={petNotes}
        onPetNameChange={setPetName}
        onPetSpeciesIdChange={handleSpeciesChange}
        onPetBreedChange={setPetBreed}
        onPetAgeChange={setPetAge}
        onPetWeightChange={setPetWeight}
        onPetNotesChange={setPetNotes}
        onSubmit={handleAddPet}
        isPetFormValid={isPetFormValid}
        isAddingPet={isAddingPet}
        addPetSuccess={addPetSuccess}
        petFeedback={petFeedback}
        hasSpeciesError={Boolean(speciesError)}
        isLoadingSpecies={isLoadingSpecies}
        hasBreedsError={Boolean(breedsError)}
        isLoadingBreeds={isLoadingBreeds}
        specieOptions={specieOptions}
        breedOptions={breedOptions}
        showBreedField={showBreedField}
      />

      <PhotoUploadCard
        isLoading={isLoading}
        photoInputRef={photoInputRef}
        selectedPhoto={selectedPhoto}
        onSelectedPhotoChange={setSelectedPhoto}
        onUpload={handlePhotoUpload}
        isUploadingPhoto={isUploadingPhoto}
        uploadSuccess={uploadSuccess}
        photoFeedback={photoFeedback}
        photos={user?.photos ?? []}
      />

      <ChangePasswordCard />

      <Dialog.Root
        open={Boolean(editingPet)}
        onOpenChange={(details) => {
          if (!details.open) {
            handleCloseEditPet();
          }
        }}
      >
        <Portal>
          <Dialog.Backdrop bg="blackAlpha.500" />
          <Dialog.Positioner p={4}>
            <Dialog.Content borderRadius="3xl" maxW="lg" w="full">
              <Dialog.Header>
                <Dialog.Title>Editar pet</Dialog.Title>
              </Dialog.Header>

              <Dialog.Body>
                <VStack align="stretch" gap={3}>
                  <Input
                    value={editPetName}
                    onChange={(event) => setEditPetName(event.target.value)}
                    placeholder="Nome"
                    h="10"
                  />

                  <NativeSelect.Root size="md" w="full">
                    <NativeSelect.Field
                      value={editPetSpeciesId}
                      onChange={(event) =>
                        handleEditSpeciesChange(event.target.value)
                      }
                      h="10"
                    >
                      <option value="">Selecione a espécie</option>
                      {specieOptions.map((option) => (
                        <option key={option.value} value={option.value}>
                          {option.label}
                        </option>
                      ))}
                    </NativeSelect.Field>
                    <NativeSelect.Indicator color="gray.500" />
                  </NativeSelect.Root>

                  {showEditBreedField ? (
                    <NativeSelect.Root
                      size="md"
                      w="full"
                      disabled={isLoadingEditBreeds || Boolean(editBreedsError)}
                    >
                      <NativeSelect.Field
                        value={editPetBreed}
                        onChange={(event) =>
                          setEditPetBreed(event.target.value)
                        }
                        h="10"
                      >
                        <option value="">
                          {isLoadingEditBreeds
                            ? "Carregando..."
                            : "Selecione a raça"}
                        </option>
                        {editBreedOptions.map((option) => (
                          <option key={option.value} value={option.label}>
                            {option.label}
                          </option>
                        ))}
                      </NativeSelect.Field>
                      <NativeSelect.Indicator color="gray.500" />
                    </NativeSelect.Root>
                  ) : null}

                  <Input
                    type="number"
                    value={editPetAge}
                    onChange={(event) => setEditPetAge(event.target.value)}
                    placeholder="Idade"
                    h="10"
                    min={0}
                  />

                  <Input
                    type="number"
                    value={editPetWeight}
                    onChange={(event) => setEditPetWeight(event.target.value)}
                    placeholder="Peso"
                    h="10"
                    min={0}
                    step="0.1"
                  />

                  <Textarea
                    value={editPetNotes}
                    onChange={(event) => setEditPetNotes(event.target.value)}
                    placeholder="Observações"
                    minH="24"
                  />

                  <Dialog.Footer px={0}>
                    <Button
                      type="button"
                      variant="outline"
                      borderRadius="full"
                      onClick={handleCloseEditPet}
                      disabled={isUpdatingPet}
                    >
                      Cancelar
                    </Button>
                    <Button
                      type="button"
                      borderRadius="full"
                      bg="green.500"
                      color="white"
                      onClick={handleEditPet}
                      disabled={isUpdatingPet}
                      _hover={{ bg: "green.600" }}
                      _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
                    >
                      {isUpdatingPet ? "Salvando..." : "Salvar"}
                    </Button>
                  </Dialog.Footer>
                </VStack>
              </Dialog.Body>
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>

      <Dialog.Root
        open={Boolean(confirmDeletePetId)}
        onOpenChange={(details) => {
          if (!details.open) {
            handleCancelDeletePet();
          }
        }}
      >
        <Portal>
          <Dialog.Backdrop bg="blackAlpha.500" />
          <Dialog.Positioner p={4}>
            <Dialog.Content borderRadius="3xl" maxW="md" w="full">
              <Dialog.Header>
                <Dialog.Title>Excluir pet?</Dialog.Title>
              </Dialog.Header>
              <Dialog.Body>
                <Text fontSize="sm" color="gray.600">
                  Esta ação remove o pet permanentemente.
                </Text>
              </Dialog.Body>
              <Dialog.Footer>
                <Button
                  type="button"
                  variant="outline"
                  borderRadius="full"
                  onClick={handleCancelDeletePet}
                  disabled={isDeletingPet}
                >
                  Cancelar
                </Button>
                <Button
                  type="button"
                  borderRadius="full"
                  bg="red.500"
                  color="white"
                  onClick={handleConfirmDeletePet}
                  disabled={isDeletingPet}
                  _hover={{ bg: "red.600" }}
                  _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
                >
                  {isDeletingPet ? "Excluindo..." : "Excluir"}
                </Button>
              </Dialog.Footer>
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>
    </PageWrapper>
  );
}
