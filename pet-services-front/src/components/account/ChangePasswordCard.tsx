"use client";

import { useMemo, useState } from "react";
import * as Form from "@radix-ui/react-form";
import * as Toggle from "@radix-ui/react-toggle";
import { Eye, EyeOff } from "lucide-react";

import { useUserChangePassword } from "@/application";

export default function ChangePasswordCard() {
  const { mutateAsync, isPending, error, isSuccess } =
    useUserChangePassword();

  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showCurrent, setShowCurrent] = useState(false);
  const [showNew, setShowNew] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);

  const feedback = useMemo(() => {
    if (!error) {
      return "";
    }

    return "Não foi possível alterar a senha. Verifique os dados e tente novamente.";
  }, [error]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (newPassword !== confirmPassword) {
      return;
    }

    await mutateAsync({
      oldPassword: currentPassword,
      newPassword,
    });

    setCurrentPassword("");
    setNewPassword("");
    setConfirmPassword("");
  };

  const isConfirmInvalid =
    confirmPassword.length > 0 && newPassword !== confirmPassword;

  return (
    <section className="rounded-4xl bg-white p-6 shadow-sm">
      <div className="mb-6">
        <p className="text-xs font-semibold uppercase text-slate-400">
          Segurança
        </p>
        <h2 className="mt-2 text-xl font-semibold text-slate-900">
          Alterar senha
        </h2>
        <p className="mt-2 text-sm text-slate-600">
          Atualize sua senha para manter sua conta protegida.
        </p>
      </div>

      <Form.Root className="space-y-5" onSubmit={handleSubmit}>
        <Form.Field className="space-y-2" name="currentPassword">
          <div className="flex items-baseline justify-between">
            <Form.Label className="text-sm font-medium">
              Senha atual
            </Form.Label>
            <Form.Message className="text-xs text-rose-500" match="valueMissing">
              Informe a senha atual
            </Form.Message>
          </div>
          <div className="relative">
            <Form.Control asChild>
              <input
                type={showCurrent ? "text" : "password"}
                value={currentPassword}
                onChange={(event) => setCurrentPassword(event.target.value)}
                className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 pr-12 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                placeholder="********"
                required
                autoComplete="current-password"
              />
            </Form.Control>
            <Toggle.Root
              pressed={showCurrent}
              onPressedChange={setShowCurrent}
              aria-label={showCurrent ? "Ocultar senha" : "Mostrar senha"}
              className="absolute inset-y-0 right-2 flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:text-slate-600"
            >
              {showCurrent ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            </Toggle.Root>
          </div>
        </Form.Field>

        <Form.Field className="space-y-2" name="newPassword">
          <div className="flex items-baseline justify-between">
            <Form.Label className="text-sm font-medium">
              Nova senha
            </Form.Label>
            <Form.Message className="text-xs text-rose-500" match="valueMissing">
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
              {showNew ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            </Toggle.Root>
          </div>
        </Form.Field>

        <Form.Field className="space-y-2" name="confirmPassword">
          <div className="flex items-baseline justify-between">
            <Form.Label className="text-sm font-medium">
              Confirmar nova senha
            </Form.Label>
            <Form.Message className="text-xs text-rose-500" match="valueMissing">
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
            Senha atualizada com sucesso.
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
    </section>
  );
}
