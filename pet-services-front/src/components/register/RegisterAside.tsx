export default function RegisterAside() {
  return (
    <div className="hidden flex-col justify-center gap-6 rounded-4xl bg-white p-10 shadow-[0_30px_80px_rgba(124,139,255,0.2)] lg:flex">
      <div className="flex items-center gap-2">
        <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-linear-to-tr from-teal-400 to-cyan-400 text-white font-semibold">
          pet
        </div>
        <span className="text-lg font-semibold">PetCare</span>
      </div>
      <h1 className="text-3xl font-semibold leading-tight">Crie sua conta</h1>
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
  );
}
