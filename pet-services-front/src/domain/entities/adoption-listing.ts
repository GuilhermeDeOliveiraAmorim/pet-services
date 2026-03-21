import type { Base } from "./base";
import type { Pet } from "./pet";

export const AdoptionListingSex = {
  Male: "male",
  Female: "female",
} as const;

export type AdoptionListingSexValue =
  (typeof AdoptionListingSex)[keyof typeof AdoptionListingSex];

export const AdoptionListingSize = {
  Small: "small",
  Medium: "medium",
  Large: "large",
} as const;

export type AdoptionListingSizeValue =
  (typeof AdoptionListingSize)[keyof typeof AdoptionListingSize];

export const AdoptionListingAgeGroup = {
  Puppy: "puppy",
  Adult: "adult",
  Senior: "senior",
} as const;

export type AdoptionListingAgeGroupValue =
  (typeof AdoptionListingAgeGroup)[keyof typeof AdoptionListingAgeGroup];

export const AdoptionListingStatus = {
  Draft: "draft",
  Published: "published",
  Paused: "paused",
  Adopted: "adopted",
  Archived: "archived",
} as const;

export type AdoptionListingStatusValue =
  (typeof AdoptionListingStatus)[keyof typeof AdoptionListingStatus];

export interface AdoptionGuardianProfileSummary {
  id: string;
  displayName: string;
  about: string;
  whatsapp: string;
  cityId: string;
  stateId: string;
}

export interface AdoptionListing extends Base {
  title: string;
  description: string;
  adoptionReason: string;
  status: string;
  sex: string;
  size: string;
  ageGroup: string;
  stateId: string;
  cityId: string;
  latitude?: number;
  longitude?: number;
  vaccinated: boolean;
  goodWithChildren: boolean;
  goodWithDogs: boolean;
  goodWithCats: boolean;
  specialNeeds: boolean;
  neutered: boolean;
  dewormed: boolean;
  requiresHouseScreening: boolean;
  petId: string;
  guardianProfileId: string;
  publishedAt?: string | null;
  adoptedAt?: string | null;
  pet?: Pet;
  guardianProfile?: AdoptionGuardianProfileSummary;
}
