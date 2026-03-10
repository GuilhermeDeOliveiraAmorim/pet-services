"use client";

import { useEffect } from "react";
import { usePathname, useRouter } from "next/navigation";

import { useAuthSession, useUserProfile } from "@/application";
import { UserTypes } from "@/domain";

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
  const { data: profileData, isLoading: isLoadingProfile } = useUserProfile({
    enabled: isAuthenticated,
  });

  const isOwnerArea = pathname?.startsWith("/owner");
  const isProviderArea = pathname?.startsWith("/provider");
  const userType = profileData?.user?.userType;

  useEffect(() => {
    if (!isHydrated) {
      return;
    }

    if (!isAuthenticated && pathname !== redirectTo) {
      router.replace(redirectTo);
      return;
    }

    if (!isAuthenticated || (!isOwnerArea && !isProviderArea)) {
      return;
    }

    if (!userType) {
      return;
    }

    if (isOwnerArea && userType === UserTypes.Provider) {
      router.replace("/provider");
      return;
    }

    if (isProviderArea && userType === UserTypes.Owner) {
      router.replace("/owner");
    }
  }, [
    isAuthenticated,
    isHydrated,
    isOwnerArea,
    isProviderArea,
    pathname,
    redirectTo,
    router,
    userType,
  ]);

  if (!isHydrated) {
    return null;
  }

  if (!isAuthenticated) {
    return null;
  }

  if ((isOwnerArea || isProviderArea) && (isLoadingProfile || !userType)) {
    return null;
  }

  return <>{children}</>;
};

export default AuthGuard;
