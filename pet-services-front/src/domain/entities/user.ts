import type { Address } from "./address";
import type { Base } from "./base";
import type { Login } from "./login";
import type { Pet } from "./pet";
import type { Phone } from "./phone";
import type { Photo } from "./photo";

export const UserTypes = {
  Owner: "owner",
  Provider: "provider",
  Admin: "admin",
} as const;

export type UserType = (typeof UserTypes)[keyof typeof UserTypes];

export interface User extends Base {
  name: string;
  userType: UserType;
  login: Login;
  phone: Phone;
  address: Address;
  emailVerified: boolean;
  photos: Photo[];
  pets: Pet[];
}
