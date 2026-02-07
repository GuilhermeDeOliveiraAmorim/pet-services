import axios, {
  AxiosHeaders,
  type AxiosInstance,
  type AxiosError,
  type InternalAxiosRequestConfig,
} from "axios";

import {
  clearAuthSession,
  getAuthSession,
  setAuthSession,
} from "@/lib/auth-session";

type RetryableRequestConfig = InternalAxiosRequestConfig & {
  _retry?: boolean;
};

export const createApiClient = (baseURL?: string): AxiosInstance => {
  const client = axios.create({
    baseURL:
      baseURL ?? process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080",
    headers: {
      "Content-Type": "application/json",
    },
  });

  client.interceptors.request.use((config) => {
    const session = getAuthSession();
    const url = config.url ?? "";
    const isAuthRoute = url.startsWith("/auth/") || url.includes("/auth/");

    if (session?.accessToken && !isAuthRoute) {
      const headers = AxiosHeaders.from(config.headers ?? {});
      headers.set("Authorization", `Bearer ${session.accessToken}`);
      config.headers = headers;
    }
    return config;
  });

  client.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
      const originalConfig = error.config as RetryableRequestConfig | undefined;
      const status = error.response?.status;

      if (!originalConfig || status !== 401 || originalConfig._retry) {
        return Promise.reject(error);
      }

      const session = getAuthSession();
      if (!session?.refreshToken) {
        clearAuthSession();
        return Promise.reject(error);
      }

      originalConfig._retry = true;

      try {
        const refreshResponse = await axios.post<{
          access_token: string;
          refresh_token: string;
          expires_in: number;
        }>(
          `${client.defaults.baseURL}/auth/refresh`,
          { refresh_token: session.refreshToken },
          {
            headers: {
              "Content-Type": "application/json",
            },
          },
        );

        const accessToken = refreshResponse.data.access_token;
        const refreshToken = refreshResponse.data.refresh_token;
        const expiresIn = refreshResponse.data.expires_in;
        const expiresAt = Date.now() + expiresIn * 1000;

        setAuthSession({ accessToken, refreshToken, expiresAt });

        client.defaults.headers.common.Authorization = `Bearer ${accessToken}`;
        const headers = AxiosHeaders.from(originalConfig.headers ?? {});
        headers.set("Authorization", `Bearer ${accessToken}`);
        originalConfig.headers = headers;

        return client(originalConfig);
      } catch (refreshError) {
        clearAuthSession();
        return Promise.reject(refreshError);
      }
    },
  );

  return client;
};
