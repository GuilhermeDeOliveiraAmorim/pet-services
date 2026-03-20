"use client";

import { Box, Flex, Text } from "@chakra-ui/react";

import { useUserProfile } from "@/application";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import { ErrorState, LoadingState } from "@/components/ui";
import ProfileForm from "./components/ProfileForm";

export default function ProfilePage() {
  const { data, isLoading } = useUserProfile();
  const user = data?.user;

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
        </Flex>

        {isLoading ? (
          <LoadingState message="Carregando perfil..." />
        ) : user ? (
          <ProfileForm key={user.id} user={user} />
        ) : (
          <ErrorState message="Não foi possível carregar o perfil." />
        )}
      </Box>
    </PageWrapper>
  );
}
