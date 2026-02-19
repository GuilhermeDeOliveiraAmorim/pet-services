import { Base } from "./base";

export interface Review extends Base {
  userId: string;
  providerId: string;
  rating: number;
  comment: string;
}
