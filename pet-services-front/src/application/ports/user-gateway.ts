import type { ChangePasswordInput, ChangePasswordOutput, RegisterUserInput, RegisterUserOutput } from "../usecases/auth";

export interface UserGateway {
  registerUser(input: RegisterUserInput): Promise<RegisterUserOutput>;
  changePassword(input: ChangePasswordInput): Promise<ChangePasswordOutput>;
}
