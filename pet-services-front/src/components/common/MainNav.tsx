import Link from "next/link";
import * as NavigationMenu from "@radix-ui/react-navigation-menu";

const navItems = [
  { label: "Home", href: "/" },
  { label: "Serviços", href: "/#services" },
  { label: "Veterinários", href: "/#doctors" },
  { label: "Depoimentos", href: "/#testimonials" },
  { label: "Contato", href: "/#contact" },
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
  return (
    <header className={`flex items-center justify-between ${className}`.trim()}>
      <Link href="/" className="flex items-center gap-2">
        <div className="flex h-9 w-9 items-center justify-center rounded-2xl bg-linear-to-tr from-teal-400 to-cyan-400 text-white font-semibold">
          pet
        </div>
        <span className="text-lg font-semibold text-slate-900">PetCare</span>
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
          <Link
            href="/login"
            className="hidden rounded-full border border-slate-200 px-4 py-2 text-sm font-medium text-slate-700 lg:inline-flex"
          >
            Login
          </Link>
          <Link
            href="/login"
            className="rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-5 py-2 text-sm font-semibold text-white shadow-lg shadow-cyan-200"
          >
            Agendar
          </Link>
        </div>
      ) : null}
    </header>
  );
}
