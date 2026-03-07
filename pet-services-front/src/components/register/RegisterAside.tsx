import { Box, Heading, HStack, Text, VStack } from "@chakra-ui/react";

export default function RegisterAside() {
  return (
    <VStack
      display={{ base: "none", lg: "flex" }}
      justify="center"
      align="stretch"
      gap={6}
      borderRadius="3xl"
      bg="white"
      p={10}
      borderWidth="1px"
      borderColor="gray.100"
      shadow="sm"
    >
      <HStack>
        <Box
          display="flex"
          alignItems="center"
          justifyContent="center"
          h="8"
          w="8"
          borderRadius="full"
          bg="green.400"
          color="white"
          fontWeight="semibold"
        >
          •
        </Box>
        <Text fontSize="lg" fontWeight="semibold">
          PetServices
        </Text>
      </HStack>

      <Heading size="4xl" lineHeight="1.15" color="gray.900">
        Crie sua conta
      </Heading>

      <Text fontSize="sm" lineHeight="1.6" color="gray.600">
        Cadastre seus dados para acessar a plataforma e acompanhar todos os
        serviços do seu pet.
      </Text>

      <Box borderRadius="3xl" bg="green.50" p={6}>
        <Text fontSize="xs" fontWeight="semibold" color="green.600">
          Segurança em primeiro lugar
        </Text>
        <Text mt={2} fontSize="sm" color="gray.600">
          Seus dados são protegidos e utilizados apenas para melhorar a sua
          experiência.
        </Text>
      </Box>
    </VStack>
  );
}
