import type { AdoptionGateway } from "@/application/ports/adoption-gateway";

export interface ReviewAdoptionApplicationInput {
  applicationId: string;
  action: "under_review" | "interview" | "approve" | "reject";
  notesInternal?: string;
}

export interface ReviewAdoptionApplicationOutput {
  id: string;
  status: string;
}

export class ReviewAdoptionApplicationUseCase {
  constructor(private readonly adoptionGateway: AdoptionGateway) {}

  execute(
    input: ReviewAdoptionApplicationInput,
  ): Promise<ReviewAdoptionApplicationOutput> {
    return this.adoptionGateway.reviewApplication(input);
  }
}
