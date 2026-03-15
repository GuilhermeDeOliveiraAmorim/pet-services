import {
  Box,
  Button,
  Grid,
  HStack,
  Input,
  NativeSelect,
  Text,
  Textarea,
  VStack,
  chakra,
} from "@chakra-ui/react";
import type { FormEvent } from "react";

type SpecieOption = {
  value: string;
  label: string;
};

type PetFormCardProps = {
  petName: string;
  petSpeciesId: string;
  petBreed: string;
  petAge: string;
  petWeight: string;
  petNotes: string;
  onPetNameChange: (value: string) => void;
  onPetSpeciesIdChange: (value: string) => void;
  onPetBreedChange: (value: string) => void;
  onPetAgeChange: (value: string) => void;
  onPetWeightChange: (value: string) => void;
  onPetNotesChange: (value: string) => void;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
  isPetFormValid: boolean;
  isAddingPet: boolean;
  addPetSuccess: boolean;
  petFeedback: string;
  hasSpeciesError: boolean;
  isLoadingSpecies: boolean;
  hasBreedsError: boolean;
  isLoadingBreeds: boolean;
  specieOptions: SpecieOption[];
  breedOptions: SpecieOption[];
  showBreedField: boolean;
};

export default function PetFormCard({
  petName,
  petSpeciesId,
  petBreed,
  petAge,
  petWeight,
  petNotes,
  onPetNameChange,
  onPetSpeciesIdChange,
  onPetBreedChange,
  onPetAgeChange,
  onPetWeightChange,
  onPetNotesChange,
  onSubmit,
  isPetFormValid,
  isAddingPet,
  addPetSuccess,
  petFeedback,
  hasSpeciesError,
  isLoadingSpecies,
  hasBreedsError,
  isLoadingBreeds,
  specieOptions,
  breedOptions,
  showBreedField,
}: PetFormCardProps) {
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
          Pets
        </Text>
        <Text
          mt={2}
          fontSize={{ base: "lg", md: "xl" }}
          fontWeight="semibold"
          color="gray.900"
        >
          Cadastrar pet
        </Text>
        <Text mt={2} fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
          Adicione um novo pet ao seu perfil.
        </Text>
      </Box>

      <chakra.form onSubmit={onSubmit}>
        <VStack align="stretch" gap={{ base: 3, md: 4 }}>
          <Grid
            gap={{ base: 3, md: 4 }}
            templateColumns={{ base: "1fr", md: "1fr 1fr" }}
          >
            <Box minW={0}>
              <Text
                fontSize={{ base: "xs", sm: "sm" }}
                fontWeight="medium"
                color="gray.700"
                mb={2}
              >
                Nome
              </Text>
              <Input
                type="text"
                value={petName}
                onChange={(event) => onPetNameChange(event.target.value)}
                h={{ base: "10", md: "11" }}
                borderRadius={{ base: "lg", md: "xl" }}
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Ex: Thor"
                fontSize={{ base: "sm", md: "base" }}
                required
              />
            </Box>

            <Box minW={0}>
              <Text
                fontSize={{ base: "xs", sm: "sm" }}
                fontWeight="medium"
                color="gray.700"
                mb={2}
              >
                Espécie
              </Text>
              <NativeSelect.Root
                size="md"
                w="full"
                minW={0}
                disabled={isLoadingSpecies || hasSpeciesError}
              >
                <NativeSelect.Field
                  name="specie"
                  value={petSpeciesId}
                  onChange={(event) => onPetSpeciesIdChange(event.target.value)}
                  h={{ base: "10", md: "11" }}
                  borderRadius={{ base: "lg", md: "xl" }}
                  bg="gray.50"
                  borderColor="gray.200"
                  borderWidth="1px"
                  focusRingColor="teal.200"
                  fontSize={{ base: "sm", md: "base" }}
                >
                  <option value="">
                    {isLoadingSpecies ? "Carregando..." : "Selecione a espécie"}
                  </option>
                  {specieOptions.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </NativeSelect.Field>
                <NativeSelect.Indicator color="gray.500" />
              </NativeSelect.Root>
            </Box>

            {showBreedField ? (
              <Box minW={0}>
                <Text
                  fontSize={{ base: "xs", sm: "sm" }}
                  fontWeight="medium"
                  color="gray.700"
                  mb={2}
                >
                  Raça
                </Text>
                <NativeSelect.Root
                  size="md"
                  w="full"
                  minW={0}
                  disabled={isLoadingBreeds || hasBreedsError}
                >
                  <NativeSelect.Field
                    name="breed"
                    value={petBreed}
                    onChange={(event) => onPetBreedChange(event.target.value)}
                    h={{ base: "10", md: "11" }}
                    borderRadius={{ base: "lg", md: "xl" }}
                    bg="gray.50"
                    borderColor="gray.200"
                    borderWidth="1px"
                    focusRingColor="teal.200"
                    fontSize={{ base: "sm", md: "base" }}
                  >
                    <option value="">
                      {isLoadingBreeds ? "Carregando..." : "Selecione a raça"}
                    </option>
                    {breedOptions.map((option) => (
                      <option key={option.value} value={option.label}>
                        {option.label}
                      </option>
                    ))}
                  </NativeSelect.Field>
                  <NativeSelect.Indicator color="gray.500" />
                </NativeSelect.Root>
              </Box>
            ) : null}

            <Box minW={0}>
              <Text
                fontSize={{ base: "xs", sm: "sm" }}
                fontWeight="medium"
                color="gray.700"
                mb={2}
              >
                Idade
              </Text>
              <Input
                type="number"
                inputMode="numeric"
                value={petAge}
                onChange={(event) => onPetAgeChange(event.target.value)}
                h={{ base: "10", md: "11" }}
                borderRadius={{ base: "lg", md: "xl" }}
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Ex: 3"
                fontSize={{ base: "sm", md: "base" }}
                min={0}
                required
              />
            </Box>

            <Box minW={0}>
              <Text
                fontSize={{ base: "xs", sm: "sm" }}
                fontWeight="medium"
                color="gray.700"
                mb={2}
              >
                Peso (kg)
              </Text>
              <Input
                type="number"
                inputMode="decimal"
                value={petWeight}
                onChange={(event) => onPetWeightChange(event.target.value)}
                h={{ base: "10", md: "11" }}
                borderRadius={{ base: "lg", md: "xl" }}
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="Ex: 12.5"
                fontSize={{ base: "sm", md: "base" }}
                min={0}
                step="0.1"
                required
              />
            </Box>
          </Grid>

          <Box minW={0}>
            <Text
              fontSize={{ base: "xs", sm: "sm" }}
              fontWeight="medium"
              color="gray.700"
              mb={2}
            >
              Observações
            </Text>
            <Textarea
              value={petNotes}
              onChange={(event) => onPetNotesChange(event.target.value)}
              minH={{ base: "20", md: "24" }}
              borderRadius={{ base: "lg", md: "xl" }}
              bg="gray.50"
              borderColor="gray.200"
              focusRingColor="teal.200"
              placeholder="Comportamento, cuidados especiais, etc."
              fontSize={{ base: "sm", md: "base" }}
            />
          </Box>

          <HStack gap={{ base: 2, md: 3 }} flexWrap="wrap">
            <Button
              type="submit"
              disabled={!isPetFormValid || isAddingPet}
              borderRadius="full"
              bg="green.400"
              color="white"
              _hover={{ bg: "green.500" }}
              _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
              fontSize={{ base: "xs", sm: "sm" }}
              h={{ base: "9", md: "10" }}
              px={{ base: 3, md: 4 }}
            >
              {isAddingPet ? "Cadastrando..." : "Cadastrar pet"}
            </Button>
            {addPetSuccess ? (
              <Text fontSize={{ base: "xs" }} color="green.600">
                Pet cadastrado com sucesso.
              </Text>
            ) : null}
          </HStack>

          {petFeedback ? (
            <Text fontSize={{ base: "xs" }} color="red.600">
              {petFeedback}
            </Text>
          ) : null}
          {hasSpeciesError ? (
            <Text fontSize={{ base: "xs" }} color="red.600">
              Não foi possível carregar as espécies.
            </Text>
          ) : null}
          {hasBreedsError ? (
            <Text fontSize={{ base: "xs" }} color="red.600">
              Não foi possível carregar as raças para a espécie selecionada.
            </Text>
          ) : null}
        </VStack>
      </chakra.form>
    </Box>
  );
}
