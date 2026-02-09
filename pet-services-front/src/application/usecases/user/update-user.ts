import type { Address, Phone, User } from "@/domain";
import type { UserGateway } from "@/application/ports";

export interface UpdateUserInput {
  name?: string;
  phone?: Partial<Phone>;
  address?: Partial<Address> & { location?: Partial<Address["location"]> };
}

export interface UpdateUserOutput {
  message?: string;
  detail?: string;
  user?: User;
}

export class UpdateUserUseCase {
  constructor(private readonly userGateway: UserGateway) {}

  execute(input: UpdateUserInput): Promise<UpdateUserOutput> {
    return this.userGateway.updateUser(input);
  }
}
