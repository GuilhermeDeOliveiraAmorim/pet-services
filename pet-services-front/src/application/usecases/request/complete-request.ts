import { Request } from "@/domain/entities";

export interface CompleteRequestInput {
  id: string;
}

export interface CompleteRequestUseCase {
  execute(input: CompleteRequestInput): Promise<Request>;
}
