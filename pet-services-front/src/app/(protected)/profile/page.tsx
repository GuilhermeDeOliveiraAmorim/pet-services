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
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import ChangePasswordCard from "@/components/account/ChangePasswordCard";
import ProfileForm from "./components/ProfileForm";

export default function ProfilePage() {
  const router = useRouter();
  const { clearSession } = useAuthSession();
  const { data, isLoading } = useUserProfile();
  const [isDeactivateOpen, setIsDeactivateOpen] = useState(false);
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
    <PageWrapper gap={10}>
      <MainNav />

      <Box
        borderRadius="3xl"
        bg="white"
        p={{ base: 5, md: 7 }}
        borderWidth="1px"
        borderColor="gray.200"
        shadow="sm"
      >
        <Flex mb={6} wrap="wrap" align="start" justify="space-between" gap={4}>
          <Box>
            <Text
              fontSize="xs"
              fontWeight="semibold"
              textTransform="uppercase"
              color="green.500"
            >
              Perfil
            </Text>
            <Text mt={2} fontSize="2xl" fontWeight="semibold" color="gray.900">
              Seus dados
            </Text>
            <Text mt={2} fontSize="sm" color="gray.600">
              Atualize suas informações para manter seu perfil completo.
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
              {user.active ? "Conta ativa" : "Conta desativada"}
            </Badge>
          ) : null}
        </Flex>

        {isLoading ? (
          <Text fontSize="sm" color="gray.500">
            Carregando perfil...
          </Text>
        ) : user ? (
          <ProfileForm key={user.id} user={user} />
        ) : (
          <Text fontSize="sm" color="gray.500">
            Não foi possível carregar o perfil.
          </Text>
        )}
      </Box>

      <ChangePasswordCard />

      <Box
        borderRadius="3xl"
        borderWidth="1px"
        borderColor="gray.200"
        bg="white"
        p={{ base: 5, md: 6 }}
        shadow="sm"
      >
        <Box mb={4}>
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

        {accountFeedback ? (
          <Box
            mb={4}
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

        <HStack gap={3} flexWrap="wrap">
          <Button
            type="button"
            onClick={() => setIsDeactivateOpen(true)}
            disabled={isDeactivating || !user || !user.active}
            borderRadius="full"
            borderWidth="1px"
            borderColor="red.200"
            color="red.600"
            bg="white"
            _hover={{ bg: "red.50" }}
            _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
          >
            {isDeactivating ? "Desativando..." : "Desativar conta"}
          </Button>

          <Button
            type="button"
            onClick={handleReactivate}
            disabled={isReactivating || !user || user.active}
            borderRadius="full"
            borderWidth="1px"
            borderColor="gray.200"
            color="gray.700"
            bg="white"
            _hover={{ bg: "gray.50" }}
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
          bg="blackAlpha.500"
          zIndex={1400}
          display="flex"
          alignItems="center"
          justifyContent="center"
          p={4}
        >
          <Box
            role="dialog"
            aria-modal="true"
            maxW="md"
            w="full"
            borderRadius="3xl"
            bg="white"
            p={6}
            shadow="xl"
          >
            <Text fontSize="lg" fontWeight="semibold" color="gray.900">
              Desativar conta?
            </Text>
            <Text mt={2} fontSize="sm" color="gray.600">
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
              >
                {isDeactivating ? "Desativando..." : "Confirmar"}
              </Button>
            </Flex>
          </Box>
        </Box>
      ) : null}
    </PageWrapper>
  );
}
