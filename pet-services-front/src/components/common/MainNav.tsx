"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import {
  Box,
  Button,
  Flex,
  HStack,
  Link as ChakraLink,
  Text,
  VStack,
} from "@chakra-ui/react";

import {
  useAuthLogout,
  useAuthSession,
  useGuardianStatus,
  useUserProfile,
} from "@/application";
import { UserTypes } from "@/domain";

import Image from "next/image";

const baseNavItems = [
  { label: "Encontre Serviços", href: "/services" },
  { label: "Adote um Pet", href: "/adoption" },
  { label: "Seja um Parceiro", href: "/partner" },
];

type MainNavProps = {
  className?: string;
  showLinks?: boolean;
  showActions?: boolean;
};

export default function MainNav({
  className = "",
  showLinks = true,
  showActions = true,
}: MainNavProps) {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const router = useRouter();
  const pathname = usePathname();
  const { session, isAuthenticated, isHydrated, clearSession } =
    useAuthSession();
  const { data: profileData } = useUserProfile({ enabled: isAuthenticated });
  const { isApprovedGuardian } = useGuardianStatus({
    enabled: isAuthenticated,
  });
  const { mutateAsync: logout, isPending } = useAuthLogout();
  const registerHref = pathname?.startsWith("/partner")
    ? "/register?user_type=provider"
    : "/register";
  const isActive = (href: string) => {
    if (!pathname) {
      return false;
    }

    if (href === "/adoption") {
      return (
        pathname === "/adoption" ||
        (/^\/adoption\/[^/]+$/.test(pathname) &&
          !pathname.startsWith("/adoption/applications"))
      );
    }

    return pathname === href || pathname.startsWith(`${href}/`);
  };
  const isOwnerUser = profileData?.user?.userType === UserTypes.Owner;
  const isProviderUser = profileData?.user?.userType === UserTypes.Provider;
  const canAccessRequests = isOwnerUser || isProviderUser;
  const navItems = [
    ...baseNavItems,
    ...(isAuthenticated
      ? [{ label: "Perfil de guardião", href: "/adoption/guardian" }]
      : []),
    ...(isAuthenticated
      ? [{ label: "Minhas candidaturas", href: "/adoption/applications" }]
      : []),
    ...(isApprovedGuardian
      ? [{ label: "Meus anúncios", href: "/adoption/listings/me" }]
      : []),
    ...(isApprovedGuardian
      ? [{ label: "Novo anúncio", href: "/adoption/listings/new" }]
      : []),
    ...(isAuthenticated ? [{ label: "Perfil", href: "/profile" }] : []),
    ...(isAuthenticated ? [{ label: "Configurações", href: "/settings" }] : []),
    ...(isOwnerUser ? [{ label: "Owner", href: "/owner" }] : []),
    ...(isProviderUser ? [{ label: "Provider", href: "/provider" }] : []),
    ...(canAccessRequests
      ? [{ label: "Minhas solicitações", href: "/requests" }]
      : []),
  ];

  const handleLogout = async () => {
    try {
      if (session) {
        await logout({ revokeAll: true });
      }
    } finally {
      clearSession();
      router.replace("/login");
    }
  };

  useEffect(() => {
    setIsMobileMenuOpen(false);
  }, [pathname]);

  return (
    <Box
      as="header"
      className={className}
      bg="white"
      borderWidth="1px"
      borderColor="gray.200"
      borderRadius="2xl"
      px={{ base: 4, md: 5 }}
      py={{ base: 3, md: 4 }}
      boxShadow="sm"
    >
      <Flex align="center" justify="space-between" gap={4}>
        <ChakraLink
          as={Link}
          href="/"
          display="inline-flex"
          alignItems="center"
          gap="2"
          _hover={{ textDecoration: "none" }}
          _focus={{ outline: "none", boxShadow: "none" }}
          _focusVisible={{ outline: "none", boxShadow: "none" }}
        >
          <Image
            src="/pawIcon.svg"
            alt="Pet Services Logo"
            width={32}
            height={32}
          />
          <Text fontSize="lg" fontWeight="semibold" color="gray.900">
            Pet Services
          </Text>
        </ChakraLink>

        {showActions ? (
          <HStack
            gap="3"
            align="center"
            display={{ base: "none", lg: "flex" }}
            flexShrink={0}
          >
            {!isHydrated ? (
              <Box w="10" h="9" />
            ) : isAuthenticated ? (
              <Button
                onClick={handleLogout}
                disabled={isPending}
                size="sm"
                h="9"
                variant="outline"
                borderRadius="full"
                borderColor="gray.300"
                color="gray.700"
                _hover={{ bg: "gray.50" }}
                _focus={{ outline: "none", boxShadow: "none" }}
                _focusVisible={{ outline: "none", boxShadow: "none" }}
                _active={{ boxShadow: "none" }}
              >
                {isPending ? "Saindo..." : "Sair"}
              </Button>
            ) : (
              <>
                <ChakraLink
                  as={Link}
                  href="/login"
                  _hover={{ textDecoration: "none" }}
                >
                  <Button
                    size="sm"
                    h="9"
                    px="4"
                    fontSize="sm"
                    fontWeight="medium"
                    variant="outline"
                    borderRadius="full"
                    borderColor="gray.200"
                    color="gray.700"
                    _hover={{ bg: "gray.50" }}
                    _focus={{ outline: "none", boxShadow: "none" }}
                    _focusVisible={{ outline: "none", boxShadow: "none" }}
                  >
                    Entrar
                  </Button>
                </ChakraLink>
                <ChakraLink
                  as={Link}
                  href={registerHref}
                  _hover={{ textDecoration: "none" }}
                >
                  <Button
                    size="sm"
                    h="9"
                    px="5"
                    fontSize="sm"
                    fontWeight="semibold"
                    borderRadius="full"
                    bgGradient="linear(to-r, teal.400, cyan.400)"
                    color="white"
                    _hover={{ opacity: 0.92 }}
                    _focus={{ outline: "none", boxShadow: "none" }}
                    _focusVisible={{ outline: "none", boxShadow: "none" }}
                  >
                    Cadastre-se
                  </Button>
                </ChakraLink>
              </>
            )}
          </HStack>
        ) : null}

        {showLinks || showActions ? (
          <Button
            display={{ base: "inline-flex", lg: "none" }}
            onClick={() => setIsMobileMenuOpen((open) => !open)}
            size="sm"
            variant="outline"
            borderRadius="full"
            borderColor="gray.300"
            color="gray.700"
            _hover={{ bg: "gray.50" }}
          >
            {isMobileMenuOpen ? "Fechar" : "Menu"}
          </Button>
        ) : null}
      </Flex>

      {showLinks ? (
        <Flex
          display={{ base: "none", lg: "flex" }}
          mt={4}
          pt={4}
          borderTopWidth="1px"
          borderTopColor="gray.100"
          gap={2}
          wrap="wrap"
        >
          {navItems.map((item) => (
            <ChakraLink
              key={item.href}
              as={Link}
              href={item.href}
              px={3}
              py={2}
              borderRadius="full"
              borderWidth="1px"
              borderColor={isActive(item.href) ? "teal.200" : "gray.200"}
              bg={isActive(item.href) ? "teal.50" : "white"}
              color={isActive(item.href) ? "teal.700" : "gray.700"}
              fontSize="sm"
              fontWeight={isActive(item.href) ? "semibold" : "medium"}
              whiteSpace="nowrap"
              _hover={{
                textDecoration: "none",
                color: "teal.700",
                borderColor: "teal.200",
                bg: "teal.50",
              }}
              _focus={{ outline: "none", boxShadow: "none" }}
              _focusVisible={{ outline: "none", boxShadow: "none" }}
            >
              {item.label}
            </ChakraLink>
          ))}
        </Flex>
      ) : null}

      {(showLinks || showActions) && isMobileMenuOpen ? (
        <VStack
          display={{ base: "flex", lg: "none" }}
          align="stretch"
          gap={2}
          pb={3}
        >
          {showLinks
            ? navItems.map((item) => (
                <ChakraLink
                  key={`mobile-${item.href}`}
                  as={Link}
                  href={item.href}
                  px={3}
                  py={2}
                  borderRadius="lg"
                  color={isActive(item.href) ? "teal.700" : "gray.700"}
                  bg={isActive(item.href) ? "teal.50" : "transparent"}
                  fontWeight={isActive(item.href) ? "semibold" : "medium"}
                >
                  {item.label}
                </ChakraLink>
              ))
            : null}

          {showActions ? (
            !isHydrated ? null : isAuthenticated ? (
              <Button
                onClick={handleLogout}
                disabled={isPending}
                size="sm"
                variant="outline"
                borderRadius="lg"
                borderColor="gray.300"
                color="gray.700"
                justifyContent="start"
              >
                {isPending ? "Saindo..." : "Sair"}
              </Button>
            ) : (
              <VStack align="stretch" gap={2}>
                <ChakraLink
                  as={Link}
                  href="/login"
                  display="inline-flex"
                  alignItems="center"
                  justifyContent="center"
                  h="9"
                  fontSize="sm"
                  fontWeight="medium"
                  color="gray.700"
                  borderWidth="1px"
                  borderColor="gray.200"
                  borderRadius="lg"
                  _hover={{ textDecoration: "none", bg: "gray.50" }}
                >
                  Entrar
                </ChakraLink>
                <ChakraLink
                  as={Link}
                  href={registerHref}
                  display="inline-flex"
                  alignItems="center"
                  justifyContent="center"
                  h="9"
                  fontSize="sm"
                  fontWeight="semibold"
                  borderRadius="lg"
                  bgGradient="linear(to-r, teal.400, cyan.400)"
                  color="white"
                  _hover={{ opacity: 0.92, textDecoration: "none" }}
                >
                  Cadastre-se
                </ChakraLink>
              </VStack>
            )
          ) : null}
        </VStack>
      ) : null}
    </Box>
  );
}
