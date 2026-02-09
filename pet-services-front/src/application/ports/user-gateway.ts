import type {
  ChangePasswordInput,
  ChangePasswordOutput,
  RegisterUserInput,
  RegisterUserOutput,
} from "../usecases/auth";
import type {
  DeactivateUserOutput,
  GetProfileOutput,
  ReactivateUserOutput,
  UpdateUserInput,
  UpdateUserOutput,
} from "../usecases/user";
import type { AddUserPhotoInput, AddUserPhotoOutput } from "../usecases/user";

export interface UserGateway {
  registerUser(input: RegisterUserInput): Promise<RegisterUserOutput>;
  changePassword(input: ChangePasswordInput): Promise<ChangePasswordOutput>;
  getProfile(): Promise<GetProfileOutput>;
  updateUser(input: UpdateUserInput): Promise<UpdateUserOutput>;
  deactivateUser(): Promise<DeactivateUserOutput>;
  reactivateUser(): Promise<ReactivateUserOutput>;
  addUserPhoto(input: AddUserPhotoInput): Promise<AddUserPhotoOutput>;
}
