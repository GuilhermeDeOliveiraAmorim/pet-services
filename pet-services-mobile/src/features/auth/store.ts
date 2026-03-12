import { create } from "zustand";

type UserType = "owner" | "provider" | "admin";

interface AuthUser {
  id: string;
  name: string;
  email: string;
  userType: UserType;
}

interface AuthSession {
  accessToken: string;
  user: AuthUser;
}

interface AuthState {
  session: AuthSession | null;
  setSession: (session: AuthSession) => void;
  clearSession: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  session: null,
  setSession: (session) => set({ session }),
  clearSession: () => set({ session: null }),
}));
