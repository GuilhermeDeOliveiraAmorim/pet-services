import type { UserGateway } from "../ports";
import { ChangePasswordUseCase, RegisterUserUseCase } from "../usecases/auth";
import {
  AddUserPhotoUseCase,
  DeactivateUserUseCase,
  GetProfileUseCase,
  ReactivateUserUseCase,
  UpdateUserUseCase,
} from "../usecases/user";

export const createUserUseCases = (gateway: UserGateway) => {
  return {
    registerUserUseCase: new RegisterUserUseCase(gateway),
    changePasswordUseCase: new ChangePasswordUseCase(gateway),
    getProfileUseCase: new GetProfileUseCase(gateway),
    updateUserUseCase: new UpdateUserUseCase(gateway),
    deactivateUserUseCase: new DeactivateUserUseCase(gateway),
    reactivateUserUseCase: new ReactivateUserUseCase(gateway),
    addUserPhotoUseCase: new AddUserPhotoUseCase(gateway),
  };
};
