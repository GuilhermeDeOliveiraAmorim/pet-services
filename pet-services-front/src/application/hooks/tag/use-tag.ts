import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { ListTagsInput } from "@/application/usecases/tag";
import { createTagCases } from "@/application/factories/tag-usecase-factory";
import { createApiContext } from "@/infra";
import { TAG_KEYS } from "./tag-query-keys";

const useTagUseCases = () => {
  return useMemo(() => {
    const { tagGateway } = createApiContext();
    return createTagCases(tagGateway);
  }, []);
};

export function useTagList(input?: ListTagsInput) {
  const { listTags } = useTagUseCases();

  return useQuery({
    queryKey: TAG_KEYS.list(input),
    queryFn: () => listTags.execute(input),
  });
}
