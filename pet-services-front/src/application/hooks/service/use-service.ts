import { useMemo } from "react";
import {
  useMutation,
  useQuery,
  type UseMutationOptions,
  type UseQueryOptions,
} from "@tanstack/react-query";

import type {
  AddServiceInput,
  AddServiceOutput,
  AddServiceCategoryOutput,
  AddServicePhotoOutput,
  AddServiceTagPayload,
  AddServiceTagOutput,
  DeleteServiceCategoryOutput,
  DeleteServiceTagOutput,
  DeleteServiceOutput,
  DeleteServicePhotoOutput,
  GetServiceOutput,
  ListServicesInput,
  ListServicesOutput,
  SearchServicesInput,
  SearchServicesOutput,
  UpdateServiceInput,
  UpdateServiceOutput,
} from "@/application";
import { createServiceCases } from "@/application/factories/service-usecase-factory";
import { createApiContext } from "@/infra";

const useServiceUseCases = () => {
  return useMemo(() => {
    const { serviceGateway } = createApiContext();
    return createServiceCases(serviceGateway);
  }, []);
};

type AddServiceOptions = Omit<
  UseMutationOptions<AddServiceOutput, Error, AddServiceInput>,
  "mutationFn"
>;

type GetServiceOptions = Omit<
  UseQueryOptions<GetServiceOutput, Error>,
  "queryKey" | "queryFn"
>;

type UpdateServiceOptions = Omit<
  UseMutationOptions<UpdateServiceOutput, Error, UpdateServiceInput>,
  "mutationFn"
>;

type DeleteServiceOptions = Omit<
  UseMutationOptions<DeleteServiceOutput, Error, string | number>,
  "mutationFn"
>;

type ListServicesOptions = Omit<
  UseQueryOptions<ListServicesOutput, Error>,
  "queryKey" | "queryFn"
> & {
  input?: ListServicesInput;
};

type SearchServicesOptions = Omit<
  UseQueryOptions<SearchServicesOutput, Error>,
  "queryKey" | "queryFn"
> & {
  input?: SearchServicesInput;
};

type AddServicePhotoOptions = Omit<
  UseMutationOptions<
    AddServicePhotoOutput,
    Error,
    { serviceId: string | number; photo: File }
  >,
  "mutationFn"
>;

type DeleteServicePhotoOptions = Omit<
  UseMutationOptions<
    DeleteServicePhotoOutput,
    Error,
    { serviceId: string | number; photoId: string | number }
  >,
  "mutationFn"
>;

type AddServiceCategoryOptions = Omit<
  UseMutationOptions<
    AddServiceCategoryOutput,
    Error,
    { serviceId: string | number; categoryId: string | number }
  >,
  "mutationFn"
>;

type AddServiceTagOptions = Omit<
  UseMutationOptions<
    AddServiceTagOutput,
    Error,
    { serviceId: string | number; payload: AddServiceTagPayload }
  >,
  "mutationFn"
>;

type DeleteServiceCategoryOptions = Omit<
  UseMutationOptions<
    DeleteServiceCategoryOutput,
    Error,
    { serviceId: string | number; categoryId: string | number }
  >,
  "mutationFn"
>;

type DeleteServiceTagOptions = Omit<
  UseMutationOptions<
    DeleteServiceTagOutput,
    Error,
    { serviceId: string | number; tagId: string | number }
  >,
  "mutationFn"
>;

export const useServiceAdd = (options?: AddServiceOptions) => {
  const { addService } = useServiceUseCases();

  return useMutation({
    mutationFn: (input) => addService.execute(input),
    ...options,
  });
};

export const useServiceGet = (
  serviceId?: string | number,
  options?: GetServiceOptions,
) => {
  const { getService } = useServiceUseCases();

  return useQuery({
    queryKey: ["service", serviceId],
    queryFn: () => getService.execute(serviceId!),
    enabled: Boolean(serviceId),
    ...options,
  });
};

export const useServiceUpdate = (options?: UpdateServiceOptions) => {
  const { updateService } = useServiceUseCases();

  return useMutation({
    mutationFn: (input) => updateService.execute(input),
    ...options,
  });
};

export const useServiceDelete = (options?: DeleteServiceOptions) => {
  const { deleteService } = useServiceUseCases();

  return useMutation({
    mutationFn: (serviceId) => deleteService.execute(serviceId),
    ...options,
  });
};

export const useServiceList = (options?: ListServicesOptions) => {
  const { listServices } = useServiceUseCases();
  const input = options?.input;

  return useQuery({
    queryKey: ["services", input],
    queryFn: () => listServices.execute(input),
    ...options,
  });
};

export const useServiceSearch = (options?: SearchServicesOptions) => {
  const { searchServices } = useServiceUseCases();
  const input = options?.input;

  return useQuery({
    queryKey: ["services-search", input],
    queryFn: () => searchServices.execute(input),
    ...options,
  });
};

export const useServiceAddPhoto = (options?: AddServicePhotoOptions) => {
  const { addServicePhoto } = useServiceUseCases();

  return useMutation({
    mutationFn: (input) => addServicePhoto.execute(input),
    ...options,
  });
};

export const useServiceDeletePhoto = (options?: DeleteServicePhotoOptions) => {
  const { deleteServicePhoto } = useServiceUseCases();

  return useMutation({
    mutationFn: (input) =>
      deleteServicePhoto.execute(input.serviceId, input.photoId),
    ...options,
  });
};

export const useServiceAddCategory = (options?: AddServiceCategoryOptions) => {
  const { addServiceCategory } = useServiceUseCases();

  return useMutation({
    mutationFn: (input) =>
      addServiceCategory.execute(input.serviceId, input.categoryId),
    ...options,
  });
};

export const useServiceAddTag = (options?: AddServiceTagOptions) => {
  const { addServiceTag } = useServiceUseCases();

  return useMutation({
    mutationFn: (input) =>
      addServiceTag.execute(input.serviceId, input.payload),
    ...options,
  });
};

export const useServiceDeleteCategory = (
  options?: DeleteServiceCategoryOptions,
) => {
  const { deleteServiceCategory } = useServiceUseCases();

  return useMutation({
    mutationFn: (input) =>
      deleteServiceCategory.execute(input.serviceId, input.categoryId),
    ...options,
  });
};

export const useServiceDeleteTag = (options?: DeleteServiceTagOptions) => {
  const { deleteServiceTag } = useServiceUseCases();

  return useMutation({
    mutationFn: (input) =>
      deleteServiceTag.execute(input.serviceId, input.tagId),
    ...options,
  });
};
