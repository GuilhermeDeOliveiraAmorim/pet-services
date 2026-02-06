export const referenceQueryKeys = {
  all: ["reference"] as const,
  countries: () => [...referenceQueryKeys.all, "countries"] as const,
  states: () => [...referenceQueryKeys.all, "states"] as const,
  cities: (stateId?: number) =>
    [...referenceQueryKeys.all, "cities", stateId ?? "all"] as const,
};
