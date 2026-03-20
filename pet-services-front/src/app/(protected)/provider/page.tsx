"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import { PROVIDER_KEYS } from "@/application/hooks/provider/provider-query-keys";
import { SERVICE_KEYS } from "@/application/hooks/service/service-query-keys";
import { TAG_KEYS } from "@/application/hooks/tag/tag-query-keys";
import { USER_KEYS } from "@/application/hooks/user/user-query-keys";
import {
  Badge,
  Box,
  Button,
  Dialog,
  Flex,
  Grid,
  HStack,
  Input,
  Portal,
  Text,
  VStack,
  chakra,
} from "@chakra-ui/react";
import {
  type ListTagsOutput,
  useProviderAdd,
  useProviderAddPhoto,
  useCategoryList,
  useProviderDeletePhoto,
  useProviderGet,
  useProviderUpdate,
  type ListServicesOutput,
  useServiceAdd,
  useServiceAddCategory,
  useServiceDeleteCategory,
  useServiceDeleteTag,
  useServiceAddPhoto,
  useServiceAddTag,
  useServiceDelete,
  useServiceDeletePhoto,
  useServiceList,
  useServiceUpdate,
  useTagList,
  useUserProfile,
} from "@/application";
import type { Service } from "@/domain";
import { getApiErrorMessage } from "@/lib/api-error";
import AddProviderForm from "./components/AddProviderForm";
import ServiceFormCard from "./components/ServiceFormCard";
import ServicesListSection from "./components/ServicesListSection";
import ChangePasswordCard from "@/components/account/ChangePasswordCard";
import { MainNav, PageWrapper, ProviderRating } from "@/components/common";

type Feedback = {
  type: "success" | "error";
  message: string;
};

const PROVIDER_PRICE_RANGE_MAX_LENGTH = 10;

