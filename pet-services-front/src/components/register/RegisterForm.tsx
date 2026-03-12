"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Checkbox, Text, VStack, chakra } from "@chakra-ui/react";
import { useReferenceCountries, useUserRegister } from "@/application";
import { UserTypes } from "@/domain";
import RegisterAccountFields from "./RegisterAccountFields";
import RegisterFormCard from "./RegisterFormCard";
import RegisterFormFooter from "./RegisterFormFooter";
import RegisterPhoneFields from "./RegisterPhoneFields";
import RegisterSubmitRow from "./RegisterSubmitRow";

const defaultPhoneCountryCode = "55";
const defaultCountryCode = "BR";
const redirectDelaySeconds = 4;

export default function RegisterForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { mutateAsync, isPending, error, isSuccess } = useUserRegister();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [selectedDialCodeValue, setSelectedDialCodeValue] = useState("");
  const [areaCode, setAreaCode] = useState("");
  const [phoneNumber, setPhoneNumber] = useState("");
  const [redirectSecondsLeft, setRedirectSecondsLeft] = useState<number | null>(
    null,
  );
  const [successFeedbackMessage, setSuccessFeedbackMessage] = useState("");
  const [isPartnerAccount, setIsPartnerAccount] = useState(
    searchParams.get("user_type") === UserTypes.Provider,
  );

  useEffect(() => {
    if (redirectSecondsLeft === null || redirectSecondsLeft <= 0) {
      return;
    }

    const timerId = window.setInterval(() => {
      setRedirectSecondsLeft((previous) => {
        if (previous === null) {
          return null;
        }

        if (previous <= 1) {
          return 0;
        }

        return previous - 1;
      });
    }, 1000);

    return () => window.clearInterval(timerId);
  }, [redirectSecondsLeft]);

  useEffect(() => {
    if (redirectSecondsLeft !== 0) {
      return;
    }

    router.replace("/login");
  }, [redirectSecondsLeft, router]);

  const { data: countriesData } = useReferenceCountries();

  const countries = useMemo(
    () => countriesData?.countries ?? [],
    [countriesData],
  );

  const dialCodeOptions = useMemo(() => {
    const seen = new Set<string>();
    return countries.reduce<Array<{ value: string; label: string }>>(
      (acc, country) => {
        const dialCode = country.dialCode
          ? String(country.dialCode).trim()
          : "";
        const digits = dialCode.replace(/\D/g, "");

        if (!digits) {
          return acc;
        }

        const label = `${country.name} (+${digits})`;
        const value = `${digits}:${country.code}`;
        const key = value;
        if (!seen.has(key)) {
          seen.add(key);
          acc.push({ value, label });
        }
        return acc;
      },
      [],
    );
  }, [countries]);

  const defaultDialCodeValue = useMemo(() => {
    if (!dialCodeOptions.length) {
      return "";
    }

    const brOption = dialCodeOptions.find((option) =>
      option.value.endsWith(`:${defaultCountryCode}`),
    );

    return brOption?.value ?? dialCodeOptions[0]?.value ?? "";
  }, [dialCodeOptions]);

  const currentDialCodeValue = useMemo(() => {
    if (!selectedDialCodeValue) {
      return defaultDialCodeValue;
    }

    const exists = dialCodeOptions.some(
      (option) => option.value === selectedDialCodeValue,
    );

    return exists ? selectedDialCodeValue : defaultDialCodeValue;
  }, [defaultDialCodeValue, dialCodeOptions, selectedDialCodeValue]);

  const selectedDialCodeDisplay = useMemo(() => {
    if (!currentDialCodeValue) {
      return "";
    }

    const [digits, countryCodeValue] = currentDialCodeValue.split(":");
    const country = countries.find((item) => item.code === countryCodeValue);
    const flag = country?.flag ?? "";
    const displayDigits = digits ?? "";

    if (!displayDigits) {
      return "";
    }

    return flag ? `${flag} +${displayDigits}` : `+${displayDigits}`;
  }, [countries, currentDialCodeValue]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const phoneDigits = phoneNumber.replace(/\D/g, "");
    const [countryCodeDigits] = currentDialCodeValue.split(":");

    const response = await mutateAsync({
      name,
      userType: isPartnerAccount ? UserTypes.Provider : UserTypes.Owner,
      login: {
        email,
        password,
      },
      phone: {
        countryCode: countryCodeDigits || defaultPhoneCountryCode,
        areaCode,
        number: phoneDigits,
      },
    });

    setSuccessFeedbackMessage(
      response.detail ||
        "Cadastro realizado! Enviamos um e-mail de confirmação para sua conta.",
    );
    setRedirectSecondsLeft(redirectDelaySeconds);
  };

  const isRedirectingAfterSuccess = redirectSecondsLeft !== null;
  const successMessageWithCountdown = isRedirectingAfterSuccess
    ? `${successFeedbackMessage} Você será redirecionado para o login em ${redirectSecondsLeft}s.`
    : successFeedbackMessage;

  return (
    <RegisterFormCard>
      <chakra.form onSubmit={handleSubmit}>
        <VStack gap={7} align="stretch">
          <Text fontSize="sm" color="gray.500">
            Etapa 1 de 2: crie sua conta com dados básicos. O perfil completo é
            finalizado após o login.
          </Text>

          <VStack mt={1} gap={7} align="stretch">
            <RegisterAccountFields
              name={name}
              onNameChange={setName}
              email={email}
              onEmailChange={setEmail}
              password={password}
              onPasswordChange={setPassword}
            />

            <RegisterPhoneFields
              countryCode={currentDialCodeValue}
              onCountryCodeChange={setSelectedDialCodeValue}
              countryCodeDisplayValue={selectedDialCodeDisplay}
              areaCode={areaCode}
              onAreaCodeChange={setAreaCode}
              phoneNumber={phoneNumber}
              onPhoneNumberChange={setPhoneNumber}
              dialCodeOptions={dialCodeOptions}
            />

            <VStack align="stretch" gap={1}>
              <Checkbox.Root
                checked={isPartnerAccount}
                onCheckedChange={() => setIsPartnerAccount((prev) => !prev)}
              >
                <Checkbox.HiddenInput />
                <Checkbox.Control />
                <Checkbox.Label>
                  Quero me cadastrar como parceiro
                </Checkbox.Label>
              </Checkbox.Root>
            </VStack>

            <RegisterSubmitRow
              isPending={isPending}
              isRedirecting={isRedirectingAfterSuccess}
            />

            <RegisterFormFooter
              error={error}
              isSuccess={isSuccess}
              successMessage={successMessageWithCountdown}
            />
          </VStack>
        </VStack>
      </chakra.form>
    </RegisterFormCard>
  );
}
