"use client";

import { type ChangeEvent, Suspense, useMemo, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import {
  Box,
  Button,
  Heading,
  IconButton,
  Input,
  Link as ChakraLink,
  Text,
  VStack,
  chakra,
} from "@chakra-ui/react";
import { Eye, EyeOff } from "lucide-react";

import { useAuthResetPassword } from "@/application";
import { getApiErrorMessage } from "@/lib/api-error";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

function ResetPasswordContent() {
  const searchParams = useSearchParams();
  const token = searchParams.get("token") ?? "";
  const { mutateAsync, isPending, error, isSuccess } = useAuthResetPassword();

  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showNew, setShowNew] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [newPasswordError, setNewPasswordError] = useState("");
  const [confirmPasswordError, setConfirmPasswordError] = useState("");

  const validateNewPassword = (value: string) => {
    if (!value.trim()) {
      return "Informe a nova senha";
    }

    return "";
  };

  const validateConfirmPassword = (newPass: string, confirmPass: string) => {
    if (!confirmPass.trim()) {
      return "Confirme a nova senha";
    }

    if (newPass !== confirmPass) {
      return "As senhas não coincidem.";
    }

    return "";
  };

  const feedback = useMemo(() => {
    if (!error) {
      return "";
    }

    return getApiErrorMessage(error, "Não foi possível redefinir a senha.");
  }, [error]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const nextNewPasswordError = validateNewPassword(newPassword);
    const nextConfirmPasswordError = validateConfirmPassword(
      newPassword,
      confirmPassword,
    );

    setNewPasswordError(nextNewPasswordError);
    setConfirmPasswordError(nextConfirmPasswordError);

    if (!token || nextNewPasswordError || nextConfirmPasswordError) {
      return;
    }

    await mutateAsync({ token, newPassword });
    setNewPassword("");
    setConfirmPassword("");
  };

  const isConfirmInvalid = Boolean(confirmPasswordError);

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
          Definir nova senha
        </Heading>
        <Text mt={2} fontSize="sm" color="gray.500">
          Crie uma nova senha para acessar sua conta.
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
            Token inválido ou ausente. Solicite uma nova redefinição.
          </Text>
        </Box>
      ) : (
        <chakra.form onSubmit={handleSubmit}>
          <VStack gap={5} align="stretch">
            <Box>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Nova senha
              </Text>
              <Box position="relative">
                <Input
                  type={showNew ? "text" : "password"}
                  value={newPassword}
                  onChange={(event: ChangeEvent<HTMLInputElement>) => {
                    setNewPassword(event.target.value);
                    if (newPasswordError) {
                      setNewPasswordError("");
                    }
                    if (confirmPasswordError && confirmPassword) {
                      setConfirmPasswordError(
                        validateConfirmPassword(
                          event.target.value,
                          confirmPassword,
                        ),
                      );
                    }
                  }}
                  onBlur={() =>
                    setNewPasswordError(validateNewPassword(newPassword))
                  }
                  placeholder="********"
                  autoComplete="new-password"
                  h="11"
                  borderRadius="xl"
                  bg="gray.50"
                  borderColor={newPasswordError ? "red.300" : "gray.200"}
                  focusRingColor={newPasswordError ? "red.200" : "teal.200"}
                  pe="12"
                />
                <IconButton
                  type="button"
                  onClick={() => setShowNew((prev) => !prev)}
                  aria-label={showNew ? "Ocultar senha" : "Mostrar senha"}
                  variant="ghost"
                  size="sm"
                  color="gray.500"
                  position="absolute"
                  right="2"
                  top="50%"
                  transform="translateY(-50%)"
                >
                  {showNew ? <EyeOff size={16} /> : <Eye size={16} />}
                </IconButton>
              </Box>
              {newPasswordError ? (
                <Text mt={1.5} fontSize="xs" color="red.500">
                  {newPasswordError}
                </Text>
              ) : null}
            </Box>

            <Box>
              <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                Confirmar nova senha
              </Text>
              <Box position="relative">
                <Input
                  type={showConfirm ? "text" : "password"}
                  value={confirmPassword}
                  onChange={(event: ChangeEvent<HTMLInputElement>) => {
                    setConfirmPassword(event.target.value);
                    if (confirmPasswordError) {
                      setConfirmPasswordError("");
                    }
                  }}
                  onBlur={() =>
                    setConfirmPasswordError(
                      validateConfirmPassword(newPassword, confirmPassword),
                    )
                  }
                  placeholder="********"
                  autoComplete="new-password"
                  h="11"
                  borderRadius="xl"
                  bg="gray.50"
                  borderColor={confirmPasswordError ? "red.300" : "gray.200"}
                  focusRingColor={confirmPasswordError ? "red.200" : "teal.200"}
                  pe="12"
                />
                <IconButton
                  type="button"
                  onClick={() => setShowConfirm((prev) => !prev)}
                  aria-label={showConfirm ? "Ocultar senha" : "Mostrar senha"}
                  variant="ghost"
                  size="sm"
                  color="gray.500"
                  position="absolute"
                  right="2"
                  top="50%"
                  transform="translateY(-50%)"
                >
                  {showConfirm ? <EyeOff size={16} /> : <Eye size={16} />}
                </IconButton>
              </Box>
              {confirmPasswordError ? (
                <Text mt={1.5} fontSize="xs" color="red.500">
                  {confirmPasswordError}
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
                  Senha redefinida com sucesso. Faça login novamente.
                </Text>
              </Box>
            ) : null}

            <Button
              type="submit"
              disabled={isPending || isConfirmInvalid}
              h="11"
              borderRadius="full"
              bgGradient="linear(to-r, green.400, teal.400)"
              color="white"
              _hover={{ opacity: 0.92 }}
              _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
            >
              {isPending ? "Atualizando..." : "Atualizar senha"}
            </Button>
          </VStack>
        </chakra.form>
      )}

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
    </Box>
  );
}

export default function ResetPasswordPage() {
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
              Definir nova senha
            </Heading>
            <Text mt={2} fontSize="sm" color="gray.500">
              Carregando...
            </Text>
          </Box>
        }
      >
        <ResetPasswordContent />
      </Suspense>
    </PageWrapper>
  );
}
