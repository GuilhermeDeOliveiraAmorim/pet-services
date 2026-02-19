import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { ListTagsInput } from "@/application/usecases/tag";
import { createTagCases } from "@/application/factories/tag-usecase-factory";
import { createApiContext } from "@/infra";

const useTagUseCases = () => {
  return useMemo(() => {
    const { tagGateway } = createApiContext();
    return createTagCases(tagGateway);
  }, []);
};

const TAG_KEYS = {
  all: ["tags"] as const,
  lists: () => [...TAG_KEYS.all, "list"] as const,
  list: (filters?: ListTagsInput) => [...TAG_KEYS.lists(), filters] as const,
};

export function useTagList(input?: ListTagsInput) {
  const { listTags } = useTagUseCases();

  return useQuery({
    queryKey: TAG_KEYS.list(input),
    queryFn: () => listTags.execute(input),
  });
}
