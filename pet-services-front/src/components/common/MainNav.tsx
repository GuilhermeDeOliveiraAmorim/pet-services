"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import * as NavigationMenu from "@radix-ui/react-navigation-menu";

import { useAuthLogout, useAuthSession } from "@/application";

import Image from "next/image";

const navItems = [
  { label: "Encontre Serviços", href: "/services" },
  { label: "Seja um Parceiro", href: "/partner" }
];

type MainNavProps = {
  className?: string;
  showLinks?: boolean;
  showActions?: boolean;
};

export default function MainNav({
  className = "",
  showLinks = true,
  showActions = true,
}: MainNavProps) {
  const router = useRouter();
  const { session, isAuthenticated, clearSession } = useAuthSession();
  const { mutateAsync: logout, isPending } = useAuthLogout();

  const handleLogout = async () => {
    try {
      if (session) {
        await logout({ revokeAll: true });
      }
    } finally {
      clearSession();
      router.replace("/login");
    }
  };

  return (
    <header className={`flex items-center justify-between ${className}`.trim()}>
      <Link href="/" className="flex items-center gap-2">
        <Image
  src="/pawIcon.svg"
  alt="Pet Services Logo"
  width={32}
  height={32}
/>
        <span className="text-lg font-semibold text-slate-900">Pet Services</span>
      </Link>

      {showLinks ? (
        <NavigationMenu.Root className="hidden lg:flex">
          <NavigationMenu.List className="flex items-center gap-6 text-sm font-medium text-slate-600">
            {navItems.map((item) => (
              <NavigationMenu.Item key={item.href}>
                <NavigationMenu.Link asChild>
                  <Link
                    href={item.href}
                    className={item.href === "/" ? "text-slate-900" : undefined}
                  >
                    {item.label}
                  </Link>
                </NavigationMenu.Link>
              </NavigationMenu.Item>
            ))}
          </NavigationMenu.List>
        </NavigationMenu.Root>
      ) : null}

      {showActions ? (
        <div className="flex items-center gap-3">
          {isAuthenticated ? (
            <button
              type="button"
              onClick={handleLogout}
              disabled={isPending}
              className="rounded-full border border-slate-200 px-4 py-2 text-sm font-medium text-slate-700 transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
            >
              {isPending ? "Saindo..." : "Sair"}
            </button>
          ) : (
            <>
              <Link
                href="/login"
                className="hidden rounded-full border border-slate-200 px-4 py-2 text-sm font-medium text-slate-700 lg:inline-flex"
              >
                Entrar
              </Link>
              <Link
                href="/register"
                className="rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-5 py-2 text-sm font-semibold text-white shadow-lg shadow-cyan-200"
              >
                Cadastre-se
              </Link>
            </>
          )}
        </div>
      ) : null}
    </header>
  );
}
