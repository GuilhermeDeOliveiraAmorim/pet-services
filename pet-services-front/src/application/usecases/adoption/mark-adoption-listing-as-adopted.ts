import type { AdoptionGateway } from "@/application/ports/adoption-gateway";

export interface MarkAdoptionListingAsAdoptedInput {
  listingId: string;
}

export interface MarkAdoptionListingAsAdoptedOutput {
  id: string;
  status: string;
}

export class MarkAdoptionListingAsAdoptedUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input: MarkAdoptionListingAsAdoptedInput,
  ): Promise<MarkAdoptionListingAsAdoptedOutput> {
    return this.adoptionGateway.markListingAsAdopted(input);
  }
}
