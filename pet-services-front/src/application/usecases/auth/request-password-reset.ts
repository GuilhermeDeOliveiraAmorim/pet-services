export interface RequestPasswordResetInput {
  email: string;
  userAgent?: string;
  ip?: string;
}

export interface RequestPasswordResetOutput {
  message?: string;
  detail?: string;
  resetToken?: string;
  expiresAt?: string;
}