export default function ProviderDashboardPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { data: userData, isLoading: isLoadingUser } = useUserProfile();
  const [createdProviderId, setCreatedProviderId] = useState<string | null>(
    null,
  );
  const providerId = userData?.providerId ?? createdProviderId ?? undefined;
  const {
    data: providerData,
    isLoading: isLoadingProvider,
    error: providerError,
  } = useProviderGet(providerId, {
    enabled: Boolean(providerId),
  });
  const { mutateAsync: addProvider, isPending: isAddingProvider } =
    useProviderAdd({
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: USER_KEYS.profile() }),
    });
  const { mutateAsync: updateProvider, isPending: isUpdatingProvider } =
    useProviderUpdate({
      onSuccess: (response, variables) => {
        if (response.provider) {
          queryClient.setQueryData(PROVIDER_KEYS.detail(response.provider.id), {
            provider: response.provider,
          });
        }

        queryClient.invalidateQueries({
          queryKey: PROVIDER_KEYS.detail(variables.providerId),
        });
      },
    });
  const { mutateAsync: addProviderPhoto, isPending: isAddingProviderPhoto } =
    useProviderAddPhoto({
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: PROVIDER_KEYS.all }),
    });
  const {
    mutateAsync: deleteProviderPhoto,
    isPending: isDeletingProviderPhoto,
  } = useProviderDeletePhoto({
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: PROVIDER_KEYS.all }),
  });

  const provider = providerData?.provider;
  const listInput = provider?.id ? { providerId: provider.id } : undefined;

  const {
    data: servicesData,
    isLoading: isLoadingServices,
    error: listServicesError,
  } = useServiceList({
    input: listInput,
    enabled: Boolean(provider?.id),
  });
  const {
    data: categoriesData,
    isLoading: isLoadingCategories,
    error: categoriesError,
  } = useCategoryList();
  const {
    data: tagsData,
    isLoading: isLoadingTags,
    error: tagsError,
  } = useTagList();

  const { mutateAsync: addService, isPending: isAddingService } = useServiceAdd(
    {
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() }),
    },
  );
  const { mutateAsync: updateService, isPending: isUpdatingService } =
    useServiceUpdate({
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() }),
    });
  const { mutateAsync: deleteService, isPending: isDeletingService } =
    useServiceDelete({
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() }),
    });
  const { mutateAsync: addServicePhoto, isPending: isAddingServicePhoto } =
    useServiceAddPhoto({
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() }),
    });
  const { mutateAsync: deleteServicePhoto, isPending: isDeletingServicePhoto } =
    useServiceDeletePhoto({
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() }),
    });
  const {
    mutateAsync: addServiceCategory,
    isPending: isAddingServiceCategory,
  } = useServiceAddCategory({
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() }),
  });
  const {
    mutateAsync: deleteServiceCategory,
    isPending: isDeletingServiceCategory,
  } = useServiceDeleteCategory({
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() }),
  });
  const { mutateAsync: deleteServiceTag, isPending: isDeletingServiceTag } =
    useServiceDeleteTag({
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() }),
    });
  const { mutateAsync: addServiceTag, isPending: isAddingServiceTag } =
    useServiceAddTag({
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() }),
    });

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [price, setPrice] = useState("");
  const [priceMinimum, setPriceMinimum] = useState("");
  const [priceMaximum, setPriceMaximum] = useState("");
  const [duration, setDuration] = useState("");
  const [editingServiceId, setEditingServiceId] = useState<string | null>(null);
  const [confirmDeleteServiceId, setConfirmDeleteServiceId] = useState<
    string | null
  >(null);
  const [deletingServiceId, setDeletingServiceId] = useState<string | null>(
    null,
  );
  const [feedback, setFeedback] = useState<Feedback | null>(null);
  const [providerFeedback, setProviderFeedback] = useState<Feedback | null>(
    null,
  );
  const [providerPhotoFeedback, setProviderPhotoFeedback] =
    useState<Feedback | null>(null);
  const [providerPhotoFile, setProviderPhotoFile] = useState<File | null>(null);
  const [deletingProviderPhotoId, setDeletingProviderPhotoId] = useState<
    string | null
  >(null);
  const [photoFilesByService, setPhotoFilesByService] = useState<
    Record<string, File | null>
  >({});
  const [uploadingPhotoServiceId, setUploadingPhotoServiceId] = useState<
    string | null
  >(null);
  const [deletingPhotoKey, setDeletingPhotoKey] = useState<string | null>(null);
  const [selectedCategoryByService, setSelectedCategoryByService] = useState<
    Record<string, string>
  >({});
  const [selectedTagByService, setSelectedTagByService] = useState<
    Record<string, string>
  >({});
  const [newTagNameByService, setNewTagNameByService] = useState<
    Record<string, string>
  >({});
  const [addingCategoryServiceId, setAddingCategoryServiceId] = useState<
    string | null
  >(null);
  const [removingCategoryKey, setRemovingCategoryKey] = useState<string | null>(
    null,
  );
  const [removingTagKey, setRemovingTagKey] = useState<string | null>(null);
  const [addingTagServiceId, setAddingTagServiceId] = useState<string | null>(
    null,
  );
  const [taxonomyFeedbackByService, setTaxonomyFeedbackByService] = useState<
    Record<string, Feedback | null>
  >({});

  const [providerBusinessName, setProviderBusinessName] = useState("");
  const [providerDescription, setProviderDescription] = useState("");
  const [providerPriceRange, setProviderPriceRange] = useState("");
  const [providerStreet, setProviderStreet] = useState("");
  const [providerAddressNumber, setProviderAddressNumber] = useState("");
  const [providerNeighborhood, setProviderNeighborhood] = useState("");
  const [providerCity, setProviderCity] = useState("");
  const [providerZipCode, setProviderZipCode] = useState("");
  const [providerState, setProviderState] = useState("");
  const [providerCountry, setProviderCountry] = useState("Brasil");
  const [providerComplement, setProviderComplement] = useState("");
  const [providerLatitude, setProviderLatitude] = useState("");
  const [providerLongitude, setProviderLongitude] = useState("");

  const services = servicesData?.services ?? [];
  const categories = categoriesData?.categories ?? [];
  const tags = tagsData?.tags ?? [];
  const isEditing = Boolean(editingServiceId);
  const isSubmitting = isAddingService || isUpdatingService;
  const isLoadingProviderContext = isLoadingUser || isLoadingProvider;
  const shouldShowProviderForm = !isLoadingProviderContext;

  const currentUser = userData?.user;

  useEffect(() => {
    if (!provider) {
      return;
    }

    setProviderBusinessName(provider.businessName ?? "");
    setProviderDescription(provider.description ?? "");
    setProviderPriceRange(provider.priceRange ?? "");
    setProviderStreet(provider.address?.street ?? "");
    setProviderAddressNumber(provider.address?.number ?? "");
    setProviderNeighborhood(provider.address?.neighborhood ?? "");
    setProviderCity(provider.address?.city ?? "");
    setProviderZipCode(provider.address?.zipCode ?? "");
    setProviderState(provider.address?.state ?? "");
    setProviderCountry(provider.address?.country ?? "Brasil");
    setProviderComplement(provider.address?.complement ?? "");
    setProviderLatitude(
      provider.address?.location?.latitude !== undefined
        ? String(provider.address.location.latitude)
        : "",
    );
    setProviderLongitude(
      provider.address?.location?.longitude !== undefined
        ? String(provider.address.location.longitude)
        : "",
    );
  }, [provider]);

  const setTaxonomyFeedback = (
    serviceId: string,
    nextFeedback: Feedback | null,
  ) => {
    setTaxonomyFeedbackByService((previous) => ({
      ...previous,
      [serviceId]: nextFeedback,
    }));
  };

  const handleCategorySelectionChange = (serviceId: string, value: string) => {
    setSelectedCategoryByService((previous) => ({
      ...previous,
      [serviceId]: value,
    }));
  };

  const handleTagSelectionChange = (serviceId: string, value: string) => {
    setSelectedTagByService((previous) => ({
      ...previous,
      [serviceId]: value,
    }));
  };

  const handleNewTagNameChange = (serviceId: string, value: string) => {
    setNewTagNameByService((previous) => ({
      ...previous,
      [serviceId]: value,
    }));
  };

  const handleAddCategoryToService = async (service: Service) => {
    const categoryId = selectedCategoryByService[service.id];

    if (!categoryId) {
      setTaxonomyFeedback(service.id, {
        type: "error",
        message: "Selecione uma categoria para associar.",
      });
      return;
    }

    const alreadyLinked = service.categories.some(
      (category) => category.id === categoryId,
    );

    if (alreadyLinked) {
      setTaxonomyFeedback(service.id, {
        type: "error",
        message: "Essa categoria já está associada ao serviço.",
      });
      return;
    }

    setAddingCategoryServiceId(service.id);

    try {
      const response = await addServiceCategory({
        serviceId: service.id,
        categoryId,
      });

      setSelectedCategoryByService((previous) => ({
        ...previous,
        [service.id]: "",
      }));
      setTaxonomyFeedback(service.id, {
        type: "success",
        message: response.message || "Categoria associada com sucesso.",
      });
    } catch (error) {
      setTaxonomyFeedback(service.id, {
        type: "error",
        message: getApiErrorMessage(
          error,
          "Não foi possível associar a categoria.",
        ),
      });
    } finally {
      setAddingCategoryServiceId(null);
    }
  };

  const handleAddTagToService = async (service: Service) => {
    const tagId = selectedTagByService[service.id];
    const tagName = newTagNameByService[service.id]?.trim() ?? "";

    if (!tagId && !tagName) {
      setTaxonomyFeedback(service.id, {
        type: "error",
        message: "Selecione uma tag existente ou informe uma nova tag.",
      });
      return;
    }

    const alreadyLinkedById =
      Boolean(tagId) && service.tags.some((tag) => tag.id === tagId);
    const alreadyLinkedByName =
      Boolean(tagName) &&
      service.tags.some(
        (tag) => tag.name.trim().toLowerCase() === tagName.toLowerCase(),
      );

    if (alreadyLinkedById || alreadyLinkedByName) {
      setTaxonomyFeedback(service.id, {
        type: "error",
        message: "Essa tag já está associada ao serviço.",
      });
      return;
    }

    setAddingTagServiceId(service.id);

    try {
      const response = await addServiceTag({
        serviceId: service.id,
        payload: {
          tagId: tagId || undefined,
          tagName: tagName || undefined,
        },
      });

      if (response.tag) {
        const addedTag = response.tag;
        queryClient.setQueriesData<ListServicesOutput>(
          { queryKey: SERVICE_KEYS.lists() },
          (previous) => {
            if (!previous) {
              return previous;
            }

            return {
              ...previous,
              services: previous.services.map((currentService) => {
                if (currentService.id !== service.id) {
                  return currentService;
                }

                const alreadyExists = currentService.tags.some(
                  (existingTag) => existingTag.id === addedTag.id,
                );

                if (alreadyExists) {
                  return currentService;
                }

                return {
                  ...currentService,
                  tags: [...currentService.tags, addedTag],
                };
              }),
            };
          },
        );

        queryClient.setQueriesData<ListTagsOutput>(
          { queryKey: TAG_KEYS.lists() },
          (previous) => {
            if (!previous) {
              return previous;
            }

            const alreadyExists = previous.tags.some(
              (currentTag) => currentTag.id === addedTag.id,
            );

            if (alreadyExists) {
              return previous;
            }

            return {
              ...previous,
              tags: [addedTag, ...previous.tags],
              total: previous.total + 1,
            };
          },
        );
      }

      queryClient.invalidateQueries({ queryKey: TAG_KEYS.lists() });

      setSelectedTagByService((previous) => ({
        ...previous,
        [service.id]: "",
      }));
      setNewTagNameByService((previous) => ({
        ...previous,
        [service.id]: "",
      }));
      setTaxonomyFeedback(service.id, {
        type: "success",
        message: response.message || "Tag adicionada/associada com sucesso.",
      });
    } catch (error) {
      setTaxonomyFeedback(service.id, {
        type: "error",
        message: getApiErrorMessage(error, "Não foi possível associar a tag."),
      });
    } finally {
      setAddingTagServiceId(null);
    }
  };

  const handleRemoveCategoryFromService = async (
    service: Service,
    categoryId: string,
  ) => {
    const categoryKey = `${service.id}:${categoryId}`;
    setRemovingCategoryKey(categoryKey);

    try {
      const response = await deleteServiceCategory({
        serviceId: service.id,
        categoryId,
      });

      setTaxonomyFeedback(service.id, {
        type: "success",
        message: response.message || "Categoria removida com sucesso.",
      });
    } catch (error) {
      setTaxonomyFeedback(service.id, {
        type: "error",
        message: getApiErrorMessage(
          error,
          "Não foi possível remover a categoria.",
        ),
      });
    } finally {
      setRemovingCategoryKey(null);
    }
  };

  const handleRemoveTagFromService = async (
    service: Service,
    tagId: string,
  ) => {
    const tagKey = `${service.id}:${tagId}`;
    setRemovingTagKey(tagKey);

    try {
      const response = await deleteServiceTag({
        serviceId: service.id,
        tagId,
      });

      setTaxonomyFeedback(service.id, {
        type: "success",
        message: response.message || "Tag removida com sucesso.",
      });
    } catch (error) {
      setTaxonomyFeedback(service.id, {
        type: "error",
        message: getApiErrorMessage(error, "Não foi possível remover a tag."),
      });
    } finally {
      setRemovingTagKey(null);
    }
  };

  const applyUserDefaultsToProviderForm = () => {
    if (!currentUser) {
      return;
    }

    setProviderBusinessName((prev) => prev || currentUser.name || "");
    setProviderDescription((prev) => prev || "");
    setProviderPriceRange((prev) => prev || "");
    setProviderStreet((prev) => prev || currentUser.address?.street || "");
    setProviderAddressNumber(
      (prev) => prev || currentUser.address?.number || "",
    );
    setProviderNeighborhood(
      (prev) => prev || currentUser.address?.neighborhood || "",
    );
    setProviderCity((prev) => prev || currentUser.address?.city || "");
    setProviderZipCode((prev) => prev || currentUser.address?.zipCode || "");
    setProviderState((prev) => prev || currentUser.address?.state || "");
    setProviderCountry(
      (prev) => prev || currentUser.address?.country || "Brasil",
    );
    setProviderComplement(
      (prev) => prev || currentUser.address?.complement || "",
    );
    setProviderLatitude((prev) => {
      if (prev) return prev;
      const value = currentUser.address?.location?.latitude;
      return typeof value === "number" ? String(value) : "";
    });
    setProviderLongitude((prev) => {
      if (prev) return prev;
      const value = currentUser.address?.location?.longitude;
      return typeof value === "number" ? String(value) : "";
    });
  };

  const handleAddOrUpdateProvider = async (
    event: React.FormEvent<HTMLFormElement>,
  ) => {
    event.preventDefault();

    const parsedLatitude = Number(providerLatitude);
    const parsedLongitude = Number(providerLongitude);

    if (
      !providerBusinessName.trim() ||
      !providerDescription.trim() ||
      !providerPriceRange.trim() ||
      !providerStreet.trim() ||
      !providerAddressNumber.trim() ||
      !providerNeighborhood.trim() ||
      !providerCity.trim() ||
      !providerZipCode.trim() ||
      !providerState.trim() ||
      !providerCountry.trim() ||
      Number.isNaN(parsedLatitude) ||
      Number.isNaN(parsedLongitude)
    ) {
      setProviderFeedback({
        type: "error",
        message:
          "Preencha todos os campos obrigatórios para cadastrar o provider.",
      });
      return;
    }

    const normalizedPriceRange = providerPriceRange.trim();

    if (normalizedPriceRange.length > PROVIDER_PRICE_RANGE_MAX_LENGTH) {
      setProviderFeedback({
        type: "error",
        message:
          "A faixa de preço deve ter no máximo 10 caracteres (ex: 80-180).",
      });
      return;
    }

    try {
      const providerPayload = {
        businessName: providerBusinessName.trim(),
        description: providerDescription.trim(),
        priceRange: normalizedPriceRange,
        address: {
          street: providerStreet.trim(),
          number: providerAddressNumber.trim(),
          neighborhood: providerNeighborhood.trim(),
          city: providerCity.trim(),
          zipCode: providerZipCode.trim(),
          state: providerState.trim(),
          country: providerCountry.trim(),
          complement: providerComplement.trim(),
          location: {
            latitude: parsedLatitude,
            longitude: parsedLongitude,
          },
        },
      };

      const response = provider?.id
        ? await updateProvider({
            providerId: provider.id,
            ...providerPayload,
          })
        : await addProvider(providerPayload);

      if (response.provider?.id) {
        setCreatedProviderId(response.provider.id);
        queryClient.setQueryData(PROVIDER_KEYS.detail(response.provider.id), {
          provider: response.provider,
        });
      }

      setProviderFeedback({
        type: "success",
        message:
          response.detail ||
          response.message ||
          (provider?.id
            ? "Provider atualizado com sucesso."
            : "Provider cadastrado com sucesso."),
      });
    } catch (error) {
      setProviderFeedback({
        type: "error",
        message: getApiErrorMessage(
          error,
          provider?.id
            ? "Não foi possível atualizar o provider."
            : "Não foi possível cadastrar o provider.",
        ),
      });
    }
  };

  const resetForm = () => {
    setName("");
    setDescription("");
    setPrice("");
    setPriceMinimum("");
    setPriceMaximum("");
    setDuration("");
    setEditingServiceId(null);
  };

  const fillFormForEdit = (service: Service) => {
    setEditingServiceId(service.id);
    setName(service.name);
    setDescription(service.description);
    setPrice(String(service.price));
    setPriceMinimum(String(service.priceMinimum));
    setPriceMaximum(String(service.priceMaximum));
    setDuration(String(service.duration));
    setFeedback(null);
  };

  const handleCancelEdit = () => {
    resetForm();
    setFeedback(null);
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!provider?.id) {
      setFeedback({
        type: "error",
        message:
          "Não foi possível identificar o provider logado para salvar o serviço.",
      });
      return;
    }

    const hasFixedPrice = price.trim() !== "";
    const hasPriceMinimum = priceMinimum.trim() !== "";
    const hasPriceMaximum = priceMaximum.trim() !== "";
    const hasPriceRange = hasPriceMinimum || hasPriceMaximum;

    const parsedPrice = Number(price);
    const parsedPriceMinimum = Number(priceMinimum);
    const parsedPriceMaximum = Number(priceMaximum);
    const parsedDuration = Number(duration);

    if (!name.trim() || !description.trim() || Number.isNaN(parsedDuration)) {
      setFeedback({
        type: "error",
        message: "Preencha todos os campos obrigatórios do serviço.",
      });
      return;
    }

    if (hasFixedPrice && hasPriceRange) {
      setFeedback({
        type: "error",
        message:
          "Use preço fixo ou faixa de preço, não ambos no mesmo serviço.",
      });
      return;
    }

    if (!hasFixedPrice && !hasPriceRange) {
      setFeedback({
        type: "error",
        message:
          "Informe preço fixo ou faixa de preço para cadastrar o serviço.",
      });
      return;
    }

    if (hasFixedPrice) {
      if (Number.isNaN(parsedPrice) || parsedPrice < 0) {
        setFeedback({
          type: "error",
          message: "Informe um preço fixo válido para o serviço.",
        });
        return;
      }
    }

    if (hasPriceRange) {
      if (!hasPriceMinimum || !hasPriceMaximum) {
        setFeedback({
          type: "error",
          message: "Para faixa de preço, informe preço mínimo e preço máximo.",
        });
        return;
      }

      if (
        Number.isNaN(parsedPriceMinimum) ||
        Number.isNaN(parsedPriceMaximum) ||
        parsedPriceMinimum < 0 ||
        parsedPriceMaximum < 0
      ) {
        setFeedback({
          type: "error",
          message: "Informe uma faixa de preço válida para o serviço.",
        });
        return;
      }

      if (parsedPriceMinimum > parsedPriceMaximum) {
        setFeedback({
          type: "error",
          message: "O preço mínimo não pode ser maior que o preço máximo.",
        });
        return;
      }
    }

    try {
      if (editingServiceId) {
        const response = await updateService({
          serviceId: editingServiceId,
          name: name.trim(),
          description: description.trim(),
          price: hasFixedPrice ? parsedPrice : 0,
          priceMinimum: hasPriceRange ? parsedPriceMinimum : 0,
          priceMaximum: hasPriceRange ? parsedPriceMaximum : 0,
          duration: parsedDuration,
        });

        setFeedback({
          type: "success",
          message:
            response.detail ||
            response.message ||
            "Serviço atualizado com sucesso.",
        });
      } else {
        const response = await addService({
          providerId: provider.id,
          name: name.trim(),
          description: description.trim(),
          price: hasFixedPrice ? parsedPrice : undefined,
          priceMinimum: hasPriceRange ? parsedPriceMinimum : undefined,
          priceMaximum: hasPriceRange ? parsedPriceMaximum : undefined,
          duration: parsedDuration,
        });

        setFeedback({
          type: "success",
          message:
            response.detail ||
            response.message ||
            "Serviço criado com sucesso.",
        });
      }

      resetForm();
    } catch (error) {
      setFeedback({
        type: "error",
        message: getApiErrorMessage(
          error,
          "Não foi possível salvar o serviço.",
        ),
      });
    }
  };

  const handleDeleteService = async (serviceId: string) => {
    setConfirmDeleteServiceId(serviceId);
  };

  const handleCancelDeleteService = () => {
    setConfirmDeleteServiceId(null);
  };

  const handleConfirmDeleteService = async () => {
    if (!confirmDeleteServiceId) {
      return;
    }

    const serviceId = confirmDeleteServiceId;
    setConfirmDeleteServiceId(null);

    setDeletingServiceId(serviceId);
    try {
      const response = await deleteService(serviceId);
      setFeedback({
        type: "success",
        message:
          response.detail ||
          response.message ||
          "Serviço excluído com sucesso.",
      });

      if (editingServiceId === serviceId) {
        resetForm();
      }
    } catch (error) {
      setFeedback({
        type: "error",
        message: getApiErrorMessage(
          error,
          "Não foi possível excluir o serviço.",
        ),
      });
    } finally {
      setDeletingServiceId(null);
    }
  };

  const handlePhotoFileChange = (
    serviceId: string,
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const file = event.target.files?.[0] ?? null;
    setPhotoFilesByService((prev) => ({
      ...prev,
      [serviceId]: file,
    }));
  };

  const handleUploadServicePhoto = async (serviceId: string) => {
    const photo = photoFilesByService[serviceId];

    if (!photo) {
      setFeedback({
        type: "error",
        message: "Selecione uma imagem antes de enviar a foto do serviço.",
      });
      return;
    }

    setUploadingPhotoServiceId(serviceId);
    try {
      const response = await addServicePhoto({
        serviceId,
        photo,
      });

      setFeedback({
        type: "success",
        message:
          response.detail || response.message || "Foto enviada com sucesso.",
      });

      setPhotoFilesByService((prev) => ({
        ...prev,
        [serviceId]: null,
      }));
    } catch (error) {
      setFeedback({
        type: "error",
        message: getApiErrorMessage(
          error,
          "Não foi possível enviar a foto do serviço.",
        ),
      });
    } finally {
      setUploadingPhotoServiceId(null);
    }
  };

  const handleDeleteServicePhoto = async (
    serviceId: string,
    photoId: string,
  ) => {
    const currentDeletingPhotoKey = `${serviceId}:${photoId}`;
    setDeletingPhotoKey(currentDeletingPhotoKey);

    try {
      const response = await deleteServicePhoto({
        serviceId,
        photoId,
      });

      setFeedback({
        type: "success",
        message:
          response.detail || response.message || "Foto removida com sucesso.",
      });
    } catch (error) {
      setFeedback({
        type: "error",
        message: getApiErrorMessage(
          error,
          "Não foi possível remover a foto do serviço.",
        ),
      });
    } finally {
      setDeletingPhotoKey(null);
    }
  };

  const handleProviderPhotoFileChange = (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const file = event.target.files?.[0] ?? null;
    setProviderPhotoFile(file);
  };

  const handleUploadProviderPhoto = async () => {
    if (!provider?.id) {
      setProviderPhotoFeedback({
        type: "error",
        message: "Não foi possível identificar o provider para envio da foto.",
      });
      return;
    }

    if (!providerPhotoFile) {
      setProviderPhotoFeedback({
        type: "error",
        message: "Selecione uma imagem antes de enviar a foto do provider.",
      });
      return;
    }

    try {
      const response = await addProviderPhoto({
        providerId: provider.id,
        photo: providerPhotoFile,
      });

      setProviderPhotoFeedback({
        type: "success",
        message:
          response.detail || response.message || "Foto enviada com sucesso.",
      });
      setProviderPhotoFile(null);
    } catch (error) {
      setProviderPhotoFeedback({
        type: "error",
        message: getApiErrorMessage(
          error,
          "Não foi possível enviar a foto do provider.",
        ),
      });
    }
  };

  const handleDeleteProviderPhoto = async (photoId: string) => {
    if (!provider?.id) {
      return;
    }

    setDeletingProviderPhotoId(photoId);
    try {
      const response = await deleteProviderPhoto({
        providerId: provider.id,
        photoId,
      });

      setProviderPhotoFeedback({
        type: "success",
        message:
          response.detail || response.message || "Foto removida com sucesso.",
      });
    } catch (error) {
      setProviderPhotoFeedback({
        type: "error",
        message: getApiErrorMessage(
          error,
          "Não foi possível remover a foto do provider.",
        ),
      });
    } finally {
      setDeletingProviderPhotoId(null);
    }
  };

  const listServicesErrorMessage = listServicesError
    ? getApiErrorMessage(
        listServicesError,
        "Não foi possível carregar os serviços.",
      )
    : null;

  const categoriesErrorMessage = categoriesError
    ? getApiErrorMessage(
        categoriesError,
        "Não foi possível carregar as categorias.",
      )
    : null;

  const tagsErrorMessage = tagsError
    ? getApiErrorMessage(tagsError, "Não foi possível carregar as tags.")
    : null;

  return (
    <PageWrapper gap={10}>
      <MainNav />

      <VStack align="stretch" gap={6}>
        <Box
          borderRadius={{ base: "2xl", md: "3xl" }}
          bg="white"
          p={{ base: 4, sm: 5, md: 7 }}
          borderWidth="1px"
          borderColor="gray.200"
          shadow="sm"
        >
          <Flex
            mb={6}
            wrap="wrap"
            align="start"
            justify="space-between"
            gap={4}
            direction={{ base: "column", sm: "row" }}
          >
            <Box>
              <Text
                fontSize={{ base: "xs" }}
                fontWeight="semibold"
                textTransform="uppercase"
                color="green.500"
              >
                Dashboard
              </Text>
              <Text
                mt={2}
                fontSize={{ base: "xl", sm: "2xl" }}
                fontWeight="semibold"
                color="gray.900"
              >
                Olá, prestador
              </Text>
              <Text mt={2} fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
                Gerencie os serviços que sua empresa oferece para os clientes.
              </Text>
            </Box>

            <Badge
              borderRadius="full"
              px={3}
              py={1}
              fontSize={{ base: "xs" }}
              colorPalette={provider?.id ? "green" : "orange"}
              variant="subtle"
            >
              {provider?.id ? "Provider ativo" : "Sem provider"}
            </Badge>
          </Flex>

          <Grid gap={4} templateColumns={{ base: "1fr", md: "1fr 1fr" }}>
            <Box
              borderRadius="2xl"
              borderWidth="1px"
              borderColor="gray.200"
              bg="gray.50"
              p={4}
            >
              <Text fontSize="sm" fontWeight="semibold" color="gray.900">
                Serviços cadastrados
              </Text>
              <Text mt={2} fontSize="2xl" fontWeight="bold" color="gray.900">
                {services.length}
              </Text>
            </Box>

            <Box
              borderRadius="2xl"
              borderWidth="1px"
              borderColor="gray.200"
              bg="gray.50"
              p={4}
            >
              <Text fontSize="sm" fontWeight="semibold" color="gray.900">
                Contexto do provider
              </Text>
              <Text mt={2} fontSize="sm" color="gray.600">
                {provider?.businessName
                  ? `Provider identificado: ${provider.businessName}`
                  : "Aguardando identificação do provider..."}
              </Text>
              {provider ? (
                <ProviderRating
                  mt={2}
                  rating={provider.averageRating}
                  labelPrefix="Avaliação média:"
                  fontSize="xs"
                />
              ) : null}
              {providerError ? (
                <Text mt={1.5} fontSize="xs" color="red.600">
                  {getApiErrorMessage(
                    providerError,
                    "Não foi possível carregar os dados do provider.",
                  )}
                </Text>
              ) : null}

              <Button
                mt={3}
                size="sm"
                borderRadius="full"
                variant="outline"
                disabled={!provider?.id}
                onClick={() => {
                  if (!provider?.id) return;
                  router.push(`/providers/${provider.id}`);
                }}
              >
                Ver página do provider
              </Button>
            </Box>
          </Grid>
        </Box>

        {shouldShowProviderForm ? (
          <AddProviderForm
            onSubmit={handleAddOrUpdateProvider}
            onPrefillFromProfile={applyUserDefaultsToProviderForm}
            canPrefillFromProfile={Boolean(currentUser)}
            isSubmitting={isAddingProvider || isUpdatingProvider}
            isLoadingProviderContext={isLoadingProviderContext}
            maxPriceRangeLength={PROVIDER_PRICE_RANGE_MAX_LENGTH}
            feedback={providerFeedback}
            mode={provider?.id ? "update" : "create"}
            businessName={providerBusinessName}
            onBusinessNameChange={setProviderBusinessName}
            priceRange={providerPriceRange}
            onPriceRangeChange={setProviderPriceRange}
            description={providerDescription}
            onDescriptionChange={setProviderDescription}
            street={providerStreet}
            onStreetChange={setProviderStreet}
            addressNumber={providerAddressNumber}
            onAddressNumberChange={setProviderAddressNumber}
            neighborhood={providerNeighborhood}
            onNeighborhoodChange={setProviderNeighborhood}
            city={providerCity}
            onCityChange={setProviderCity}
            zipCode={providerZipCode}
            onZipCodeChange={setProviderZipCode}
            state={providerState}
            onStateChange={setProviderState}
            country={providerCountry}
            onCountryChange={setProviderCountry}
            latitude={providerLatitude}
            onLatitudeChange={setProviderLatitude}
            longitude={providerLongitude}
            onLongitudeChange={setProviderLongitude}
            complement={providerComplement}
            onComplementChange={setProviderComplement}
          />
        ) : null}

        {provider?.id ? (
          <Box
            borderRadius="3xl"
            bg="white"
            p={{ base: 5, md: 6 }}
            borderWidth="1px"
            borderColor="gray.200"
            shadow="sm"
          >
            <Box mb={4}>
              <Text
                fontSize="xs"
                fontWeight="semibold"
                textTransform="uppercase"
                color="green.500"
              >
                Provider
              </Text>
              <Text mt={2} fontSize="xl" fontWeight="semibold" color="gray.900">
                Fotos do provider
              </Text>
              <Text mt={2} fontSize="sm" color="gray.600">
                Gerencie as imagens que aparecem na vitrine do seu perfil.
              </Text>
            </Box>

            {(provider.photos ?? []).length ? (
              <Flex mt={3} gap={3} wrap="wrap">
                {(provider.photos ?? []).map((photo) => {
                  const photoId = String(photo.id);
                  const isCurrentDeletingPhoto =
                    deletingProviderPhotoId === photoId;

                  return (
                    <Box
                      key={photoId}
                      borderWidth="1px"
                      borderColor="gray.200"
                      borderRadius="xl"
                      bg="gray.50"
                      p={2}
                      w={{ base: "full", sm: "140px" }}
                    >
                      <chakra.img
                        src={photo.url}
                        alt={`Foto do provider ${provider.businessName || ""}`}
                        w="full"
                        h="124px"
                        objectFit="cover"
                        borderRadius="lg"
                      />
                      <Button
                        mt={2}
                        size="xs"
                        borderRadius="full"
                        colorPalette="red"
                        variant="subtle"
                        w="full"
                        onClick={() => handleDeleteProviderPhoto(photoId)}
                        disabled={
                          isDeletingProviderPhoto ||
                          isCurrentDeletingPhoto ||
                          isAddingProviderPhoto
                        }
                      >
                        {isCurrentDeletingPhoto ? "Removendo..." : "Remover"}
                      </Button>
                    </Box>
                  );
                })}
              </Flex>
            ) : (
              <Text mt={3} fontSize="sm" color="gray.500">
                Nenhuma foto cadastrada.
              </Text>
            )}

            <HStack
              mt={4}
              align={{ base: "stretch", md: "end" }}
              gap={3}
              flexWrap="wrap"
              flexDir={{ base: "column", md: "row" }}
            >
              <Box minW={{ base: "full", md: "420px" }}>
                <Input
                  type="file"
                  accept="image/*"
                  onChange={handleProviderPhotoFileChange}
                  h="10"
                  bg="white"
                  borderRadius="xl"
                  borderColor="gray.200"
                  p={1}
                  disabled={isAddingProviderPhoto || isDeletingProviderPhoto}
                />
                {providerPhotoFile ? (
                  <Text mt={1} fontSize="xs" color="gray.500">
                    Arquivo selecionado: {providerPhotoFile.name}
                  </Text>
                ) : null}
              </Box>

              <Button
                size="sm"
                borderRadius="full"
                bg="green.400"
                color="white"
                _hover={{ bg: "green.500" }}
                onClick={handleUploadProviderPhoto}
                disabled={
                  isAddingProviderPhoto ||
                  isDeletingProviderPhoto ||
                  !providerPhotoFile
                }
              >
                {isAddingProviderPhoto ? "Enviando..." : "Enviar foto"}
              </Button>
            </HStack>

            {providerPhotoFeedback ? (
              <Text
                mt={3}
                fontSize="xs"
                color={
                  providerPhotoFeedback.type === "error"
                    ? "red.600"
                    : "green.600"
                }
              >
                {providerPhotoFeedback.message}
              </Text>
            ) : null}
          </Box>
        ) : null}

        {provider?.id ? (
          <ServiceFormCard
            isEditing={isEditing}
            onSubmit={handleSubmit}
            name={name}
            onNameChange={setName}
            duration={duration}
            onDurationChange={setDuration}
            price={price}
            onPriceChange={setPrice}
            priceMinimum={priceMinimum}
            onPriceMinimumChange={setPriceMinimum}
            priceMaximum={priceMaximum}
            onPriceMaximumChange={setPriceMaximum}
            description={description}
            onDescriptionChange={setDescription}
            isSubmitting={isSubmitting}
            canSubmit={Boolean(provider?.id) && !isLoadingProviderContext}
            onCancelEdit={handleCancelEdit}
            feedback={feedback}
          />
        ) : null}

        <ServicesListSection
          isLoadingProviderContext={isLoadingProviderContext}
          isLoadingServices={isLoadingServices}
          hasProvider={Boolean(provider?.id)}
          listServicesErrorMessage={listServicesErrorMessage}
          services={services}
          deletingServiceId={deletingServiceId}
          isSubmitting={isSubmitting}
          isDeletingService={isDeletingService}
          confirmDeleteServiceId={confirmDeleteServiceId}
          onViewDetails={(serviceId) =>
            router.push(`/services/${serviceId}?from=/provider`)
          }
          onEditService={fillFormForEdit}
          onDeleteService={handleDeleteService}
          photoFilesByService={photoFilesByService}
          uploadingPhotoServiceId={uploadingPhotoServiceId}
          deletingPhotoKey={deletingPhotoKey}
          isAddingServicePhoto={isAddingServicePhoto}
          isDeletingServicePhoto={isDeletingServicePhoto}
          onPhotoFileChange={handlePhotoFileChange}
          onUploadServicePhoto={handleUploadServicePhoto}
          onDeleteServicePhoto={handleDeleteServicePhoto}
          selectedCategoryByService={selectedCategoryByService}
          selectedTagByService={selectedTagByService}
          newTagNameByService={newTagNameByService}
          addingCategoryServiceId={addingCategoryServiceId}
          addingTagServiceId={addingTagServiceId}
          removingCategoryKey={removingCategoryKey}
          removingTagKey={removingTagKey}
          taxonomyFeedbackByService={taxonomyFeedbackByService}
          categories={categories}
          tags={tags}
          isLoadingCategories={isLoadingCategories}
          isLoadingTags={isLoadingTags}
          isAddingServiceCategory={isAddingServiceCategory}
          isDeletingServiceCategory={isDeletingServiceCategory}
          isAddingServiceTag={isAddingServiceTag}
          isDeletingServiceTag={isDeletingServiceTag}
          categoriesErrorMessage={categoriesErrorMessage}
          tagsErrorMessage={tagsErrorMessage}
          onCategorySelectionChange={handleCategorySelectionChange}
          onTagSelectionChange={handleTagSelectionChange}
          onNewTagNameChange={handleNewTagNameChange}
          onAddCategoryToService={handleAddCategoryToService}
          onRemoveCategoryFromService={handleRemoveCategoryFromService}
          onAddTagToService={handleAddTagToService}
          onRemoveTagFromService={handleRemoveTagFromService}
        />
      </VStack>

      <Dialog.Root
        open={Boolean(confirmDeleteServiceId)}
        onOpenChange={(details) => {
          if (!details.open) {
            handleCancelDeleteService();
          }
        }}
      >
        <Portal>
          <Dialog.Backdrop bg="blackAlpha.500" />
          <Dialog.Positioner p={4}>
            <Dialog.Content borderRadius="3xl" maxW="md" w="full">
              <Dialog.Header>
                <Dialog.Title>Excluir serviço?</Dialog.Title>
              </Dialog.Header>
              <Dialog.Body>
                <Text fontSize="sm" color="gray.600">
                  Esta ação remove o serviço permanentemente.
                </Text>
              </Dialog.Body>
              <Dialog.Footer>
                <Button
                  type="button"
                  variant="outline"
                  borderRadius="full"
                  onClick={handleCancelDeleteService}
                  disabled={isDeletingService}
                >
                  Cancelar
                </Button>
                <Button
                  type="button"
                  borderRadius="full"
                  bg="red.500"
                  color="white"
                  onClick={handleConfirmDeleteService}
                  disabled={isDeletingService}
                  _hover={{ bg: "red.600" }}
                  _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
                >
                  {isDeletingService ? "Excluindo..." : "Excluir"}
                </Button>
              </Dialog.Footer>
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>

      <ChangePasswordCard />
    </PageWrapper>
  );
}
