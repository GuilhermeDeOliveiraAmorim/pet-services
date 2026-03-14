"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import {
  Button,
  Flex,
  HStack,
  Link as ChakraLink,
  Text,
} from "@chakra-ui/react";

import { useAuthLogout, useAuthSession, useUserProfile } from "@/application";
import { UserTypes } from "@/domain";

import Image from "next/image";

const baseNavItems = [
  { label: "Encontre Serviços", href: "/services" },
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
  const router = useRouter();
  const pathname = usePathname();
  const { session, isAuthenticated, clearSession } = useAuthSession();
  const { data: profileData } = useUserProfile({ enabled: isAuthenticated });
  const { mutateAsync: logout, isPending } = useAuthLogout();
  const registerHref = pathname?.startsWith("/partner")
    ? "/register?user_type=provider"
    : "/register";
  const isActive = (href: string) =>
    pathname === href || pathname?.startsWith(`${href}/`);
  const isOwnerUser = profileData?.user?.userType === UserTypes.Owner;
  const isProviderUser = profileData?.user?.userType === UserTypes.Provider;
  const navItems = [
    ...baseNavItems,
    ...(isAuthenticated ? [{ label: "Perfil", href: "/profile" }] : []),
    ...(isOwnerUser ? [{ label: "Owner", href: "/owner" }] : []),
    ...(isProviderUser ? [{ label: "Provider", href: "/provider" }] : []),
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

  return (
    <Flex
      as="header"
      align="center"
      justify="space-between"
      py={{ base: "2", md: "3" }}
      className={className}
    >
      <ChakraLink
        as={Link}
        href="/"
        display="inline-flex"
        alignItems="center"
        gap="2"
        _hover={{ textDecoration: "none" }}
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

      {showLinks ? (
        <HStack
          display={{ base: "none", lg: "flex" }}
          gap="6"
          fontSize="sm"
          fontWeight="medium"
          color="gray.600"
        >
          {navItems.map((item) => (
            <ChakraLink
              key={item.href}
              as={Link}
              href={item.href}
              color={isActive(item.href) ? "teal.600" : "gray.600"}
              fontWeight={isActive(item.href) ? "semibold" : "medium"}
              _hover={{
                textDecoration: "none",
                color: "gray.800",
              }}
            >
              {item.label}
            </ChakraLink>
          ))}
        </HStack>
      ) : null}

      {showActions ? (
        <HStack gap="3">
          {isAuthenticated ? (
            <Button
              onClick={handleLogout}
              disabled={isPending}
              size="sm"
              variant="outline"
              borderRadius="full"
              borderColor="gray.300"
              color="gray.700"
              _hover={{ bg: "gray.50" }}
            >
              {isPending ? "Saindo..." : "Sair"}
            </Button>
          ) : (
            <>
              <ChakraLink
                as={Link}
                href="/login"
                display={{ base: "none", lg: "inline-flex" }}
                alignItems="center"
                justifyContent="center"
                h="9"
                px="4"
                fontSize="sm"
                fontWeight="medium"
                color="gray.700"
                borderWidth="1px"
                borderColor="gray.200"
                borderRadius="full"
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
                px="5"
                fontSize="sm"
                fontWeight="semibold"
                borderRadius="full"
                bgGradient="linear(to-r, teal.400, cyan.400)"
                color="white"
                _hover={{ opacity: 0.92, textDecoration: "none" }}
              >
                Cadastre-se
              </ChakraLink>
            </>
          )}
        </HStack>
      ) : null}
    </Flex>
  );
}
