export const USER_KEYS = {
  all: ["user"] as const,

  // Profile queries
  profiles: () => [...USER_KEYS.all, "profile"] as const,
  profile: () => [...USER_KEYS.profiles()] as const,
} as const;
