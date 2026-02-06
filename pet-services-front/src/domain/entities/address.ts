import type { Location } from "./location";

export interface Address {
  street: string;
  number: string;
  neighborhood: string;
  city: string;
  zipCode: string;
  state: string;
  country: string;
  complement: string;
  location: Location;
}
