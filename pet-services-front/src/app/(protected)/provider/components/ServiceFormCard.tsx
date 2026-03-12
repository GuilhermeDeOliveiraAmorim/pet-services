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

type ServiceFormCardProps = {
  isEditing: boolean;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
  name: string;
  onNameChange: (value: string) => void;
  duration: string;
  onDurationChange: (value: string) => void;
  price: string;
  onPriceChange: (value: string) => void;
  priceMinimum: string;
  onPriceMinimumChange: (value: string) => void;
  priceMaximum: string;
  onPriceMaximumChange: (value: string) => void;
  description: string;
  onDescriptionChange: (value: string) => void;
  isSubmitting: boolean;
  canSubmit: boolean;
  onCancelEdit: () => void;
  feedback: Feedback | null;
};

export default function ServiceFormCard({
  isEditing,
  onSubmit,
  name,
  onNameChange,
  duration,
  onDurationChange,
  price,
  onPriceChange,
  priceMinimum,
  onPriceMinimumChange,
  priceMaximum,
  onPriceMaximumChange,
  description,
  onDescriptionChange,
  isSubmitting,
  canSubmit,
  onCancelEdit,
  feedback,
}: ServiceFormCardProps) {
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
          Serviços
        </Text>
        <Text mt={2} fontSize="xl" fontWeight="semibold" color="gray.900">
          {isEditing ? "Editar serviço" : "Cadastrar serviço"}
        </Text>
      </Box>

      <chakra.form onSubmit={onSubmit}>
        <VStack align="stretch" gap={4}>
          <Grid gap={4} templateColumns={{ base: "1fr", md: "1fr 1fr" }}>
            <Box minW={0}>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Nome
              </Text>
              <Input
                value={name}
                onChange={(event) => onNameChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Duração (minutos)
              </Text>
              <Input
                type="number"
                inputMode="numeric"
                value={duration}
                onChange={(event) => onDurationChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Preço base (R$)
              </Text>
              <Input
                type="number"
                inputMode="decimal"
                value={price}
                onChange={(event) => onPriceChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Preço mínimo (R$)
              </Text>
              <Input
                type="number"
                inputMode="decimal"
                value={priceMinimum}
                onChange={(event) => onPriceMinimumChange(event.target.value)}
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
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Preço máximo (R$)
              </Text>
              <Input
                type="number"
                inputMode="decimal"
                value={priceMaximum}
                onChange={(event) => onPriceMaximumChange(event.target.value)}
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
              onChange={(event) => onDescriptionChange(event.target.value)}
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
              disabled={isSubmitting || !canSubmit}
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
                onClick={onCancelEdit}
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
  );
}
