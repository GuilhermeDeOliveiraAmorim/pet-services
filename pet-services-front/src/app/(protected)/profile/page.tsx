"use client";

import { useMemo, useState } from "react";
import Image from "next/image";
import { useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import * as Dialog from "@radix-ui/react-dialog";
import * as Form from "@radix-ui/react-form";
import { isAxiosError } from "axios";

import {
  type ProblemDetailsResponse,
  useAuthSession,
  useUserAddPhoto,
  useUserDeactivate,
  useUserProfile,
  useUserReactivate,
  useUserUpdate,
} from "@/application";
import type { UpdateUserInput } from "@/application";
import ChangePasswordCard from "@/components/account/ChangePasswordCard";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

export default function ProfilePage() {
  const router = useRouter();
  const { clearSession } = useAuthSession();
  const { data, isLoading } = useUserProfile();
  const [isDeactivateOpen, setIsDeactivateOpen] = useState(false);
  const {
    mutateAsync: deactivateUser,
    isPending: isDeactivating,
    error: deactivateError,
  } = useUserDeactivate();
  const {
    mutateAsync: reactivateUser,
    isPending: isReactivating,
    error: reactivateError,
  } = useUserReactivate();

  const user = data?.user;

  const accountFeedback = useMemo(() => {
    const error = deactivateError ?? reactivateError;
    if (!error) {
      return "";
    }

    if (isAxiosError<ProblemDetailsResponse>(error)) {
      const problem = error.response?.data?.errors?.[0];
      return (
        problem?.detail ||
        problem?.title ||
        "Não foi possível atualizar o status da conta."
      );
    }

    return "Não foi possível atualizar o status da conta.";
  }, [deactivateError, reactivateError]);

  const handleConfirmDeactivate = async () => {
    await deactivateUser();
    clearSession();
    router.replace("/login");
  };

  const handleReactivate = async () => {
    await reactivateUser();
  };

  return (
    <PageWrapper gap={10}>
      <MainNav />

      <section className="rounded-4xl bg-white p-6 shadow-sm">
        <div className="mb-6 flex flex-wrap items-start justify-between gap-4">
          <div>
            <p className="text-xs font-semibold uppercase text-cyan-400">
              Perfil
            </p>
            <h1 className="mt-2 text-2xl font-semibold text-slate-900">
              Seus dados
            </h1>
            <p className="mt-2 text-sm text-slate-600">
              Atualize suas informações de contato.
            </p>
          </div>
          {user ? (
            <span
              className={`inline-flex items-center gap-2 rounded-full px-3 py-1 text-xs font-semibold ${
                user.active
                  ? "bg-emerald-50 text-emerald-600"
                  : "bg-amber-50 text-amber-600"
              }`}
            >
              {user.active ? "Conta ativa" : "Conta desativada"}
            </span>
          ) : null}
        </div>

        {isLoading ? (
          <p className="text-sm text-slate-500">Carregando perfil...</p>
        ) : user ? (
          <ProfileForm key={user.id} user={user} />
        ) : (
          <p className="text-sm text-slate-500">
            Não foi possível carregar o perfil.
          </p>
        )}
      </section>

      <ChangePasswordCard />

      <section className="rounded-4xl border border-rose-200/70 bg-white p-6 shadow-sm">
        <div className="mb-4">
          <p className="text-xs font-semibold uppercase text-rose-400">Conta</p>
          <h2 className="mt-2 text-xl font-semibold text-slate-900">
            Status da conta
          </h2>
          <p className="mt-2 text-sm text-slate-600">
            Você pode desativar sua conta e reativá-la depois.
          </p>
        </div>

        {accountFeedback ? (
          <div className="mb-4 rounded-3xl border border-rose-200/80 bg-rose-50/70 px-4 py-3 text-sm text-rose-600 shadow-sm">
            {accountFeedback}
          </div>
        ) : null}

        <div className="flex flex-wrap gap-3">
          <Dialog.Root
            open={isDeactivateOpen}
            onOpenChange={setIsDeactivateOpen}
          >
            <Dialog.Trigger asChild>
              <button
                type="button"
                disabled={isDeactivating || !user || !user.active}
                className="rounded-full border border-rose-200 px-4 py-2 text-sm font-medium text-rose-600 transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
              >
                {isDeactivating ? "Desativando..." : "Desativar conta"}
              </button>
            </Dialog.Trigger>
            <Dialog.Portal>
              <Dialog.Overlay className="fixed inset-0 bg-slate-900/30" />
              <Dialog.Content className="fixed left-1/2 top-1/2 w-[92vw] max-w-md -translate-x-1/2 -translate-y-1/2 rounded-3xl bg-white p-6 shadow-xl">
                <Dialog.Title className="text-lg font-semibold text-slate-900">
                  Desativar conta?
                </Dialog.Title>
                <Dialog.Description className="mt-2 text-sm text-slate-600">
                  Ao desativar sua conta, você será desconectado e precisará
                  reativá-la para acessar novamente.
                </Dialog.Description>

                <div className="mt-6 flex flex-wrap justify-end gap-3">
                  <Dialog.Close asChild>
                    <button
                      type="button"
                      className="rounded-full border border-slate-200 px-4 py-2 text-sm font-medium text-slate-700"
                    >
                      Cancelar
                    </button>
                  </Dialog.Close>
                  <button
                    type="button"
                    onClick={handleConfirmDeactivate}
                    disabled={isDeactivating}
                    className="rounded-full bg-rose-500 px-4 py-2 text-sm font-semibold text-white transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
                  >
                    {isDeactivating ? "Desativando..." : "Confirmar"}
                  </button>
                </div>
              </Dialog.Content>
            </Dialog.Portal>
          </Dialog.Root>
          <button
            type="button"
            onClick={handleReactivate}
            disabled={isReactivating || !user || user.active}
            className="rounded-full border border-slate-200 px-4 py-2 text-sm font-medium text-slate-700 transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
          >
            {isReactivating ? "Reativando..." : "Reativar conta"}
          </button>
        </div>
      </section>
    </PageWrapper>
  );
}

type ProfileFormProps = {
  user: NonNullable<ReturnType<typeof useUserProfile>["data"]>["user"];
};

const ProfileForm = ({ user }: ProfileFormProps) => {
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

    if (isAxiosError<ProblemDetailsResponse>(updateError)) {
      const problem = updateError.response?.data?.errors?.[0];
      return (
        problem?.detail ||
        problem?.title ||
        "Não foi possível atualizar o perfil."
      );
    }

    return "Não foi possível atualizar o perfil.";
  }, [updateError]);

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
    <Form.Root className="space-y-6" onSubmit={handleSubmit}>
      <div className="rounded-3xl border border-slate-100 bg-slate-50/60 p-4">
        <div className="flex flex-wrap items-center gap-4">
          <div className="h-20 w-20 overflow-hidden rounded-3xl bg-slate-200">
            {currentPhotoUrl ? (
              <Image
                src={currentPhotoUrl}
                alt="Foto do usuário"
                width={80}
                height={80}
                className="h-full w-full object-cover"
                unoptimized
              />
            ) : null}
          </div>
          <div className="flex flex-1 flex-col gap-2">
            <p className="text-sm font-medium text-slate-800">Foto do perfil</p>
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
          </div>
        </div>
      </div>
      <div className="grid gap-4 md:grid-cols-2">
        <Form.Field className="space-y-2" name="name">
          <Form.Label className="text-sm font-medium">Nome</Form.Label>
          <Form.Control asChild>
            <input
              value={name}
              onChange={(event) => setName(event.target.value)}
              className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
              placeholder="Seu nome"
              required
            />
          </Form.Control>
        </Form.Field>

        <Form.Field className="space-y-2" name="email">
          <Form.Label className="text-sm font-medium">Email</Form.Label>
          <Form.Control asChild>
            <input
              value={email}
              readOnly
              className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-100 px-4 text-sm text-slate-500"
            />
          </Form.Control>
          <p className="text-xs text-slate-400">
            O email é fixo e não pode ser alterado aqui.
          </p>
        </Form.Field>
      </div>

      <div>
        <p className="text-sm font-semibold text-slate-900">Telefone</p>
        <div className="mt-3 grid gap-4 md:grid-cols-3">
          <Form.Field className="space-y-2" name="countryCode">
            <Form.Label className="text-sm font-medium">DDI</Form.Label>
            <Form.Control asChild>
              <input
                value={phoneCountryCode}
                onChange={(event) => setPhoneCountryCode(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="55"
              />
            </Form.Control>
          </Form.Field>

          <Form.Field className="space-y-2" name="areaCode">
            <Form.Label className="text-sm font-medium">DDD</Form.Label>
            <Form.Control asChild>
              <input
                value={phoneAreaCode}
                onChange={(event) => setPhoneAreaCode(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="82"
              />
            </Form.Control>
          </Form.Field>

          <Form.Field className="space-y-2" name="phoneNumber">
            <Form.Label className="text-sm font-medium">Número</Form.Label>
            <Form.Control asChild>
              <input
                value={phoneNumber}
                onChange={(event) => setPhoneNumber(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="999999999"
              />
            </Form.Control>
          </Form.Field>
        </div>
      </div>

      <div>
        <p className="text-sm font-semibold text-slate-900">Endereço</p>
        <div className="mt-3 grid gap-4 md:grid-cols-2">
          <Form.Field className="space-y-2" name="street">
            <Form.Label className="text-sm font-medium">Rua</Form.Label>
            <Form.Control asChild>
              <input
                value={street}
                onChange={(event) => setStreet(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="Rua"
              />
            </Form.Control>
          </Form.Field>

          <Form.Field className="space-y-2" name="addressNumber">
            <Form.Label className="text-sm font-medium">Número</Form.Label>
            <Form.Control asChild>
              <input
                value={addressNumber}
                onChange={(event) => setAddressNumber(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="123"
              />
            </Form.Control>
          </Form.Field>

          <Form.Field className="space-y-2" name="neighborhood">
            <Form.Label className="text-sm font-medium">Bairro</Form.Label>
            <Form.Control asChild>
              <input
                value={neighborhood}
                onChange={(event) => setNeighborhood(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="Centro"
              />
            </Form.Control>
          </Form.Field>

          <Form.Field className="space-y-2" name="city">
            <Form.Label className="text-sm font-medium">Cidade</Form.Label>
            <Form.Control asChild>
              <input
                value={city}
                onChange={(event) => setCity(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="Maceió"
              />
            </Form.Control>
          </Form.Field>

          <Form.Field className="space-y-2" name="zipCode">
            <Form.Label className="text-sm font-medium">CEP</Form.Label>
            <Form.Control asChild>
              <input
                value={zipCode}
                onChange={(event) => setZipCode(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="00000-000"
              />
            </Form.Control>
          </Form.Field>

          <Form.Field className="space-y-2" name="state">
            <Form.Label className="text-sm font-medium">Estado</Form.Label>
            <Form.Control asChild>
              <input
                value={state}
                onChange={(event) => setState(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="AL"
              />
            </Form.Control>
          </Form.Field>

          <Form.Field className="space-y-2" name="country">
            <Form.Label className="text-sm font-medium">País</Form.Label>
            <Form.Control asChild>
              <input
                value={country}
                onChange={(event) => setCountry(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="Brasil"
              />
            </Form.Control>
          </Form.Field>

          <Form.Field className="space-y-2" name="complement">
            <Form.Label className="text-sm font-medium">Complemento</Form.Label>
            <Form.Control asChild>
              <input
                value={complement}
                onChange={(event) => setComplement(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="Apartamento"
              />
            </Form.Control>
          </Form.Field>

          <Form.Field className="space-y-2" name="latitude">
            <Form.Label className="text-sm font-medium">Latitude</Form.Label>
            <Form.Control asChild>
              <input
                value={latitude}
                onChange={(event) => setLatitude(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="-10.000"
              />
            </Form.Control>
          </Form.Field>

          <Form.Field className="space-y-2" name="longitude">
            <Form.Label className="text-sm font-medium">Longitude</Form.Label>
            <Form.Control asChild>
              <input
                value={longitude}
                onChange={(event) => setLongitude(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="-37.000"
              />
            </Form.Control>
          </Form.Field>
        </div>
      </div>

      {feedback ? (
        <div className="rounded-3xl border border-rose-200/80 bg-rose-50/70 px-4 py-3 text-sm text-rose-600 shadow-sm">
          {feedback}
        </div>
      ) : null}

      {updateSuccess ? (
        <div className="rounded-3xl border border-emerald-200/80 bg-emerald-50/70 px-4 py-3 text-sm text-emerald-700 shadow-sm">
          Perfil atualizado com sucesso.
        </div>
      ) : null}

      <button
        type="submit"
        disabled={isSaving || !hasChanges}
        className="inline-flex h-11 w-full items-center justify-center rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-4 text-sm font-semibold text-white shadow-lg shadow-cyan-200 transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
      >
        {isSaving ? "Salvando..." : "Salvar alterações"}
      </button>
    </Form.Root>
  );
};
