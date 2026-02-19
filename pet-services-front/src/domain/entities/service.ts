import type { Base } from "./base";
import type { Photo } from "./photo";
import type { Category } from "./category";
import type { Tag } from "./tag";

export interface Service extends Base {
  providerId: string;
  name: string;
  description: string;
  price: number;
  priceMinimum: number;
  priceMaximum: number;
  duration: number;
  photos: Photo[];
  categories: Category[];
  tags: Tag[];
}
