"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

import { type UserType, UserTypes } from "@/domain";
import { useUserRegister } from "@/application";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

const defaultCountryCode = "55";
const defaultCountry = "Brasil";

export default function RegisterPage() {
  const router = useRouter();
  const { mutateAsync, isPending, error, isSuccess } = useUserRegister();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [userType, setUserType] = useState<UserType>(UserTypes.Owner);
  const [countryCode, setCountryCode] = useState(defaultCountryCode);
  const [areaCode, setAreaCode] = useState("");
  const [phoneNumber, setPhoneNumber] = useState("");

  const [street, setStreet] = useState("");
  const [addressNumber, setAddressNumber] = useState("");
  const [neighborhood, setNeighborhood] = useState("");
  const [city, setCity] = useState("");
  const [state, setState] = useState("");
  const [zipCode, setZipCode] = useState("");
  const [country, setCountry] = useState(defaultCountry);
  const [complement, setComplement] = useState("");
  const [latitude, setLatitude] = useState("0");
  const [longitude, setLongitude] = useState("0");

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    await mutateAsync({
      name,
      userType,
      login: {
        email,
        password,
      },
      phone: {
        countryCode,
        areaCode,
        number: phoneNumber,
      },
      address: {
        street,
        number: addressNumber,
        neighborhood,
        city,
        zipCode,
        state,
        country,
        complement,
        location: {
          latitude: Number(latitude),
          longitude: Number(longitude),
        },
      },
    });

    router.replace("/login");
  };

  return (
    <PageWrapper className="gap-16">
      <MainNav showLinks={false} showActions={false} />
      <div className="grid w-full gap-10 lg:grid-cols-[1.1fr_0.9fr]">
        <div className="hidden flex-col justify-center gap-6 rounded-4xl bg-white p-10 shadow-[0_30px_80px_rgba(124,139,255,0.2)] lg:flex">
          <div className="flex items-center gap-2">
            <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-linear-to-tr from-teal-400 to-cyan-400 text-white font-semibold">
              pet
            </div>
            <span className="text-lg font-semibold">PetCare</span>
          </div>
          <h1 className="text-3xl font-semibold leading-tight">
            Crie sua conta
          </h1>
          <p className="text-sm leading-6 text-slate-600">
            Cadastre seus dados para acessar a plataforma e acompanhar todos os
            serviços do seu pet.
          </p>
          <div className="rounded-4xl bg-linear-to-br from-teal-50 to-cyan-50 p-6">
            <p className="text-xs font-semibold text-teal-500">
              Segurança em primeiro lugar
            </p>
            <p className="mt-2 text-sm text-slate-600">
              Seus dados são protegidos e utilizados apenas para melhorar a sua
              experiência.
            </p>
          </div>
        </div>

        <div className="flex items-center justify-center">
          <div className="w-full rounded-4xl bg-white p-8 shadow-[0_30px_80px_rgba(124,139,255,0.15)]">
            <div className="mb-6 flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-semibold">Cadastro</h2>
                <p className="mt-2 text-sm text-slate-500">
                  Preencha os campos para criar sua conta.
                </p>
              </div>
              <div className="flex h-12 w-12 items-center justify-center rounded-3xl bg-linear-to-br from-cyan-100 to-blue-100 text-sm font-semibold text-slate-700">
                Novo
              </div>
            </div>

            <form className="space-y-6" onSubmit={handleSubmit}>
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="name">
                    Nome completo
                  </label>
                  <input
                    id="name"
                    value={name}
                    onChange={(event) => setName(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="Seu nome"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="userType">
                    Tipo de usuário
                  </label>
                  <select
                    id="userType"
                    value={userType}
                    onChange={(event) =>
                      setUserType(event.target.value as UserType)
                    }
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                  >
                    <option value={UserTypes.Owner}>Tutor</option>
                    <option value={UserTypes.Provider}>Prestador</option>
                  </select>
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="email">
                    Email
                  </label>
                  <input
                    id="email"
                    type="email"
                    value={email}
                    onChange={(event) => setEmail(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="voce@email.com"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="password">
                    Senha
                  </label>
                  <input
                    id="password"
                    type="password"
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="********"
                    required
                    autoComplete="new-password"
                  />
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-3">
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="countryCode">
                    DDI
                  </label>
                  <input
                    id="countryCode"
                    value={countryCode}
                    onChange={(event) => setCountryCode(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="55"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="areaCode">
                    DDD
                  </label>
                  <input
                    id="areaCode"
                    value={areaCode}
                    onChange={(event) => setAreaCode(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="11"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="phoneNumber">
                    Telefone
                  </label>
                  <input
                    id="phoneNumber"
                    value={phoneNumber}
                    onChange={(event) => setPhoneNumber(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="99999-9999"
                    required
                  />
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="street">
                    Rua
                  </label>
                  <input
                    id="street"
                    value={street}
                    onChange={(event) => setStreet(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="Rua Exemplo"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label
                    className="text-sm font-medium"
                    htmlFor="addressNumber"
                  >
                    Número
                  </label>
                  <input
                    id="addressNumber"
                    value={addressNumber}
                    onChange={(event) => setAddressNumber(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="123"
                    required
                  />
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="neighborhood">
                    Bairro
                  </label>
                  <input
                    id="neighborhood"
                    value={neighborhood}
                    onChange={(event) => setNeighborhood(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="Centro"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="city">
                    Cidade
                  </label>
                  <input
                    id="city"
                    value={city}
                    onChange={(event) => setCity(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="São Paulo"
                    required
                  />
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-3">
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="state">
                    Estado
                  </label>
                  <input
                    id="state"
                    value={state}
                    onChange={(event) => setState(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="SP"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="zipCode">
                    CEP
                  </label>
                  <input
                    id="zipCode"
                    value={zipCode}
                    onChange={(event) => setZipCode(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="00000-000"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="country">
                    País
                  </label>
                  <input
                    id="country"
                    value={country}
                    onChange={(event) => setCountry(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="Brasil"
                    required
                  />
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="complement">
                    Complemento
                  </label>
                  <input
                    id="complement"
                    value={complement}
                    onChange={(event) => setComplement(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="Apartamento, bloco, etc"
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="latitude">
                    Latitude
                  </label>
                  <input
                    id="latitude"
                    type="number"
                    step="0.000001"
                    value={latitude}
                    onChange={(event) => setLatitude(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="-23.550520"
                    required
                  />
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label className="text-sm font-medium" htmlFor="longitude">
                    Longitude
                  </label>
                  <input
                    id="longitude"
                    type="number"
                    step="0.000001"
                    value={longitude}
                    onChange={(event) => setLongitude(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="-46.633308"
                    required
                  />
                </div>
                <div className="flex items-end">
                  <button
                    type="submit"
                    disabled={isPending}
                    className="inline-flex h-11 w-full items-center justify-center rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-4 text-sm font-semibold text-white shadow-lg shadow-cyan-200 transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
                  >
                    {isPending ? "Criando..." : "Criar conta"}
                  </button>
                </div>
              </div>

              {error ? (
                <p className="text-sm text-rose-500">
                  Não foi possível criar sua conta. Verifique os dados.
                </p>
              ) : null}

              {isSuccess ? (
                <p className="text-sm text-emerald-600">
                  Cadastro realizado com sucesso! Redirecionando...
                </p>
              ) : null}

              <p className="text-center text-xs text-slate-500">
                Já tem conta?{" "}
                <Link href="/login" className="text-cyan-600">
                  Entrar
                </Link>
              </p>
            </form>
          </div>
        </div>
      </div>
    </PageWrapper>
  );
}
