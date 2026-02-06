export interface PasswordResetToken {
  token: string;
  userId: string;
  expiresAt: string;
  userAgent: string;
  ip: string;
  revokedAt?: string | null;
}
