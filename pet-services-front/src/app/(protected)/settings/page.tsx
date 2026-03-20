"use client";

import { Box, Grid, Text } from "@chakra-ui/react";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import ChangePasswordCard from "@/components/account/ChangePasswordCard";
import AccountStatusCard from "@/components/account/AccountStatusCard";

export default function SettingsPage() {
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
        <Text
          fontSize={{ base: "xs" }}
          fontWeight="semibold"
          textTransform="uppercase"
          color="teal.500"
        >
          Configurações
        </Text>
        <Text
          mt={2}
          fontSize={{ base: "xl", sm: "2xl" }}
          fontWeight="semibold"
          color="gray.900"
        >
          Gerencie sua conta
        </Text>
        <Text mt={2} fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
          Altere sua senha e status da conta.
        </Text>
      </Box>
      <Grid
        gap={6}
        alignItems="start"
        templateColumns={{ base: "1fr", lg: "1fr 1fr" }}
      >
        <ChangePasswordCard />
        <AccountStatusCard />
      </Grid>
    </PageWrapper>
  );
}
