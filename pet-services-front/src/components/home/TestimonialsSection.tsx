import { Box, Grid, Heading, HStack, Text, VStack } from "@chakra-ui/react";

const highlights = [
  "Atendimento incrível",
  "Equipe cuidadosa",
  "Ambiente seguro",
  "Veterinários dedicados",
];

export default function TestimonialsSection() {
  return (
    <Grid id="testimonials" gap="10" templateColumns={{ base: "1fr", lg: "1fr 1fr" }}>
      <VStack align="stretch" gap="6">
        <Text fontSize="xs" fontWeight="semibold" textTransform="uppercase" color="pink.400">
          Depoimentos
        </Text>
        <Heading size="xl" color="gray.900">
          O que nossos clientes dizem?
        </Heading>

        <Box borderRadius="3xl" bg="white" p="6" boxShadow="sm">
          <Text fontSize="sm" color="gray.600">
            “Meu cachorro voltou super feliz e o acompanhamento foi perfeito.
            Atendimento rápido e equipe cuidadosa.”
          </Text>

          <HStack mt="4" gap="3">
            <Box
              h="10"
              w="10"
              borderRadius="full"
              bgGradient="linear(to-br, orange.200, pink.200)"
              display="flex"
              alignItems="center"
              justifyContent="center"
              fontWeight="semibold"
              color="gray.700"
            >
              MT
            </Box>
            <Box>
              <Text fontSize="sm" fontWeight="semibold" color="gray.900">
                Mariana Torres
              </Text>
              <Text fontSize="xs" color="gray.500">
                Tutora
              </Text>
            </Box>
          </HStack>
        </Box>

        <Grid gap="4" templateColumns={{ base: "1fr", sm: "1fr 1fr" }}>
          {highlights.map((text) => (
            <Box key={text} borderRadius="2xl" bg="white" p="4" fontSize="sm" color="gray.600" boxShadow="sm">
              {text}
            </Box>
          ))}
        </Grid>
      </VStack>

      <Box position="relative" display="flex" alignItems="center" justifyContent="center">
        <Box
          position="absolute"
          inset="6"
          borderRadius="40px"
          bgGradient="linear(to-br, teal.100, cyan.100)"
        />
        <VStack
          position="relative"
          zIndex="1"
          w="full"
          maxW="sm"
          align="stretch"
          gap="4"
          borderRadius="3xl"
          bg="white"
          p="6"
          boxShadow="0 30px 80px rgba(124,139,255,0.2)"
        >
          <HStack justify="space-between">
            <Text fontSize="sm" fontWeight="semibold" color="gray.900">
              Consulta
            </Text>
            <Text borderRadius="full" bg="pink.100" px="3" py="1" fontSize="xs" fontWeight="semibold" color="pink.500">
              24/7
            </Text>
          </HStack>
          <Box h="48" borderRadius="3xl" bgGradient="linear(to-br, yellow.100, orange.100)" />
          <Box borderRadius="2xl" bg="gray.50" px="4" py="3" fontSize="sm" color="gray.600">
            Cuidamos com dedicação e tecnologia em cada detalhe.
          </Box>
        </VStack>
      </Box>
    </Grid>
  );
}