import type { ChangePasswordInput, ChangePasswordOutput, RegisterUserInput, RegisterUserOutput, UserGateway } from "@/application";
import type { AxiosInstance } from "axios";

export class UserGatewayAxios implements UserGateway {
  constructor(private readonly http: AxiosInstance) {}

  setAccessToken(token?: string) {
    if (token) {
      this.http.defaults.headers.common.Authorization = `Bearer ${token}`;
      return;
    }

    delete this.http.defaults.headers.common.Authorization;
  }

  async registerUser(input: RegisterUserInput): Promise<RegisterUserOutput> {
    const { data } = await this.http.post<RegisterUserOutput>("/users/register", input);
    return data;
  }

  async changePassword(input: ChangePasswordInput): Promise<ChangePasswordOutput> {
    const { data } = await this.http.post<ChangePasswordOutput>("/users/change-password", {
      old_password: input.oldPassword,
      new_password: input.newPassword,
    });

    return data;
  }
}
