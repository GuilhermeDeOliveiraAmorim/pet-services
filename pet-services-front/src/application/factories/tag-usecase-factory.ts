import { TagGateway } from "@/application/ports/tag-gateway";
import {
  ListTagsUseCase,
  ListTagsInput,
  ListTagsOutput,
} from "@/application/usecases/tag";

class ListTagsUseCaseImpl implements ListTagsUseCase {
  constructor(private tagGateway: TagGateway) {}

  async execute(input?: ListTagsInput): Promise<ListTagsOutput> {
    return this.tagGateway.listTags(input);
  }
}

export function createTagCases(tagGateway: TagGateway) {
  return {
    listTags: new ListTagsUseCaseImpl(tagGateway),
  };
}
