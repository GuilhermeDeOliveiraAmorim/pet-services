import { RequestGateway } from "@/application/ports/request-gateway";
import {
  AddRequestUseCase,
  AddRequestInput,
  GetRequestUseCase,
  GetRequestInput,
  ListRequestsUseCase,
  ListRequestsInput,
  ListRequestsOutput,
  AcceptRequestUseCase,
  AcceptRequestInput,
  RejectRequestUseCase,
  RejectRequestInput,
  CompleteRequestUseCase,
  CompleteRequestInput,
} from "@/application/usecases/request";

class AddRequestUseCaseImpl implements AddRequestUseCase {
  constructor(private requestGateway: RequestGateway) {}

  async execute(input: AddRequestInput) {
    return this.requestGateway.addRequest(input);
  }
}

class GetRequestUseCaseImpl implements GetRequestUseCase {
  constructor(private requestGateway: RequestGateway) {}

  async execute(input: GetRequestInput) {
    return this.requestGateway.getRequest(input);
  }
}

class ListRequestsUseCaseImpl implements ListRequestsUseCase {
  constructor(private requestGateway: RequestGateway) {}

  async execute(input?: ListRequestsInput): Promise<ListRequestsOutput> {
    return this.requestGateway.listRequests(input);
  }
}

class AcceptRequestUseCaseImpl implements AcceptRequestUseCase {
  constructor(private requestGateway: RequestGateway) {}

  async execute(input: AcceptRequestInput) {
    return this.requestGateway.acceptRequest(input);
  }
}

class RejectRequestUseCaseImpl implements RejectRequestUseCase {
  constructor(private requestGateway: RequestGateway) {}

  async execute(input: RejectRequestInput) {
    return this.requestGateway.rejectRequest(input);
  }
}

class CompleteRequestUseCaseImpl implements CompleteRequestUseCase {
  constructor(private requestGateway: RequestGateway) {}

  async execute(input: CompleteRequestInput) {
    return this.requestGateway.completeRequest(input);
  }
}

export function createRequestCases(requestGateway: RequestGateway) {
  return {
    addRequest: new AddRequestUseCaseImpl(requestGateway),
    getRequest: new GetRequestUseCaseImpl(requestGateway),
    listRequests: new ListRequestsUseCaseImpl(requestGateway),
    acceptRequest: new AcceptRequestUseCaseImpl(requestGateway),
    rejectRequest: new RejectRequestUseCaseImpl(requestGateway),
    completeRequest: new CompleteRequestUseCaseImpl(requestGateway),
  };
}
