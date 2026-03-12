import { apiClient } from "../../lib/http/client";

export interface LoginInput {
  email: string;
  password: string;
}

export interface RegisterInput {
  name: string;
  userType: "owner" | "provider";
  login: {
    email: string;
    password: string;
  };
  phone: {
    countryCode: string;
    areaCode: string;
    number: string;
  };
}

export async function login(input: LoginInput) {
  const response = await apiClient.post("/auth/login", input);
  return response.data;
}

export async function register(input: RegisterInput) {
  const response = await apiClient.post("/users", input);
  return response.data;
}
