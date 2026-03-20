import { Box, Text } from "@chakra-ui/react";
import type { ReactNode } from "react";

type ErrorStateProps = {
  message?: string;
  children?: ReactNode;
};

export default function ErrorState({
  message = "Ocorreu um erro. Tente novamente.",
  children,
}: ErrorStateProps) {
  return (
    <Box
      borderRadius={{ base: "2xl", md: "3xl" }}
      borderWidth="1px"
      borderColor="red.200"
      bg="red.50"
      p={{ base: 4, md: 6 }}
    >
      <Text fontSize={{ base: "xs", sm: "sm" }} color="red.700">
        {message}
      </Text>
      {children}
    </Box>
  );
}
