import Link from "next/link";
import { Box, Link as ChakraLink, Text } from "@chakra-ui/react";

type RegisterFormFooterProps = {
  error: Error | null;
  isSuccess: boolean;
};

export default function RegisterFormFooter({
  error,
  isSuccess,
}: RegisterFormFooterProps) {
  return (
    <Box>
      {error ? (
        <Text fontSize="sm" color="red.500">
          Não foi possível criar sua conta. Verifique os dados.
        </Text>
      ) : null}

      {isSuccess ? (
        <Text fontSize="sm" color="green.600">
          Cadastro inicial realizado! Entre para completar seu perfil.
        </Text>
      ) : null}

      <Text mt={2} textAlign="center" fontSize="xs" color="gray.500">
        Já tem conta?{" "}
        <ChakraLink
          as={Link}
          href="/login"
          color="green.500"
          fontWeight="medium"
        >
          Entrar
        </ChakraLink>
      </Text>
    </Box>
  );
}
