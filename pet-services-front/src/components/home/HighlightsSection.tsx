import {
  Box,
  Heading,
  SimpleGrid,
  VStack,
  Text,
  Image,
} from "@chakra-ui/react";

const highlights = [
  {
    title: "Pet para adoção",
    img: "/img/pet1.jpg",
    desc: "Encontre um novo amigo!",
  },
  {
    title: "Serviço popular",
    img: "/img/service1.jpg",
    desc: "Banho & Tosa com desconto",
  },
];

export default function HighlightsSection() {
  return (
    <Box as="section" py={10} maxW="6xl" mx="auto">
      <Heading as="h2" size="lg" mb={6} color="teal.700">
        Destaques
      </Heading>
      <SimpleGrid columns={{ base: 1, md: 2 }} gap={6}>
        {highlights.map((item) => (
          <VStack
            key={item.title}
            bg="white"
            borderRadius="2xl"
            boxShadow="md"
            p={6}
            align="start"
          >
            <Image
              src={item.img}
              alt={item.title}
              borderRadius="xl"
              boxSize="120px"
              objectFit="cover"
            />
            <Text fontWeight="bold" color="teal.800">
              {item.title}
            </Text>
            <Text color="gray.600">{item.desc}</Text>
          </VStack>
        ))}
      </SimpleGrid>
    </Box>
  );
}
