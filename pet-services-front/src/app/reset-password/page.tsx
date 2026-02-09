"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import * as Form from "@radix-ui/react-form";
import * as Toggle from "@radix-ui/react-toggle";
import { Eye, EyeOff } from "lucide-react";
import { isAxiosError } from "axios";

import {
  type ProblemDetailsResponse,
  useAuthResetPassword,
} from "@/application";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

export default function ResetPasswordPage() {
  const searchParams = useSearchParams();
  const token = searchParams.get("token") ?? "";
  const { mutateAsync, isPending, error, isSuccess } = useAuthResetPassword();

  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showNew, setShowNew] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);

  const feedback = useMemo(() => {
    if (!error) {
      return "";
    }

    if (isAxiosError<ProblemDetailsResponse>(error)) {
      const problem = error.response?.data?.errors?.[0];
      return (
        problem?.detail ||
        problem?.title ||
        "Não foi possível redefinir a senha."
      );
    }

    return "Não foi possível redefinir a senha.";
  }, [error]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!token || newPassword !== confirmPassword) {
      return;
    }

    await mutateAsync({ token, newPassword });
    setNewPassword("");
    setConfirmPassword("");
  };

  const isConfirmInvalid =
    confirmPassword.length > 0 && newPassword !== confirmPassword;

  return (
    <PageWrapper className="gap-16">
      <MainNav showLinks={false} showActions={false} />

      <div className="mx-auto w-full max-w-xl rounded-4xl bg-white p-8 shadow-[0_30px_80px_rgba(124,139,255,0.15)]">
        <div className="mb-6">
          <h1 className="text-2xl font-semibold">Definir nova senha</h1>
          <p className="mt-2 text-sm text-slate-500">
            Crie uma nova senha para acessar sua conta.
          </p>
        </div>

        {!token ? (
          <div className="rounded-3xl border border-rose-200/80 bg-rose-50/70 px-4 py-3 text-sm text-rose-600 shadow-sm">
            Token inválido ou ausente. Solicite uma nova redefinição.
          </div>
        ) : (
          <Form.Root className="space-y-5" onSubmit={handleSubmit}>
            <Form.Field className="space-y-2" name="newPassword">
              <div className="flex items-baseline justify-between">
                <Form.Label className="text-sm font-medium">
                  Nova senha
                </Form.Label>
                <Form.Message
                  className="text-xs text-rose-500"
                  match="valueMissing"
                >
                  Informe a nova senha
                </Form.Message>
              </div>
              <div className="relative">
                <Form.Control asChild>
                  <input
                    type={showNew ? "text" : "password"}
                    value={newPassword}
                    onChange={(event) => setNewPassword(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 pr-12 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="********"
                    required
                    autoComplete="new-password"
                  />
                </Form.Control>
                <Toggle.Root
                  pressed={showNew}
                  onPressedChange={setShowNew}
                  aria-label={showNew ? "Ocultar senha" : "Mostrar senha"}
                  className="absolute inset-y-0 right-2 flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:text-slate-600"
                >
                  {showNew ? (
                    <EyeOff className="h-4 w-4" />
                  ) : (
                    <Eye className="h-4 w-4" />
                  )}
                </Toggle.Root>
              </div>
            </Form.Field>

            <Form.Field className="space-y-2" name="confirmPassword">
              <div className="flex items-baseline justify-between">
                <Form.Label className="text-sm font-medium">
                  Confirmar nova senha
                </Form.Label>
                <Form.Message
                  className="text-xs text-rose-500"
                  match="valueMissing"
                >
                  Confirme a nova senha
                </Form.Message>
              </div>
              <div className="relative">
                <Form.Control asChild>
                  <input
                    type={showConfirm ? "text" : "password"}
                    value={confirmPassword}
                    onChange={(event) => setConfirmPassword(event.target.value)}
                    className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 pr-12 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    placeholder="********"
                    required
                    autoComplete="new-password"
                  />
                </Form.Control>
                <Toggle.Root
                  pressed={showConfirm}
                  onPressedChange={setShowConfirm}
                  aria-label={showConfirm ? "Ocultar senha" : "Mostrar senha"}
                  className="absolute inset-y-0 right-2 flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:text-slate-600"
                >
                  {showConfirm ? (
                    <EyeOff className="h-4 w-4" />
                  ) : (
                    <Eye className="h-4 w-4" />
                  )}
                </Toggle.Root>
              </div>
              {isConfirmInvalid ? (
                <p className="text-xs text-rose-500">
                  As senhas não coincidem.
                </p>
              ) : null}
            </Form.Field>

            {feedback ? (
              <div className="rounded-3xl border border-rose-200/80 bg-rose-50/70 px-4 py-3 text-sm text-rose-600 shadow-sm">
                {feedback}
              </div>
            ) : null}

            {isSuccess ? (
              <div className="rounded-3xl border border-emerald-200/80 bg-emerald-50/70 px-4 py-3 text-sm text-emerald-700 shadow-sm">
                Senha redefinida com sucesso. Faça login novamente.
              </div>
            ) : null}

            <button
              type="submit"
              disabled={isPending || isConfirmInvalid}
              className="inline-flex h-11 w-full items-center justify-center rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-4 text-sm font-semibold text-white shadow-lg shadow-cyan-200 transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
            >
              {isPending ? "Atualizando..." : "Atualizar senha"}
            </button>
          </Form.Root>
        )}

        <div className="mt-6 text-center text-xs text-slate-500">
          <Link href="/login" className="text-cyan-600">
            Voltar para o login
          </Link>
        </div>
      </div>
    </PageWrapper>
  );
}
