export interface ResendVerificationEmailInput {
  email: string;
  userAgent?: string;
  ip?: string;
}

export interface ResendVerificationEmailOutput {
  message?: string;
  detail?: string;
  verifyToken?: string;
  expiresAt?: string;
}
