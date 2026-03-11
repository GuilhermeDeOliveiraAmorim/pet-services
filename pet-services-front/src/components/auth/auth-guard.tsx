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
  const isProfileArea = pathname?.startsWith("/profile");
  const userType = profileData?.user?.userType;
  const profileComplete = profileData?.user?.profileComplete;

  useEffect(() => {
    if (!isHydrated) {
      return;
    }

    if (!isAuthenticated && pathname !== redirectTo) {
      router.replace(redirectTo);
      return;
    }

    if (!isAuthenticated) {
      return;
    }

    if (!isLoadingProfile && profileComplete === false && !isProfileArea) {
      router.replace("/profile");
      return;
    }

    if (!isOwnerArea && !isProviderArea) {
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
    isLoadingProfile,
    isOwnerArea,
    isProfileArea,
    isProviderArea,
    pathname,
    profileComplete,
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

  if (isLoadingProfile) {
    return null;
  }

  if (profileComplete === false && !isProfileArea) {
    return null;
  }

  if ((isOwnerArea || isProviderArea) && !userType) {
    return null;
  }

  return <>{children}</>;
};

export default AuthGuard;
