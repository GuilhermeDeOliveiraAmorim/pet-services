"use client";

import { useMemo, useRef, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";

import {
  usePetAdd,
  useSpeciesList,
  useUserAddPhoto,
  useUserProfile,
} from "@/application";
import ChangePasswordCard from "@/components/account/ChangePasswordCard";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import { getApiErrorMessage } from "@/lib/api-error";

import DashboardIntro from "./components/DashboardIntro";
import PetFormCard from "./components/PetFormCard";
import PhotoUploadCard from "./components/PhotoUploadCard";

export default function OwnerDashboardPage() {
  const queryClient = useQueryClient();
  const { data, isLoading } = useUserProfile();
  const user = data?.user;

  const [selectedPhoto, setSelectedPhoto] = useState<File | null>(null);
  const [petName, setPetName] = useState("");
  const [petSpeciesId, setPetSpeciesId] = useState("");
  const [petAge, setPetAge] = useState("");
  const [petWeight, setPetWeight] = useState("");
  const [petNotes, setPetNotes] = useState("");

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
      setPetName("");
      setPetSpeciesId("");
      setPetAge("");
      setPetWeight("");
      setPetNotes("");
    },
  });

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

    return getApiErrorMessage(uploadError, "Não foi possível enviar a foto.");
  }, [uploadError]);

  const petFeedback = useMemo(() => {
    if (!addPetError) {
      return "";
    }

    return getApiErrorMessage(addPetError, "Não foi possível cadastrar o pet.");
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

      <DashboardIntro />

      <PetFormCard
        petName={petName}
        petSpeciesId={petSpeciesId}
        petAge={petAge}
        petWeight={petWeight}
        petNotes={petNotes}
        onPetNameChange={setPetName}
        onPetSpeciesIdChange={setPetSpeciesId}
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
        specieOptions={specieOptions}
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
    </PageWrapper>
  );
}
