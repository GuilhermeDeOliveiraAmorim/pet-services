import ChangePasswordCard from "@/components/account/ChangePasswordCard";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";
import { Box, Flex, Grid, Text, VStack } from "@chakra-ui/react";

export default function ProviderDashboardPage() {
  return (
    <PageWrapper gap={10}>
      <MainNav />

      <VStack align="stretch" gap={6}>
        <Box>
          <Text
            fontSize="xs"
            fontWeight="semibold"
            textTransform="uppercase"
            color="green.500"
          >
            Dashboard
          </Text>
          <Text mt={2} fontSize="2xl" fontWeight="semibold" color="gray.900">
            Olá, prestador
          </Text>
          <Text mt={2} fontSize="sm" color="gray.600">
            Este é o seu painel inicial. Aqui vão aparecer seus serviços e
            agendamentos.
          </Text>
        </Box>

        <Box
          borderRadius="3xl"
          borderWidth="1px"
          borderStyle="dashed"
          borderColor="gray.300"
          bg="white"
          px={6}
          py={16}
          textAlign="center"
        >
          <Flex
            mx="auto"
            h="12"
            w="12"
            align="center"
            justify="center"
            borderRadius="2xl"
            bg="green.50"
            color="green.500"
            fontSize="lg"
            fontWeight="semibold"
          >
            🩺
          </Flex>
          <Text mt={4} fontSize="lg" fontWeight="semibold" color="gray.900">
            Sem dados ainda
          </Text>
          <Text mt={2} fontSize="sm" color="gray.600">
            Quando você cadastrar seu primeiro serviço ou receber um
            agendamento, as informações vão aparecer aqui.
          </Text>
        </Box>

        <Grid
          gap={4}
          templateColumns={{ base: "1fr", md: "repeat(3, minmax(0, 1fr))" }}
        >
          <Box
            borderRadius="2xl"
            borderWidth="1px"
            borderColor="gray.200"
            bg="white"
            p={4}
          >
            <Text fontSize="sm" fontWeight="semibold" color="gray.900">
              Serviços
            </Text>
            <Text mt={2} fontSize="sm" color="gray.600">
              Estrutura pronta para listagem, criação e edição de serviços.
            </Text>
          </Box>

          <Box
            borderRadius="2xl"
            borderWidth="1px"
            borderColor="gray.200"
            bg="white"
            p={4}
          >
            <Text fontSize="sm" fontWeight="semibold" color="gray.900">
              Agenda
            </Text>
            <Text mt={2} fontSize="sm" color="gray.600">
              Estrutura pronta para próximos agendamentos e confirmações.
            </Text>
          </Box>

          <Box
            borderRadius="2xl"
            borderWidth="1px"
            borderColor="gray.200"
            bg="white"
            p={4}
          >
            <Text fontSize="sm" fontWeight="semibold" color="gray.900">
              Indicadores
            </Text>
            <Text mt={2} fontSize="sm" color="gray.600">
              Estrutura pronta para métricas de atendimento e avaliações.
            </Text>
          </Box>
        </Grid>
      </VStack>

      <ChangePasswordCard />
    </PageWrapper>
  );
}
