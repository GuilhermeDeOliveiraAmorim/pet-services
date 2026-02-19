import { Request } from "@/domain/entities";

export interface GetRequestInput {
  id: string;
}

export interface GetRequestUseCase {
  execute(input: GetRequestInput): Promise<Request>;
}
