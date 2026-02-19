import { Tag } from "@/domain/entities";
import {
  ListTagsInput,
  ListTagsOutput,
} from "@/application/usecases/tag";

export interface TagGateway {
  listTags(input?: ListTagsInput): Promise<ListTagsOutput>;
}
