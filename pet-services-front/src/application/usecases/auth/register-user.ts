import type { Address, Login, Phone, UserType } from "@/domain";

export interface RegisterUserInput {
  name: string;
  userType: UserType;
  login: Login;
  phone: Phone;
  address: Address;
}

export interface RegisterUserOutput {
  message?: string;
  detail?: string;
}
