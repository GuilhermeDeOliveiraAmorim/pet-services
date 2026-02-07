"use client";

import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import * as Form from "@radix-ui/react-form";
import * as Tabs from "@radix-ui/react-tabs";
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
        const label = `${country.name} (+${country.dialCode})`;
        const key = `${country.dialCode}-${country.name}`;
        if (!seen.has(key)) {
          seen.add(key);
          acc.push({ value: country.dialCode, label });
        }
        return acc;
      },
      [],
    );
  }, [countries]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

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
        number: phoneNumber,
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
      <Form.Root className="space-y-6" onSubmit={handleSubmit}>
        <Tabs.Root value={step} onValueChange={setStep}>
          <Tabs.List className="grid grid-cols-2 rounded-full bg-slate-100 p-1 text-sm">
            <Tabs.Trigger
              value="account"
              className="rounded-full px-3 py-2 font-medium text-slate-500 data-[state=active]:bg-white data-[state=active]:text-slate-900 data-[state=active]:shadow-sm"
            >
              Conta
            </Tabs.Trigger>
            <Tabs.Trigger
              value="address"
              className="rounded-full px-3 py-2 font-medium text-slate-500 data-[state=active]:bg-white data-[state=active]:text-slate-900 data-[state=active]:shadow-sm"
            >
              Endereço
            </Tabs.Trigger>
          </Tabs.List>

          <Tabs.Content value="account" className="mt-6 space-y-6">
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
              countryCode={phoneCountryCode}
              onCountryCodeChange={setPhoneCountryCode}
              areaCode={areaCode}
              onAreaCodeChange={setAreaCode}
              phoneNumber={phoneNumber}
              onPhoneNumberChange={setPhoneNumber}
              dialCodeOptions={dialCodeOptions}
            />

            <div className="flex justify-end">
              <button
                type="button"
                onClick={() => setStep("address")}
                className="rounded-full border border-slate-200 px-6 py-2 text-sm font-semibold text-slate-700"
              >
                Continuar
              </button>
            </div>
          </Tabs.Content>

          <Tabs.Content value="address" className="mt-6 space-y-6">
            <RegisterAddressFields
              street={street}
              onStreetChange={setStreet}
              addressNumber={addressNumber}
              onAddressNumberChange={setAddressNumber}
              neighborhood={neighborhood}
              onNeighborhoodChange={setNeighborhood}
              zipCode={zipCode}
              onZipCodeChange={setZipCode}
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
              complement={complement}
              onComplementChange={setComplement}
              latitude={latitude}
              onLatitudeChange={setLatitude}
            />

            <RegisterSubmitRow
              longitude={longitude}
              onLongitudeChange={setLongitude}
              isPending={isPending}
            />

            <RegisterFormFooter error={error} isSuccess={isSuccess} />
          </Tabs.Content>
        </Tabs.Root>
      </Form.Root>
    </RegisterFormCard>
  );
}
