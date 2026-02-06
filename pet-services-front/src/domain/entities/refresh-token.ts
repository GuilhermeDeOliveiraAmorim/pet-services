import type { Base } from "./base";

export const TokenTypes = {
  Refresh: "refresh",
  Access: "access",
} as const;

export type TokenType = (typeof TokenTypes)[keyof typeof TokenTypes];

export interface RefreshToken extends Base {
  userId: string;
  token: string;
  expiresAt: string;
  revokedAt?: string | null;
  userAgent: string;
  ipAddress: string;
  tokenType: TokenType;
}
