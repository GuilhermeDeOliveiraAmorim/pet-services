"use client";

import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import * as Checkbox from "@radix-ui/react-checkbox";
import * as Form from "@radix-ui/react-form";
import * as Toggle from "@radix-ui/react-toggle";
import { Check, Eye, EyeOff } from "lucide-react";
import { isAxiosError } from "axios";

import {
  type ProblemDetailsResponse,
  useAuthLogin,
  useAuthResendVerificationEmail,
  useAuthSession,
} from "@/application";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

export default function LoginPage() {
  const router = useRouter();
  const { setSession } = useAuthSession();
  const { mutateAsync, isPending, error } = useAuthLogin();
  const {
    mutateAsync: resendVerification,
    isPending: isResendingVerification,
    data: resendVerificationResult,
  } = useAuthResendVerificationEmail();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);

  const feedback = useMemo(() => {
    if (!error) {
      return { message: "", canResend: false };
    }

    if (isAxiosError<ProblemDetailsResponse>(error)) {
      const errorPayload = error.response?.data;
      const problem = errorPayload?.errors?.[0];
      const title = problem?.title?.toLowerCase() ?? "";
      const detail = problem?.detail ?? "";
      const status = problem?.status ?? error.response?.status;

      const isUnverifiedEmail =
        status === 403 &&
        (title.includes("email não verificado") ||
          detail.toLowerCase().includes("email") ||
          detail.toLowerCase().includes("verifica"));

      return {
        message:
          detail || "Não foi possível fazer login. Verifique suas credenciais.",
        canResend: isUnverifiedEmail,
      };
    }

    return {
      message: "Não foi possível fazer login. Verifique suas credenciais.",
      canResend: false,
    };
  }, [error]);

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
            Bem-vindo de volta
          </h1>
          <p className="text-sm leading-6 text-slate-600">
            Gerencie seus pets, acompanhe consultas e mantenha o histórico
            sempre atualizado em um só lugar.
          </p>
          <div className="relative mt-4 rounded-4xl bg-linear-to-br from-teal-50 to-cyan-50 p-6">
            <div className="absolute -right-6 -top-6 h-16 w-16 rounded-full bg-pink-200/60 blur-xl" />
            <div className="rounded-3xl bg-white p-4 shadow-sm">
              <p className="text-xs font-semibold text-teal-500">Dica do dia</p>
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
          <div className="w-full rounded-4xl bg-white p-8 shadow-[0_30px_80px_rgba(124,139,255,0.15)]">
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

            <Form.Root className="space-y-5" onSubmit={handleSubmit}>
              <Form.Field className="space-y-2" name="email">
                <div className="flex items-baseline justify-between">
                  <Form.Label className="text-sm font-medium">Email</Form.Label>
                  <Form.Message
                    className="text-xs text-rose-500"
                    match="valueMissing"
                  >
                    Informe o email
                  </Form.Message>
                  <Form.Message
                    className="text-xs text-rose-500"
                    match="typeMismatch"
                  >
                    Email inválido
                  </Form.Message>
                </div>
                <Form.Control asChild>
                  <input
                    type="email"
                    value={email}
                    onChange={(event) => setEmail(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="voce@email.com"
                    required
                    autoComplete="email"
                  />
                </Form.Control>
              </Form.Field>

              <Form.Field className="space-y-2" name="password">
                <div className="flex items-baseline justify-between">
                  <Form.Label className="text-sm font-medium">Senha</Form.Label>
                  <Form.Message
                    className="text-xs text-rose-500"
                    match="valueMissing"
                  >
                    Informe a senha
                  </Form.Message>
                </div>
                <div className="relative">
                  <Form.Control asChild>
                    <input
                      type={showPassword ? "text" : "password"}
                      value={password}
                      onChange={(event) => setPassword(event.target.value)}
                      className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 pr-12 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                      placeholder="********"
                      required
                      autoComplete="current-password"
                    />
                  </Form.Control>
                  <Toggle.Root
                    pressed={showPassword}
                    onPressedChange={setShowPassword}
                    aria-label={
                      showPassword ? "Ocultar senha" : "Mostrar senha"
                    }
                    className="absolute inset-y-0 right-2 flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:text-slate-600"
                  >
                    {showPassword ? (
                      <EyeOff className="h-4 w-4" />
                    ) : (
                      <Eye className="h-4 w-4" />
                    )}
                  </Toggle.Root>
                </div>
              </Form.Field>

              <div className="flex items-center justify-between text-xs text-slate-500">
                <label className="flex items-center gap-2" htmlFor="remember">
                  <Checkbox.Root
                    id="remember"
                    className="flex h-4 w-4 items-center justify-center rounded border border-slate-300 bg-white shadow-sm data-[state=checked]:bg-cyan-500 data-[state=checked]:border-cyan-500"
                  >
                    <Checkbox.Indicator className="text-white">
                      <Check className="h-3 w-3" />
                    </Checkbox.Indicator>
                  </Checkbox.Root>
                  Lembrar de mim
                </label>
                <button type="button" className="text-cyan-600">
                  Esqueci minha senha
                </button>
              </div>

              {error ? (
                <div className="rounded-3xl border border-rose-200/80 bg-rose-50/70 px-4 py-3 text-sm text-rose-600 shadow-sm">
                  <div className="flex items-start gap-3">
                    <span className="mt-0.5 inline-flex h-6 w-6 items-center justify-center rounded-full bg-rose-100 text-xs font-semibold text-rose-600">
                      !
                    </span>
                    <div className="space-y-2">
                      <p className="leading-5 text-rose-600">
                        {feedback.message}
                      </p>
                      {feedback.canResend ? (
                        <button
                          type="button"
                          disabled={!email || isResendingVerification}
                          onClick={() => resendVerification({ email })}
                          className="inline-flex items-center justify-center rounded-full border border-rose-200 bg-white px-3 py-1 text-xs font-semibold text-rose-600 transition hover:bg-rose-100 disabled:cursor-not-allowed disabled:opacity-60"
                        >
                          {isResendingVerification
                            ? "Reenviando..."
                            : "Reenviar email de verificação"}
                        </button>
                      ) : null}
                    </div>
                  </div>
                </div>
              ) : null}

              {resendVerificationResult?.message ? (
                <div className="rounded-3xl border border-emerald-200/80 bg-emerald-50/70 px-4 py-3 text-xs text-emerald-600 shadow-sm">
                  <div className="flex items-start gap-3">
                    <span className="mt-0.5 inline-flex h-6 w-6 items-center justify-center rounded-full bg-emerald-100 text-xs font-semibold text-emerald-600">
                      ✓
                    </span>
                    <p className="leading-5">
                      {resendVerificationResult.detail ??
                        resendVerificationResult.message}
                    </p>
                  </div>
                </div>
              ) : null}

              <Form.Submit asChild>
                <button
                  type="submit"
                  disabled={isPending}
                  className="inline-flex h-11 w-full items-center justify-center rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-4 text-sm font-semibold text-white shadow-lg shadow-cyan-200 transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
                >
                  {isPending ? "Entrando..." : "Entrar"}
                </button>
              </Form.Submit>

              <p className="text-center text-xs text-slate-500">
                Ainda não tem conta?{" "}
                <Link href="/register" className="text-cyan-600">
                  Criar conta
                </Link>
              </p>
            </Form.Root>
          </div>
        </div>
      </div>
    </PageWrapper>
  );
}
