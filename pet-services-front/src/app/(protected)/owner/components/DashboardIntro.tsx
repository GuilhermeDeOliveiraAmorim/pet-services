import { Box, Flex, Text, VStack } from "@chakra-ui/react";

type DashboardIntroProps = {
  hasPets: boolean;
  petsCount: number;
};

export default function DashboardIntro({
  hasPets,
  petsCount,
}: DashboardIntroProps) {
  return (
    <VStack align="stretch" gap={{ base: 4, md: 6 }}>
      <Box>
        <Text
          fontSize="xs"
          fontWeight="semibold"
          textTransform="uppercase"
          color="green.500"
        >
          Dashboard
        </Text>
        <Text
          mt={2}
          fontSize={{ base: "xl", md: "2xl" }}
          fontWeight="semibold"
          color="gray.900"
        >
          Olá, tutor
        </Text>
        <Text mt={2} fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
          Este é o seu painel inicial. Aqui vão aparecer seus pets e
          agendamentos.
        </Text>
      </Box>

      <Box
        borderRadius={{ base: "2xl", md: "3xl" }}
        borderWidth="1px"
        borderStyle="dashed"
        borderColor="gray.300"
        bg="white"
        px={{ base: 4, sm: 6 }}
        py={{ base: 8, md: 16 }}
        textAlign="center"
      >
        <Flex
          mx="auto"
          h={{ base: "10", md: "12" }}
          w={{ base: "10", md: "12" }}
          align="center"
          justify="center"
          borderRadius={{ base: "lg", md: "2xl" }}
          bg="green.50"
          color="green.500"
          fontSize={{ base: "base", md: "lg" }}
          fontWeight="semibold"
        >
          {hasPets ? "✅" : "🐾"}
        </Flex>
        <Text
          mt={3}
          fontSize={{ base: "sm", md: "lg" }}
          fontWeight="semibold"
          color="gray.900"
        >
          {hasPets ? "Tudo pronto" : "Comece por aqui"}
        </Text>
        <Text mt={2} fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
          {hasPets
            ? `Você já tem ${petsCount} pet${petsCount > 1 ? "s" : ""} cadastrado${petsCount > 1 ? "s" : ""}. Edite os dados, adicione fotos e acompanhe seus agendamentos.`
            : "Cadastre seu primeiro pet no formulário abaixo para liberar agendamentos e histórico no painel."}
        </Text>
      </Box>
    </VStack>
  );
}
