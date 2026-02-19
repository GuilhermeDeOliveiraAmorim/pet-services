import { Base } from "./base";
import { Pet } from "./pet";

export interface Request extends Base {
  userId: string;
  providerId: string;
  serviceId: string;
  petId: string;
  pet?: Pet;
  notes: string;
  status: "pending" | "accepted" | "rejected" | "completed" | "cancelled";
  rejectReason?: string;
}
