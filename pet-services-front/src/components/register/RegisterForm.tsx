"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { Box, Button, HStack, VStack, chakra } from "@chakra-ui/react";
import { type UserType, UserTypes } from "@/domain";
import {
  useReferenceCities,
  useReferenceCountries,
  useReferenceStates,
  useUserRegister,
} from "@/application";
import RegisterAccountFields from "./RegisterAccountFields";
import RegisterAdditionalFields from "./RegisterAdditionalFields";
import RegisterAddressFields from "./RegisterAddressFields";
import RegisterFormCard from "./RegisterFormCard";
import RegisterFormFooter from "./RegisterFormFooter";
import RegisterPhoneFields from "./RegisterPhoneFields";
import RegisterSubmitRow from "./RegisterSubmitRow";

const defaultPhoneCountryCode = "55";
const defaultCountryCode = "BR";

export default function RegisterForm() {
  const router = useRouter();
  const { mutateAsync, isPending, error, isSuccess } = useUserRegister();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [userType, setUserType] = useState<UserType>(UserTypes.Owner);
  const [phoneCountryCode, setPhoneCountryCode] = useState(
    defaultPhoneCountryCode,
  );
  const [areaCode, setAreaCode] = useState("");
  const [phoneNumber, setPhoneNumber] = useState("");

  const [street, setStreet] = useState("");
  const [addressNumber, setAddressNumber] = useState("");
  const [neighborhood, setNeighborhood] = useState("");
  const [countryCode, setCountryCode] = useState(defaultCountryCode);
  const [stateId, setStateId] = useState<number | undefined>();
  const [cityId, setCityId] = useState<number | undefined>();
  const [zipCode, setZipCode] = useState("");
  const [complement, setComplement] = useState("");
  const [latitude, setLatitude] = useState("0");
  const [longitude, setLongitude] = useState("0");
  const [step, setStep] = useState("account");
  const [isGeocoding, setIsGeocoding] = useState(false);
  const [geocodeError, setGeocodeError] = useState<string | undefined>();

  const { data: countriesData } = useReferenceCountries();
  const { data: statesData } = useReferenceStates();
  const { data: citiesData } = useReferenceCities(
    { stateId },
    { enabled: Boolean(stateId) },
  );

  const countries = useMemo(
    () => countriesData?.countries ?? [],
    [countriesData],
  );
  const states = useMemo(() => statesData?.states ?? [], [statesData]);
  const cities = useMemo(() => citiesData?.cities ?? [], [citiesData]);

  const countryOptions = useMemo(() => {
    const map = new Map<string, { value: string; label: string }>();
    countries.forEach((country) => {
      const key = `${country.code}:${country.name}`;
      if (!map.has(key)) {
        map.set(key, {
          value: country.code,
          label: country.name,
        });
      }
    });
    return Array.from(map.values());
  }, [countries]);

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

  useEffect(() => {
    if (!dialCodeOptions.length) {
      return;
    }

    const brOption = dialCodeOptions.find((option) =>
      option.value.endsWith(`:${defaultCountryCode}`),
    );

    if (!brOption) {
      return;
    }

    const hasMatch = dialCodeOptions.some((option) =>
      option.value.startsWith(`${phoneCountryCode}:`),
    );

    if (!phoneCountryCode || !hasMatch) {
      const [digits] = brOption.value.split(":");
      setPhoneCountryCode(digits ?? defaultPhoneCountryCode);
    }
  }, [dialCodeOptions, phoneCountryCode]);

  const selectedDialCodeValue = useMemo(() => {
    if (!phoneCountryCode) {
      return "";
    }

    const match = dialCodeOptions.find((option) =>
      option.value.startsWith(`${phoneCountryCode}:`),
    );

    return match?.value ?? "";
  }, [dialCodeOptions, phoneCountryCode]);

  const selectedDialCodeDisplay = useMemo(() => {
    if (!selectedDialCodeValue) {
      return "";
    }

    const [digits, countryCodeValue] = selectedDialCodeValue.split(":");
    const country = countries.find((item) => item.code === countryCodeValue);
    const flag = country?.flag ?? "";
    const displayDigits = digits ?? "";

    if (!displayDigits) {
      return "";
    }

    return flag ? `${flag} +${displayDigits}` : `+${displayDigits}`;
  }, [countries, selectedDialCodeValue]);

  const geocodeAddress = useMemo(() => zipCode.trim(), [zipCode]);

  const canGeocode = Boolean(geocodeAddress);

  const handleGeocode = async () => {
    if (!geocodeAddress) {
      setGeocodeError("Informe o endereço completo ou o CEP.");
      return;
    }

    try {
      setIsGeocoding(true);
      setGeocodeError(undefined);

      const params = new URLSearchParams({ address: geocodeAddress });
      const response = await fetch(`/api/geocode?${params.toString()}`);
      const data = await response.json();

      if (!response.ok) {
        throw new Error(data?.error || "Não foi possível obter coordenadas");
      }

      const lat = Number.parseFloat(String(data.latitude ?? ""));
      const lng = Number.parseFloat(String(data.longitude ?? ""));

      if (!Number.isNaN(lat)) {
        setLatitude(lat.toFixed(2));
      }
      if (!Number.isNaN(lng)) {
        setLongitude(lng.toFixed(2));
      }
    } catch (err) {
      const message =
        err instanceof Error
          ? err.message
          : "Não foi possível obter coordenadas";
      setGeocodeError(message);
    } finally {
      setIsGeocoding(false);
    }
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const phoneDigits = phoneNumber.replace(/\D/g, "");

    await mutateAsync({
      name,
      userType,
      login: {
        email,
        password,
      },
      phone: {
        countryCode: phoneCountryCode,
        areaCode,
        number: phoneDigits,
      },
      address: {
        street,
        number: addressNumber,
        neighborhood,
        city: cities.find((item) => item.id === cityId)?.name ?? "",
        zipCode,
        state: states.find((item) => item.id === stateId)?.code ?? "",
        country:
          countries.find((item) => item.code === countryCode)?.name ?? "",
        complement,
        location: {
          latitude: Number(latitude),
          longitude: Number(longitude),
        },
      },
    });

    router.replace("/login");
  };

  return (
    <RegisterFormCard>
      <chakra.form onSubmit={handleSubmit}>
        <VStack gap={7} align="stretch">
          <HStack
            bg="gray.100"
            p={1}
            borderRadius="full"
            gap={1}
            align="stretch"
            borderWidth="1px"
            borderColor="gray.200"
          >
            <Button
              type="button"
              flex={1}
              borderRadius="full"
              variant={step === "account" ? "solid" : "ghost"}
              bg={step === "account" ? "white" : "transparent"}
              color={step === "account" ? "gray.900" : "gray.600"}
              boxShadow={step === "account" ? "sm" : "none"}
              onClick={() => setStep("account")}
            >
              Conta
            </Button>
            <Button
              type="button"
              flex={1}
              borderRadius="full"
              variant={step === "address" ? "solid" : "ghost"}
              bg={step === "address" ? "white" : "transparent"}
              color={step === "address" ? "gray.900" : "gray.600"}
              boxShadow={step === "address" ? "sm" : "none"}
              onClick={() => setStep("address")}
            >
              Endereço
            </Button>
          </HStack>

          {step === "account" ? (
            <VStack mt={1} gap={7} align="stretch">
              <RegisterAccountFields
                name={name}
                onNameChange={setName}
                userType={userType}
                onUserTypeChange={setUserType}
                email={email}
                onEmailChange={setEmail}
                password={password}
                onPasswordChange={setPassword}
              />

              <RegisterPhoneFields
                countryCode={selectedDialCodeValue}
                onCountryCodeChange={(value) => {
                  const [digits] = value.split(":");
                  setPhoneCountryCode(digits ?? "");
                }}
                countryCodeDisplayValue={selectedDialCodeDisplay}
                areaCode={areaCode}
                onAreaCodeChange={setAreaCode}
                phoneNumber={phoneNumber}
                onPhoneNumberChange={setPhoneNumber}
                dialCodeOptions={dialCodeOptions}
              />

              <Box display="flex" justifyContent="flex-end">
                <Button
                  type="button"
                  onClick={() => setStep("address")}
                  borderRadius="full"
                  borderWidth="1px"
                  borderColor="gray.200"
                  color="gray.700"
                  bg="white"
                  _hover={{ bg: "gray.50", borderColor: "green.300" }}
                >
                  Continuar
                </Button>
              </Box>
            </VStack>
          ) : null}

          {step === "address" ? (
            <VStack mt={1} gap={7} align="stretch">
              <RegisterAddressFields
                street={street}
                onStreetChange={setStreet}
                addressNumber={addressNumber}
                onAddressNumberChange={setAddressNumber}
                neighborhood={neighborhood}
                onNeighborhoodChange={setNeighborhood}
                zipCode={zipCode}
                onZipCodeChange={setZipCode}
                complement={complement}
                onComplementChange={setComplement}
                countryCode={countryCode}
                onCountryCodeChange={(value) => {
                  setCountryCode(value);
                  setStateId(undefined);
                  setCityId(undefined);
                }}
                stateId={stateId}
                onStateIdChange={(value) => {
                  setStateId(value);
                  setCityId(undefined);
                }}
                cityId={cityId}
                onCityIdChange={setCityId}
                countryOptions={countryOptions}
                states={states}
                cities={cities}
                cityDisabled={!stateId}
              />

              <RegisterAdditionalFields
                latitude={latitude}
                onLatitudeChange={setLatitude}
                longitude={longitude}
                onLongitudeChange={setLongitude}
                onGeocode={handleGeocode}
                isGeocoding={isGeocoding}
                geocodeError={geocodeError}
                canGeocode={canGeocode}
                latitudeDisabled
                longitudeDisabled
              />

              <RegisterSubmitRow isPending={isPending} />

              <RegisterFormFooter error={error} isSuccess={isSuccess} />
            </VStack>
          ) : null}
        </VStack>
      </chakra.form>
    </RegisterFormCard>
  );
}
