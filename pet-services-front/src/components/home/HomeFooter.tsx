import { Box, HStack, Text } from "@chakra-ui/react";

export default function HomeFooter() {
  return (
    <>
      <Box h="1px" w="full" bg="gray.200" />
      <HStack
        as="footer"
        flexDir={{ base: "column", lg: "row" }}
        align="center"
        justify="space-between"
        gap="6"
        pt="8"
        fontSize="sm"
        color="gray.500"
      >
        <HStack gap="2">
          <Box
            display="flex"
            h="9"
            w="9"
            alignItems="center"
            justifyContent="center"
            borderRadius="2xl"
            bgGradient="linear(to-tr, teal.400, cyan.400)"
            color="white"
            fontWeight="semibold"
          >
            pet
          </Box>
          <Text color="gray.700">PetServices</Text>
        </HStack>
        <Text>© 2026 PetServices. Todos os direitos reservados.</Text>
      </HStack>
    </>
  );
}