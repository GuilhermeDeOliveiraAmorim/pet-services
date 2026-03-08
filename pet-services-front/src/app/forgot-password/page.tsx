"use client";

import { type ChangeEvent, useMemo, useState } from "react";
import Link from "next/link";
import {
  Box,
  Button,
  Heading,
  Input,
  Link as ChakraLink,
  Text,
  VStack,
  chakra,
} from "@chakra-ui/react";

import { useAuthRequestPasswordReset } from "@/application";
import { getApiErrorMessage } from "@/lib/api-error";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

export default function ForgotPasswordPage() {
  const { mutateAsync, isPending, error, isSuccess } =
    useAuthRequestPasswordReset();
  const [email, setEmail] = useState("");
  const [emailError, setEmailError] = useState("");

  const validateEmail = (value: string) => {
    const parsed = value.trim();

    if (!parsed) {
      return "Informe o email";
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(parsed)) {
      return "Email inválido";
    }

    return "";
  };

  const feedback = useMemo(() => {
    if (!error) {
      return "";
    }

    return getApiErrorMessage(
      error,
      "Não foi possível solicitar a redefinição de senha.",
    );
  }, [error]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const nextEmailError = validateEmail(email);
    setEmailError(nextEmailError);

    if (nextEmailError) {
      return;
    }

    await mutateAsync({ email });
  };

  return (
    <PageWrapper gap={16}>
      <MainNav showLinks={false} showActions={false} />

      <chakra.form
        onSubmit={handleSubmit}
        mx="auto"
        w="full"
        maxW="xl"
        borderRadius="3xl"
        bg="white"
        p={{ base: 6, md: 8 }}
        boxShadow="0 30px 80px rgba(35, 49, 82, 0.12)"
      >
        <Box mb={6}>
          <Heading size="2xl" color="gray.900">
            Redefinir senha
          </Heading>
          <Text mt={2} fontSize="sm" color="gray.500">
            Enviaremos um link para redefinir sua senha.
          </Text>
        </Box>

        <VStack gap={5} align="stretch">
          <Box>
            <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
              Email
            </Text>
            <Input
              type="email"
              value={email}
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                setEmail(event.target.value);
                if (emailError) {
                  setEmailError("");
                }
              }}
              onBlur={() => setEmailError(validateEmail(email))}
              placeholder="voce@email.com"
              autoComplete="email"
              h="11"
              borderRadius="xl"
              bg="gray.50"
              borderColor={emailError ? "red.300" : "gray.200"}
              focusRingColor={emailError ? "red.200" : "teal.200"}
            />
            {emailError ? (
              <Text mt={1.5} fontSize="xs" color="red.500">
                {emailError}
              </Text>
            ) : null}
          </Box>

          {feedback ? (
            <Box
              borderRadius="2xl"
              borderWidth="1px"
              borderColor="red.200"
              bg="red.50"
              p={4}
            >
              <Text fontSize="sm" color="red.600">
                {feedback}
              </Text>
            </Box>
          ) : null}

          {isSuccess ? (
            <Box
              borderRadius="2xl"
              borderWidth="1px"
              borderColor="green.200"
              bg="green.50"
              p={4}
            >
              <Text fontSize="sm" color="green.700">
                Se o email existir, enviaremos instruções para redefinição.
              </Text>
            </Box>
          ) : null}

          <Button
            type="submit"
            disabled={isPending}
            h="11"
            borderRadius="full"
            bgGradient="linear(to-r, green.400, teal.400)"
            color="white"
            _hover={{ opacity: 0.92 }}
            _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
          >
            {isPending ? "Enviando..." : "Enviar link"}
          </Button>
        </VStack>

        <Text mt={6} textAlign="center" fontSize="xs" color="gray.500">
          <ChakraLink
            as={Link}
            href="/login"
            color="teal.500"
            fontWeight="medium"
          >
            Voltar para o login
          </ChakraLink>
        </Text>
      </chakra.form>
    </PageWrapper>
  );
}
