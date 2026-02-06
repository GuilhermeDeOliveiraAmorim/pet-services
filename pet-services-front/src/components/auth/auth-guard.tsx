"use client";

import { useEffect } from "react";
import { usePathname, useRouter } from "next/navigation";

import { useAuthSession } from "@/application";

type AuthGuardProps = {
  children: React.ReactNode;
  redirectTo?: string;
};

const DEFAULT_REDIRECT = "/login";

export const AuthGuard = ({
  children,
  redirectTo = DEFAULT_REDIRECT,
}: AuthGuardProps) => {
  const router = useRouter();
  const pathname = usePathname();
  const { isAuthenticated, isHydrated } = useAuthSession();

  useEffect(() => {
    if (!isHydrated) {
      return;
    }

    if (!isAuthenticated && pathname !== redirectTo) {
      router.replace(redirectTo);
    }
  }, [isAuthenticated, isHydrated, pathname, redirectTo, router]);

  if (!isHydrated) {
    return null;
  }

  if (!isAuthenticated) {
    return null;
  }

  return <>{children}</>;
};

export default AuthGuard;
