import { Box, HStack, Link, Text } from "@chakra-ui/react";

export default function InstitutionalSection() {
  return (
    <Box as="section" py={8} maxW="6xl" mx="auto" textAlign="center">
      <HStack gap={8} justify="center" flexWrap="wrap">
        <Link href="/sobre" color="teal.700" fontWeight="bold">
          Sobre
        </Link>
        <Link href="/blog" color="teal.700" fontWeight="bold">
          Blog
        </Link>
        <Link href="/beneficios" color="teal.700" fontWeight="bold">
          Benefícios
        </Link>
        <Link href="/contato" color="teal.700" fontWeight="bold">
          Contato
        </Link>
        <Link
          href="https://instagram.com"
          color="teal.700"
          fontWeight="bold"
          target="_blank"
          rel="noopener noreferrer"
        >
          Instagram
        </Link>
      </HStack>
      <Text mt={4} color="gray.500" fontSize="sm">
        © {new Date().getFullYear()} Pet Services. Todos os direitos reservados.
      </Text>
    </Box>
  );
}
