import {
  Box,
  SimpleGrid,
  VStack,
  Heading,
  Button,
  Text,
} from "@chakra-ui/react";

const calls = [
  {
    title: "Quero oferecer serviços",
    desc: "Cadastre-se como prestador e alcance mais clientes!",
  },
  { title: "É uma ONG?", desc: "Divulgue pets para adoção e receba apoio." },
  { title: "Seja voluntário", desc: "Ajude ONGs e pets em sua região." },
];

export default function PartnerCallToAction() {
  return (
    <Box as="section" py={10} bg="teal.50">
      <SimpleGrid columns={{ base: 1, md: 3 }} gap={6} maxW="6xl" mx="auto">
        {calls.map((call) => (
          <VStack
            key={call.title}
            bg="white"
            borderRadius="2xl"
            boxShadow="md"
            p={6}
            align="center"
          >
            <Heading as="h3" size="md" color="teal.700">
              {call.title}
            </Heading>
            <Text color="gray.700">{call.desc}</Text>
            <Button colorScheme="teal" variant="solid">
              Saiba mais
            </Button>
          </VStack>
        ))}
      </SimpleGrid>
    </Box>
  );
}
