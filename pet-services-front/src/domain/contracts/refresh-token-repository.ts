import type { PasswordResetToken, RefreshToken } from "../entities";

export interface RefreshTokenRepository {
  create(token: RefreshToken): Promise<void>;
  findById(id: string): Promise<RefreshToken | null>;
  findByToken(token: string): Promise<RefreshToken | null>;
  findActiveByUserId(userId: string): Promise<RefreshToken[]>;
  update(token: RefreshToken): Promise<void>;
  revoke(tokenId: string): Promise<void>;
  revokeAllByUserId(userId: string): Promise<void>;
  deleteExpired(): Promise<void>;
  isValid(token: string): Promise<boolean>;
  createPasswordReset(token: PasswordResetToken): Promise<void>;
  revokeAllPasswordResetByUserId(userId: string): Promise<void>;
  findValidPasswordResetByToken(token: string): Promise<PasswordResetToken | null>;
  revokePasswordResetByToken(token: string): Promise<void>;
}
