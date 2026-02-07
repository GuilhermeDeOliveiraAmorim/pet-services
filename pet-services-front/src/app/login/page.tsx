"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

import { useAuthLogin, useAuthSession } from "@/application";
import MainNav from "@/components/common/MainNav";

export default function LoginPage() {
  const router = useRouter();
  const { setSession } = useAuthSession();
  const { mutateAsync, isPending, error } = useAuthLogin();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const response = await mutateAsync({ email, password });
    const expiresAt = Date.now() + response.expiresIn * 1000;

    setSession({
      accessToken: response.accessToken,
      refreshToken: response.refreshToken,
      expiresAt,
    });

    router.replace("/");
  };

  return (
    <div className="min-h-screen bg-[#f7f9ff] px-4 text-slate-900">
      <div className="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-10 py-12">
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
              Bem-vindo de volta
            </h1>
            <p className="text-sm leading-6 text-slate-600">
              Gerencie seus pets, acompanhe consultas e mantenha o histórico
              sempre atualizado em um só lugar.
            </p>
            <div className="relative mt-4 rounded-4xl bg-linear-to-br from-teal-50 to-cyan-50 p-6">
              <div className="absolute -right-6 -top-6 h-16 w-16 rounded-full bg-pink-200/60 blur-xl" />
              <div className="rounded-3xl bg-white p-4 shadow-sm">
                <p className="text-xs font-semibold text-teal-500">
                  Dica do dia
                </p>
                <p className="mt-2 text-sm text-slate-600">
                  Mantenha a carteirinha de vacinação em dia para garantir
                  proteção completa.
                </p>
              </div>
              <div className="mt-4 flex items-center justify-between rounded-2xl bg-white/80 px-4 py-3 text-sm text-slate-600">
                <span>Pets ativos</span>
                <span className="font-semibold text-slate-900">85k+</span>
              </div>
            </div>
          </div>

          <div className="flex items-center justify-center">
            <div className="w-full max-w-md rounded-4xl bg-white p-8 shadow-[0_30px_80px_rgba(124,139,255,0.15)]">
              <div className="mb-6 flex items-center justify-between">
                <div>
                  <h2 className="text-2xl font-semibold">Login</h2>
                  <p className="mt-2 text-sm text-slate-500">
                    Acesse sua conta para continuar.
                  </p>
                </div>
                <div className="flex h-12 w-12 items-center justify-center rounded-3xl bg-linear-to-br from-cyan-100 to-blue-100 text-sm font-semibold text-slate-700">
                  24/7
                </div>
              </div>

              <form className="space-y-5" onSubmit={handleSubmit}>
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
                    autoComplete="email"
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
                    autoComplete="current-password"
                  />
                </div>

                <div className="flex items-center justify-between text-xs text-slate-500">
                  <label className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      className="h-4 w-4 rounded border-slate-300"
                    />
                    Lembrar de mim
                  </label>
                  <button type="button" className="text-cyan-600">
                    Esqueci minha senha
                  </button>
                </div>

                {error ? (
                  <p className="text-sm text-rose-500">
                    Não foi possível fazer login. Verifique suas credenciais.
                  </p>
                ) : null}

                <button
                  type="submit"
                  disabled={isPending}
                  className="inline-flex h-11 w-full items-center justify-center rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-4 text-sm font-semibold text-white shadow-lg shadow-cyan-200 transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
                >
                  {isPending ? "Entrando..." : "Entrar"}
                </button>

                <p className="text-center text-xs text-slate-500">
                  Ainda não tem conta?{" "}
                  <span className="text-cyan-600">Criar conta</span>
                </p>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
