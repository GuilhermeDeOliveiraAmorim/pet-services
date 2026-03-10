"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import {
  Box,
  Button,
  Dialog,
  Flex,
  Grid,
  HStack,
  Input,
  NativeSelect,
  Portal,
  Text,
  Textarea,
  VStack,
  chakra,
} from "@chakra-ui/react";
import {
  useProviderAdd,
  useCategoryList,
  useProviderGet,
  useServiceAdd,
  useServiceAddCategory,
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
import ChangePasswordCard from "@/components/account/ChangePasswordCard";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

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
        queryClient.invalidateQueries({ queryKey: ["user-profile"] }),
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
        queryClient.invalidateQueries({ queryKey: ["services"] }),
    },
  );
  const { mutateAsync: updateService, isPending: isUpdatingService } =
    useServiceUpdate({
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: ["services"] }),
    });
  const { mutateAsync: deleteService, isPending: isDeletingService } =
    useServiceDelete({
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: ["services"] }),
    });
  const { mutateAsync: addServicePhoto, isPending: isAddingServicePhoto } =
    useServiceAddPhoto({
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: ["services"] }),
    });
  const { mutateAsync: deleteServicePhoto, isPending: isDeletingServicePhoto } =
    useServiceDeletePhoto({
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: ["services"] }),
    });
  const {
    mutateAsync: addServiceCategory,
    isPending: isAddingServiceCategory,
  } = useServiceAddCategory({
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["services"] }),
  });
  const { mutateAsync: addServiceTag, isPending: isAddingServiceTag } =
    useServiceAddTag({
      onSuccess: () =>
        queryClient.invalidateQueries({ queryKey: ["services"] }),
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
  const [addingCategoryServiceId, setAddingCategoryServiceId] = useState<
    string | null
  >(null);
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
  const shouldShowAddProviderForm =
    !isLoadingProviderContext && !provider?.id && !providerId;

  const currentUser = userData?.user;

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

    if (!tagId) {
      setTaxonomyFeedback(service.id, {
        type: "error",
        message: "Selecione uma tag para associar.",
      });
      return;
    }

    const alreadyLinked = service.tags.some((tag) => tag.id === tagId);

    if (alreadyLinked) {
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
        tagId,
      });

      setSelectedTagByService((previous) => ({
        ...previous,
        [service.id]: "",
      }));
      setTaxonomyFeedback(service.id, {
        type: "success",
        message: response.message || "Tag associada com sucesso.",
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

  const handleAddProvider = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (provider?.id) {
      setProviderFeedback({
        type: "error",
        message: "Você já possui um provider cadastrado.",
      });
      return;
    }

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
      const response = await addProvider({
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
      });

      if (response.provider?.id) {
        setCreatedProviderId(response.provider.id);
        queryClient.invalidateQueries({
          queryKey: ["provider", response.provider.id],
        });
      }

      setProviderFeedback({
        type: "success",
        message:
          response.detail ||
          response.message ||
          "Provider cadastrado com sucesso.",
      });
    } catch (error) {
      setProviderFeedback({
        type: "error",
        message: getApiErrorMessage(
          error,
          "Não foi possível cadastrar o provider.",
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

  return (
    <PageWrapper gap={10}>
      <MainNav />

      <VStack align="stretch" gap={6}>
        <Box>
          <Text
            fontSize="xs"
            fontWeight="semibold"
            textTransform="uppercase"
            color="green.500"
          >
            Dashboard
          </Text>
          <Text mt={2} fontSize="2xl" fontWeight="semibold" color="gray.900">
            Olá, prestador
          </Text>
          <Text mt={2} fontSize="sm" color="gray.600">
            Gerencie os serviços que sua empresa oferece para os clientes.
          </Text>
        </Box>

        <Grid gap={4} templateColumns={{ base: "1fr", md: "1fr 1fr" }}>
          <Box
            borderRadius="2xl"
            borderWidth="1px"
            borderColor="gray.200"
            bg="white"
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
            bg="white"
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
            {providerError ? (
              <Text mt={1.5} fontSize="xs" color="red.600">
                {getApiErrorMessage(
                  providerError,
                  "Não foi possível carregar os dados do provider.",
                )}
              </Text>
            ) : null}
          </Box>
        </Grid>

        {shouldShowAddProviderForm ? (
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
                Cadastrar provider
              </Text>
              <Text mt={2} fontSize="sm" color="gray.600">
                Antes de criar serviços, você precisa cadastrar os dados do seu
                provider.
              </Text>
            </Box>

            <chakra.form onSubmit={handleAddProvider}>
              <VStack align="stretch" gap={4}>
                <HStack justify="space-between" flexWrap="wrap">
                  <Text fontSize="sm" color="gray.600">
                    Use seus dados do perfil como base e ajuste se necessário.
                  </Text>
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    borderRadius="full"
                    onClick={applyUserDefaultsToProviderForm}
                    disabled={!currentUser || isAddingProvider}
                  >
                    Preencher com dados do perfil
                  </Button>
                </HStack>

                <Grid gap={4} templateColumns={{ base: "1fr", md: "1fr 1fr" }}>
                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      Nome comercial
                    </Text>
                    <Input
                      value={providerBusinessName}
                      onChange={(event) =>
                        setProviderBusinessName(event.target.value)
                      }
                      h="11"
                      borderRadius="xl"
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="Ex: Clínica Pet Saúde"
                      required
                    />
                  </Box>

                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      Faixa de preço
                    </Text>
                    <Input
                      value={providerPriceRange}
                      onChange={(event) =>
                        setProviderPriceRange(
                          event.target.value.slice(
                            0,
                            PROVIDER_PRICE_RANGE_MAX_LENGTH,
                          ),
                        )
                      }
                      maxLength={PROVIDER_PRICE_RANGE_MAX_LENGTH}
                      h="11"
                      borderRadius="xl"
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="Ex: 80-180"
                      required
                    />
                    <Text mt={1.5} fontSize="xs" color="gray.500">
                      Máximo de {PROVIDER_PRICE_RANGE_MAX_LENGTH} caracteres.
                    </Text>
                  </Box>

                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      Descrição
                    </Text>
                    <Textarea
                      value={providerDescription}
                      onChange={(event) =>
                        setProviderDescription(event.target.value)
                      }
                      minH="20"
                      borderRadius="xl"
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="Descreva o seu negócio e os serviços prestados"
                      required
                    />
                  </Box>

                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      Rua
                    </Text>
                    <Input
                      value={providerStreet}
                      onChange={(event) =>
                        setProviderStreet(event.target.value)
                      }
                      h="11"
                      borderRadius="xl"
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="Rua"
                      required
                    />
                  </Box>

                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      Número do endereço
                    </Text>
                    <Input
                      value={providerAddressNumber}
                      onChange={(event) =>
                        setProviderAddressNumber(event.target.value)
                      }
                      h="11"
                      borderRadius="xl"
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="123"
                      required
                    />
                  </Box>

                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      Bairro
                    </Text>
                    <Input
                      value={providerNeighborhood}
                      onChange={(event) =>
                        setProviderNeighborhood(event.target.value)
                      }
                      h="11"
                      borderRadius="xl"
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="Bairro"
                      required
                    />
                  </Box>

                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      Cidade
                    </Text>
                    <Input
                      value={providerCity}
                      onChange={(event) => setProviderCity(event.target.value)}
                      h="11"
                      borderRadius="xl"
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="Cidade"
                      required
                    />
                  </Box>

                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      CEP
                    </Text>
                    <Input
                      value={providerZipCode}
                      onChange={(event) =>
                        setProviderZipCode(event.target.value)
                      }
                      h="11"
                      borderRadius="xl"
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="57000-000"
                      required
                    />
                  </Box>

                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      Estado
                    </Text>
                    <Input
                      value={providerState}
                      onChange={(event) => setProviderState(event.target.value)}
                      h="11"
                      borderRadius="xl"
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="AL"
                      required
                    />
                  </Box>

                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      País
                    </Text>
                    <Input
                      value={providerCountry}
                      onChange={(event) =>
                        setProviderCountry(event.target.value)
                      }
                      h="11"
                      borderRadius="xl"
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="Brasil"
                      required
                    />
                  </Box>

                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      Latitude
                    </Text>
                    <Input
                      type="number"
                      inputMode="decimal"
                      value={providerLatitude}
                      onChange={(event) =>
                        setProviderLatitude(event.target.value)
                      }
                      h="11"
                      borderRadius="xl"
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="-9.6498"
                      step="0.000001"
                      required
                    />
                  </Box>

                  <Box minW={0}>
                    <Text
                      fontSize="sm"
                      fontWeight="medium"
                      color="gray.700"
                      mb={2}
                    >
                      Longitude
                    </Text>
                    <Input
                      type="number"
                      inputMode="decimal"
                      value={providerLongitude}
                      onChange={(event) =>
                        setProviderLongitude(event.target.value)
                      }
                      h="11"
                      borderRadius="xl"
                      bg="gray.50"
                      borderColor="gray.200"
                      focusRingColor="teal.200"
                      placeholder="-35.7089"
                      step="0.000001"
                      required
                    />
                  </Box>
                </Grid>

                <Box minW={0}>
                  <Text
                    fontSize="sm"
                    fontWeight="medium"
                    color="gray.700"
                    mb={2}
                  >
                    Complemento (opcional)
                  </Text>
                  <Input
                    value={providerComplement}
                    onChange={(event) =>
                      setProviderComplement(event.target.value)
                    }
                    h="11"
                    borderRadius="xl"
                    bg="gray.50"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                    placeholder="Sala, bloco, referência"
                  />
                </Box>

                <HStack gap={3} flexWrap="wrap">
                  <Button
                    type="submit"
                    borderRadius="full"
                    bg="green.400"
                    color="white"
                    _hover={{ bg: "green.500" }}
                    _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
                    disabled={isAddingProvider || isLoadingProviderContext}
                  >
                    {isAddingProvider
                      ? "Cadastrando provider..."
                      : "Cadastrar provider"}
                  </Button>
                </HStack>

                {providerFeedback ? (
                  <Text
                    fontSize="xs"
                    color={
                      providerFeedback.type === "error"
                        ? "red.600"
                        : "green.600"
                    }
                  >
                    {providerFeedback.message}
                  </Text>
                ) : null}
              </VStack>
            </chakra.form>
          </Box>
        ) : null}

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
              Serviços
            </Text>
            <Text mt={2} fontSize="xl" fontWeight="semibold" color="gray.900">
              {isEditing ? "Editar serviço" : "Cadastrar serviço"}
            </Text>
          </Box>

          <chakra.form onSubmit={handleSubmit}>
            <VStack align="stretch" gap={4}>
              <Grid gap={4} templateColumns={{ base: "1fr", md: "1fr 1fr" }}>
                <Box minW={0}>
                  <Text
                    fontSize="sm"
                    fontWeight="medium"
                    color="gray.700"
                    mb={2}
                  >
                    Nome
                  </Text>
                  <Input
                    value={name}
                    onChange={(event) => setName(event.target.value)}
                    h="11"
                    borderRadius="xl"
                    bg="gray.50"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                    placeholder="Ex: Consulta veterinária"
                    required
                  />
                </Box>

                <Box minW={0}>
                  <Text
                    fontSize="sm"
                    fontWeight="medium"
                    color="gray.700"
                    mb={2}
                  >
                    Duração (minutos)
                  </Text>
                  <Input
                    type="number"
                    inputMode="numeric"
                    value={duration}
                    onChange={(event) => setDuration(event.target.value)}
                    h="11"
                    borderRadius="xl"
                    bg="gray.50"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                    placeholder="Ex: 60"
                    min={1}
                    required
                  />
                </Box>

                <Box minW={0}>
                  <Text
                    fontSize="sm"
                    fontWeight="medium"
                    color="gray.700"
                    mb={2}
                  >
                    Preço base (R$)
                  </Text>
                  <Input
                    type="number"
                    inputMode="decimal"
                    value={price}
                    onChange={(event) => setPrice(event.target.value)}
                    h="11"
                    borderRadius="xl"
                    bg="gray.50"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                    placeholder="Ex: 120"
                    min={0}
                    step="0.01"
                  />
                </Box>

                <Box minW={0}>
                  <Text
                    fontSize="sm"
                    fontWeight="medium"
                    color="gray.700"
                    mb={2}
                  >
                    Preço mínimo (R$)
                  </Text>
                  <Input
                    type="number"
                    inputMode="decimal"
                    value={priceMinimum}
                    onChange={(event) => setPriceMinimum(event.target.value)}
                    h="11"
                    borderRadius="xl"
                    bg="gray.50"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                    placeholder="Ex: 100"
                    min={0}
                    step="0.01"
                  />
                </Box>

                <Box minW={0}>
                  <Text
                    fontSize="sm"
                    fontWeight="medium"
                    color="gray.700"
                    mb={2}
                  >
                    Preço máximo (R$)
                  </Text>
                  <Input
                    type="number"
                    inputMode="decimal"
                    value={priceMaximum}
                    onChange={(event) => setPriceMaximum(event.target.value)}
                    h="11"
                    borderRadius="xl"
                    bg="gray.50"
                    borderColor="gray.200"
                    focusRingColor="teal.200"
                    placeholder="Ex: 180"
                    min={0}
                    step="0.01"
                  />
                </Box>
              </Grid>

              <Text fontSize="xs" color="gray.500">
                Use preço fixo ou faixa (mínimo e máximo), não ambos.
              </Text>

              <Box minW={0}>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  Descrição
                </Text>
                <Textarea
                  value={description}
                  onChange={(event) => setDescription(event.target.value)}
                  minH="24"
                  borderRadius="xl"
                  bg="gray.50"
                  borderColor="gray.200"
                  focusRingColor="teal.200"
                  placeholder="Detalhes do serviço, público-alvo e diferenciais"
                  required
                />
              </Box>

              <HStack gap={3} flexWrap="wrap">
                <Button
                  type="submit"
                  disabled={
                    isSubmitting || !provider?.id || isLoadingProviderContext
                  }
                  borderRadius="full"
                  bg="green.400"
                  color="white"
                  _hover={{ bg: "green.500" }}
                  _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
                >
                  {isSubmitting
                    ? isEditing
                      ? "Salvando..."
                      : "Cadastrando..."
                    : isEditing
                      ? "Salvar alterações"
                      : "Cadastrar serviço"}
                </Button>

                {isEditing ? (
                  <Button
                    type="button"
                    variant="outline"
                    borderRadius="full"
                    onClick={handleCancelEdit}
                    disabled={isSubmitting}
                  >
                    Cancelar edição
                  </Button>
                ) : null}
              </HStack>

              {feedback ? (
                <Text
                  fontSize="xs"
                  color={feedback.type === "error" ? "red.600" : "green.600"}
                >
                  {feedback.message}
                </Text>
              ) : null}
            </VStack>
          </chakra.form>
        </Box>

        <Box
          borderRadius="3xl"
          bg="white"
          p={{ base: 5, md: 6 }}
          borderWidth="1px"
          borderColor="gray.200"
          shadow="sm"
        >
          <Text fontSize="xl" fontWeight="semibold" color="gray.900" mb={4}>
            Seus serviços
          </Text>

          {isLoadingProviderContext || isLoadingServices ? (
            <Text fontSize="sm" color="gray.500">
              Carregando serviços...
            </Text>
          ) : !provider?.id ? (
            <Text fontSize="sm" color="red.600">
              Não encontramos um provider vinculado ao seu usuário.
            </Text>
          ) : listServicesError ? (
            <Text fontSize="sm" color="red.600">
              {getApiErrorMessage(
                listServicesError,
                "Não foi possível carregar os serviços.",
              )}
            </Text>
          ) : services.length === 0 ? (
            <Text fontSize="sm" color="gray.500">
              Você ainda não cadastrou serviços.
            </Text>
          ) : (
            <VStack align="stretch" gap={3}>
              {services.map((service) => {
                const isCurrentDeleting = deletingServiceId === service.id;
                const currentPhotoFile =
                  photoFilesByService[service.id] ?? null;
                const isCurrentUploadingPhoto =
                  uploadingPhotoServiceId === service.id;
                const selectedCategoryId =
                  selectedCategoryByService[service.id] ?? "";
                const selectedTagId = selectedTagByService[service.id] ?? "";
                const isCurrentAddingCategory =
                  addingCategoryServiceId === service.id;
                const isCurrentAddingTag = addingTagServiceId === service.id;
                const taxonomyFeedback =
                  taxonomyFeedbackByService[service.id] ?? null;

                return (
                  <Box
                    key={service.id}
                    borderRadius="2xl"
                    borderWidth="1px"
                    borderColor="gray.200"
                    bg="gray.50"
                    p={4}
                  >
                    <Flex
                      align={{ base: "start", md: "center" }}
                      justify="space-between"
                      direction={{ base: "column", md: "row" }}
                      gap={3}
                    >
                      <Box>
                        <Text
                          fontSize="md"
                          fontWeight="semibold"
                          color="gray.900"
                        >
                          {service.name}
                        </Text>
                        <Text mt={1} fontSize="sm" color="gray.600">
                          {service.description}
                        </Text>
                        <Text mt={2} fontSize="xs" color="gray.500">
                          Preço base: R$ {service.price.toFixed(2)} · Faixa: R${" "}
                          {service.priceMinimum.toFixed(2)} - R${" "}
                          {service.priceMaximum.toFixed(2)} · Duração:{" "}
                          {service.duration} min
                        </Text>
                      </Box>

                      <HStack gap={2}>
                        <Button
                          size="sm"
                          borderRadius="full"
                          variant="subtle"
                          onClick={() =>
                            router.push(
                              `/services/${service.id}?from=/provider`,
                            )
                          }
                        >
                          Ver detalhes
                        </Button>
                        <Button
                          size="sm"
                          borderRadius="full"
                          variant="outline"
                          onClick={() => fillFormForEdit(service)}
                          disabled={isSubmitting || isCurrentDeleting}
                        >
                          Editar
                        </Button>
                        <Button
                          size="sm"
                          borderRadius="full"
                          colorPalette="red"
                          variant="subtle"
                          onClick={() => handleDeleteService(service.id)}
                          disabled={
                            isDeletingService ||
                            isSubmitting ||
                            Boolean(confirmDeleteServiceId)
                          }
                        >
                          {isCurrentDeleting ? "Excluindo..." : "Excluir"}
                        </Button>
                      </HStack>
                    </Flex>

                    <Box mt={4}>
                      <Text fontSize="sm" fontWeight="medium" color="gray.700">
                        Fotos do serviço
                      </Text>

                      {service.photos.length === 0 ? (
                        <Text mt={1.5} fontSize="xs" color="gray.500">
                          Nenhuma foto cadastrada.
                        </Text>
                      ) : (
                        <Flex mt={2} gap={2} wrap="wrap">
                          {service.photos.map((photo) => {
                            const photoDeleteKey = `${service.id}:${photo.id}`;
                            const isCurrentDeletingPhoto =
                              deletingPhotoKey === photoDeleteKey;

                            return (
                              <Box
                                key={photo.id}
                                borderWidth="1px"
                                borderColor="gray.200"
                                borderRadius="lg"
                                bg="white"
                                p={2}
                              >
                                <chakra.img
                                  src={photo.url}
                                  alt={`Foto do serviço ${service.name}`}
                                  w="84px"
                                  h="84px"
                                  objectFit="cover"
                                  borderRadius="md"
                                />
                                <Button
                                  mt={2}
                                  size="xs"
                                  borderRadius="full"
                                  colorPalette="red"
                                  variant="subtle"
                                  onClick={() =>
                                    handleDeleteServicePhoto(
                                      service.id,
                                      photo.id,
                                    )
                                  }
                                  disabled={
                                    isDeletingServicePhoto ||
                                    isCurrentDeletingPhoto ||
                                    isAddingServicePhoto
                                  }
                                >
                                  {isCurrentDeletingPhoto
                                    ? "Removendo..."
                                    : "Remover"}
                                </Button>
                              </Box>
                            );
                          })}
                        </Flex>
                      )}

                      <HStack mt={3} align="end" gap={2} flexWrap="wrap">
                        <Box>
                          <Input
                            type="file"
                            accept="image/*"
                            onChange={(event) =>
                              handlePhotoFileChange(service.id, event)
                            }
                            h="10"
                            bg="white"
                            borderRadius="xl"
                            borderColor="gray.200"
                            p={1}
                            disabled={
                              isAddingServicePhoto || isCurrentUploadingPhoto
                            }
                          />
                          {currentPhotoFile ? (
                            <Text mt={1} fontSize="xs" color="gray.500">
                              Arquivo selecionado: {currentPhotoFile.name}
                            </Text>
                          ) : null}
                        </Box>

                        <Button
                          size="sm"
                          borderRadius="full"
                          bg="green.400"
                          color="white"
                          _hover={{ bg: "green.500" }}
                          onClick={() => handleUploadServicePhoto(service.id)}
                          disabled={
                            isAddingServicePhoto ||
                            isCurrentUploadingPhoto ||
                            !currentPhotoFile ||
                            isDeletingServicePhoto
                          }
                        >
                          {isCurrentUploadingPhoto
                            ? "Enviando..."
                            : "Enviar foto"}
                        </Button>
                      </HStack>

                      <Box mt={4}>
                        <Text
                          fontSize="sm"
                          fontWeight="medium"
                          color="gray.700"
                          mb={2}
                        >
                          Categorias e tags
                        </Text>

                        <Grid
                          templateColumns={{ base: "1fr", md: "1fr 1fr" }}
                          gap={3}
                        >
                          <Box>
                            <Text fontSize="xs" color="gray.500" mb={1}>
                              Categorias atuais
                            </Text>
                            <Text fontSize="sm" color="gray.700">
                              {service.categories.length > 0
                                ? service.categories
                                    .map((category) => category.name)
                                    .join(", ")
                                : "Nenhuma categoria associada."}
                            </Text>

                            <HStack mt={2} align="end" gap={2} flexWrap="wrap">
                              <NativeSelect.Root
                                size="sm"
                                minW="220px"
                                disabled={
                                  isLoadingCategories ||
                                  Boolean(categoriesError) ||
                                  isCurrentAddingCategory ||
                                  isAddingServiceTag ||
                                  isAddingServiceCategory
                                }
                              >
                                <NativeSelect.Field
                                  value={selectedCategoryId}
                                  onChange={(event) =>
                                    handleCategorySelectionChange(
                                      service.id,
                                      event.target.value,
                                    )
                                  }
                                  borderRadius="xl"
                                  bg="white"
                                  borderColor="gray.200"
                                >
                                  <option value="">
                                    {isLoadingCategories
                                      ? "Carregando categorias..."
                                      : "Selecione uma categoria"}
                                  </option>
                                  {categories.map((category) => (
                                    <option
                                      key={category.id}
                                      value={category.id}
                                    >
                                      {category.name}
                                    </option>
                                  ))}
                                </NativeSelect.Field>
                                <NativeSelect.Indicator color="gray.500" />
                              </NativeSelect.Root>

                              <Button
                                size="sm"
                                borderRadius="full"
                                variant="outline"
                                onClick={() =>
                                  handleAddCategoryToService(service)
                                }
                                disabled={
                                  !selectedCategoryId ||
                                  isCurrentAddingCategory ||
                                  isAddingServiceCategory ||
                                  isLoadingCategories ||
                                  Boolean(categoriesError)
                                }
                              >
                                {isCurrentAddingCategory
                                  ? "Associando..."
                                  : "Associar categoria"}
                              </Button>
                            </HStack>

                            {categoriesError ? (
                              <Text mt={1} fontSize="xs" color="red.600">
                                {getApiErrorMessage(
                                  categoriesError,
                                  "Não foi possível carregar as categorias.",
                                )}
                              </Text>
                            ) : null}
                          </Box>

                          <Box>
                            <Text fontSize="xs" color="gray.500" mb={1}>
                              Tags atuais
                            </Text>
                            <Text fontSize="sm" color="gray.700">
                              {service.tags.length > 0
                                ? service.tags.map((tag) => tag.name).join(", ")
                                : "Nenhuma tag associada."}
                            </Text>

                            <HStack mt={2} align="end" gap={2} flexWrap="wrap">
                              <NativeSelect.Root
                                size="sm"
                                minW="220px"
                                disabled={
                                  isLoadingTags ||
                                  Boolean(tagsError) ||
                                  isCurrentAddingTag ||
                                  isAddingServiceTag ||
                                  isAddingServiceCategory
                                }
                              >
                                <NativeSelect.Field
                                  value={selectedTagId}
                                  onChange={(event) =>
                                    handleTagSelectionChange(
                                      service.id,
                                      event.target.value,
                                    )
                                  }
                                  borderRadius="xl"
                                  bg="white"
                                  borderColor="gray.200"
                                >
                                  <option value="">
                                    {isLoadingTags
                                      ? "Carregando tags..."
                                      : "Selecione uma tag"}
                                  </option>
                                  {tags.map((tag) => (
                                    <option key={tag.id} value={tag.id}>
                                      {tag.name}
                                    </option>
                                  ))}
                                </NativeSelect.Field>
                                <NativeSelect.Indicator color="gray.500" />
                              </NativeSelect.Root>

                              <Button
                                size="sm"
                                borderRadius="full"
                                variant="outline"
                                onClick={() => handleAddTagToService(service)}
                                disabled={
                                  !selectedTagId ||
                                  isCurrentAddingTag ||
                                  isAddingServiceTag ||
                                  isLoadingTags ||
                                  Boolean(tagsError)
                                }
                              >
                                {isCurrentAddingTag
                                  ? "Associando..."
                                  : "Associar tag"}
                              </Button>
                            </HStack>

                            {tagsError ? (
                              <Text mt={1} fontSize="xs" color="red.600">
                                {getApiErrorMessage(
                                  tagsError,
                                  "Não foi possível carregar as tags.",
                                )}
                              </Text>
                            ) : null}
                          </Box>
                        </Grid>

                        {taxonomyFeedback ? (
                          <Text
                            mt={2}
                            fontSize="xs"
                            color={
                              taxonomyFeedback.type === "error"
                                ? "red.600"
                                : "green.600"
                            }
                          >
                            {taxonomyFeedback.message}
                          </Text>
                        ) : null}
                      </Box>
                    </Box>
                  </Box>
                );
              })}
            </VStack>
          )}
        </Box>
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
