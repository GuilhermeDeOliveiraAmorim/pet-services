import { useMemo } from "react";
import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseQueryOptions,
} from "@tanstack/react-query";
import {
  AddRequestInput,
  GetRequestInput,
  ListRequestsInput,
  ListRequestsOutput,
  AcceptRequestInput,
  RejectRequestInput,
  CompleteRequestInput,
} from "@/application/usecases/request";
import { createRequestCases } from "@/application/factories/request-usecase-factory";
import { createApiContext } from "@/infra";

const useRequestUseCases = () => {
  return useMemo(() => {
    const { requestGateway } = createApiContext();
    return createRequestCases(requestGateway);
  }, []);
};

const REQUEST_KEYS = {
  all: ["requests"] as const,
  lists: () => [...REQUEST_KEYS.all, "list"] as const,
  list: (filters?: ListRequestsInput) =>
    [...REQUEST_KEYS.lists(), filters] as const,
  details: () => [...REQUEST_KEYS.all, "detail"] as const,
  detail: (id: string) => [...REQUEST_KEYS.details(), id] as const,
};

type ListRequestsOptions = Omit<
  UseQueryOptions<ListRequestsOutput, Error>,
  "queryKey" | "queryFn"
>;

export function useRequestAdd() {
  const queryClient = useQueryClient();
  const { addRequest } = useRequestUseCases();

  return useMutation({
    mutationFn: (input: AddRequestInput) => addRequest.execute(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: REQUEST_KEYS.lists() });
    },
  });
}

export function useRequestGet(input: GetRequestInput) {
  const { getRequest } = useRequestUseCases();

  return useQuery({
    queryKey: REQUEST_KEYS.detail(input.id),
    queryFn: () => getRequest.execute(input),
    enabled: !!input.id,
  });
}

export function useRequestList(
  input?: ListRequestsInput,
  options?: ListRequestsOptions,
) {
  const { listRequests } = useRequestUseCases();

  return useQuery({
    queryKey: REQUEST_KEYS.list(input),
    queryFn: () => listRequests.execute(input),
    ...options,
  });
}

export function useRequestAccept() {
  const queryClient = useQueryClient();
  const { acceptRequest } = useRequestUseCases();

  return useMutation({
    mutationFn: (input: AcceptRequestInput) => acceptRequest.execute(input),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: REQUEST_KEYS.lists() });
      queryClient.invalidateQueries({
        queryKey: REQUEST_KEYS.detail(variables.id),
      });
    },
  });
}

export function useRequestReject() {
  const queryClient = useQueryClient();
  const { rejectRequest } = useRequestUseCases();

  return useMutation({
    mutationFn: (input: RejectRequestInput) => rejectRequest.execute(input),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: REQUEST_KEYS.lists() });
      queryClient.invalidateQueries({
        queryKey: REQUEST_KEYS.detail(variables.id),
      });
    },
  });
}

export function useRequestComplete() {
  const queryClient = useQueryClient();
  const { completeRequest } = useRequestUseCases();

  return useMutation({
    mutationFn: (input: CompleteRequestInput) => completeRequest.execute(input),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: REQUEST_KEYS.lists() });
      queryClient.invalidateQueries({
        queryKey: REQUEST_KEYS.detail(variables.id),
      });
    },
  });
}
