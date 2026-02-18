import type { Base } from "./base";
import type { Photo } from "./photo";
import type { Species } from "./specie";

export interface Pet extends Base {
  userId: string;
  name: string;
  specie: Species;
  age: number;
  weight: number;
  notes: string;
  photos: Photo[];
}
