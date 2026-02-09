import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

export default function OwnerDashboardPage() {
  return (
    <PageWrapper className="gap-10">
      <MainNav />

      <section className="flex flex-col gap-6">
        <div>
          <p className="text-xs font-semibold uppercase text-cyan-400">
            Dashboard
          </p>
          <h1 className="mt-2 text-2xl font-semibold text-slate-900">
            Olá, tutor
          </h1>
          <p className="mt-2 text-sm text-slate-600">
            Este é o seu painel inicial. Aqui vão aparecer seus pets e
            agendamentos.
          </p>
        </div>

        <div className="rounded-4xl border border-dashed border-slate-200 bg-white px-6 py-16 text-center">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-2xl bg-cyan-50 text-cyan-500">
            <span className="text-lg font-semibold">🐾</span>
          </div>
          <h2 className="mt-4 text-lg font-semibold text-slate-900">
            Sem dados ainda
          </h2>
          <p className="mt-2 text-sm text-slate-600">
            Quando você cadastrar seu primeiro pet ou agendar um serviço, as
            informações vão aparecer aqui.
          </p>
        </div>
      </section>
    </PageWrapper>
  );
}
