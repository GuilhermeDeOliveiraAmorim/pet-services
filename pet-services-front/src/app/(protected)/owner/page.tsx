"use client";

import { useMemo, useState } from "react";
import Image from "next/image";
import { useQueryClient } from "@tanstack/react-query";
import { isAxiosError } from "axios";
import * as Form from "@radix-ui/react-form";

import {
  type ProblemDetailsResponse,
  usePetAdd,
  useSpeciesList,
  useUserAddPhoto,
  useUserProfile,
} from "@/application";
import ChangePasswordCard from "@/components/account/ChangePasswordCard";
import MainNav from "@/components/common/MainNav";
import RadixSelectField from "@/components/common/RadixSelectField";
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
    petName.trim() &&
    petSpeciesId.trim() &&
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

      <section className="flex flex-col gap-6">
        <div>
          <p className="text-xs font-semibold uppercase text-cyan-400">
            Dashboard
          </p>
          <h1 className="mt-2 text-2xl font-semibold text-slate-900">
            Olá, tutor
          </h1>
          <p className="mt-2 text-sm text-slate-600">
            Este é o seu painel inicial. Aqui vão aparecer seus pets e
            agendamentos.
          </p>
        </div>

        <div className="rounded-4xl border border-dashed border-slate-200 bg-white px-6 py-16 text-center">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-2xl bg-cyan-50 text-cyan-500">
            <span className="text-lg font-semibold">🐾</span>
          </div>
          <h2 className="mt-4 text-lg font-semibold text-slate-900">
            Sem dados ainda
          </h2>
          <p className="mt-2 text-sm text-slate-600">
            Quando você cadastrar seu primeiro pet ou agendar um serviço, as
            informações vão aparecer aqui.
          </p>
        </div>
      </section>

      <section className="rounded-4xl bg-white p-6 shadow-sm">
        <div className="mb-4">
          <p className="text-xs font-semibold uppercase text-cyan-400">Pets</p>
          <h2 className="mt-2 text-xl font-semibold text-slate-900">
            Cadastrar pet
          </h2>
          <p className="mt-2 text-sm text-slate-600">
            Adicione um novo pet ao seu perfil.
          </p>
        </div>

        <Form.Root onSubmit={handleAddPet} className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <label className="space-y-2 text-sm text-slate-700">
              <span className="font-medium">Nome</span>
              <input
                type="text"
                value={petName}
                onChange={(event) => setPetName(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="Ex: Thor"
                required
              />
            </label>
            <RadixSelectField
              name="specie"
              label="Espécie"
              value={petSpeciesId}
              onValueChange={setPetSpeciesId}
              options={specieOptions}
              placeholder={
                isLoadingSpecies
                  ? "Carregando espécies..."
                  : "Selecione a espécie"
              }
              required
              disabled={isLoadingSpecies || Boolean(speciesError)}
              searchable
            />
            <label className="space-y-2 text-sm text-slate-700">
              <span className="font-medium">Idade</span>
              <input
                type="number"
                inputMode="numeric"
                value={petAge}
                onChange={(event) => setPetAge(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="Ex: 3"
                min={0}
                required
              />
            </label>
            <label className="space-y-2 text-sm text-slate-700">
              <span className="font-medium">Peso (kg)</span>
              <input
                type="number"
                inputMode="decimal"
                value={petWeight}
                onChange={(event) => setPetWeight(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="Ex: 12.5"
                min={0}
                step="0.1"
                required
              />
            </label>
          </div>

          <label className="space-y-2 text-sm text-slate-700">
            <span className="font-medium">Observações</span>
            <textarea
              value={petNotes}
              onChange={(event) => setPetNotes(event.target.value)}
              className="min-h-24 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
              placeholder="Comportamento, cuidados especiais, etc."
            />
          </label>

          <div className="flex flex-wrap items-center gap-3">
            <button
              type="submit"
              disabled={!isPetFormValid || isAddingPet}
              className="rounded-full bg-cyan-500 px-4 py-2 text-xs font-semibold text-white transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
            >
              {isAddingPet ? "Cadastrando..." : "Cadastrar pet"}
            </button>
            {addPetSuccess ? (
              <p className="text-xs text-emerald-600">
                Pet cadastrado com sucesso.
              </p>
            ) : null}
          </div>

          {petFeedback ? (
            <p className="text-xs text-rose-600">{petFeedback}</p>
          ) : null}
          {speciesError ? (
            <p className="text-xs text-rose-600">
              Não foi possível carregar as espécies.
            </p>
          ) : null}
        </Form.Root>
      </section>

      <section className="rounded-4xl bg-white p-6 shadow-sm">
        <div className="mb-4">
          <p className="text-xs font-semibold uppercase text-cyan-400">Fotos</p>
          <h2 className="mt-2 text-xl font-semibold text-slate-900">
            Fotos do perfil
          </h2>
          <p className="mt-2 text-sm text-slate-600">
            Envie uma foto para aparecer no seu perfil.
          </p>
        </div>

        {isLoading ? (
          <p className="text-sm text-slate-500">Carregando fotos...</p>
        ) : (
          <div className="space-y-4">
            <div className="flex flex-wrap items-center gap-3">
              <input
                type="file"
                accept="image/*"
                onChange={(event) =>
                  setSelectedPhoto(event.target.files?.[0] ?? null)
                }
                className="text-sm text-slate-600"
              />
              <button
                type="button"
                onClick={handlePhotoUpload}
                disabled={!selectedPhoto || isUploadingPhoto}
                className="rounded-full bg-cyan-500 px-4 py-2 text-xs font-semibold text-white transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
              >
                {isUploadingPhoto ? "Enviando..." : "Enviar foto"}
              </button>
            </div>

            {uploadSuccess ? (
              <p className="text-xs text-emerald-600">
                Foto enviada com sucesso.
              </p>
            ) : null}

            {photoFeedback ? (
              <p className="text-xs text-rose-600">{photoFeedback}</p>
            ) : null}

            <div className="grid gap-3 sm:grid-cols-3 lg:grid-cols-4">
              {user?.photos?.length ? (
                user.photos.map((photo) => (
                  <div
                    key={photo.id}
                    className="h-28 w-full overflow-hidden rounded-2xl bg-slate-100"
                  >
                    <Image
                      src={photo.url}
                      alt="Foto do usuário"
                      width={112}
                      height={112}
                      className="h-full w-full object-cover"
                      unoptimized
                    />
                  </div>
                ))
              ) : (
                <p className="text-sm text-slate-500">
                  Nenhuma foto enviada ainda.
                </p>
              )}
            </div>
          </div>
        )}
      </section>

      <ChangePasswordCard />
    </PageWrapper>
  );
}
