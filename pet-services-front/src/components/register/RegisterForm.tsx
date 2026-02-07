"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import * as Form from "@radix-ui/react-form";
import { type UserType, UserTypes } from "@/domain";
import { useUserRegister } from "@/application";
import RegisterAccountFields from "./RegisterAccountFields";
import RegisterAdditionalFields from "./RegisterAdditionalFields";
import RegisterAddressFields from "./RegisterAddressFields";
import RegisterFormCard from "./RegisterFormCard";
import RegisterFormFooter from "./RegisterFormFooter";
import RegisterPhoneFields from "./RegisterPhoneFields";
import RegisterSubmitRow from "./RegisterSubmitRow";

const defaultCountryCode = "55";
const defaultCountry = "Brasil";

export default function RegisterForm() {
  const router = useRouter();
  const { mutateAsync, isPending, error, isSuccess } = useUserRegister();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [userType, setUserType] = useState<UserType>(UserTypes.Owner);
  const [countryCode, setCountryCode] = useState(defaultCountryCode);
  const [areaCode, setAreaCode] = useState("");
  const [phoneNumber, setPhoneNumber] = useState("");

  const [street, setStreet] = useState("");
  const [addressNumber, setAddressNumber] = useState("");
  const [neighborhood, setNeighborhood] = useState("");
  const [city, setCity] = useState("");
  const [state, setState] = useState("");
  const [zipCode, setZipCode] = useState("");
  const [country, setCountry] = useState(defaultCountry);
  const [complement, setComplement] = useState("");
  const [latitude, setLatitude] = useState("0");
  const [longitude, setLongitude] = useState("0");

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
        countryCode,
        areaCode,
        number: phoneNumber,
      },
      address: {
        street,
        number: addressNumber,
        neighborhood,
        city,
        zipCode,
        state,
        country,
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
          countryCode={countryCode}
          onCountryCodeChange={setCountryCode}
          areaCode={areaCode}
          onAreaCodeChange={setAreaCode}
          phoneNumber={phoneNumber}
          onPhoneNumberChange={setPhoneNumber}
        />

        <RegisterAddressFields
          street={street}
          onStreetChange={setStreet}
          addressNumber={addressNumber}
          onAddressNumberChange={setAddressNumber}
          neighborhood={neighborhood}
          onNeighborhoodChange={setNeighborhood}
          city={city}
          onCityChange={setCity}
          state={state}
          onStateChange={setState}
          zipCode={zipCode}
          onZipCodeChange={setZipCode}
          country={country}
          onCountryChange={setCountry}
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
      </Form.Root>
    </RegisterFormCard>
  );
}
