export type AuthSession = {
  accessToken: string;
  refreshToken: string;
  expiresAt: number;
};

const STORAGE_KEY = "pet-services:auth";
const EXPIRY_SKEW_MS = 30_000;

const isBrowser = () => typeof window !== "undefined";

export const getAuthSession = (): AuthSession | null => {
  if (!isBrowser()) {
    return null;
  }

  const raw = window.localStorage.getItem(STORAGE_KEY);
  if (!raw) {
    return null;
  }

  try {
    const parsed = JSON.parse(raw) as AuthSession;
    if (!parsed?.accessToken || !parsed?.refreshToken || !parsed?.expiresAt) {
      return null;
    }

    return parsed;
  } catch {
    return null;
  }
};

export const setAuthSession = (session: AuthSession) => {
  if (!isBrowser()) {
    return;
  }

  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(session));
};

export const clearAuthSession = () => {
  if (!isBrowser()) {
    return;
  }

  window.localStorage.removeItem(STORAGE_KEY);
};

export const isAuthSessionValid = (session: AuthSession | null) => {
  if (!session) {
    return false;
  }

  return session.expiresAt - EXPIRY_SKEW_MS > Date.now();
};
