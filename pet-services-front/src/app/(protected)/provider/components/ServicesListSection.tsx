import type { ChangeEvent } from "react";
import { Box, Text, VStack } from "@chakra-ui/react";

import type { Service } from "@/domain";
import ServiceCard from "./ServiceCard";

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
  onPhotoFileChange: (
    serviceId: string,
    event: ChangeEvent<HTMLInputElement>,
  ) => void;
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
            const selectedCategoryId =
              selectedCategoryByService[service.id] ?? "";
            const selectedTagId = selectedTagByService[service.id] ?? "";
            const isCurrentAddingCategory =
              addingCategoryServiceId === service.id;
            const isCurrentAddingTag = addingTagServiceId === service.id;
            const taxonomyFeedback =
              taxonomyFeedbackByService[service.id] ?? null;
            const newTagName = newTagNameByService[service.id] ?? "";

            return (
              <ServiceCard
                key={service.id}
                service={service}
                isCurrentDeleting={isCurrentDeleting}
                isSubmitting={isSubmitting}
                isDeletingService={isDeletingService}
                confirmDeleteServiceId={confirmDeleteServiceId}
                onViewDetails={onViewDetails}
                onEditService={onEditService}
                onDeleteService={onDeleteService}
                currentPhotoFile={currentPhotoFile}
                isCurrentUploadingPhoto={isCurrentUploadingPhoto}
                deletingPhotoKey={deletingPhotoKey}
                isAddingServicePhoto={isAddingServicePhoto}
                isDeletingServicePhoto={isDeletingServicePhoto}
                onPhotoFileChange={onPhotoFileChange}
                onUploadServicePhoto={onUploadServicePhoto}
                onDeleteServicePhoto={onDeleteServicePhoto}
                selectedCategoryId={selectedCategoryId}
                selectedTagId={selectedTagId}
                newTagName={newTagName}
                isCurrentAddingCategory={isCurrentAddingCategory}
                isCurrentAddingTag={isCurrentAddingTag}
                removingCategoryKey={removingCategoryKey}
                removingTagKey={removingTagKey}
                taxonomyFeedback={taxonomyFeedback}
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
                onCategorySelectionChange={onCategorySelectionChange}
                onTagSelectionChange={onTagSelectionChange}
                onNewTagNameChange={onNewTagNameChange}
                onAddCategoryToService={onAddCategoryToService}
                onRemoveCategoryFromService={onRemoveCategoryFromService}
                onAddTagToService={onAddTagToService}
                onRemoveTagFromService={onRemoveTagFromService}
              />
            );
          })}
        </VStack>
      )}
    </Box>
  );
}
