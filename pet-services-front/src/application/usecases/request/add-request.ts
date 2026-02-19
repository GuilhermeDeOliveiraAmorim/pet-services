import { Request } from "@/domain/entities";

export interface AddRequestInput {
  providerId: string;
  serviceId: string;
  petId: string;
  notes: string;
}

export interface AddRequestUseCase {
  execute(input: AddRequestInput): Promise<Request>;
}
