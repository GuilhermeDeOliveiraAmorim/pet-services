"use client";

import { useMemo, useState } from "react";
import { Eye, EyeOff } from "lucide-react";
import {
  Box,
  Button,
  IconButton,
  Input,
  InputGroup,
  Text,
  VStack,
  chakra,
} from "@chakra-ui/react";

import { useUserChangePassword } from "@/application";
import { getApiErrorMessage } from "@/lib/api-error";

export default function ChangePasswordCard() {
  const { mutateAsync, isPending, error, isSuccess } = useUserChangePassword();

  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showCurrent, setShowCurrent] = useState(false);
  const [showNew, setShowNew] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);

  const feedback = useMemo(() => {
    if (!error) {
      return "";
    }

    return getApiErrorMessage(
      error,
      "Não foi possível alterar a senha. Verifique os dados e tente novamente.",
    );
  }, [error]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (newPassword !== confirmPassword) {
      return;
    }

    await mutateAsync({
      oldPassword: currentPassword,
      newPassword,
    });

    setCurrentPassword("");
    setNewPassword("");
    setConfirmPassword("");
  };

  const isConfirmInvalid =
    confirmPassword.length > 0 && newPassword !== confirmPassword;

  return (
    <Box
      borderRadius="3xl"
      bg="white"
      p={{ base: 5, md: 6 }}
      borderWidth="1px"
      borderColor="gray.200"
      shadow="sm"
    >
      <Box mb={6}>
        <Text
          fontSize="xs"
          fontWeight="semibold"
          textTransform="uppercase"
          color="gray.500"
        >
          Segurança
        </Text>
        <Text mt={2} fontSize="xl" fontWeight="semibold" color="gray.900">
          Alterar senha
        </Text>
        <Text mt={2} fontSize="sm" color="gray.600">
          Atualize sua senha para manter sua conta protegida.
        </Text>
      </Box>

      <chakra.form onSubmit={handleSubmit}>
        <VStack align="stretch" gap={5}>
          <Box>
            <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
              Senha atual
            </Text>
            <InputGroup
              endElement={
                <IconButton
                  aria-label={showCurrent ? "Ocultar senha" : "Mostrar senha"}
                  variant="ghost"
                  size="sm"
                  borderRadius="full"
                  color="gray.500"
                  onClick={() => setShowCurrent((prev) => !prev)}
                >
                  {showCurrent ? <EyeOff size={16} /> : <Eye size={16} />}
                </IconButton>
              }
            >
              <Input
                type={showCurrent ? "text" : "password"}
                value={currentPassword}
                onChange={(event) => setCurrentPassword(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="********"
                required
                autoComplete="current-password"
              />
            </InputGroup>
          </Box>

          <Box>
            <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
              Nova senha
            </Text>
            <InputGroup
              endElement={
                <IconButton
                  aria-label={showNew ? "Ocultar senha" : "Mostrar senha"}
                  variant="ghost"
                  size="sm"
                  borderRadius="full"
                  color="gray.500"
                  onClick={() => setShowNew((prev) => !prev)}
                >
                  {showNew ? <EyeOff size={16} /> : <Eye size={16} />}
                </IconButton>
              }
            >
              <Input
                type={showNew ? "text" : "password"}
                value={newPassword}
                onChange={(event) => setNewPassword(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="********"
                required
                autoComplete="new-password"
              />
            </InputGroup>
          </Box>

          <Box>
            <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
              Confirmar nova senha
            </Text>
            <InputGroup
              endElement={
                <IconButton
                  aria-label={showConfirm ? "Ocultar senha" : "Mostrar senha"}
                  variant="ghost"
                  size="sm"
                  borderRadius="full"
                  color="gray.500"
                  onClick={() => setShowConfirm((prev) => !prev)}
                >
                  {showConfirm ? <EyeOff size={16} /> : <Eye size={16} />}
                </IconButton>
              }
            >
              <Input
                type={showConfirm ? "text" : "password"}
                value={confirmPassword}
                onChange={(event) => setConfirmPassword(event.target.value)}
                h="11"
                borderRadius="xl"
                bg="gray.50"
                borderColor="gray.200"
                focusRingColor="teal.200"
                placeholder="********"
                required
                autoComplete="new-password"
              />
            </InputGroup>
            {isConfirmInvalid ? (
              <Text mt={1.5} fontSize="xs" color="red.500">
                As senhas não coincidem.
              </Text>
            ) : null}
          </Box>

          {feedback ? (
            <Box
              borderRadius="2xl"
              borderWidth="1px"
              borderColor="red.200"
              bg="red.50"
              px={4}
              py={3}
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
              px={4}
              py={3}
            >
              <Text fontSize="sm" color="green.700">
                Senha atualizada com sucesso.
              </Text>
            </Box>
          ) : null}

          <Button
            type="submit"
            disabled={isPending || isConfirmInvalid}
            h="11"
            w="full"
            borderRadius="full"
            bg="green.400"
            color="white"
            _hover={{ bg: "green.500" }}
            _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
          >
            {isPending ? "Atualizando..." : "Atualizar senha"}
          </Button>
        </VStack>
      </chakra.form>
    </Box>
  );
}
