import { Box, HStack, Link, Text } from "@chakra-ui/react";

import { serviceOptions } from "./service-options";

export default function ServicesSection() {
  return (
    <HStack
      id="services"
      gap="3"
      overflowX="auto"
      overflowY="hidden"
      flexWrap="nowrap"
      align="stretch"
      pb="2"
      scrollbarWidth="none"
      css={{
        "&::-webkit-scrollbar": {
          height: "0px",
        },
        "&:hover::-webkit-scrollbar": {
          height: "6px",
        },
        "&:hover::-webkit-scrollbar-thumb": {
          background: "rgba(148, 163, 184, 0.6)",
          borderRadius: "9999px",
        },
      }}
      _hover={{ scrollbarWidth: "thin" }}
    >
      {serviceOptions.map((item) => {
        const ItemIcon = item.icon;

        return (
          <Link
            key={item.label}
            href={item.href}
            display="inline-flex"
            alignItems="center"
            cursor="pointer"
            borderRadius="full"
            bg="gray.100"
            borderWidth="1px"
            borderColor="gray.200"
            px="3"
            py="2"
            gap="2"
            flex="0 0 auto"
            minW="max-content"
            transition="all 0.2s ease"
            _hover={{ bg: "gray.200", textDecoration: "none" }}
          >
            <Box
              display="inline-flex"
              alignItems="center"
              justifyContent="center"
              h="7"
              w="7"
              borderRadius="full"
              bg={item.iconBg}
              color={item.iconColor}
            >
              <ItemIcon size={16} />
            </Box>
            <Text fontSize="xs" fontWeight="medium" color="gray.700" whiteSpace="nowrap">
              {item.label}
            </Text>
          </Link>
        );
      })}
    </HStack>
  );
}