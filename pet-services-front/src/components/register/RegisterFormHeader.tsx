export default function RegisterFormHeader() {
  return (
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
  );
}
