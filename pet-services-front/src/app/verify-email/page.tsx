"use client";

import { Suspense, useEffect, useMemo, useRef } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import {
  Box,
  Button,
  Heading,
  Link as ChakraLink,
  Text,
  VStack,
} from "@chakra-ui/react";

import { useAuthVerifyEmail } from "@/application";
import { getApiErrorMessage } from "@/lib/api-error";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

function VerifyEmailContent() {
  const searchParams = useSearchParams();
  const token = searchParams.get("token")?.trim() ?? "";

  const hasTriggeredRef = useRef(false);
  const { mutate, isPending, isSuccess, error, reset } = useAuthVerifyEmail();

  useEffect(() => {
    if (!token || hasTriggeredRef.current) {
      return;
    }

    hasTriggeredRef.current = true;
    mutate({ token });
  }, [mutate, token]);

  const feedback = useMemo(() => {
    if (!error) {
      return "";
    }

    return getApiErrorMessage(
      error,
      "Não foi possível verificar seu e-mail. Tente novamente.",
    );
  }, [error]);

  const handleRetry = () => {
    if (!token) {
      return;
    }

    reset();
    mutate({ token });
  };

  return (
    <Box
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
          Verificar e-mail
        </Heading>
        <Text mt={2} fontSize="sm" color="gray.500">
          Estamos confirmando o token recebido no link.
        </Text>
      </Box>

      {!token ? (
        <Box
          borderRadius="2xl"
          borderWidth="1px"
          borderColor="red.200"
          bg="red.50"
          p={4}
        >
          <Text fontSize="sm" color="red.600">
            Token inválido ou ausente na URL. Solicite um novo e-mail de
            verificação.
          </Text>
        </Box>
      ) : null}

      {token ? (
        <VStack align="stretch" gap={4}>
          {isPending ? (
            <Box
              borderRadius="2xl"
              borderWidth="1px"
              borderColor="blue.200"
              bg="blue.50"
              p={4}
            >
              <Text fontSize="sm" color="blue.700">
                Validando seu e-mail...
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
                E-mail verificado com sucesso. Você já pode fazer login.
              </Text>
            </Box>
          ) : null}

          {!isPending && !isSuccess && feedback ? (
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

          {!isPending && !isSuccess ? (
            <Button
              onClick={handleRetry}
              h="11"
              borderRadius="full"
              variant="outline"
            >
              Tentar novamente
            </Button>
          ) : null}
        </VStack>
      ) : null}

      <Text mt={6} textAlign="center" fontSize="xs" color="gray.500">
        <ChakraLink as={Link} href="/login" color="teal.500" fontWeight="medium">
          Ir para login
        </ChakraLink>
      </Text>
    </Box>
  );
}

export default function VerifyEmailPage() {
  return (
    <PageWrapper gap={16}>
      <MainNav showLinks={false} showActions={false} />
      <Suspense
        fallback={
          <Box
            mx="auto"
            w="full"
            maxW="xl"
            borderRadius="3xl"
            bg="white"
            p={{ base: 6, md: 8 }}
            boxShadow="0 30px 80px rgba(35, 49, 82, 0.12)"
          >
            <Heading size="2xl" color="gray.900">
              Verificar e-mail
            </Heading>
            <Text mt={2} fontSize="sm" color="gray.500">
              Carregando...
            </Text>
          </Box>
        }
      >
        <VerifyEmailContent />
      </Suspense>
    </PageWrapper>
  );
}
