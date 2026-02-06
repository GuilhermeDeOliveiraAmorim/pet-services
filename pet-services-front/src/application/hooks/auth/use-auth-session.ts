import { useCallback, useEffect, useMemo, useState } from "react";

import {
  type AuthSession,
  clearAuthSession,
  getAuthSession,
  isAuthSessionValid,
  setAuthSession,
} from "@/lib/auth-session";

export const useAuthSession = () => {
  const [session, setSessionState] = useState<AuthSession | null>(null);
  const [isHydrated, setIsHydrated] = useState(false);

  useEffect(() => {
    const stored = getAuthSession();
    setSessionState(stored);
    setIsHydrated(true);
  }, []);

  const updateSession = useCallback((next: AuthSession) => {
    setAuthSession(next);
    setSessionState(next);
  }, []);

  const clearSession = useCallback(() => {
    clearAuthSession();
    setSessionState(null);
  }, []);

  const isAuthenticated = useMemo(
    () => isAuthSessionValid(session),
    [session],
  );

  return {
    session,
    isAuthenticated,
    isHydrated,
    setSession: updateSession,
    clearSession,
  };
};
