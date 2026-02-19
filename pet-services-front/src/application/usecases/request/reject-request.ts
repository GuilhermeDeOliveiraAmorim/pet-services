import { Request } from "@/domain/entities";

export interface RejectRequestInput {
  id: string;
  rejectReason: string;
}

export interface RejectRequestUseCase {
  execute(input: RejectRequestInput): Promise<Request>;
}
