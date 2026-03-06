import { Box, Button, Grid, Heading, HStack, Text, VStack } from "@chakra-ui/react";

const services = [
  {
    title: "Banho & Tosa",
    desc: "Profissionais com carinho e técnica.",
  },
  { title: "Pet Care", desc: "Planos mensais personalizados." },
  {
    title: "Tratamentos",
    desc: "Exames e acompanhamento clínico.",
  },
  { title: "Vacinação", desc: "Calendário completo de vacinas." },
];

export default function ServicesSection() {
  return (
    <VStack id="services" align="stretch" gap="8">
      <HStack justify="space-between" align="center">
        <Box>
          <Text fontSize="xs" fontWeight="semibold" textTransform="uppercase" color="teal.400">
            Serviços
          </Text>
          <Heading size="xl" color="gray.900" mt="1">
            Nossos serviços para o seu pet
          </Heading>
        </Box>
        <Button
          borderRadius="full"
          variant="outline"
          colorPalette="gray"
          size="sm"
          fontWeight="medium"
        >
          Ver todos
        </Button>
      </HStack>

      <Grid gap="4" templateColumns={{ base: "1fr", sm: "1fr 1fr", lg: "repeat(4, 1fr)" }}>
        {services.map((item) => (
          <Box key={item.title} borderRadius="3xl" bg="white" p="6" boxShadow="sm">
            <Box mb="4" h="10" w="10" borderRadius="2xl" bg="teal.100" />
            <Text fontSize="sm" fontWeight="semibold" color="gray.900">
              {item.title}
            </Text>
            <Text mt="2" fontSize="sm" color="gray.600">
              {item.desc}
            </Text>
          </Box>
        ))}
      </Grid>
    </VStack>
  );
}