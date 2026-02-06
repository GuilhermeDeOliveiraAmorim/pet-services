export interface ChangePasswordInput {
  userId: string;
  oldPassword: string;
  newPassword: string;
}

export interface ChangePasswordInputBody {
  oldPassword: string;
  newPassword: string;
}

export interface ChangePasswordOutput {
  message?: string;
  detail?: string;
}
