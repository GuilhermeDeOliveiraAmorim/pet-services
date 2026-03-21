"use client";

import { useMemo, useState } from "react";
import {
  Box,
  Button,
  Checkbox,
  Input,
  NativeSelect,
  SimpleGrid,
  Stack,
  Text,
  Textarea,
} from "@chakra-ui/react";

import {
  usePetList,
  useReferenceCities,
  useReferenceStates,
} from "@/application";
import { type AdoptionListing } from "@/domain";

type AdoptionListingFormValues = {
  petId: string;
  title: string;
  description: string;
  adoptionReason: string;
  sex: string;
  size: string;
  ageGroup: string;
  cityId: string;
  stateId: string;
  latitude?: number;
  longitude?: number;
  vaccinated: boolean;
  neutered: boolean;
  dewormed: boolean;
  specialNeeds: boolean;
  goodWithChildren: boolean;
  goodWithDogs: boolean;
  goodWithCats: boolean;
  requiresHouseScreening: boolean;
};

type AdoptionListingFormProps = {
  initialValues?: Partial<AdoptionListingFormValues>;
  listing?: AdoptionListing;
  isSubmitting?: boolean;
  submitLabel: string;
  onSubmit: (values: AdoptionListingFormValues) => Promise<void> | void;
};

const sexOptions = [
  { value: "male", label: "Macho" },
  { value: "female", label: "Fêmea" },
];

const sizeOptions = [
  { value: "small", label: "Pequeno" },
  { value: "medium", label: "Médio" },
  { value: "large", label: "Grande" },
];

const ageGroupOptions = [
  { value: "puppy", label: "Filhote" },
  { value: "adult", label: "Adulto" },
  { value: "senior", label: "Sênior" },
];

