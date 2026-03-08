import type { ReactNode } from "react";
import { Box, Flex } from "@chakra-ui/react";
import RegisterFormHeader from "./RegisterFormHeader";

type RegisterFormCardProps = {
  children: ReactNode;
};

export default function RegisterFormCard({ children }: RegisterFormCardProps) {
  return (
    <Flex align="center" justify="center">
      <Box
        w="full"
        maxW="2xl"
        borderRadius="3xl"
        bg="white"
        p={{ base: 6, md: 10 }}
        borderWidth="1px"
        borderColor="gray.100"
        shadow="sm"
      >
        <RegisterFormHeader />
        {children}
      </Box>
    </Flex>
  );
}
