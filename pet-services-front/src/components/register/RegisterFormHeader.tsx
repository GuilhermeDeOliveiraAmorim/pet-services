import { Box, Flex, Heading, Text } from "@chakra-ui/react";

export default function RegisterFormHeader() {
  return (
    <Flex align="center" justify="space-between" mb={6}>
      <Box>
        <Heading size="2xl" color="gray.900">
          Cadastro
        </Heading>
        <Text mt={2} fontSize="sm" color="gray.500">
          Preencha os campos para criar sua conta.
        </Text>
      </Box>
      <Flex
        h="12"
        px={4}
        align="center"
        justify="center"
        borderRadius="full"
        bg="gray.100"
        borderWidth="1px"
        borderColor="gray.200"
        fontSize="sm"
        fontWeight="semibold"
        color="gray.700"
      >
        Novo
      </Flex>
    </Flex>
  );
}
