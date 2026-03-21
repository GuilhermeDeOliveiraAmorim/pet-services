"use client";
import { Box, HStack, VStack, Text, Button } from "@chakra-ui/react";
import { FaBath, FaDog, FaHotel, FaWalking } from "react-icons/fa";

const actions = [
  { icon: <FaBath size={28} />, label: "Banho" },
  { icon: <FaDog size={28} />, label: "Tosa" },
  { icon: <FaWalking size={28} />, label: "Passeio" },
  { icon: <FaHotel size={28} />, label: "Hotel" },
];

export default function QuickActionsBar() {
  return (
    <Box as="section" py={6} bg="gray.50">
      <HStack gap={6} justify="center">
        {actions.map((action) => (
          <VStack key={action.label} gap={2}>
            <Button
              size="lg"
              variant="ghost"
              colorScheme="teal"
              borderRadius="full"
              display="flex"
              alignItems="center"
              justifyContent="center"
              p={0}
            >
              {action.icon}
            </Button>
            <Text fontSize="sm" color="gray.700">
              {action.label}
            </Text>
          </VStack>
        ))}
      </HStack>
    </Box>
  );
}
