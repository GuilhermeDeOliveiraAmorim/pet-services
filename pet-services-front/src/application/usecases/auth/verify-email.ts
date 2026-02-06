export interface VerifyEmailInput {
  token: string;
}

export interface VerifyEmailOutput {
  message?: string;
  detail?: string;
}
