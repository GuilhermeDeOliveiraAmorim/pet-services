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

type ServiceCardProps = {
  service: Service;
  isCurrentDeleting: boolean;
  isSubmitting: boolean;
  isDeletingService: boolean;
  confirmDeleteServiceId: string | null;
  onViewDetails: (serviceId: string) => void;
  onEditService: (service: Service) => void;
  onDeleteService: (serviceId: string) => void;
  currentPhotoFile: File | null;
  isCurrentUploadingPhoto: boolean;
  deletingPhotoKey: string | null;
  isAddingServicePhoto: boolean;
  isDeletingServicePhoto: boolean;
  onPhotoFileChange: (
    serviceId: string,
    event: ChangeEvent<HTMLInputElement>,
  ) => void;
  onUploadServicePhoto: (serviceId: string) => void;
  onDeleteServicePhoto: (serviceId: string, photoId: string) => void;
  selectedCategoryId: string;
  selectedTagId: string;
  newTagName: string;
  isCurrentAddingCategory: boolean;
  isCurrentAddingTag: boolean;
  removingCategoryKey: string | null;
  removingTagKey: string | null;
  taxonomyFeedback: Feedback | null;
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

export default function ServiceCard({
  service,
  isCurrentDeleting,
  isSubmitting,
  isDeletingService,
  confirmDeleteServiceId,
  onViewDetails,
  onEditService,
  onDeleteService,
  currentPhotoFile,
  isCurrentUploadingPhoto,
  deletingPhotoKey,
  isAddingServicePhoto,
  isDeletingServicePhoto,
  onPhotoFileChange,
  onUploadServicePhoto,
  onDeleteServicePhoto,
  selectedCategoryId,
  selectedTagId,
  newTagName,
  isCurrentAddingCategory,
  isCurrentAddingTag,
  removingCategoryKey,
  removingTagKey,
  taxonomyFeedback,
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
}: ServiceCardProps) {
  const hasFixedPrice = service.price > 0;
  const hasRangePrice = service.priceMinimum > 0 || service.priceMaximum > 0;

  return (
    <Box
      borderRadius="2xl"
      borderWidth="1px"
      borderColor="gray.200"
      bg="white"
      p={{ base: 4, md: 5 }}
    >
      <Flex
        align={{ base: "start", md: "center" }}
        justify="space-between"
        direction={{ base: "column", md: "row" }}
        gap={4}
      >
        <Box minW={0}>
          <Text fontSize="lg" fontWeight="semibold" color="gray.900">
            {service.name}
          </Text>
          <Text mt={1} fontSize="sm" color="gray.600">
            {service.description}
          </Text>

          <Flex mt={3} gap={2} wrap="wrap">
            <Box
              borderRadius="full"
              borderWidth="1px"
              borderColor="gray.200"
              bg="gray.50"
              px={3}
              py={1.5}
            >
              <Text fontSize="xs" fontWeight="medium" color="gray.700">
                Duração: {service.duration} min
              </Text>
            </Box>

            {hasFixedPrice ? (
              <Box
                borderRadius="full"
                borderWidth="1px"
                borderColor="green.200"
                bg="green.50"
                px={3}
                py={1.5}
              >
                <Text fontSize="xs" fontWeight="medium" color="green.700">
                  Preço base: R$ {service.price.toFixed(2)}
                </Text>
              </Box>
            ) : null}

            {hasRangePrice ? (
              <Box
                borderRadius="full"
                borderWidth="1px"
                borderColor="teal.200"
                bg="teal.50"
                px={3}
                py={1.5}
              >
                <Text fontSize="xs" fontWeight="medium" color="teal.700">
                  Faixa: R$ {service.priceMinimum.toFixed(2)} a R${" "}
                  {service.priceMaximum.toFixed(2)}
                </Text>
              </Box>
            ) : null}
          </Flex>
        </Box>

        <HStack gap={2} flexWrap="wrap" justify={{ base: "start", md: "end" }}>
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

      <VStack mt={5} gap={4} align="stretch">
        <Box
          borderRadius="2xl"
          borderWidth="1px"
          borderColor="gray.200"
          bg="gray.50"
          p={{ base: 3, md: 4 }}
        >
          <Text fontSize="sm" fontWeight="semibold" color="gray.900">
            Fotos do serviço
          </Text>
          <Text mt={1} fontSize="xs" color="gray.500">
            Adicione imagens para enriquecer a apresentação do serviço.
          </Text>

          {(service.photos ?? []).length === 0 ? (
            <Text mt={3} fontSize="xs" color="gray.500">
              Nenhuma foto cadastrada.
            </Text>
          ) : (
            <Flex mt={3} gap={3} wrap="wrap">
              {(service.photos ?? []).map((photo) => {
                const photoDeleteKey = `${service.id}:${photo.id}`;
                const isCurrentDeletingPhoto =
                  deletingPhotoKey === photoDeleteKey;

                return (
                  <Box
                    key={photo.id}
                    borderWidth="1px"
                    borderColor="gray.200"
                    borderRadius="xl"
                    bg="white"
                    p={2}
                    w="112px"
                  >
                    <chakra.img
                      src={photo.url}
                      alt={`Foto do serviço ${service.name}`}
                      w="96px"
                      h="96px"
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
                      onClick={() => onDeleteServicePhoto(service.id, photo.id)}
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

          <Box
            mt={4}
            borderRadius="xl"
            borderWidth="1px"
            borderColor="gray.200"
            bg="white"
            p={3}
          >
            <Text fontSize="xs" fontWeight="medium" color="gray.600" mb={2}>
              Nova foto
            </Text>
            <Input
              type="file"
              accept="image/*"
              onChange={(event) => onPhotoFileChange(service.id, event)}
              h="10"
              bg="white"
              borderRadius="xl"
              borderColor="gray.200"
              p={1}
              disabled={isAddingServicePhoto || isCurrentUploadingPhoto}
            />
            {currentPhotoFile ? (
              <Text mt={1} fontSize="xs" color="gray.500">
                Arquivo selecionado: {currentPhotoFile.name}
              </Text>
            ) : null}
            <Button
              mt={3}
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
          </Box>
        </Box>

        <Box
          borderRadius="2xl"
          borderWidth="1px"
          borderColor="gray.200"
          bg="gray.50"
          p={{ base: 3, md: 4 }}
        >
          <Text fontSize="sm" fontWeight="semibold" color="gray.900" mb={1}>
            Categorias e tags
          </Text>
          <Text fontSize="xs" color="gray.500">
            Organize como o serviço será classificado e encontrado.
          </Text>

          <Grid mt={3} templateColumns={{ base: "1fr", md: "1fr 1fr" }} gap={4}>
            <Box>
              <Box
                borderRadius="xl"
                borderWidth="1px"
                borderColor="gray.200"
                bg="white"
                p={3}
              >
                <Text fontSize="xs" fontWeight="medium" color="gray.600" mb={2}>
                  Associar categoria
                </Text>
                <NativeSelect.Root
                  size="sm"
                  w="full"
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
                      onCategorySelectionChange(service.id, event.target.value)
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
                  mt={3}
                  size="sm"
                  borderRadius="full"
                  variant="outline"
                  w={{ base: "full", sm: "auto" }}
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
              </Box>

              <Text
                fontSize="xs"
                fontWeight="medium"
                color="gray.500"
                mt={3}
                mb={2}
              >
                Categorias atuais
              </Text>
              {service.categories.length > 0 ? (
                <VStack align="stretch" gap={2}>
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
                        borderRadius="xl"
                        px={3}
                        py={2}
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
                            onRemoveCategoryFromService(service, category.id)
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

              {categoriesErrorMessage ? (
                <Text mt={1} fontSize="xs" color="red.600">
                  {categoriesErrorMessage}
                </Text>
              ) : null}
            </Box>

            <Box>
              <Box
                borderRadius="xl"
                borderWidth="1px"
                borderColor="gray.200"
                bg="white"
                p={3}
              >
                <Text fontSize="xs" fontWeight="medium" color="gray.600" mb={2}>
                  Associar ou criar tag
                </Text>
                <NativeSelect.Root
                  size="sm"
                  w="full"
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
                      onTagSelectionChange(service.id, event.target.value)
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
                  mt={3}
                  value={newTagName}
                  onChange={(event) =>
                    onNewTagNameChange(service.id, event.target.value)
                  }
                  placeholder="Nova tag (ex: urgente)"
                  h="9"
                  w="full"
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
                  mt={3}
                  size="sm"
                  borderRadius="full"
                  variant="outline"
                  w={{ base: "full", sm: "auto" }}
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
              </Box>

              <Text
                fontSize="xs"
                fontWeight="medium"
                color="gray.500"
                mt={3}
                mb={2}
              >
                Tags atuais
              </Text>
              {service.tags.length > 0 ? (
                <VStack align="stretch" gap={2}>
                  {service.tags.map((tag) => {
                    const tagKey = `${service.id}:${tag.id}`;
                    const isCurrentRemovingTag = removingTagKey === tagKey;

                    return (
                      <HStack
                        key={tag.id}
                        justify="space-between"
                        borderWidth="1px"
                        borderColor="gray.200"
                        bg="white"
                        borderRadius="xl"
                        px={3}
                        py={2}
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
                          {isCurrentRemovingTag ? "Removendo..." : "Remover"}
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
                taxonomyFeedback.type === "error" ? "red.600" : "green.600"
              }
            >
              {taxonomyFeedback.message}
            </Text>
          ) : null}
        </Box>
      </VStack>
    </Box>
  );
}
