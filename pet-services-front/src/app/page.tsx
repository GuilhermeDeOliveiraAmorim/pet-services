import * as Avatar from "@radix-ui/react-avatar";
import * as Separator from "@radix-ui/react-separator";
import * as Tooltip from "@radix-ui/react-tooltip";

import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import Banner from "../../public/banner.png";

export default function Home() {
  return (
    <PageWrapper className="gap-16">
      <MainNav 
       />
   <section
        id="home"
        className="relative flex items-center justify-center py-24 bg-cover bg-center rounded-4xl overflow-hidden"
        style={{ backgroundImage: `url(${Banner.src})` }}
      >
        <div className="absolute inset-0 bg-black/40" />
        <div className="relative z-10 flex w-full max-w-3xl flex-col items-center gap-8 text-center">
          <h1 className="text-4xl font-semibold leading-tight text-white sm:text-5xl">
            Cuidado Profissional para quem você Ama
          </h1>

          <p className="max-w-xl text-base leading-7 text-white">
            Conecte-se com os melhores veterinários, pet shops e cuidadores da sua
            região.
          </p>
          <div className="flex w-full max-w-xl items-stretch rounded-full border border-slate-300 overflow-hidden">
            <input
              type="text"
              placeholder="Onde você está?"
              className="flex-1 bg-white px-4 py-2 text-sm text-slate-700 focus:outline-none"
            />
            <select
              className="bg-white px-4 py-2 text-sm text-slate-700 focus:outline-none"
              defaultValue=""
            >
              <option value="" disabled>
                Qual serviço seu pet precisa?
              </option>
              <option>Clínica Veterinária</option>
              <option>Pet Shop</option>
              <option>Banho e Tosa</option>
              <option>Hotelzinho e Creche</option>
              <option>Passeador(a)</option>
              <option>Pet Sitter</option>
              <option>Adestrador</option>
            </select>
            <button className="bg-teal-400 px-6 py-2 text-sm font-semibold text-white">
              Buscar
            </button>
          </div>

          <Tooltip.Provider delayDuration={200}>
            <div className="flex flex-wrap items-center gap-4">
              {/* …restante dos botões… */}
            </div>
          </Tooltip.Provider>

          <div className="grid grid-cols-3 gap-6 pt-6 text-sm">
            {/* … */}
          </div>
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

      <section id="testimonials" className="grid gap-10 lg:grid-cols-[1fr_1fr]">
        <div className="flex flex-col gap-6">
          <p className="text-xs font-semibold uppercase text-pink-400">
            Depoimentos
          </p>
          <h2 className="text-2xl font-semibold text-slate-900">
            O que nossos clientes dizem?
          </h2>
          <div className="rounded-3xl bg-white p-6 shadow-sm">
            <p className="text-sm text-slate-600">
              “Meu cachorro voltou super feliz e o acompanhamento foi perfeito.
              Atendimento rápido e equipe cuidadosa.”
            </p>
            <div className="mt-4 flex items-center gap-3">
              <Avatar.Root className="h-10 w-10 overflow-hidden rounded-full bg-linear-to-br from-orange-200 to-pink-200">
                <Avatar.Fallback>MT</Avatar.Fallback>
              </Avatar.Root>
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
        <p className="text-xs font-semibold uppercase text-cyan-400">Contato</p>
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

      <Separator.Root className="h-px w-full bg-slate-200" />
      <footer className="flex flex-col items-center justify-between gap-6 pt-8 text-sm text-slate-500 lg:flex-row">
        <div className="flex items-center gap-2">
          <div className="flex h-9 w-9 items-center justify-center rounded-2xl bg-linear-to-tr from-teal-400 to-cyan-400 text-white font-semibold">
            pet
          </div>
          <span className="text-slate-700">PetCare</span>
        </div>
        <p>© 2026 PetCare. Todos os direitos reservados.</p>
      </footer>
    </PageWrapper>
  );
}
