import type { UserGateway } from "../ports";
import { ChangePasswordUseCase, RegisterUserUseCase } from "../usecases/auth";

export const createUserUseCases = (gateway: UserGateway) => {
  return {
    registerUserUseCase: new RegisterUserUseCase(gateway),
    changePasswordUseCase: new ChangePasswordUseCase(gateway),
  };
};
