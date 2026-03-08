import { type ChangeEvent, useMemo, useState } from "react";
import { Box, Grid, Input, NativeSelect, Text } from "@chakra-ui/react";

type RegisterAddressFieldsProps = {
  street: string;
  onStreetChange: (value: string) => void;
  addressNumber: string;
  onAddressNumberChange: (value: string) => void;
  neighborhood: string;
  onNeighborhoodChange: (value: string) => void;
  zipCode: string;
  onZipCodeChange: (value: string) => void;
  complement: string;
  onComplementChange: (value: string) => void;
  countryCode: string;
  onCountryCodeChange: (value: string) => void;
  stateId?: number;
  onStateIdChange: (value?: number) => void;
  cityId?: number;
  onCityIdChange: (value?: number) => void;
  countryOptions: Array<{ value: string; label: string }>;
  states: Array<{ id: number; code: string; name: string }>;
  cities: Array<{ id: number; name: string }>;
  cityDisabled?: boolean;
};

export default function RegisterAddressFields({
  street,
  onStreetChange,
  addressNumber,
  onAddressNumberChange,
  neighborhood,
  onNeighborhoodChange,
  countryCode,
  onCountryCodeChange,
  stateId,
  onStateIdChange,
  cityId,
  onCityIdChange,
  countryOptions,
  states,
  cities,
  cityDisabled,
  zipCode,
  onZipCodeChange,
  complement,
  onComplementChange,
}: RegisterAddressFieldsProps) {
  const [countrySearch, setCountrySearch] = useState("");
  const [stateSearch, setStateSearch] = useState("");
  const [citySearch, setCitySearch] = useState("");

  const filteredCountries = useMemo(() => {
    const query = countrySearch.trim().toLowerCase();
    if (!query) return countryOptions;
    return countryOptions.filter((country) =>
      country.label.toLowerCase().includes(query),
    );
  }, [countryOptions, countrySearch]);

  const filteredStates = useMemo(() => {
    const query = stateSearch.trim().toLowerCase();
    if (!query) return states;
    return states.filter((state) => state.name.toLowerCase().includes(query));
  }, [states, stateSearch]);

  const filteredCities = useMemo(() => {
    const query = citySearch.trim().toLowerCase();
    if (!query) return cities;
    return cities.filter((city) => city.name.toLowerCase().includes(query));
  }, [cities, citySearch]);

  return (
    <>
      <Grid gap={4} templateColumns={{ base: "1fr", sm: "1.3fr 0.7fr" }}>
        <Box minW={0}>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Rua
          </Text>
          <Input
            id="street"
            name="street"
            value={street}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              onStreetChange(event.target.value)
            }
            h="11"
            borderRadius="xl"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            placeholder="Rua Exemplo"
            required
            w="full"
          />
        </Box>

        <Box minW={0}>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Número
          </Text>
          <Input
            id="addressNumber"
            name="addressNumber"
            value={addressNumber}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              onAddressNumberChange(event.target.value)
            }
            h="11"
            borderRadius="xl"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            placeholder="123"
            required
            w="full"
          />
        </Box>
      </Grid>

      <Grid gap={4} templateColumns={{ base: "1fr", sm: "1fr 1fr" }}>
        <Box minW={0}>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Bairro
          </Text>
          <Input
            id="neighborhood"
            name="neighborhood"
            value={neighborhood}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              onNeighborhoodChange(event.target.value)
            }
            h="11"
            borderRadius="xl"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            placeholder="Centro"
            required
            w="full"
          />
        </Box>

        <Box minW={0}>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Cidade
          </Text>
          <Input
            value={citySearch}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              setCitySearch(event.target.value)
            }
            placeholder="Buscar cidade"
            h="9"
            mb={2}
            borderRadius="lg"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            disabled={cityDisabled}
            w="full"
          />
          <NativeSelect.Root
            size="md"
            h="11"
            borderRadius="xl"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            disabled={cityDisabled}
            w="full"
            minW={0}
          >
            <NativeSelect.Field
              id="cityId"
              name="cityId"
              value={cityId ? String(cityId) : ""}
              onChange={(event: ChangeEvent<HTMLSelectElement>) =>
                onCityIdChange(
                  event.target.value ? Number(event.target.value) : undefined,
                )
              }
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
      </Grid>

      <Grid
        gap={4}
        templateColumns={{ base: "1fr", sm: "repeat(3, minmax(0, 1fr))" }}
      >
        <Box minW={0}>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Estado
          </Text>
          <Input
            value={stateSearch}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              setStateSearch(event.target.value)
            }
            placeholder="Buscar estado"
            h="9"
            mb={2}
            borderRadius="lg"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            w="full"
          />
          <NativeSelect.Root
            size="md"
            h="11"
            borderRadius="xl"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            w="full"
            minW={0}
          >
            <NativeSelect.Field
              id="stateId"
              name="stateId"
              value={stateId ? String(stateId) : ""}
              onChange={(event: ChangeEvent<HTMLSelectElement>) =>
                onStateIdChange(
                  event.target.value ? Number(event.target.value) : undefined,
                )
              }
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

        <Box minW={0}>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            CEP
          </Text>
          <Input
            id="zipCode"
            name="zipCode"
            value={zipCode}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              onZipCodeChange(event.target.value)
            }
            h="11"
            borderRadius="xl"
            bg="white"
            borderColor="gray.200"
            focusRingColor="teal.200"
            placeholder="00000-000"
            inputMode="numeric"
            required
            w="full"
          />
        </Box>

        <Box minW={0}>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            País
          </Text>
          <Input
            value={countrySearch}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              setCountrySearch(event.target.value)
            }
            placeholder="Buscar país"
            h="9"
            mb={2}
            borderRadius="lg"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            w="full"
          />
          <NativeSelect.Root
            size="md"
            h="11"
            borderRadius="xl"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            w="full"
            minW={0}
          >
            <NativeSelect.Field
              id="countryCode"
              name="countryCode"
              value={countryCode}
              onChange={(event: ChangeEvent<HTMLSelectElement>) =>
                onCountryCodeChange(event.target.value)
              }
            >
              {filteredCountries.map((country) => (
                <option key={country.value} value={country.value}>
                  {country.label}
                </option>
              ))}
            </NativeSelect.Field>
            <NativeSelect.Indicator />
          </NativeSelect.Root>
        </Box>
      </Grid>

      <Box minW={0}>
        <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
          Complemento
        </Text>
        <Input
          id="complement"
          name="complement"
          value={complement}
          onChange={(event: ChangeEvent<HTMLInputElement>) =>
            onComplementChange(event.target.value)
          }
          h="11"
          borderRadius="xl"
          bg="white"
          borderColor="gray.200"
          focusRingColor="teal.200"
          placeholder="Apartamento, bloco, etc"
          w="full"
        />
      </Box>
    </>
  );
}
