/**
 * Global Cache Configuration
 *
 * Define staleTime e gcTime (cacheTime) para cada domínio
 * Baseado em frequência de mudança e criticidade dos dados
 */

export const CACHE_CONFIG = {
  // Static reference data - rarely changes
  reference: {
    staleTime: Infinity, // Never consider stale
    gcTime: 30 * 60 * 1000, // Keep 30 min in garbage
  },

  // Species and Breeds - rarely change
  species: {
    staleTime: 60 * 60 * 1000, // 1 hour
    gcTime: 60 * 60 * 1000, // 1 hour
  },
  breeds: {
    staleTime: 60 * 60 * 1000, // 1 hour
    gcTime: 60 * 60 * 1000, // 1 hour
  },

  // Categories, Tags - rarely change
  categories: {
    staleTime: 30 * 60 * 1000, // 30 min
    gcTime: 60 * 60 * 1000, // 1 hour
  },
  tags: {
    staleTime: 30 * 60 * 1000, // 30 min
    gcTime: 60 * 60 * 1000, // 1 hour
  },

  // User Profile - medium priority
  user: {
    staleTime: 5 * 60 * 1000, // 5 min
    gcTime: 30 * 60 * 1000, // 30 min (keep longer to avoid refetch on back)
    refetchOnWindowFocus: false,
  },

  // Services - data can change, but not frequently
  services: {
    staleTime: 60 * 1000, // 1 min
    gcTime: 10 * 60 * 1000, // 10 min
    refetchOnWindowFocus: false,
  },

  // Providers - data can change
  providers: {
    staleTime: 60 * 1000, // 1 min
    gcTime: 10 * 60 * 1000, // 10 min
    refetchOnWindowFocus: false,
  },

  // Reviews - frequently read, medium priority
  reviews: {
    staleTime: 5 * 60 * 1000, // 5 min (more dynamic)
    gcTime: 30 * 60 * 1000, // 30 min
    refetchOnWindowFocus: true, // Refetch when user returns (new reviews possible)
  },

  // Requests - CRITICAL, changes status frequently
  requests: {
    staleTime: 30 * 1000, // 30 sec (critical, status changes)
    gcTime: 10 * 60 * 1000, // Keep 10 min cached
    refetchOnWindowFocus: true, // Refetch when user returns
    refetchInterval: 30 * 1000, // Poll every 30 sec (optional, can be enabled in page)
  },
} as const;

export type CacheConfigKey = keyof typeof CACHE_CONFIG;
