import Link from "next/link";
import { Box, Link as ChakraLink, Text } from "@chakra-ui/react";
import { getApiErrorMessage } from "@/lib/api-error";

type RegisterFormFooterProps = {
  error: Error | null;
  isSuccess: boolean;
  successMessage?: string;
};

export default function RegisterFormFooter({
  error,
  isSuccess,
  successMessage,
}: RegisterFormFooterProps) {
  const feedbackMessage = (() => {
    if (!error) {
      return "";
    }

    return getApiErrorMessage(
      error,
      "Não foi possível criar sua conta. Verifique os dados.",
    );
  })();

  return (
    <Box>
      {feedbackMessage ? (
        <Text fontSize="sm" color="red.500">
          {feedbackMessage}
        </Text>
      ) : null}

      {isSuccess ? (
        <Text fontSize="sm" color="green.600">
          {successMessage ||
            "Cadastro inicial realizado! Entre para completar seu perfil."}
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
