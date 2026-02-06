import axios, { type AxiosInstance } from "axios";

export const createApiClient = (baseURL?: string): AxiosInstance => {
  return axios.create({
    baseURL: baseURL ?? process.env.NEXT_PUBLIC_API_URL,
    headers: {
      "Content-Type": "application/json",
    },
  });
};
