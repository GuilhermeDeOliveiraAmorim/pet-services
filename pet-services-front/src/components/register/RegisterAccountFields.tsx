import * as Form from "@radix-ui/react-form";

import { type UserType, UserTypes } from "@/domain";

type RegisterAccountFieldsProps = {
  name: string;
  onNameChange: (value: string) => void;
  userType: UserType;
  onUserTypeChange: (value: UserType) => void;
  email: string;
  onEmailChange: (value: string) => void;
  password: string;
  onPasswordChange: (value: string) => void;
};

export default function RegisterAccountFields({
  name,
  onNameChange,
  userType,
  onUserTypeChange,
  email,
  onEmailChange,
  password,
  onPasswordChange,
}: RegisterAccountFieldsProps) {
  return (
    <>
      <div className="grid gap-4 sm:grid-cols-2">
        <Form.Field className="space-y-2" name="name">
          <div className="flex items-baseline justify-between">
            <Form.Label className="text-sm font-medium">
              Nome completo
            </Form.Label>
            <Form.Message
              className="text-xs text-rose-500"
              match="valueMissing"
            >
              Informe o nome
            </Form.Message>
          </div>
          <Form.Control asChild>
            <input
              id="name"
              value={name}
              onChange={(event) => onNameChange(event.target.value)}
              className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
              placeholder="Seu nome"
              required
            />
          </Form.Control>
        </Form.Field>

        <Form.Field className="space-y-2" name="userType">
          <Form.Label className="text-sm font-medium">
            Tipo de usuário
          </Form.Label>
          <div className="relative">
            <Form.Control asChild>
              <select
                id="userType"
                value={userType}
                onChange={(event) =>
                  onUserTypeChange(event.target.value as UserType)
                }
                className="h-11 w-full appearance-none rounded-2xl border border-slate-200 bg-slate-50 px-4 pr-10 text-sm text-slate-700 shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                required
              >
                <option value={UserTypes.Owner}>Tutor</option>
                <option value={UserTypes.Provider}>Prestador</option>
              </select>
            </Form.Control>
            <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3 text-slate-400">
              <svg
                aria-hidden="true"
                className="h-4 w-4"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fillRule="evenodd"
                  d="M5.23 7.21a.75.75 0 0 1 1.06.02L10 11.126l3.71-3.895a.75.75 0 1 1 1.08 1.04l-4.24 4.45a.75.75 0 0 1-1.08 0l-4.24-4.45a.75.75 0 0 1 .02-1.06Z"
                  clipRule="evenodd"
                />
              </svg>
            </div>
          </div>
        </Form.Field>
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
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
              id="email"
              type="email"
              value={email}
              onChange={(event) => onEmailChange(event.target.value)}
              className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
              placeholder="voce@email.com"
              required
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
          <Form.Control asChild>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(event) => onPasswordChange(event.target.value)}
              className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
              placeholder="********"
              required
              autoComplete="new-password"
            />
          </Form.Control>
        </Form.Field>
      </div>
    </>
  );
}
