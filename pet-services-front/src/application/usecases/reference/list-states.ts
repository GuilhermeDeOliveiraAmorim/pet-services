import type { State } from "@/domain";
import type { ReferenceGateway } from "@/application/ports";

export interface ListStatesOutput {
  states: State[];
}

export class ListStatesUseCase {
  constructor(private readonly gateway: ReferenceGateway) {}

  execute(): Promise<ListStatesOutput> {
    return this.gateway.listStates();
  }
}
