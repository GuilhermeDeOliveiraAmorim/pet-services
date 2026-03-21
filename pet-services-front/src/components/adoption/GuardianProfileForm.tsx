"use client";

import { useMemo, useState } from "react";
import {
  Box,
  Button,
  Input,
  NativeSelect,
  Stack,
  Text,
  Textarea,
} from "@chakra-ui/react";

import { useReferenceCities, useReferenceStates } from "@/application";

type GuardianProfileFormValues = {
  displayName: string;
  guardianType: string;
  document: string;
  phone: string;
  whatsapp: string;
  about: string;
  cityId: string;
  stateId: string;
};

type GuardianProfileFormProps = {
  initialValues?: Partial<GuardianProfileFormValues>;
  isSubmitting?: boolean;
  submitLabel: string;
  onSubmit: (values: GuardianProfileFormValues) => Promise<void> | void;
};

const guardianTypeOptions = [
  { value: "ngo", label: "ONG" },
  { value: "independent", label: "Protetor independente" },
  { value: "owner", label: "Tutor responsável" },
];

export default function GuardianProfileForm({
  initialValues,
  isSubmitting,
  submitLabel,
  onSubmit,
}: GuardianProfileFormProps) {
  const [displayName, setDisplayName] = useState(
    initialValues?.displayName ?? "",
  );
  const [guardianType, setGuardianType] = useState(
    initialValues?.guardianType ?? "independent",
  );
  const [document, setDocument] = useState(initialValues?.document ?? "");
  const [phone, setPhone] = useState(initialValues?.phone ?? "");
  const [whatsapp, setWhatsapp] = useState(initialValues?.whatsapp ?? "");
  const [about, setAbout] = useState(initialValues?.about ?? "");
  const [stateId, setStateId] = useState(initialValues?.stateId ?? "");
  const [cityId, setCityId] = useState(initialValues?.cityId ?? "");
  const [stateSearch, setStateSearch] = useState("");
  const [citySearch, setCitySearch] = useState("");

  const { data: statesData } = useReferenceStates();
  const numericStateId = stateId ? Number(stateId) : undefined;
  const { data: citiesData } = useReferenceCities(
    numericStateId ? { stateId: numericStateId } : {},
    { enabled: Boolean(numericStateId) },
  );

  const filteredStates = useMemo(() => {
    const states = statesData?.states ?? [];
    const query = stateSearch.trim().toLowerCase();
    if (!query) {
      return states;
    }
    return states.filter((state) => state.name.toLowerCase().includes(query));
  }, [stateSearch, statesData?.states]);

  const filteredCities = useMemo(() => {
    const cities = citiesData?.cities ?? [];
    const query = citySearch.trim().toLowerCase();
    if (!query) {
      return cities;
    }
    return cities.filter((city) => city.name.toLowerCase().includes(query));
  }, [citiesData?.cities, citySearch]);

  const isValid =
    Boolean(displayName.trim()) &&
    Boolean(guardianType) &&
    Boolean(document.trim()) &&
    Boolean(phone.trim()) &&
    Boolean(stateId) &&
    Boolean(cityId);

  const handleSubmit = async () => {
    if (!isValid) {
      return;
    }

    await onSubmit({
      displayName: displayName.trim(),
      guardianType,
      document: document.trim(),
      phone: phone.trim(),
      whatsapp: whatsapp.trim(),
      about: about.trim(),
      stateId,
      cityId,
    });
  };

  return (
    <Stack gap={4}>
      <Box>
        <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
          Nome público
        </Text>
        <Input
          value={displayName}
          onChange={(event) => setDisplayName(event.target.value)}
          placeholder="Como o perfil deve aparecer nos anúncios"
          bg="white"
          borderColor="gray.200"
          focusRingColor="teal.200"
        />
      </Box>

      <Box>
        <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
          Tipo de guardião
        </Text>
        <NativeSelect.Root>
          <NativeSelect.Field
            value={guardianType}
            onChange={(event) => setGuardianType(event.target.value)}
            bg="white"
            borderColor="gray.200"
          >
            {guardianTypeOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </NativeSelect.Field>
          <NativeSelect.Indicator />
        </NativeSelect.Root>
      </Box>

      <Box>
        <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
          Documento
        </Text>
        <Input
          value={document}
          onChange={(event) => setDocument(event.target.value)}
          placeholder="CPF ou CNPJ"
          bg="white"
          borderColor="gray.200"
          focusRingColor="teal.200"
        />
      </Box>

      <Stack direction={{ base: "column", md: "row" }} gap={4}>
        <Box flex="1">
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Telefone
          </Text>
          <Input
            value={phone}
            onChange={(event) => setPhone(event.target.value)}
            placeholder="(11) 99999-9999"
            bg="white"
            borderColor="gray.200"
            focusRingColor="teal.200"
          />
        </Box>

        <Box flex="1">
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            WhatsApp
          </Text>
          <Input
            value={whatsapp}
            onChange={(event) => setWhatsapp(event.target.value)}
            placeholder="Opcional"
            bg="white"
            borderColor="gray.200"
            focusRingColor="teal.200"
          />
        </Box>
      </Stack>

      <Stack direction={{ base: "column", md: "row" }} gap={4}>
        <Box flex="1">
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Buscar estado
          </Text>
          <Input
            value={stateSearch}
            onChange={(event) => setStateSearch(event.target.value)}
            placeholder="Digite o nome do estado"
            bg="white"
            borderColor="gray.200"
            focusRingColor="teal.200"
            mb={2}
          />
          <NativeSelect.Root>
            <NativeSelect.Field
              value={stateId}
              onChange={(event) => {
                setStateId(event.target.value);
                setCityId("");
              }}
              bg="white"
              borderColor="gray.200"
            >
              <option value="">Selecione</option>
              {filteredStates.map((state) => (
                <option key={state.id} value={String(state.id)}>
                  {state.name}
                </option>
              ))}
            </NativeSelect.Field>
            <NativeSelect.Indicator />
          </NativeSelect.Root>
        </Box>

        <Box flex="1">
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Buscar cidade
          </Text>
          <Input
            value={citySearch}
            onChange={(event) => setCitySearch(event.target.value)}
            placeholder="Digite o nome da cidade"
            bg="white"
            borderColor="gray.200"
            focusRingColor="teal.200"
            mb={2}
            disabled={!stateId}
          />
          <NativeSelect.Root disabled={!stateId}>
            <NativeSelect.Field
              value={cityId}
              onChange={(event) => setCityId(event.target.value)}
              bg="white"
              borderColor="gray.200"
            >
              <option value="">Selecione</option>
              {filteredCities.map((city) => (
                <option key={city.id} value={String(city.id)}>
                  {city.name}
                </option>
              ))}
            </NativeSelect.Field>
            <NativeSelect.Indicator />
          </NativeSelect.Root>
        </Box>
      </Stack>

      <Box>
        <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
          Sobre você ou sua organização
        </Text>
        <Textarea
          value={about}
          onChange={(event) => setAbout(event.target.value)}
          placeholder="Conte sua experiência com resgate, adoção e critérios de acompanhamento."
          minH="32"
          bg="white"
          borderColor="gray.200"
          focusRingColor="teal.200"
        />
      </Box>

      <Button
        borderRadius="full"
        onClick={handleSubmit}
        disabled={!isValid || isSubmitting}
        alignSelf="start"
      >
        {isSubmitting ? "Salvando..." : submitLabel}
      </Button>
    </Stack>
  );
}
