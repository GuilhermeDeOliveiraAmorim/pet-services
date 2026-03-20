import { Flex, Spinner, Text } from "@chakra-ui/react";

type LoadingStateProps = {
  message?: string;
};

export default function LoadingState({
  message = "Carregando...",
}: LoadingStateProps) {
  return (
    <Flex
      borderRadius={{ base: "2xl", md: "3xl" }}
      borderWidth="1px"
      borderColor="gray.200"
      bg="white"
      py={12}
      justify="center"
      align="center"
      gap={3}
      px={4}
    >
      <Spinner color="teal.500" size="sm" />
      <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
        {message}
      </Text>
    </Flex>
  );
}
