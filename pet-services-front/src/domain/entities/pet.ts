import type { Base } from "./base";
import type { Breed } from "./breed";
import type { Photo } from "./photo";
import type { Specie } from "./specie";

export interface Pet extends Base {
  userId: string;
  name: string;
  specie: Specie;
  breed: Breed;
  age: number;
  weight: number;
  notes: string;
  photos: Photo[];
}
