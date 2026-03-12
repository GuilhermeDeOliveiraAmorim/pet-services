import type { FormEvent } from "react";
import {
  Box,
  Button,
  Grid,
  HStack,
  Input,
  Text,
  Textarea,
  VStack,
  chakra,
} from "@chakra-ui/react";

type Feedback = {
  type: "success" | "error";
  message: string;
};

type AddProviderFormProps = {
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
  onPrefillFromProfile: () => void;
  canPrefillFromProfile: boolean;
  isSubmitting: boolean;
  isLoadingProviderContext: boolean;
  maxPriceRangeLength: number;
  feedback: Feedback | null;
  businessName: string;
  onBusinessNameChange: (value: string) => void;
  priceRange: string;
  onPriceRangeChange: (value: string) => void;
  description: string;
  onDescriptionChange: (value: string) => void;
  street: string;
  onStreetChange: (value: string) => void;
  addressNumber: string;
  onAddressNumberChange: (value: string) => void;
  neighborhood: string;
  onNeighborhoodChange: (value: string) => void;
  city: string;
  onCityChange: (value: string) => void;
  zipCode: string;
  onZipCodeChange: (value: string) => void;
  state: string;
  onStateChange: (value: string) => void;
  country: string;
  onCountryChange: (value: string) => void;
  latitude: string;
  onLatitudeChange: (value: string) => void;
  longitude: string;
  onLongitudeChange: (value: string) => void;
  complement: string;
  onComplementChange: (value: string) => void;
};

export default function AddProviderForm({
  onSubmit,
  onPrefillFromProfile,
  canPrefillFromProfile,
  isSubmitting,
  isLoadingProviderContext,
  maxPriceRangeLength,
  feedback,
  businessName,
  onBusinessNameChange,
  priceRange,
  onPriceRangeChange,
  description,
  onDescriptionChange,
  street,
  onStreetChange,
  addressNumber,
  onAddressNumberChange,
  neighborhood,
  onNeighborhoodChange,
  city,
  onCityChange,
  zipCode,
  onZipCodeChange,
  state,
  onStateChange,
  country,
  onCountryChange,
  latitude,
  onLatitudeChange,
  longitude,
  onLongitudeChange,
  complement,
  onComplementChange,
}: AddProviderFormProps) {
  return (
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

      <chakra.form onSubmit={onSubmit}>
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
              onClick={onPrefillFromProfile}
              disabled={!canPrefillFromProfile || isSubmitting}
            >
              Preencher com dados do perfil
            </Button>
          </HStack>

          <Grid gap={4} templateColumns={{ base: "1fr", md: "1fr 1fr" }}>
            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Nome comercial
              </Text>
              <Input
                value={businessName}
                onChange={(event) => onBusinessNameChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Faixa de preço
              </Text>
              <Input
                value={priceRange}
                onChange={(event) =>
                  onPriceRangeChange(
                    event.target.value.slice(0, maxPriceRangeLength),
                  )
                }
                maxLength={maxPriceRangeLength}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Ex: 80-180"
                required
              />
              <Text mt={1.5} fontSize="xs" color="gray.500">
                Máximo de {maxPriceRangeLength} caracteres.
              </Text>
            </Box>

            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Descrição
              </Text>
              <Textarea
                value={description}
                onChange={(event) => onDescriptionChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Rua
              </Text>
              <Input
                value={street}
                onChange={(event) => onStreetChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Número do endereço
              </Text>
              <Input
                value={addressNumber}
                onChange={(event) => onAddressNumberChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Bairro
              </Text>
              <Input
                value={neighborhood}
                onChange={(event) => onNeighborhoodChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Cidade
              </Text>
              <Input
                value={city}
                onChange={(event) => onCityChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                CEP
              </Text>
              <Input
                value={zipCode}
                onChange={(event) => onZipCodeChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Estado
              </Text>
              <Input
                value={state}
                onChange={(event) => onStateChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                País
              </Text>
              <Input
                value={country}
                onChange={(event) => onCountryChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Latitude
              </Text>
              <Input
                type="number"
                inputMode="decimal"
                value={latitude}
                onChange={(event) => onLatitudeChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Longitude
              </Text>
              <Input
                type="number"
                inputMode="decimal"
                value={longitude}
                onChange={(event) => onLongitudeChange(event.target.value)}
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
            <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
              Complemento (opcional)
            </Text>
            <Input
              value={complement}
              onChange={(event) => onComplementChange(event.target.value)}
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
              disabled={isSubmitting || isLoadingProviderContext}
            >
              {isSubmitting ? "Cadastrando provider..." : "Cadastrar provider"}
            </Button>
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
  );
}
