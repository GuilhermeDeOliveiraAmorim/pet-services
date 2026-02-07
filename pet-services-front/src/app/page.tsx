import MainNav from "@/components/common/MainNav";

export default function Home() {
  return (
    <div className="min-h-screen bg-[#f7f9ff] text-slate-900">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-16 px-6 py-10 lg:px-8">
        <MainNav />

        <section
          id="home"
          className="grid items-center gap-12 lg:grid-cols-[1.1fr_0.9fr]"
        >
          <div className="flex flex-col gap-6">
            <span className="w-fit rounded-full bg-pink-100 px-4 py-1 text-xs font-semibold text-pink-500">
              Plataforma de cuidados
            </span>
            <h1 className="text-4xl font-semibold leading-tight text-slate-900 sm:text-5xl">
              Cuidado excelente para deixar seu pet feliz
            </h1>
            <p className="max-w-xl text-base leading-7 text-slate-600">
              Consultas, banho e tosa, vacinas e acompanhamentos em um só lugar.
              Faça a gestão completa do bem-estar do seu pet com agilidade e
              carinho.
            </p>
            <div className="flex flex-wrap items-center gap-4">
              <button className="rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-cyan-200">
                Agendar agora
              </button>
              <button className="rounded-full border border-slate-200 px-6 py-3 text-sm font-semibold text-slate-700">
                Ver serviços
              </button>
            </div>
            <div className="grid grid-cols-3 gap-6 pt-6 text-sm">
              <div>
                <p className="text-2xl font-semibold text-slate-900">99%</p>
                <p className="text-slate-500">Satisfação</p>
              </div>
              <div>
                <p className="text-2xl font-semibold text-slate-900">85k+</p>
                <p className="text-slate-500">Clientes</p>
              </div>
              <div>
                <p className="text-2xl font-semibold text-slate-900">90k+</p>
                <p className="text-slate-500">Atendimentos</p>
              </div>
            </div>
          </div>

          <div className="relative">
            <div className="absolute -left-6 top-10 hidden h-16 w-16 rounded-full bg-teal-300/40 blur-xl lg:block" />
            <div className="absolute -right-6 bottom-8 hidden h-20 w-20 rounded-full bg-pink-300/40 blur-xl lg:block" />
            <div className="relative rounded-4xl bg-white p-6 shadow-[0_30px_80px_rgba(124,139,255,0.25)]">
              <div className="flex items-start justify-between">
                <div className="rounded-2xl bg-teal-50 px-4 py-2 text-xs font-semibold text-teal-500">
                  Melhor atendimento
                </div>
                <div className="flex items-center gap-2">
                  <span className="inline-flex h-8 w-8 items-center justify-center rounded-full bg-cyan-100 text-xs font-semibold text-cyan-600">
                    2
                  </span>
                  <span className="inline-flex h-8 w-8 items-center justify-center rounded-full bg-pink-100 text-xs font-semibold text-pink-500">
                    4
                  </span>
                </div>
              </div>
              <div className="mt-6 flex items-center gap-6">
                <div className="relative flex h-36 w-36 items-center justify-center rounded-full bg-linear-to-tr from-yellow-200 via-orange-200 to-pink-200">
                  <div className="h-24 w-24 rounded-[28px] bg-white shadow-xl" />
                </div>
                <div className="flex flex-col gap-3">
                  <div className="rounded-2xl border border-slate-100 bg-white px-4 py-3 text-sm shadow-sm">
                    <p className="text-slate-600">Consulta veterinária</p>
                    <p className="font-semibold text-slate-900">Dra. Sophia</p>
                  </div>
                  <div className="rounded-2xl border border-slate-100 bg-white px-4 py-3 text-sm shadow-sm">
                    <p className="text-slate-600">Dica do dia</p>
                    <p className="font-semibold text-slate-900">Hidratação</p>
                  </div>
                </div>
              </div>
              <div className="mt-6 flex items-center justify-between rounded-2xl bg-slate-50 px-4 py-3 text-sm">
                <span className="text-slate-600">Pets atendidos</span>
                <span className="font-semibold text-slate-900">85k</span>
              </div>
            </div>
          </div>
        </section>

        <section className="grid gap-8 rounded-4xl bg-white px-6 py-10 shadow-sm lg:grid-cols-[1.1fr_0.9fr]">
          <div className="flex flex-col gap-4">
            <h2 className="text-2xl font-semibold text-slate-900">
              Nossos melhores profissionais
            </h2>
            <p className="text-slate-600">
              Veterinários e especialistas certificados, prontos para cuidar do
              seu animal com carinho e tecnologia.
            </p>
          </div>
          <div className="grid gap-4 sm:grid-cols-3">
            {["Dra. Livia Amaral", "Dr. João Esporte", "Dra. Emery Roser"].map(
              (name, index) => (
                <div
                  key={name}
                  className="flex flex-col items-center gap-3 rounded-3xl bg-slate-50 px-4 py-6 text-center"
                >
                  <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-linear-to-br from-cyan-200 to-blue-200 text-lg font-semibold text-slate-700">
                    {index + 1}
                  </div>
                  <p className="text-sm font-semibold text-slate-900">{name}</p>
                  <p className="text-xs text-slate-500">Veterinário</p>
                </div>
              ),
            )}
          </div>
        </section>

        <section id="services" className="flex flex-col gap-8">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs font-semibold uppercase text-teal-400">
                Serviços
              </p>
              <h2 className="text-2xl font-semibold text-slate-900">
                Nossos serviços para o seu pet
              </h2>
            </div>
            <button className="rounded-full border border-slate-200 px-4 py-2 text-sm font-medium text-slate-700">
              Ver todos
            </button>
          </div>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {[
              {
                title: "Banho & Tosa",
                desc: "Profissionais com carinho e técnica.",
              },
              { title: "Pet Care", desc: "Planos mensais personalizados." },
              {
                title: "Tratamentos",
                desc: "Exames e acompanhamento clínico.",
              },
              { title: "Vacinação", desc: "Calendário completo de vacinas." },
            ].map((item) => (
              <div
                key={item.title}
                className="rounded-3xl bg-white p-6 shadow-sm"
              >
                <div className="mb-4 h-10 w-10 rounded-2xl bg-teal-100" />
                <h3 className="text-sm font-semibold text-slate-900">
                  {item.title}
                </h3>
                <p className="mt-2 text-sm text-slate-600">{item.desc}</p>
              </div>
            ))}
          </div>
        </section>

        <section
          id="testimonials"
          className="grid gap-10 lg:grid-cols-[1fr_1fr]"
        >
          <div className="flex flex-col gap-6">
            <p className="text-xs font-semibold uppercase text-pink-400">
              Depoimentos
            </p>
            <h2 className="text-2xl font-semibold text-slate-900">
              O que nossos clientes dizem?
            </h2>
            <div className="rounded-3xl bg-white p-6 shadow-sm">
              <p className="text-sm text-slate-600">
                “Meu cachorro voltou super feliz e o acompanhamento foi
                perfeito. Atendimento rápido e equipe cuidadosa.”
              </p>
              <div className="mt-4 flex items-center gap-3">
                <div className="h-10 w-10 rounded-full bg-linear-to-br from-orange-200 to-pink-200" />
                <div>
                  <p className="text-sm font-semibold text-slate-900">
                    Mariana Torres
                  </p>
                  <p className="text-xs text-slate-500">Tutora</p>
                </div>
              </div>
            </div>
            <div className="grid gap-4 sm:grid-cols-2">
              {[
                "Atendimento incrível",
                "Equipe cuidadosa",
                "Ambiente seguro",
                "Veterinários dedicados",
              ].map((text) => (
                <div
                  key={text}
                  className="rounded-2xl bg-white p-4 text-sm text-slate-600 shadow-sm"
                >
                  {text}
                </div>
              ))}
            </div>
          </div>
          <div className="relative flex items-center justify-center">
            <div className="absolute inset-6 rounded-[40px] bg-linear-to-br from-teal-100 to-cyan-100" />
            <div className="relative z-10 flex w-full max-w-sm flex-col gap-4 rounded-4xl bg-white p-6 shadow-[0_30px_80px_rgba(124,139,255,0.2)]">
              <div className="flex items-center justify-between">
                <span className="text-sm font-semibold text-slate-900">
                  Consulta
                </span>
                <span className="rounded-full bg-pink-100 px-3 py-1 text-xs font-semibold text-pink-500">
                  24/7
                </span>
              </div>
              <div className="h-48 rounded-3xl bg-linear-to-br from-yellow-100 to-orange-100" />
              <div className="rounded-2xl bg-slate-50 px-4 py-3 text-sm text-slate-600">
                Cuidamos com dedicação e tecnologia em cada detalhe.
              </div>
            </div>
          </div>
        </section>

        <section
          id="contact"
          className="rounded-4xl bg-linear-to-br from-cyan-50 to-teal-50 px-6 py-12 text-center"
        >
          <p className="text-xs font-semibold uppercase text-cyan-400">
            Contato
          </p>
          <h2 className="mt-3 text-2xl font-semibold text-slate-900">
            Tem alguma dúvida?
          </h2>
          <p className="mt-2 text-sm text-slate-600">
            Fale com a nossa equipe e receba atendimento rápido e humanizado.
          </p>
          <button className="mt-6 rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-cyan-200">
            Falar com a equipe
          </button>
        </section>

        <footer className="flex flex-col items-center justify-between gap-6 border-t border-slate-200 pt-8 text-sm text-slate-500 lg:flex-row">
          <div className="flex items-center gap-2">
            <div className="flex h-9 w-9 items-center justify-center rounded-2xl bg-linear-to-tr from-teal-400 to-cyan-400 text-white font-semibold">
              pet
            </div>
            <span className="text-slate-700">PetCare</span>
          </div>
          <p>© 2026 PetCare. Todos os direitos reservados.</p>
        </footer>
      </div>
    </div>
  );
}
