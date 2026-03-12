"use client";

import Link from "next/link";
import {
  Box,
  Button,
  Heading,
  SimpleGrid,
  Text,
  VStack,
} from "@chakra-ui/react";

const steps = [
  {
    title: "1. Crie sua conta",
    description:
      "Cadastre-se em poucos minutos para acessar seu painel de parceiro.",
  },
  {
    title: "2. Complete seu perfil",
    description:
      "Preencha dados do negócio, endereço e informações que ajudam clientes a encontrar você.",
  },
  {
    title: "3. Publique serviços",
    description:
      "Adicione categorias, preços e fotos para destacar seus atendimentos.",
  },
];

export default function PartnerOnboardingSection() {
  return (
    <VStack
      borderRadius="3xl"
      bg="white"
      borderWidth="1px"
      borderColor="gray.200"
      px={{ base: "5", md: "8" }}
      py={{ base: "8", md: "10" }}
      align="stretch"
      gap="8"
    >
      <VStack align="stretch" gap="2">
        <Text
          fontSize="xs"
          fontWeight="semibold"
          textTransform="uppercase"
          color="teal.500"
        >
          Onboarding de parceiro
        </Text>
        <Heading size={{ base: "lg", md: "xl" }} color="gray.900">
          Comece agora e ative seu perfil profissional
        </Heading>
        <Text fontSize="sm" color="gray.600">
          Estruturamos o processo para você ficar visível rapidamente e com um
          perfil completo para atrair mais clientes.
        </Text>
      </VStack>

      <SimpleGrid columns={{ base: 1, md: 3 }} gap="4">
        {steps.map((step) => (
          <Box
            key={step.title}
            borderRadius="2xl"
            borderWidth="1px"
            borderColor="gray.200"
            bg="gray.50"
            p="5"
          >
            <Text fontWeight="semibold" color="gray.900">
              {step.title}
            </Text>
            <Text mt="2" fontSize="sm" color="gray.600">
              {step.description}
            </Text>
          </Box>
        ))}
      </SimpleGrid>

      <Button
        as={Link}
        href="/register?user_type=provider"
        alignSelf={{ base: "stretch", md: "flex-start" }}
        borderRadius="full"
        bgGradient="linear(to-r, teal.400, cyan.400)"
        px="6"
        py="3"
        fontSize="sm"
        fontWeight="semibold"
        color="white"
        _hover={{ opacity: 0.9 }}
      >
        Iniciar onboarding
      </Button>
    </VStack>
  );
}
