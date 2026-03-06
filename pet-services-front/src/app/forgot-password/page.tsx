"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import * as Form from "@radix-ui/react-form";
import { isAxiosError } from "axios";

import {
  type ProblemDetailsResponse,
  useAuthRequestPasswordReset,
} from "@/application";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

export default function ForgotPasswordPage() {
  const { mutateAsync, isPending, error, isSuccess } =
    useAuthRequestPasswordReset();
  const [email, setEmail] = useState("");

  const feedback = useMemo(() => {
    if (!error) {
      return "";
    }

    if (isAxiosError<ProblemDetailsResponse>(error)) {
      const problem = error.response?.data?.errors?.[0];
      return (
        problem?.detail ||
        problem?.title ||
        "Não foi possível solicitar a redefinição de senha."
      );
    }

    return "Não foi possível solicitar a redefinição de senha.";
  }, [error]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    await mutateAsync({ email });
  };

  return (
    <PageWrapper gap={16}>
      <MainNav showLinks={false} showActions={false} />

      <div className="mx-auto w-full max-w-xl rounded-4xl bg-white p-8 shadow-[0_30px_80px_rgba(124,139,255,0.15)]">
        <div className="mb-6">
          <h1 className="text-2xl font-semibold">Redefinir senha</h1>
          <p className="mt-2 text-sm text-slate-500">
            Enviaremos um link para redefinir sua senha.
          </p>
        </div>

        <Form.Root className="space-y-5" onSubmit={handleSubmit}>
          <Form.Field className="space-y-2" name="email">
            <div className="flex items-baseline justify-between">
              <Form.Label className="text-sm font-medium">Email</Form.Label>
              <Form.Message className="text-xs text-rose-500" match="valueMissing">
                Informe o email
              </Form.Message>
              <Form.Message className="text-xs text-rose-500" match="typeMismatch">
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

          {feedback ? (
            <div className="rounded-3xl border border-rose-200/80 bg-rose-50/70 px-4 py-3 text-sm text-rose-600 shadow-sm">
              {feedback}
            </div>
          ) : null}

          {isSuccess ? (
            <div className="rounded-3xl border border-emerald-200/80 bg-emerald-50/70 px-4 py-3 text-sm text-emerald-700 shadow-sm">
              Se o email existir, enviaremos instruções para redefinição.
            </div>
          ) : null}

          <button
            type="submit"
            disabled={isPending}
            className="inline-flex h-11 w-full items-center justify-center rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-4 text-sm font-semibold text-white shadow-lg shadow-cyan-200 transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
          >
            {isPending ? "Enviando..." : "Enviar link"}
          </button>
        </Form.Root>

        <div className="mt-6 text-center text-xs text-slate-500">
          <Link href="/login" className="text-cyan-600">
            Voltar para o login
          </Link>
        </div>
      </div>
    </PageWrapper>
  );
}
