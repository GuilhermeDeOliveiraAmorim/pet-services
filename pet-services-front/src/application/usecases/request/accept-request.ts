import { Request } from "@/domain/entities";

export interface AcceptRequestInput {
  id: string;
}

export interface AcceptRequestUseCase {
  execute(input: AcceptRequestInput): Promise<Request>;
}