export default function AdoptionListingForm({
  initialValues,
  listing,
  isSubmitting,
  submitLabel,
  onSubmit,
}: AdoptionListingFormProps) {
  const [petId, setPetId] = useState(
    initialValues?.petId ?? listing?.petId ?? "",
  );
  const [title, setTitle] = useState(
    initialValues?.title ?? listing?.title ?? "",
  );
  const [description, setDescription] = useState(
    initialValues?.description ?? listing?.description ?? "",
  );
  const [adoptionReason, setAdoptionReason] = useState(
    initialValues?.adoptionReason ?? listing?.adoptionReason ?? "",
  );
  const [sex, setSex] = useState(initialValues?.sex ?? listing?.sex ?? "male");
  const [size, setSize] = useState(
    initialValues?.size ?? listing?.size ?? "medium",
  );
  const [ageGroup, setAgeGroup] = useState(
    initialValues?.ageGroup ?? listing?.ageGroup ?? "adult",
  );
  const [stateId, setStateId] = useState(
    initialValues?.stateId ?? listing?.stateId ?? "",
  );
  const [cityId, setCityId] = useState(
    initialValues?.cityId ?? listing?.cityId ?? "",
  );
  const [latitude, setLatitude] = useState(
    initialValues?.latitude ?? listing?.latitude ?? undefined,
  );
  const [longitude, setLongitude] = useState(
    initialValues?.longitude ?? listing?.longitude ?? undefined,
  );
  const [vaccinated, setVaccinated] = useState(
    initialValues?.vaccinated ?? false,
  );
  const [neutered, setNeutered] = useState(
    initialValues?.neutered ?? listing?.neutered ?? false,
  );
  const [dewormed, setDewormed] = useState(
    initialValues?.dewormed ?? listing?.dewormed ?? false,
  );
  const [specialNeeds, setSpecialNeeds] = useState(
    initialValues?.specialNeeds ?? listing?.specialNeeds ?? false,
  );
  const [goodWithChildren, setGoodWithChildren] = useState(
    initialValues?.goodWithChildren ?? listing?.goodWithChildren ?? false,
  );
  const [goodWithDogs, setGoodWithDogs] = useState(
    initialValues?.goodWithDogs ?? listing?.goodWithDogs ?? false,
  );
  const [goodWithCats, setGoodWithCats] = useState(
    initialValues?.goodWithCats ?? listing?.goodWithCats ?? false,
  );
  const [requiresHouseScreening, setRequiresHouseScreening] = useState(
    initialValues?.requiresHouseScreening ??
      listing?.requiresHouseScreening ??
      false,
  );
  const [stateSearch, setStateSearch] = useState("");
  const [citySearch, setCitySearch] = useState("");

  const { data: petsData } = usePetList();
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

  const pets = petsData?.pets ?? [];

  const isValid =
    Boolean(petId) &&
    Boolean(title.trim()) &&
    Boolean(description.trim()) &&
    Boolean(adoptionReason.trim()) &&
    Boolean(stateId) &&
    Boolean(cityId);

  const handleSubmit = async () => {
    if (!isValid) {
      return;
    }

    await onSubmit({
      petId,
      title: title.trim(),
      description: description.trim(),
      adoptionReason: adoptionReason.trim(),
      sex,
      size,
      ageGroup,
      stateId,
      cityId,
      latitude,
      longitude,
      vaccinated,
      neutered,
      dewormed,
      specialNeeds,
      goodWithChildren,
      goodWithDogs,
      goodWithCats,
      requiresHouseScreening,
    });
  };

  return (
    <Stack gap={4}>
      <Box>
        <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
          Pet vinculado
        </Text>
        <NativeSelect.Root disabled={Boolean(listing?.petId)}>
          <NativeSelect.Field
            value={petId}
            onChange={(event) => setPetId(event.target.value)}
            bg="white"
            borderColor="gray.200"
          >
            <option value="">Selecione</option>
            {pets.map((pet) => (
              <option key={pet.id} value={pet.id}>
                {pet.name} {pet.breed ? `• ${pet.breed}` : ""}
              </option>
            ))}
          </NativeSelect.Field>
          <NativeSelect.Indicator />
        </NativeSelect.Root>
      </Box>

      <Box>
        <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
          Título do anúncio
        </Text>
        <Input
          value={title}
          onChange={(event) => setTitle(event.target.value)}
          placeholder="Ex.: Luna procura uma família tranquila"
          bg="white"
          borderColor="gray.200"
          focusRingColor="teal.200"
        />
      </Box>

      <Box>
        <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
          Descrição
        </Text>
        <Textarea
          value={description}
          onChange={(event) => setDescription(event.target.value)}
          minH="28"
          placeholder="Descreva rotina, energia, comportamento e cuidados do pet."
          bg="white"
          borderColor="gray.200"
          focusRingColor="teal.200"
        />
      </Box>

      <Box>
        <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
          Motivo da adoção
        </Text>
        <Textarea
          value={adoptionReason}
          onChange={(event) => setAdoptionReason(event.target.value)}
          minH="24"
          placeholder="Explique o contexto da adoção e o tipo de lar ideal."
          bg="white"
          borderColor="gray.200"
          focusRingColor="teal.200"
        />
      </Box>

      <SimpleGrid columns={{ base: 1, md: 3 }} gap={4}>
        <Box>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Sexo
          </Text>
          <NativeSelect.Root>
            <NativeSelect.Field
              value={sex}
              onChange={(event) => setSex(event.target.value)}
              bg="white"
              borderColor="gray.200"
            >
              {sexOptions.map((option) => (
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
            Porte
          </Text>
          <NativeSelect.Root>
            <NativeSelect.Field
              value={size}
              onChange={(event) => setSize(event.target.value)}
              bg="white"
              borderColor="gray.200"
            >
              {sizeOptions.map((option) => (
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
            Faixa etária
          </Text>
          <NativeSelect.Root>
            <NativeSelect.Field
              value={ageGroup}
              onChange={(event) => setAgeGroup(event.target.value)}
              bg="white"
              borderColor="gray.200"
            >
              {ageGroupOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </NativeSelect.Field>
            <NativeSelect.Indicator />
          </NativeSelect.Root>
        </Box>
      </SimpleGrid>

      <SimpleGrid columns={{ base: 1, md: 2 }} gap={4}>
        <Box>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Buscar estado
          </Text>
          <Input
            value={stateSearch}
            onChange={(event) => setStateSearch(event.target.value)}
            placeholder="Digite o estado"
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

        <Box>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Buscar cidade
          </Text>
          <Input
            value={citySearch}
            onChange={(event) => setCitySearch(event.target.value)}
            placeholder="Digite a cidade"
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
      </SimpleGrid>

      <SimpleGrid columns={{ base: 1, md: 2 }} gap={4}>
        <Box>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Latitude
          </Text>
          <Input
            type="number"
            value={latitude ?? ""}
            onChange={(event) =>
              setLatitude(
                event.target.value === ""
                  ? undefined
                  : Number(event.target.value),
              )
            }
            bg="white"
            borderColor="gray.200"
            focusRingColor="teal.200"
          />
        </Box>
        <Box>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Longitude
          </Text>
          <Input
            type="number"
            value={longitude ?? ""}
            onChange={(event) =>
              setLongitude(
                event.target.value === ""
                  ? undefined
                  : Number(event.target.value),
              )
            }
            bg="white"
            borderColor="gray.200"
            focusRingColor="teal.200"
          />
        </Box>
      </SimpleGrid>

      <SimpleGrid columns={{ base: 1, md: 2 }} gap={3}>
        {[
          {
            checked: vaccinated,
            setChecked: setVaccinated,
            label: "Vacinado",
          },
          {
            checked: neutered,
            setChecked: setNeutered,
            label: "Castrado",
          },
          {
            checked: dewormed,
            setChecked: setDewormed,
            label: "Vermifugado",
          },
          {
            checked: specialNeeds,
            setChecked: setSpecialNeeds,
            label: "Possui necessidades especiais",
          },
          {
            checked: goodWithChildren,
            setChecked: setGoodWithChildren,
            label: "Bom com crianças",
          },
          {
            checked: goodWithDogs,
            setChecked: setGoodWithDogs,
            label: "Bom com cães",
          },
          {
            checked: goodWithCats,
            setChecked: setGoodWithCats,
            label: "Bom com gatos",
          },
          {
            checked: requiresHouseScreening,
            setChecked: setRequiresHouseScreening,
            label: "Exige triagem/visita domiciliar",
          },
        ].map(({ checked, setChecked, label }) => (
          <Checkbox.Root
            key={label}
            checked={checked}
            onCheckedChange={(details) => setChecked(Boolean(details.checked))}
          >
            <Checkbox.HiddenInput />
            <Checkbox.Control />
            <Checkbox.Label>{label}</Checkbox.Label>
          </Checkbox.Root>
        ))}
      </SimpleGrid>

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
