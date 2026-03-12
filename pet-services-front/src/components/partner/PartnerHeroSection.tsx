"use client";

import Link from "next/link";
import { Box, Button, Heading, HStack, Text, VStack } from "@chakra-ui/react";

import Banner from "../../../public/banner.png";

export default function PartnerHeroSection() {
  return (
    <Box
      id="partner"
      position="relative"
      display="flex"
      alignItems="center"
      justifyContent="center"
      py={{ base: "16", md: "24" }}
      bgImage={`url(${Banner.src})`}
      bgSize="cover"
      bgPos="center"
      borderRadius="3xl"
      overflow="hidden"
    >
      <Box position="absolute" inset="0" bg="blackAlpha.600" />

      <VStack
        position="relative"
        zIndex="1"
        w="full"
        maxW="3xl"
        gap="8"
        textAlign="center"
        px="4"
      >
        <Heading
          size={{ base: "2xl", md: "4xl" }}
          color="white"
          lineHeight="1.2"
        >
          Seja um Parceiro Pet Services
        </Heading>

        <Text maxW="2xl" fontSize="md" lineHeight="7" color="white">
          Cadastre seu negócio, publique seus serviços e conecte-se com tutores
          que procuram atendimento de confiança na sua região.
        </Text>

        <HStack
          w="full"
          maxW="xl"
          gap="3"
          justify="center"
          flexWrap={{ base: "wrap", md: "nowrap" }}
        >
          <Link href="/register?user_type=provider">
            <Button
              borderRadius="full"
              px="6"
              bg="teal.400"
              color="white"
              fontSize="sm"
              fontWeight="semibold"
              _hover={{ bg: "teal.500" }}
            >
              Quero me cadastrar
            </Button>
          </Link>

          <Link href="/login">
            <Button
              borderRadius="full"
              px="6"
              variant="outline"
              borderColor="whiteAlpha.700"
              color="white"
              fontSize="sm"
              fontWeight="semibold"
              _hover={{ bg: "whiteAlpha.200" }}
            >
              Já tenho conta
            </Button>
          </Link>
        </HStack>

        <Text fontSize="xs" color="whiteAlpha.800">
          Depois do login, conclua o onboarding no painel do parceiro.
        </Text>
      </VStack>
    </Box>
  );
}
