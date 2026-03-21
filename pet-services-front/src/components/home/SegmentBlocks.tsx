import { Box, SimpleGrid, Button, Heading, VStack } from "@chakra-ui/react";

const segments = [
  { label: "Adoção", color: "pink.400", action: "Ver pets" },
  { label: "Serviços", color: "teal.400", action: "Ver opções" },
  { label: "Marketplace", color: "yellow.400", action: "Ver lojas" },
  { label: "Veterinário", color: "green.400", action: "Agendar" },
];

export default function SegmentBlocks() {
  return (
    <Box as="section" py={10}>
      <SimpleGrid
        columns={{ base: 1, sm: 2, md: 4 }}
        gap={6}
        maxW="6xl"
        mx="auto"
      >
        {segments.map((seg) => (
          <VStack
            key={seg.label}
            bg={seg.color}
            borderRadius="2xl"
            p={8}
            gap={4}
            align="center"
            boxShadow="md"
          >
            <Heading as="h2" size="md" color="white">
              {seg.label}
            </Heading>
            <Button colorScheme="whiteAlpha" variant="outline" size="md">
              {seg.action}
            </Button>
          </VStack>
        ))}
      </SimpleGrid>
    </Box>
  );
}
