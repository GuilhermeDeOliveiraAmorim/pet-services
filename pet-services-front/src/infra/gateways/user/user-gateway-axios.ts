import type {
  ChangePasswordInput,
  ChangePasswordOutput,
  RegisterUserInput,
  RegisterUserOutput,
  UserGateway,
} from "@/application";
import type {
  AddUserPhotoInput,
  AddUserPhotoOutput,
  DeactivateUserOutput,
  GetProfileOutput,
  ReactivateUserOutput,
  UpdateUserInput,
  UpdateUserOutput,
} from "@/application";
import type { AxiosInstance } from "axios";
import { mapUserFromApi } from "@/infra/mappers/user-mapper";

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
      login: {
        email: input.login.email,
        password: input.login.password,
      },
      phone: {
        country_code: input.phone.countryCode,
        area_code: input.phone.areaCode,
        number: input.phone.number,
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

  async getProfile(): Promise<GetProfileOutput> {
    const { data } = await this.http.get<{ user: unknown; provider_id?: string }>(
      "/users/profile",
    );

    return {
      user: mapUserFromApi(data.user as Record<string, unknown>),
      providerId: data.provider_id,
    };
  }

  async updateUser(input: UpdateUserInput): Promise<UpdateUserOutput> {
    const payload: {
      name?: string;
      phone?: {
        country_code?: string;
        area_code?: string;
        number?: string;
      };
      address?: {
        street?: string;
        number?: string;
        neighborhood?: string;
        city?: string;
        zip_code?: string;
        state?: string;
        country?: string;
        complement?: string;
        location?: {
          latitude: number;
          longitude: number;
        };
      };
    } = {};

    if (typeof input.name === "string") {
      payload.name = input.name;
    }

    if (input.phone) {
      payload.phone = {
        country_code: input.phone.countryCode,
        area_code: input.phone.areaCode,
        number: input.phone.number,
      };
    }

    if (input.address) {
      payload.address = {
        street: input.address.street,
        number: input.address.number,
        neighborhood: input.address.neighborhood,
        city: input.address.city,
        zip_code: input.address.zipCode,
        state: input.address.state,
        country: input.address.country,
        complement: input.address.complement,
        location: input.address.location
          ? {
              latitude: input.address.location.latitude,
              longitude: input.address.location.longitude,
            }
          : undefined,
      };
    }

    const { data } = await this.http.put<{
      message?: string;
      detail?: string;
      user?: unknown;
    }>("/users", payload);

    return {
      message: data.message,
      detail: data.detail,
      user: data.user
        ? mapUserFromApi(data.user as Record<string, unknown>)
        : undefined,
    };
  }

  async deactivateUser(): Promise<DeactivateUserOutput> {
    const { data } =
      await this.http.post<DeactivateUserOutput>("/users/deactivate");
    return data;
  }

  async reactivateUser(): Promise<ReactivateUserOutput> {
    const { data } =
      await this.http.post<ReactivateUserOutput>("/users/reactivate");
    return data;
  }

  async deleteUser(): Promise<DeactivateUserOutput> {
    const { data } = await this.http.delete<DeactivateUserOutput>("/users");
    return data;
  }

  async addUserPhoto(input: AddUserPhotoInput): Promise<AddUserPhotoOutput> {
    const formData = new FormData();
    formData.append("file", input.file);

    const { data } = await this.http.post<AddUserPhotoOutput>(
      "/users/photos",
      formData,
      {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      },
    );

    return data;
  }
}
