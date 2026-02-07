import * as Form from "@radix-ui/react-form";

import { type UserType, UserTypes } from "@/domain";
import RadixSelectField from "@/components/common/RadixSelectField";

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

        <RadixSelectField
          name="userType"
          label="Tipo de usuário"
          value={userType}
          onValueChange={(value) => onUserTypeChange(value as UserType)}
          options={[
            { value: UserTypes.Owner, label: "Tutor" },
            { value: UserTypes.Provider, label: "Prestador" },
          ]}
          required
        />
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
