import type { User } from "@/domain";

export interface LoginInput {
  email: string;
  password: string;
  userAgent?: string;
  ip?: string;
}

export interface LoginOutput {
  user: User;
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}
