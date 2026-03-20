import { Box, Text } from "@chakra-ui/react";
import type { ReactNode } from "react";

type EmptyStateProps = {
  message?: string;
  children?: ReactNode;
};

export default function EmptyState({
  message = "Nenhum item encontrado.",
  children,
}: EmptyStateProps) {
  return (
    <Box
      borderRadius={{ base: "2xl", md: "3xl" }}
      borderWidth="1px"
      borderColor="gray.200"
      bg="white"
      p={{ base: 4, md: 6 }}
    >
      <Text fontSize={{ base: "xs", sm: "sm" }} color="gray.600">
        {message}
      </Text>
      {children}
    </Box>
  );
}
