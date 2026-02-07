import { type UserType, UserTypes } from "@/domain";
import SelectField from "@/components/common/SelectField";

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
        <div className="space-y-2">
          <label className="text-sm font-medium" htmlFor="name">
            Nome completo
          </label>
          <input
            id="name"
            value={name}
            onChange={(event) => onNameChange(event.target.value)}
            className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
            placeholder="Seu nome"
            required
          />
        </div>
        <SelectField
          id="userType"
          label="Tipo de usuário"
          value={userType}
          onChange={(value) => onUserTypeChange(value as UserType)}
          options={[
            { value: UserTypes.Owner, label: "Tutor" },
            { value: UserTypes.Provider, label: "Prestador" },
          ]}
        />
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
            onChange={(event) => onEmailChange(event.target.value)}
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
            onChange={(event) => onPasswordChange(event.target.value)}
            className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
            placeholder="********"
            required
            autoComplete="new-password"
          />
        </div>
      </div>
    </>
  );
}
