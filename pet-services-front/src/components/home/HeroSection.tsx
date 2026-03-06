import {
  Box,
  Button,
  Heading,
  HStack,
  Input,
  Text,
} from "@chakra-ui/react";

import Banner from "../../../public/banner.png";

const serviceOptions = [
  "Clínica Veterinária",
  "Pet Shop",
  "Banho e Tosa",
  "Hotelzinho e Creche",
  "Passeador(a)",
  "Pet Sitter",
  "Adestrador",
];

export default function HeroSection() {
  return (
    <Box
      id="home"
      position="relative"
      display="flex"
      alignItems="center"
      justifyContent="center"
      py={{ base: "16", md: "24" }}
      bgImage={`url(${Banner.src})`}
      bgSize="cover"
      bgPos="center"
      borderRadius="3xl"
      overflow="hidden"
    >
      <Box position="absolute" inset="0" bg="blackAlpha.600" />

      <Box
        position="relative"
        zIndex="1"
        display="flex"
        w="full"
        maxW="3xl"
        flexDir="column"
        alignItems="center"
        gap="8"
        textAlign="center"
        px="4"
      >
        <Heading size={{ base: "2xl", md: "4xl" }} color="white" lineHeight="1.2">
          Cuidado Profissional para quem você Ama
        </Heading>

        <Text maxW="xl" fontSize="md" lineHeight="7" color="white">
          Conecte-se com os melhores veterinários, pet shops e cuidadores da sua
          região.
        </Text>

        <HStack
          w="full"
          maxW="xl"
          align="stretch"
          gap="0"
          borderWidth="1px"
          borderColor="gray.300"
          borderRadius="full"
          overflow="hidden"
          bg="white"
          flexWrap={{ base: "wrap", md: "nowrap" }}
        >
          <Input
            flex="1"
            minW={{ base: "full", md: "0" }}
            border="none"
            borderRadius="0"
            px="4"
            py="2"
            fontSize="sm"
            color="gray.700"
            placeholder="Onde você está?"
            _focusVisible={{ boxShadow: "none" }}
          />

          <Box
            as="select"
            bg="white"
            px="4"
            py="2"
            fontSize="sm"
            color="gray.700"
            border="none"
            outline="none"
            minW={{ base: "full", md: "280px" }}
          >
            <option value="">Qual serviço seu pet precisa?</option>
            {serviceOptions.map((option) => (
              <option key={option} value={option}>
                {option}
              </option>
            ))}
          </Box>

          <Button
            borderRadius="0"
            px="6"
            bg="teal.400"
            color="white"
            fontSize="sm"
            fontWeight="semibold"
            _hover={{ bg: "teal.500" }}
          >
            Buscar
          </Button>
        </HStack>
      </Box>
    </Box>
  );
}