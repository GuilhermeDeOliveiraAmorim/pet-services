import { Box, Text } from "@chakra-ui/react";

export default function Footer() {
  return (
    <Box as="footer" py={6} bg="gray.100" textAlign="center">
      <Text color="gray.600" fontSize="sm">
        Pet Services © {new Date().getFullYear()} — Plataforma para o
        ecossistema pet
      </Text>
    </Box>
  );
}
