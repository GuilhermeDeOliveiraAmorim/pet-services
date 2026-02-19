import type {
  AddServiceInput,
  AddServiceOutput,
  AddServiceCategoryOutput,
  AddServicePhotoOutput,
  AddServiceTagOutput,
  DeleteServiceOutput,
  DeleteServicePhotoOutput,
  GetServiceOutput,
  ListServicesInput,
  ListServicesOutput,
  SearchServicesInput,
  SearchServicesOutput,
  ServiceGateway,
  UpdateServiceInput,
  UpdateServiceOutput,
} from "@/application";
import { mapServiceFromApi } from "@/infra/mappers/service-mapper";
import type { AxiosInstance } from "axios";

export class ServiceGatewayAxios implements ServiceGateway {
  constructor(private readonly http: AxiosInstance) {}

  async addService(input: AddServiceInput): Promise<AddServiceOutput> {
    const payload = {
      provider_id: input.providerId,
      name: input.name,
      description: input.description,
      price: input.price,
      price_minimum: input.priceMinimum,
      price_maximum: input.priceMaximum,
      duration: input.duration,
    };

    const { data } = await this.http.post<AddServiceOutput>(
      "/services",
      payload,
    );

    return data;
  }

  async getService(serviceId: string | number): Promise<GetServiceOutput> {
    const { data } = await this.http.get<{ service: unknown }>(
      `/services/${serviceId}`,
    );

    return {
      service: mapServiceFromApi(data.service as Record<string, unknown>),
    };
  }

  async updateService(
    input: UpdateServiceInput,
  ): Promise<UpdateServiceOutput> {
    const payload: {
      name?: string;
      description?: string;
      price?: number;
      price_minimum?: number;
      price_maximum?: number;
      duration?: number;
    } = {};

    if (input.name) payload.name = input.name;
    if (input.description) payload.description = input.description;
    if (typeof input.price === "number") payload.price = input.price;
    if (typeof input.priceMinimum === "number")
      payload.price_minimum = input.priceMinimum;
    if (typeof input.priceMaximum === "number")
      payload.price_maximum = input.priceMaximum;
    if (typeof input.duration === "number") payload.duration = input.duration;

    const { data } = await this.http.put<{
      message?: string;
      detail?: string;
      service?: unknown;
    }>(`/services/${input.serviceId}`, payload);

    return {
      message: data.message,
      detail: data.detail,
      service: data.service
        ? mapServiceFromApi(data.service as Record<string, unknown>)
        : undefined,
    };
  }

  async deleteService(
    serviceId: string | number,
  ): Promise<DeleteServiceOutput> {
    const { data } = await this.http.delete<DeleteServiceOutput>(
      `/services/${serviceId}`,
    );
    return data;
  }

  async listServices(input?: ListServicesInput): Promise<ListServicesOutput> {
    const params = new URLSearchParams();
    if (input?.providerId) params.append("provider_id", input.providerId);
    if (input?.categoryId) params.append("category_id", input.categoryId);
    if (input?.tagId) params.append("tag_id", input.tagId);
    if (typeof input?.priceMin === "number")
      params.append("price_min", input.priceMin.toString());
    if (typeof input?.priceMax === "number")
      params.append("price_max", input.priceMax.toString());
    if (typeof input?.page === "number")
      params.append("page", input.page.toString());
    if (typeof input?.pageSize === "number")
      params.append("page_size", input.pageSize.toString());

    const { data } = await this.http.get<{
      services: unknown[];
      total: number;
    }>("/services", { params });

    return {
      services: Array.isArray(data.services)
        ? data.services.map((s) =>
            mapServiceFromApi(s as Record<string, unknown>),
          )
        : [],
      total: data.total ?? 0,
    };
  }

  async searchServices(
    input?: SearchServicesInput,
  ): Promise<SearchServicesOutput> {
    const params = new URLSearchParams();
    if (input?.query) params.append("query", input.query);
    if (input?.categoryId) params.append("category_id", input.categoryId);
    if (input?.tagId) params.append("tag_id", input.tagId);
    if (typeof input?.latitude === "number")
      params.append("latitude", input.latitude.toString());
    if (typeof input?.longitude === "number")
      params.append("longitude", input.longitude.toString());
    if (typeof input?.radiusKm === "number")
      params.append("radius_km", input.radiusKm.toString());
    if (typeof input?.priceMin === "number")
      params.append("price_min", input.priceMin.toString());
    if (typeof input?.priceMax === "number")
      params.append("price_max", input.priceMax.toString());
    if (typeof input?.page === "number")
      params.append("page", input.page.toString());
    if (typeof input?.pageSize === "number")
      params.append("page_size", input.pageSize.toString());

    const { data } = await this.http.get<{
      services: unknown[];
      total: number;
    }>("/services/search", { params });

    return {
      services: Array.isArray(data.services)
        ? data.services.map((s) =>
            mapServiceFromApi(s as Record<string, unknown>),
          )
        : [],
      total: data.total ?? 0,
    };
  }

  async addServicePhoto(
    serviceId: string | number,
    photo: File,
  ): Promise<AddServicePhotoOutput> {
    const formData = new FormData();
    formData.append("file", photo);

    const { data } = await this.http.post<{
      message?: string;
      detail?: string;
      photo?: { id: string; url: string };
    }>(`/services/${serviceId}/photos`, formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });

    return {
      message: data.message,
      detail: data.detail,
      photo: data.photo,
    };
  }

  async deleteServicePhoto(
    serviceId: string | number,
    photoId: string | number,
  ): Promise<DeleteServicePhotoOutput> {
    const { data } = await this.http.delete<{
      message?: string;
      detail?: string;
    }>(`/services/${serviceId}/photos/${photoId}`);

    return {
      message: data.message,
      detail: data.detail,
    };
  }

  async addServiceCategory(
    serviceId: string | number,
    categoryId: string | number,
  ): Promise<AddServiceCategoryOutput> {
    const { data } = await this.http.post<{
      message?: string;
      detail?: string;
    }>(`/services/${serviceId}/categories`, {
      category_id: categoryId,
    });

    return {
      message: data.message,
      detail: data.detail,
    };
  }

  async addServiceTag(
    serviceId: string | number,
    tagId: string | number,
  ): Promise<AddServiceTagOutput> {
    const { data } = await this.http.post<{
      message?: string;
      detail?: string;
    }>(`/services/${serviceId}/tags`, {
      tag_id: tagId,
    });

    return {
      message: data.message,
      detail: data.detail,
    };
  }
}
