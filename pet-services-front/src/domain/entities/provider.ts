import type { Base } from "./base";
import type { Address } from "./address";
import type { Photo } from "./photo";

export interface Provider extends Base {
  userId: string;
  businessName: string;
  address: Address;
  description: string;
  priceRange: string;
  averageRating: number;
  photos: Photo[];
}
