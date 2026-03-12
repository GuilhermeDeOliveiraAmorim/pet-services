import { Box, HStack, Text } from "@chakra-ui/react";
import {
  CheckCircle2,
  MapPin,
  ShieldCheck,
  Sparkles,
  Users,
} from "lucide-react";

const benefits = [
  {
    label: "Mais visibilidade local",
    icon: MapPin,
    iconBg: "cyan.100",
    iconColor: "cyan.600",
  },
  {
    label: "Conquiste novos clientes",
    icon: Users,
    iconBg: "teal.100",
    iconColor: "teal.600",
  },
  {
    label: "Perfil profissional verificado",
    icon: ShieldCheck,
    iconBg: "green.100",
    iconColor: "green.600",
  },
  {
    label: "Catálogo de serviços",
    icon: Sparkles,
    iconBg: "orange.100",
    iconColor: "orange.600",
  },
  {
    label: "Onboarding simples",
    icon: CheckCircle2,
    iconBg: "purple.100",
    iconColor: "purple.600",
  },
];

export default function PartnerBenefitsSection() {
  return (
    <HStack
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
      {benefits.map((item) => {
        const ItemIcon = item.icon;

        return (
          <Box
            key={item.label}
            display="inline-flex"
            alignItems="center"
            borderRadius="full"
            bg="gray.100"
            borderWidth="1px"
            borderColor="gray.200"
            px="3"
            py="2"
            gap="2"
            flex="0 0 auto"
            minW="max-content"
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
            <Text
              fontSize="xs"
              fontWeight="medium"
              color="gray.700"
              whiteSpace="nowrap"
            >
              {item.label}
            </Text>
          </Box>
        );
      })}
    </HStack>
  );
}
