import { Request } from "@/domain/entities";
import {
  AddRequestInput,
  GetRequestInput,
  ListRequestsInput,
  ListRequestsOutput,
  AcceptRequestInput,
  RejectRequestInput,
  CompleteRequestInput,
} from "@/application/usecases/request";

export interface RequestGateway {
  addRequest(input: AddRequestInput): Promise<Request>;
  getRequest(input: GetRequestInput): Promise<Request>;
  listRequests(input?: ListRequestsInput): Promise<ListRequestsOutput>;
  acceptRequest(input: AcceptRequestInput): Promise<Request>;
  rejectRequest(input: RejectRequestInput): Promise<Request>;
  completeRequest(input: CompleteRequestInput): Promise<Request>;
}
