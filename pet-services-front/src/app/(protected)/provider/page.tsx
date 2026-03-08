"use client";

import { useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import {
  Box,
  Button,
  Flex,
  Grid,
  HStack,
  Input,
  Text,
  Textarea,
  VStack,
  chakra,
} from "@chakra-ui/react";
import {
  useProviderList,
  useServiceAdd,
  useServiceDelete,
  useServiceList,
  useServiceUpdate,
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

export default function ProviderDashboardPage() {
  const queryClient = useQueryClient();
  const { data: userData, isLoading: isLoadingUser } = useUserProfile();
  const { data: providersData, isLoading: isLoadingProviders } =
    useProviderList();

  const userId = userData?.user?.id;
  const provider = useMemo(
    () => providersData?.providers?.find((item) => item.userId === userId),
    [providersData?.providers, userId],
  );

  const listInput = useMemo(
    () => (provider?.id ? { providerId: provider.id } : undefined),
    [provider?.id],
  );

  const {
    data: servicesData,
    isLoading: isLoadingServices,
    error: listServicesError,
  } = useServiceList({
    input: listInput,
    enabled: Boolean(provider?.id),
  });

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

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [price, setPrice] = useState("");
  const [priceMinimum, setPriceMinimum] = useState("");
  const [priceMaximum, setPriceMaximum] = useState("");
  const [duration, setDuration] = useState("");
  const [editingServiceId, setEditingServiceId] = useState<string | null>(null);
  const [deletingServiceId, setDeletingServiceId] = useState<string | null>(
    null,
  );
  const [feedback, setFeedback] = useState<Feedback | null>(null);

  const services = servicesData?.services ?? [];
  const isEditing = Boolean(editingServiceId);
  const isSubmitting = isAddingService || isUpdatingService;
  const isLoadingProviderContext = isLoadingUser || isLoadingProviders;

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

    const parsedPrice = Number(price);
    const parsedPriceMinimum = Number(priceMinimum);
    const parsedPriceMaximum = Number(priceMaximum);
    const parsedDuration = Number(duration);

    if (
      !name.trim() ||
      !description.trim() ||
      Number.isNaN(parsedPrice) ||
      Number.isNaN(parsedPriceMinimum) ||
      Number.isNaN(parsedPriceMaximum) ||
      Number.isNaN(parsedDuration)
    ) {
      setFeedback({
        type: "error",
        message: "Preencha todos os campos obrigatórios do serviço.",
      });
      return;
    }

    try {
      if (editingServiceId) {
        const response = await updateService({
          serviceId: editingServiceId,
          name: name.trim(),
          description: description.trim(),
          price: parsedPrice,
          priceMinimum: parsedPriceMinimum,
          priceMaximum: parsedPriceMaximum,
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
          price: parsedPrice,
          priceMinimum: parsedPriceMinimum,
          priceMaximum: parsedPriceMaximum,
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
    const confirmed = window.confirm("Deseja realmente excluir este serviço?");

    if (!confirmed) {
      return;
    }

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
          </Box>
        </Grid>

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
                    required
                  />
                </Box>
              </Grid>

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
                          disabled={isDeletingService || isSubmitting}
                        >
                          {isCurrentDeleting ? "Excluindo..." : "Excluir"}
                        </Button>
                      </HStack>
                    </Flex>
                  </Box>
                );
              })}
            </VStack>
          )}
        </Box>
      </VStack>

      <ChangePasswordCard />
    </PageWrapper>
  );
}
