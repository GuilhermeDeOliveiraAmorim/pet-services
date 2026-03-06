import { Box, Flex } from "@chakra-ui/react";
import type { ReactNode } from "react";

type PageWrapperProps = {
  children: ReactNode;
  gap?: number | string;
};

export default function PageWrapper({
  children,
  gap = 10,
}: PageWrapperProps) {
  return (
    <Box minH="100vh" bg="gray.50" color="gray.900">
      <Flex
        minH="100vh"
        w="full"
        maxW="6xl"
        mx="auto"
        direction="column"
        gap={gap}
        px={{ base: 6, lg: 8 }}
        py="10"
      >
        {children}
      </Flex>
    </Box>
  );
}
