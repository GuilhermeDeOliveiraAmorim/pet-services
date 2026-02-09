import type { Photo } from "@/domain";
import type { UserGateway } from "@/application/ports";

export interface AddUserPhotoInput {
  file: File;
}

export interface AddUserPhotoOutput {
  message?: string;
  detail?: string;
  photo?: Photo;
}

export class AddUserPhotoUseCase {
  constructor(private readonly userGateway: UserGateway) {}

  execute(input: AddUserPhotoInput): Promise<AddUserPhotoOutput> {
    return this.userGateway.addUserPhoto(input);
  }
}
