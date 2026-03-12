import Constants from "expo-constants";

type ExtraConfig = {
  apiUrl?: string;
};

type ExpoConfigWithExtra = {
  extra?: ExtraConfig;
};

const expoConfig = (Constants.expoConfig ?? {}) as ExpoConfigWithExtra;

export const env = {
  apiUrl: expoConfig.extra?.apiUrl ?? "http://localhost:8080",
};
