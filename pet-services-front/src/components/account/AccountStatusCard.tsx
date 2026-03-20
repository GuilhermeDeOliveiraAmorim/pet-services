"use client";

import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { Badge, Box, Button, Flex, HStack, Text } from "@chakra-ui/react";
import {
  useAuthSession,
  useUserDeactivate,
  useUserProfile,
  useUserReactivate,
} from "@/application";
import { getApiErrorMessage } from "@/lib/api-error";

export default function AccountStatusCard() {
  const router = useRouter();
  const [isDeactivateOpen, setIsDeactivateOpen] = useState(false);
  const { clearSession } = useAuthSession();
  const { data } = useUserProfile();
  const {
    mutateAsync: deactivateUser,
    isPending: isDeactivating,
    error: deactivateError,
  } = useUserDeactivate();
  const {
    mutateAsync: reactivateUser,
    isPending: isReactivating,
    error: reactivateError,
  } = useUserReactivate();

  const user = data?.user;

  const accountFeedback = useMemo(() => {
    const error = deactivateError ?? reactivateError;
    if (!error) {
      return "";
    }
    return getApiErrorMessage(
      error,
      "Não foi possível atualizar o status da conta.",
    );
  }, [deactivateError, reactivateError]);

  const handleConfirmDeactivate = async () => {
    await deactivateUser();
    clearSession();
    router.replace("/login");
  };

  const handleReactivate = async () => {
    await reactivateUser();
  };

  return (
    <>
      <Box
        borderRadius="3xl"
        bg="white"
        p={{ base: 5, md: 6 }}
        borderWidth="1px"
        borderColor="gray.200"
        shadow="sm"
      >
        <Flex
          mb={6}
          wrap="wrap"
          align="start"
          justify="space-between"
          gap={4}
          direction={{ base: "column", sm: "row" }}
        >
          <Box>
            <Text
              fontSize="xs"
              fontWeight="semibold"
              textTransform="uppercase"
              color="gray.500"
            >
              Conta
            </Text>
            <Text mt={2} fontSize="xl" fontWeight="semibold" color="gray.900">
              Status da conta
            </Text>
            <Text mt={2} fontSize="sm" color="gray.600">
              Você pode desativar sua conta e reativá-la depois.
            </Text>
          </Box>
          {user ? (
            <Badge
              borderRadius="full"
              px={3}
              py={1}
              fontSize="xs"
              colorPalette={user.active ? "green" : "orange"}
              variant="subtle"
            >
              {user.active ? "Ativo" : "Inativo"}
            </Badge>
          ) : null}
        </Flex>

        {accountFeedback ? (
          <Box
            mb={5}
            borderRadius="2xl"
            borderWidth="1px"
            borderColor="red.200"
            bg="red.50"
            px={4}
            py={3}
          >
            <Text fontSize="sm" color="red.600">
              {accountFeedback}
            </Text>
          </Box>
        ) : null}

        <HStack
          gap={3}
          flexWrap="wrap"
          justify={{ base: "center", sm: "flex-start" }}
        >
          <Button
            type="button"
            onClick={() => setIsDeactivateOpen(true)}
            disabled={isDeactivating || !user || !user.active}
            h="10"
            borderRadius="full"
            borderWidth="1px"
            borderColor="orange.300"
            color="orange.700"
            bg="white"
            _hover={{ bg: "orange.50" }}
            _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
          >
            {isDeactivating ? "Desativando..." : "Desativar conta"}
          </Button>
          <Button
            type="button"
            onClick={handleReactivate}
            disabled={isReactivating || !user || user.active}
            h="10"
            borderRadius="full"
            bg="green.400"
            color="white"
            _hover={{ bg: "green.500" }}
            _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
          >
            {isReactivating ? "Reativando..." : "Reativar conta"}
          </Button>
        </HStack>
      </Box>

      {isDeactivateOpen ? (
        <Box
          position="fixed"
          inset={0}
          zIndex={30}
          bg="blackAlpha.400"
          display="flex"
          alignItems="center"
          justifyContent="center"
          px={4}
        >
          <Box
            role="dialog"
            aria-modal="true"
            maxW="md"
            w="full"
            borderRadius={{ base: "2xl", md: "3xl" }}
            bg="white"
            p={{ base: 4, md: 6 }}
            shadow="xl"
          >
            <Text
              fontSize={{ base: "base", md: "lg" }}
              fontWeight="semibold"
              color="gray.900"
            >
              Desativar conta?
            </Text>
            <Text mt={2} fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
              Ao desativar sua conta, você será desconectado e precisará
              reativá-la para acessar novamente.
            </Text>

            <Flex mt={6} justify="flex-end" gap={3} wrap="wrap">
              <Button
                type="button"
                borderRadius="full"
                borderWidth="1px"
                borderColor="gray.200"
                bg="white"
                color="gray.700"
                onClick={() => setIsDeactivateOpen(false)}
                h={{ base: "9", md: "10" }}
              >
                Cancelar
              </Button>
              <Button
                type="button"
                borderRadius="full"
                bg="red.500"
                color="white"
                onClick={handleConfirmDeactivate}
                disabled={isDeactivating}
                _hover={{ bg: "red.600" }}
                _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
                h={{ base: "9", md: "10" }}
              >
                {isDeactivating ? "Desativando..." : "Confirmar"}
              </Button>
            </Flex>
          </Box>
        </Box>
      ) : null}
    </>
  );
}
