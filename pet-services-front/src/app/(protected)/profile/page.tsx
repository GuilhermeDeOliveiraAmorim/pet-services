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
import { ErrorState, LoadingState } from "@/components/ui";
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
        borderRadius={{ base: "2xl", md: "3xl" }}
        bg="white"
        p={{ base: 4, sm: 5, md: 7 }}
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
              fontSize={{ base: "xs" }}
              fontWeight="semibold"
              textTransform="uppercase"
              color="green.500"
            >
              Perfil
            </Text>
            <Text
              mt={2}
              fontSize={{ base: "xl", sm: "2xl" }}
              fontWeight="semibold"
              color="gray.900"
            >
              Seus dados
            </Text>
            <Text mt={2} fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
              Atualize suas informações para manter seu perfil completo.
            </Text>
          </Box>
          {user ? (
            <Badge
              borderRadius="full"
              px={3}
              py={1}
              fontSize={{ base: "xs" }}
              colorPalette={user.active ? "green" : "orange"}
              variant="subtle"
            >
              {user.active ? "Ativo" : "Inativo"}
            </Badge>
          ) : null}
        </Flex>

        {isLoading ? (
          <LoadingState message="Carregando perfil..." />
        ) : user ? (
          <ProfileForm key={user.id} user={user} />
        ) : (
          <ErrorState message="Não foi possível carregar o perfil." />
        )}
      </Box>

      <ChangePasswordCard />

      <Box
        borderRadius={{ base: "2xl", md: "3xl" }}
        borderWidth="1px"
        borderColor="gray.200"
        bg="white"
        p={{ base: 4, sm: 5, md: 6 }}
        shadow="sm"
      >
        <Box mb={4}>
          <Text
            fontSize={{ base: "xs" }}
            fontWeight="semibold"
            textTransform="uppercase"
            color="gray.500"
          >
            Conta
          </Text>
          <Text
            mt={2}
            fontSize={{ base: "lg", md: "xl" }}
            fontWeight="semibold"
            color="gray.900"
          >
            Status da conta
          </Text>
          <Text mt={2} fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
            Você pode desativar sua conta e reativá-la depois.
          </Text>
        </Box>

        {accountFeedback ? (
          <Box
            mb={4}
            borderRadius={{ base: "xl", md: "2xl" }}
            borderWidth="1px"
            borderColor="red.200"
            bg="red.50"
            px={{ base: 3, md: 4 }}
            py={3}
          >
            <Text fontSize={{ base: "xs", sm: "sm" }} color="red.600">
              {accountFeedback}
            </Text>
          </Box>
        ) : null}

        <HStack
          gap={{ base: 2, sm: 3 }}
          flexWrap="wrap"
          justify={{ base: "center", sm: "flex-start" }}
        >
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
            fontSize={{ base: "xs", sm: "sm" }}
            h={{ base: "9", md: "10" }}
          >
            {isDeactivating ? "Desativando..." : "Desativar"}
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
            fontSize={{ base: "xs", sm: "sm" }}
            h={{ base: "9", md: "10" }}
          >
            {isReactivating ? "Reativando..." : "Reativar"}
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
                fontSize={{ base: "xs", sm: "sm" }}
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
                fontSize={{ base: "xs", sm: "sm" }}
                h={{ base: "9", md: "10" }}
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
