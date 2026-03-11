"use client";

import { type ChangeEvent, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import {
  Box,
  Button,
  Checkbox,
  chakra,
  Flex,
  Grid,
  Heading,
  HStack,
  IconButton,
  Input,
  Link as ChakraLink,
  Text,
  VStack,
} from "@chakra-ui/react";
import { Eye, EyeOff } from "lucide-react";

import {
  useAuthLogin,
  useAuthResendVerificationEmail,
  useAuthSession,
} from "@/application";
import { UserTypes } from "@/domain";
import { getApiErrorMessage, isUnverifiedEmailError } from "@/lib/api-error";
import MainNav from "@/components/common/MainNav";
import PageWrapper from "@/components/common/PageWrapper";

export default function LoginPage() {
  const router = useRouter();
  const { setSession } = useAuthSession();
  const { mutateAsync, isPending, error } = useAuthLogin();
  const {
    mutateAsync: resendVerification,
    isPending: isResendingVerification,
    data: resendVerificationResult,
  } = useAuthResendVerificationEmail();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [rememberMe, setRememberMe] = useState(false);
  const [emailError, setEmailError] = useState("");
  const [passwordError, setPasswordError] = useState("");

  const validateEmail = (value: string) => {
    const parsed = value.trim();
    if (!parsed) {
      return "Informe o email";
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(parsed)) {
      return "Email inválido";
    }

    return "";
  };

  const validatePassword = (value: string) => {
    if (!value.trim()) {
      return "Informe a senha";
    }

    return "";
  };

  const feedback = useMemo(() => {
    if (!error) {
      return { message: "", canResend: false };
    }

    return {
      message: getApiErrorMessage(
        error,
        "Não foi possível fazer login. Verifique suas credenciais.",
      ),
      canResend: isUnverifiedEmailError(error),
    };
  }, [error]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const nextEmailError = validateEmail(email);
    const nextPasswordError = validatePassword(password);

    setEmailError(nextEmailError);
    setPasswordError(nextPasswordError);

    if (nextEmailError || nextPasswordError) {
      return;
    }

    const response = await mutateAsync({ email, password });
    const expiresAt = Date.now() + response.expiresIn * 1000;

    setSession({
      accessToken: response.accessToken,
      refreshToken: response.refreshToken,
      expiresAt,
    });

    if (!response.user.profileComplete) {
      router.replace("/profile");
      return;
    }

    if (response.user.userType === UserTypes.Owner) {
      router.replace("/owner");
      return;
    }

    if (response.user.userType === UserTypes.Provider) {
      router.replace("/provider");
      return;
    }

    router.replace("/");
  };

  return (
    <PageWrapper gap={16}>
      <MainNav showLinks={false} showActions={false} />
      <Grid
        w="full"
        gap={10}
        templateColumns={{ base: "1fr", lg: "1.1fr 0.9fr" }}
      >
        <Box
          display={{ base: "none", lg: "flex" }}
          borderRadius="3xl"
          overflow="hidden"
          position="relative"
          minH="640px"
          bgImage="url('/banner.png')"
          bgSize="cover"
          backgroundPosition="center"
        >
          <Box position="absolute" inset={0} bg="blackAlpha.650" />
          <VStack
            position="relative"
            zIndex={1}
            align="stretch"
            justify="space-between"
            p={10}
            w="full"
          >
            <Box>
              <Heading size="4xl" lineHeight="1.1" color="white">
                Cuidado Profissional para quem você Ama
              </Heading>
              <Text mt={4} color="whiteAlpha.900" fontSize="md" maxW="lg">
                Conecte-se com os melhores veterinários, pet shops e cuidadores
                da sua região.
              </Text>
            </Box>

            <HStack
              borderRadius="2xl"
              p={5}
              bg="whiteAlpha.180"
              backdropFilter="blur(8px)"
              justify="space-between"
            >
              <Text color="white" fontWeight="medium">
                Login seguro
              </Text>
              <Text color="green.200" fontWeight="bold">
                24/7
              </Text>
            </HStack>
          </VStack>
        </Box>

        <Flex align="center" justify="center">
          <chakra.form
            onSubmit={handleSubmit}
            w="full"
            borderRadius="3xl"
            bg="white"
            p={{ base: 6, md: 8 }}
            boxShadow="0 30px 80px rgba(35, 49, 82, 0.12)"
          >
            <Flex align="center" justify="space-between" mb={6}>
              <Box>
                <Heading size="2xl" color="gray.900">
                  Login
                </Heading>
                <Text mt={2} fontSize="sm" color="gray.500">
                  Acesse sua conta para continuar.
                </Text>
              </Box>
              <Flex
                h="12"
                w="12"
                align="center"
                justify="center"
                borderRadius="2xl"
                bgGradient="linear(to-br, teal.100, cyan.100)"
                color="gray.700"
                fontSize="sm"
                fontWeight="semibold"
              >
                24/7
              </Flex>
            </Flex>

            <VStack gap={5} align="stretch">
              <Box>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  Email
                </Text>
                <Input
                  type="email"
                  value={email}
                  onChange={(event: ChangeEvent<HTMLInputElement>) => {
                    setEmail(event.target.value);
                    if (emailError) {
                      setEmailError("");
                    }
                  }}
                  onBlur={() => setEmailError(validateEmail(email))}
                  placeholder="voce@email.com"
                  autoComplete="email"
                  h="11"
                  borderRadius="xl"
                  bg="gray.50"
                  borderColor={emailError ? "red.300" : "gray.200"}
                  focusRingColor={emailError ? "red.200" : "teal.200"}
                />
                {emailError ? (
                  <Text mt={1.5} fontSize="xs" color="red.500">
                    {emailError}
                  </Text>
                ) : null}
              </Box>

              <Box>
                <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
                  Senha
                </Text>
                <Box position="relative">
                  <Input
                    type={showPassword ? "text" : "password"}
                    value={password}
                    onChange={(event: ChangeEvent<HTMLInputElement>) => {
                      setPassword(event.target.value);
                      if (passwordError) {
                        setPasswordError("");
                      }
                    }}
                    onBlur={() => setPasswordError(validatePassword(password))}
                    placeholder="********"
                    autoComplete="current-password"
                    h="11"
                    borderRadius="xl"
                    bg="gray.50"
                    borderColor={passwordError ? "red.300" : "gray.200"}
                    focusRingColor={passwordError ? "red.200" : "teal.200"}
                    pe="12"
                  />
                  <IconButton
                    type="button"
                    onClick={() => setShowPassword((prev) => !prev)}
                    aria-label={
                      showPassword ? "Ocultar senha" : "Mostrar senha"
                    }
                    variant="ghost"
                    size="sm"
                    color="gray.500"
                    position="absolute"
                    right="2"
                    top="50%"
                    transform="translateY(-50%)"
                  >
                    {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                  </IconButton>
                </Box>
                {passwordError ? (
                  <Text mt={1.5} fontSize="xs" color="red.500">
                    {passwordError}
                  </Text>
                ) : null}
              </Box>

              <Flex
                align="center"
                justify="space-between"
                fontSize="xs"
                color="gray.500"
              >
                <Checkbox.Root
                  checked={rememberMe}
                  onCheckedChange={() => setRememberMe((prev) => !prev)}
                >
                  <Checkbox.HiddenInput />
                  <Checkbox.Control />
                  <Checkbox.Label>Lembrar de mim</Checkbox.Label>
                </Checkbox.Root>

                <ChakraLink
                  as={Link}
                  href="/forgot-password"
                  color="teal.500"
                  fontWeight="medium"
                >
                  Esqueci minha senha
                </ChakraLink>
              </Flex>

              {error ? (
                <Box
                  borderRadius="2xl"
                  borderWidth="1px"
                  borderColor="red.200"
                  bg="red.50"
                  p={4}
                >
                  <Text fontSize="sm" color="red.600">
                    {feedback.message}
                  </Text>
                  {feedback.canResend ? (
                    <Button
                      mt={3}
                      size="xs"
                      type="button"
                      variant="outline"
                      colorPalette="red"
                      borderRadius="full"
                      disabled={!email || isResendingVerification}
                      onClick={() => resendVerification({ email })}
                    >
                      {isResendingVerification
                        ? "Reenviando..."
                        : "Reenviar email de verificação"}
                    </Button>
                  ) : null}
                </Box>
              ) : null}

              {resendVerificationResult?.message ? (
                <Box
                  borderRadius="2xl"
                  borderWidth="1px"
                  borderColor="green.200"
                  bg="green.50"
                  p={4}
                >
                  <Text fontSize="sm" color="green.700">
                    {resendVerificationResult.detail ??
                      resendVerificationResult.message}
                  </Text>
                </Box>
              ) : null}

              <Button
                type="submit"
                disabled={isPending}
                h="11"
                borderRadius="full"
                bgGradient="linear(to-r, green.400, teal.400)"
                color="white"
                _hover={{ opacity: 0.92 }}
                _disabled={{ opacity: 0.7, cursor: "not-allowed" }}
              >
                {isPending ? "Entrando..." : "Entrar"}
              </Button>

              <Text textAlign="center" fontSize="xs" color="gray.500">
                Ainda não tem conta?{" "}
                <ChakraLink
                  as={Link}
                  href="/register"
                  color="teal.500"
                  fontWeight="medium"
                >
                  Criar conta
                </ChakraLink>
              </Text>
            </VStack>
          </chakra.form>
        </Flex>
      </Grid>
    </PageWrapper>
  );
}
