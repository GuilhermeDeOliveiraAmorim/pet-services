import type {
  ChangePasswordInput,
  ChangePasswordOutput,
  RegisterUserInput,
  RegisterUserOutput,
  UserGateway,
} from "@/application";
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
    const payload = {
      name: input.name,
      user_type: input.userType,
      login: {
        email: input.login.email,
        password: input.login.password,
      },
      phone: {
        country_code: input.phone.countryCode,
        area_code: input.phone.areaCode,
        number: input.phone.number,
      },
      address: {
        street: input.address.street,
        number: input.address.number,
        neighborhood: input.address.neighborhood,
        city: input.address.city,
        zip_code: input.address.zipCode,
        state: input.address.state,
        country: input.address.country,
        complement: input.address.complement,
        location: {
          latitude: input.address.location.latitude,
          longitude: input.address.location.longitude,
        },
      },
    };

    const { data } = await this.http.post<RegisterUserOutput>(
      "/users/register",
      payload,
    );
    return data;
  }

  async changePassword(
    input: ChangePasswordInput,
  ): Promise<ChangePasswordOutput> {
    const { data } = await this.http.post<ChangePasswordOutput>(
      "/users/change-password",
      {
        old_password: input.oldPassword,
        new_password: input.newPassword,
      },
    );

    return data;
  }
}
