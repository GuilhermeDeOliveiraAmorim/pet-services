import type { ChangeEvent } from "react";
import {
  Box,
  Button,
  Flex,
  Grid,
  HStack,
  Input,
  NativeSelect,
  Text,
  VStack,
  chakra,
} from "@chakra-ui/react";

import type { Service } from "@/domain";

type Feedback = {
  type: "success" | "error";
  message: string;
};

type OptionItem = {
  id: string;
  name: string;
};

type ServicesListSectionProps = {
  isLoadingProviderContext: boolean;
  isLoadingServices: boolean;
  hasProvider: boolean;
  listServicesErrorMessage: string | null;
  services: Service[];
  deletingServiceId: string | null;
  isSubmitting: boolean;
  isDeletingService: boolean;
  confirmDeleteServiceId: string | null;
  onViewDetails: (serviceId: string) => void;
  onEditService: (service: Service) => void;
  onDeleteService: (serviceId: string) => void;
  photoFilesByService: Record<string, File | null>;
  uploadingPhotoServiceId: string | null;
  deletingPhotoKey: string | null;
  isAddingServicePhoto: boolean;
  isDeletingServicePhoto: boolean;
  onPhotoFileChange: (serviceId: string, event: ChangeEvent<HTMLInputElement>) => void;
  onUploadServicePhoto: (serviceId: string) => void;
  onDeleteServicePhoto: (serviceId: string, photoId: string) => void;
  selectedCategoryByService: Record<string, string>;
  selectedTagByService: Record<string, string>;
  newTagNameByService: Record<string, string>;
  addingCategoryServiceId: string | null;
  addingTagServiceId: string | null;
  removingCategoryKey: string | null;
  removingTagKey: string | null;
  taxonomyFeedbackByService: Record<string, Feedback | null>;
  categories: OptionItem[];
  tags: OptionItem[];
  isLoadingCategories: boolean;
  isLoadingTags: boolean;
  isAddingServiceCategory: boolean;
  isDeletingServiceCategory: boolean;
  isAddingServiceTag: boolean;
  isDeletingServiceTag: boolean;
  categoriesErrorMessage: string | null;
  tagsErrorMessage: string | null;
  onCategorySelectionChange: (serviceId: string, value: string) => void;
  onTagSelectionChange: (serviceId: string, value: string) => void;
  onNewTagNameChange: (serviceId: string, value: string) => void;
  onAddCategoryToService: (service: Service) => void;
  onRemoveCategoryFromService: (service: Service, categoryId: string) => void;
  onAddTagToService: (service: Service) => void;
  onRemoveTagFromService: (service: Service, tagId: string) => void;
};

