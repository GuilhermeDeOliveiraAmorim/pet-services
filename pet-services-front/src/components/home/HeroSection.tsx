import {
  Box,
  Button,
  Heading,
  HStack,
  Input,
  NativeSelect,
  Text,
} from "@chakra-ui/react";

import Banner from "../../../public/banner.png";
import { serviceOptions } from "./service-options";

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

        <form action="/services" method="get" style={{ width: "100%" }}>
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
              name="zip_code"
              flex="1"
              minW={{ base: "full", md: "220px" }}
              border="none"
              borderRadius="0"
              px="4"
              py="2"
              fontSize="sm"
              color="gray.700"
              placeholder="CEP (ex: 01310-100)"
              _focusVisible={{ boxShadow: "none" }}
            />

            <NativeSelect.Root
              flexShrink={0}
              w={{ base: "full", md: "280px" }}
              minW={{ base: "full", md: "280px" }}
              border="none"
              borderRadius="0"
            >
              <NativeSelect.Field
                name="q"
                bg="white"
                px="4"
                py="2"
                fontSize="sm"
                color="gray.700"
                border="none"
              >
                <option value="">Qual serviço seu pet precisa?</option>
                {serviceOptions.map((option) => (
                  <option key={option.label} value={option.label}>
                    {option.label}
                  </option>
                ))}
              </NativeSelect.Field>
              <NativeSelect.Indicator />
            </NativeSelect.Root>

            <Button
              type="submit"
              flexShrink={0}
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
        </form>
      </Box>
    </Box>
  );
}