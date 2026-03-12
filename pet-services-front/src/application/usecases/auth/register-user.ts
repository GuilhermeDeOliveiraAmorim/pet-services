import type { Login, Phone, UserType } from "@/domain";
import type { UserGateway } from "@/application/ports";

export interface RegisterUserInput {
  name: string;
  userType: UserType;
  login: Login;
  phone: Phone;
}

export interface RegisterUserOutput {
  message?: string;
  detail?: string;
}

export class RegisterUserUseCase {
  constructor(private readonly userGateway: UserGateway) {}

  execute(input: RegisterUserInput): Promise<RegisterUserOutput> {
    return this.userGateway.registerUser(input);
  }
}