export default function ServicesListSection({
  isLoadingProviderContext,
  isLoadingServices,
  hasProvider,
  listServicesErrorMessage,
  services,
  deletingServiceId,
  isSubmitting,
  isDeletingService,
  confirmDeleteServiceId,
  onViewDetails,
  onEditService,
  onDeleteService,
  photoFilesByService,
  uploadingPhotoServiceId,
  deletingPhotoKey,
  isAddingServicePhoto,
  isDeletingServicePhoto,
  onPhotoFileChange,
  onUploadServicePhoto,
  onDeleteServicePhoto,
  selectedCategoryByService,
  selectedTagByService,
  newTagNameByService,
  addingCategoryServiceId,
  addingTagServiceId,
  removingCategoryKey,
  removingTagKey,
  taxonomyFeedbackByService,
  categories,
  tags,
  isLoadingCategories,
  isLoadingTags,
  isAddingServiceCategory,
  isDeletingServiceCategory,
  isAddingServiceTag,
  isDeletingServiceTag,
  categoriesErrorMessage,
  tagsErrorMessage,
  onCategorySelectionChange,
  onTagSelectionChange,
  onNewTagNameChange,
  onAddCategoryToService,
  onRemoveCategoryFromService,
  onAddTagToService,
  onRemoveTagFromService,
}: ServicesListSectionProps) {
  return (
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
      ) : !hasProvider ? (
        <Text fontSize="sm" color="red.600">
          Não encontramos um provider vinculado ao seu usuário.
        </Text>
      ) : listServicesErrorMessage ? (
        <Text fontSize="sm" color="red.600">
          {listServicesErrorMessage}
        </Text>
      ) : services.length === 0 ? (
        <Text fontSize="sm" color="gray.500">
          Você ainda não cadastrou serviços.
        </Text>
      ) : (
        <VStack align="stretch" gap={3}>
          {services.map((service) => {
            const isCurrentDeleting = deletingServiceId === service.id;
            const currentPhotoFile = photoFilesByService[service.id] ?? null;
            const isCurrentUploadingPhoto =
              uploadingPhotoServiceId === service.id;
            const selectedCategoryId = selectedCategoryByService[service.id] ?? "";
            const selectedTagId = selectedTagByService[service.id] ?? "";
            const isCurrentAddingCategory =
              addingCategoryServiceId === service.id;
            const isCurrentAddingTag = addingTagServiceId === service.id;
            const taxonomyFeedback =
              taxonomyFeedbackByService[service.id] ?? null;
            const newTagName = newTagNameByService[service.id] ?? "";

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
                    <Text fontSize="md" fontWeight="semibold" color="gray.900">
                      {service.name}
                    </Text>
                    <Text mt={1} fontSize="sm" color="gray.600">
                      {service.description}
                    </Text>
                    <Text mt={2} fontSize="xs" color="gray.500">
                      Preço base: R$ {service.price.toFixed(2)} · Faixa: R${" "}
                      {service.priceMinimum.toFixed(2)} - R${" "}
                      {service.priceMaximum.toFixed(2)} · Duração: {service.duration} min
                    </Text>
                  </Box>

                  <HStack gap={2}>
                    <Button
                      size="sm"
                      borderRadius="full"
                      variant="subtle"
                      onClick={() => onViewDetails(service.id)}
                    >
                      Ver detalhes
                    </Button>
                    <Button
                      size="sm"
                      borderRadius="full"
                      variant="outline"
                      onClick={() => onEditService(service)}
                      disabled={isSubmitting || isCurrentDeleting}
                    >
                      Editar
                    </Button>
                    <Button
                      size="sm"
                      borderRadius="full"
                      colorPalette="red"
                      variant="subtle"
                      onClick={() => onDeleteService(service.id)}
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
                                onDeleteServicePhoto(service.id, photo.id)
                              }
                              disabled={
                                isDeletingServicePhoto ||
                                isCurrentDeletingPhoto ||
                                isAddingServicePhoto
                              }
                            >
                              {isCurrentDeletingPhoto ? "Removendo..." : "Remover"}
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
                        onChange={(event) => onPhotoFileChange(service.id, event)}
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
                      onClick={() => onUploadServicePhoto(service.id)}
                      disabled={
                        isAddingServicePhoto ||
                        isCurrentUploadingPhoto ||
                        !currentPhotoFile ||
                        isDeletingServicePhoto
                      }
                    >
                      {isCurrentUploadingPhoto ? "Enviando..." : "Enviar foto"}
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

                    <Grid templateColumns={{ base: "1fr", md: "1fr 1fr" }} gap={3}>
                      <Box>
                        <Text fontSize="xs" color="gray.500" mb={1}>
                          Categorias atuais
                        </Text>
                        {service.categories.length > 0 ? (
                          <VStack align="stretch" gap={1}>
                            {service.categories.map((category) => {
                              const categoryKey = `${service.id}:${category.id}`;
                              const isCurrentRemovingCategory =
                                removingCategoryKey === categoryKey;

                              return (
                                <HStack
                                  key={category.id}
                                  justify="space-between"
                                  borderWidth="1px"
                                  borderColor="gray.200"
                                  bg="white"
                                  borderRadius="md"
                                  px={2}
                                  py={1}
                                >
                                  <Text fontSize="sm" color="gray.700">
                                    {category.name}
                                  </Text>
                                  <Button
                                    size="xs"
                                    borderRadius="full"
                                    colorPalette="red"
                                    variant="subtle"
                                    onClick={() =>
                                      onRemoveCategoryFromService(
                                        service,
                                        category.id,
                                      )
                                    }
                                    disabled={
                                      isCurrentRemovingCategory ||
                                      isDeletingServiceCategory ||
                                      isAddingServiceCategory
                                    }
                                  >
                                    {isCurrentRemovingCategory
                                      ? "Removendo..."
                                      : "Remover"}
                                  </Button>
                                </HStack>
                              );
                            })}
                          </VStack>
                        ) : (
                          <Text fontSize="sm" color="gray.700">
                            Nenhuma categoria associada.
                          </Text>
                        )}

                        <HStack mt={2} align="end" gap={2} flexWrap="wrap">
                          <NativeSelect.Root
                            size="sm"
                            minW="220px"
                            disabled={
                              isLoadingCategories ||
                              Boolean(categoriesErrorMessage) ||
                              isCurrentAddingCategory ||
                              isAddingServiceTag ||
                              isAddingServiceCategory
                            }
                          >
                            <NativeSelect.Field
                              value={selectedCategoryId}
                              onChange={(event) =>
                                onCategorySelectionChange(
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
                                <option key={category.id} value={category.id}>
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
                            onClick={() => onAddCategoryToService(service)}
                            disabled={
                              !selectedCategoryId ||
                              isCurrentAddingCategory ||
                              isAddingServiceCategory ||
                              isDeletingServiceCategory ||
                              isLoadingCategories ||
                              Boolean(categoriesErrorMessage)
                            }
                          >
                            {isCurrentAddingCategory
                              ? "Associando..."
                              : "Associar categoria"}
                          </Button>
                        </HStack>

                        {categoriesErrorMessage ? (
                          <Text mt={1} fontSize="xs" color="red.600">
                            {categoriesErrorMessage}
                          </Text>
                        ) : null}
                      </Box>

                      <Box>
                        <Text fontSize="xs" color="gray.500" mb={1}>
                          Tags atuais
                        </Text>
                        {service.tags.length > 0 ? (
                          <VStack align="stretch" gap={1}>
                            {service.tags.map((tag) => {
                              const tagKey = `${service.id}:${tag.id}`;
                              const isCurrentRemovingTag =
                                removingTagKey === tagKey;

                              return (
                                <HStack
                                  key={tag.id}
                                  justify="space-between"
                                  borderWidth="1px"
                                  borderColor="gray.200"
                                  bg="white"
                                  borderRadius="md"
                                  px={2}
                                  py={1}
                                >
                                  <Text fontSize="sm" color="gray.700">
                                    {tag.name}
                                  </Text>
                                  <Button
                                    size="xs"
                                    borderRadius="full"
                                    colorPalette="red"
                                    variant="subtle"
                                    onClick={() =>
                                      onRemoveTagFromService(service, tag.id)
                                    }
                                    disabled={
                                      isCurrentRemovingTag ||
                                      isDeletingServiceTag ||
                                      isAddingServiceTag
                                    }
                                  >
                                    {isCurrentRemovingTag
                                      ? "Removendo..."
                                      : "Remover"}
                                  </Button>
                                </HStack>
                              );
                            })}
                          </VStack>
                        ) : (
                          <Text fontSize="sm" color="gray.700">
                            Nenhuma tag associada.
                          </Text>
                        )}

                        <HStack mt={2} align="end" gap={2} flexWrap="wrap">
                          <NativeSelect.Root
                            size="sm"
                            minW="220px"
                            disabled={
                              isLoadingTags ||
                              Boolean(tagsErrorMessage) ||
                              isCurrentAddingTag ||
                              isAddingServiceTag ||
                              isAddingServiceCategory ||
                              isDeletingServiceTag
                            }
                          >
                            <NativeSelect.Field
                              value={selectedTagId}
                              onChange={(event) =>
                                onTagSelectionChange(
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
                                  : "Selecione uma tag (ou adicione)"}
                              </option>
                              {tags.map((tag) => (
                                <option key={tag.id} value={tag.id}>
                                  {tag.name}
                                </option>
                              ))}
                            </NativeSelect.Field>
                            <NativeSelect.Indicator color="gray.500" />
                          </NativeSelect.Root>

                          <Input
                            value={newTagName}
                            onChange={(event) =>
                              onNewTagNameChange(service.id, event.target.value)
                            }
                            placeholder="Nova tag (ex: urgente)"
                            h="9"
                            minW="220px"
                            borderRadius="xl"
                            bg="white"
                            borderColor="gray.200"
                            disabled={
                              isCurrentAddingTag ||
                              isAddingServiceTag ||
                              isLoadingTags ||
                              isDeletingServiceTag
                            }
                          />

                          <Button
                            size="sm"
                            borderRadius="full"
                            variant="outline"
                            onClick={() => onAddTagToService(service)}
                            disabled={
                              (!selectedTagId && !newTagName.trim()) ||
                              isCurrentAddingTag ||
                              isAddingServiceTag ||
                              isDeletingServiceTag ||
                              isLoadingTags ||
                              Boolean(tagsErrorMessage)
                            }
                          >
                            {isCurrentAddingTag
                              ? "Associando..."
                              : "Adicionar/associar tag"}
                          </Button>
                        </HStack>

                        {tagsErrorMessage ? (
                          <Text mt={1} fontSize="xs" color="red.600">
                            {tagsErrorMessage}
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
  );
}
