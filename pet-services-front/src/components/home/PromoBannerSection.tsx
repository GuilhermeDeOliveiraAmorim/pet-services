import { Box, Image } from "@chakra-ui/react";

export default function PromoBannerSection() {
  return (
    <Box as="section" py={6} maxW="6xl" mx="auto">
      <Image
        src="/img/banner-promo.jpg"
        alt="Promoção"
        borderRadius="xl"
        w="full"
        h="120px"
        objectFit="cover"
      />
    </Box>
  );
}
