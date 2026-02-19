import { Tag } from "@/domain/entities";

export interface ListTagsInput {
  name?: string;
  page?: number;
  pageSize?: number;
}

export interface ListTagsOutput {
  tags: Tag[];
  total: number;
  page: number;
  pageSize: number;
}

export interface ListTagsUseCase {
  execute(input?: ListTagsInput): Promise<ListTagsOutput>;
}
