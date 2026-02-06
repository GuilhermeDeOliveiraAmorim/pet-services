export interface RefreshTokenInput {
  refreshToken: string;
  userAgent?: string;
  ip?: string;
}

export interface RefreshTokenOutput {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}
