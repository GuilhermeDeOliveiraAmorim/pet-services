import { Box, Button, Heading, Text, VStack } from "@chakra-ui/react";

export default function ContactSection() {
  return (
    <VStack
      id="contact"
      borderRadius="3xl"
      bgGradient="linear(to-br, cyan.50, teal.50)"
      px="6"
      py="12"
      textAlign="center"
      gap="0"
    >
      <Text fontSize="xs" fontWeight="semibold" textTransform="uppercase" color="cyan.400">
        Contato
      </Text>
      <Heading mt="3" size="xl" color="gray.900">
        Tem alguma dúvida?
      </Heading>
      <Text mt="2" fontSize="sm" color="gray.600">
        Fale com a nossa equipe e receba atendimento rápido e humanizado.
      </Text>
      <Button
        mt="6"
        borderRadius="full"
        bgGradient="linear(to-r, teal.400, cyan.400)"
        px="6"
        py="3"
        fontSize="sm"
        fontWeight="semibold"
        color="white"
        boxShadow="lg"
        _hover={{ opacity: 0.9 }}
      >
        Falar com a equipe
      </Button>
    </VStack>
  );
}