"use client";

import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import * as Form from "@radix-ui/react-form";
import { isAxiosError } from "axios";

import {
  type ProblemDetailsResponse,
  useAuthSession,
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

  const handleDeactivate = async () => {
    await deactivateUser();
    clearSession();
    router.replace("/login");
  };

  const handleReactivate = async () => {
    await reactivateUser();
  };

  return (
    <PageWrapper className="gap-10">
      <MainNav />

      <section className="rounded-4xl bg-white p-6 shadow-sm">
        <div className="mb-6">
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
          <button
            type="button"
            onClick={handleDeactivate}
            disabled={isDeactivating}
            className="rounded-full border border-rose-200 px-4 py-2 text-sm font-medium text-rose-600 transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
          >
            {isDeactivating ? "Desativando..." : "Desativar conta"}
          </button>
          <button
            type="button"
            onClick={handleReactivate}
            disabled={isReactivating}
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
  const {
    mutateAsync: updateUser,
    isPending: isSaving,
    error: updateError,
    isSuccess: updateSuccess,
  } = useUserUpdate();

  const [name, setName] = useState(user.name ?? "");
  const [email] = useState(user.login?.email ?? "");
  const [phoneCountryCode, setPhoneCountryCode] = useState(
    user.phone?.countryCode ?? "",
  );
  const [phoneAreaCode, setPhoneAreaCode] = useState(
    user.phone?.areaCode ?? "",
  );
  const [phoneNumber, setPhoneNumber] = useState(user.phone?.number ?? "");

  const [street, setStreet] = useState(user.address?.street ?? "");
  const [addressNumber, setAddressNumber] = useState(
    user.address?.number ?? "",
  );
  const [neighborhood, setNeighborhood] = useState(
    user.address?.neighborhood ?? "",
  );
  const [city, setCity] = useState(user.address?.city ?? "");
  const [zipCode, setZipCode] = useState(user.address?.zipCode ?? "");
  const [state, setState] = useState(user.address?.state ?? "");
  const [country, setCountry] = useState(user.address?.country ?? "");
  const [complement, setComplement] = useState(user.address?.complement ?? "");
  const [latitude, setLatitude] = useState(
    user.address?.location?.latitude !== undefined
      ? String(user.address.location.latitude)
      : "",
  );
  const [longitude, setLongitude] = useState(
    user.address?.location?.longitude !== undefined
      ? String(user.address.location.longitude)
      : "",
  );

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
        disabled={isSaving}
        className="inline-flex h-11 w-full items-center justify-center rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-4 text-sm font-semibold text-white shadow-lg shadow-cyan-200 transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
      >
        {isSaving ? "Salvando..." : "Salvar alterações"}
      </button>
    </Form.Root>
  );
};
