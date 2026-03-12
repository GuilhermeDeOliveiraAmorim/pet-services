import Image from "next/image";
import type { RefObject } from "react";
import {
  Box,
  Button,
  Grid,
  HStack,
  Input,
  Text,
  VStack,
} from "@chakra-ui/react";

type PhotoItem = {
  id: string;
  url: string;
};

type PhotoUploadCardProps = {
  isLoading: boolean;
  photoInputRef: RefObject<HTMLInputElement | null>;
  selectedPhoto: File | null;
  onSelectedPhotoChange: (file: File | null) => void;
  onUpload: () => void;
  isUploadingPhoto: boolean;
  uploadSuccess: boolean;
  photoFeedback: string;
  photos: PhotoItem[];
};

export default function PhotoUploadCard({
  isLoading,
  photoInputRef,
  selectedPhoto,
  onSelectedPhotoChange,
  onUpload,
  isUploadingPhoto,
  uploadSuccess,
  photoFeedback,
  photos,
}: PhotoUploadCardProps) {
  return (
    <Box
      borderRadius={{ base: "2xl", md: "3xl" }}
      bg="white"
      p={{ base: 4, sm: 5, md: 7 }}
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
          Fotos
        </Text>
        <Text
          mt={2}
          fontSize={{ base: "lg", md: "xl" }}
          fontWeight="semibold"
          color="gray.900"
        >
          Fotos do perfil
        </Text>
        <Text mt={2} fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
          Envie uma foto para aparecer no seu perfil.
        </Text>
      </Box>

      {isLoading ? (
        <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.500">
          Carregando fotos...
        </Text>
      ) : (
        <VStack align="stretch" gap={{ base: 3, md: 4 }}>
          <Input
            ref={photoInputRef}
            type="file"
            accept="image/*"
            display="none"
            onChange={(event) =>
              onSelectedPhotoChange(event.target.files?.[0] ?? null)
            }
          />

          <HStack gap={{ base: 2, md: 3 }} flexWrap="wrap">
            <Button
              type="button"
              borderRadius="full"
              variant="outline"
              borderColor="gray.300"
              color="gray.700"
              onClick={() => photoInputRef.current?.click()}
              fontSize={{ base: "xs", sm: "sm" }}
              h={{ base: "9", md: "10" }}
              px={{ base: 3, md: 4 }}
            >
              Escolher foto
            </Button>
            <Text
              fontSize={{ base: "xs" }}
              color="gray.500"
              maxW={{ base: "150px", sm: "250px", md: "300px" }}
              overflow="hidden"
              textOverflow="ellipsis"
              whiteSpace="nowrap"
            >
              {selectedPhoto ? selectedPhoto.name : "Nenhum arquivo"}
            </Text>
            <Button
              type="button"
              onClick={onUpload}
              disabled={!selectedPhoto || isUploadingPhoto}
              borderRadius="full"
              bg="green.400"
              color="white"
              _hover={{ bg: "green.500" }}
              _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
              fontSize={{ base: "xs", sm: "sm" }}
              h={{ base: "9", md: "10" }}
              px={{ base: 3, md: 4 }}
            >
              {isUploadingPhoto ? "Enviando..." : "Enviar foto"}
            </Button>
          </HStack>

          {uploadSuccess ? (
            <Text fontSize={{ base: "xs" }} color="green.600">
              Foto enviada com sucesso.
            </Text>
          ) : null}

          {photoFeedback ? (
            <Text fontSize={{ base: "xs" }} color="red.600">
              {photoFeedback}
            </Text>
          ) : null}

          <Grid
            gap={{ base: 2, md: 3 }}
            templateColumns={{
              base: "1fr 1fr",
              sm: "repeat(3, 1fr)",
              lg: "repeat(4, 1fr)",
            }}
          >
            {photos.length ? (
              photos.map((photo) => (
                <Box
                  key={photo.id}
                  h={{ base: "24", sm: "28" }}
                  w="full"
                  overflow="hidden"
                  borderRadius={{ base: "lg", md: "2xl" }}
                  bg="gray.100"
                >
                  <Image
                    src={photo.url}
                    alt="Foto do usuário"
                    width={112}
                    height={112}
                    style={{
                      width: "100%",
                      height: "100%",
                      objectFit: "cover",
                    }}
                    unoptimized
                  />
                </Box>
              ))
            ) : (
              <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.500">
                Nenhuma foto enviada ainda.
              </Text>
            )}
          </Grid>
        </VStack>
      )}
    </Box>
  );
}
