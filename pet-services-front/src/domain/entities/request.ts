import { Base } from "./base";
import { Pet } from "./pet";

export interface Request extends Base {
  userId: string;
  userName?: string;
  providerId: string;
  businessName?: string;
  serviceId: string;
  serviceName?: string;
  petId: string;
  pet?: Pet;
  notes: string;
  status: "pending" | "accepted" | "rejected" | "completed" | "cancelled";
  rejectReason?: string;
}
