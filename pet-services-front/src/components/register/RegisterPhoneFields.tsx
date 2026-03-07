import { type ChangeEvent } from "react";
import { Box, Grid, Input, NativeSelect, Text } from "@chakra-ui/react";

type RegisterPhoneFieldsProps = {
  countryCode: string;
  onCountryCodeChange: (value: string) => void;
  countryCodeDisplayValue?: string;
  areaCode: string;
  onAreaCodeChange: (value: string) => void;
  phoneNumber: string;
  onPhoneNumberChange: (value: string) => void;
  dialCodeOptions: Array<{ value: string; label: string }>;
};

export default function RegisterPhoneFields({
  countryCode,
  onCountryCodeChange,
  countryCodeDisplayValue,
  areaCode,
  onAreaCodeChange,
  phoneNumber,
  onPhoneNumberChange,
  dialCodeOptions,
}: RegisterPhoneFieldsProps) {
  const formatPhoneNumber = (value: string) => {
    const digits = value.replace(/\D/g, "");
    if (!digits) {
      return "";
    }

    const isMobile = digits.length > 8;
    const maxDigits = isMobile ? 9 : 8;
    const splitAt = isMobile ? 5 : 4;
    const trimmed = digits.slice(0, maxDigits);

    return `${trimmed.slice(0, splitAt)}-${trimmed.slice(splitAt)}`;
  };

  return (
    <Grid
      gap={4}
      templateColumns={{ base: "1fr", sm: "repeat(3, minmax(0, 1fr))" }}
    >
      <Box minW={0}>
        <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
          DDI
        </Text>
        <NativeSelect.Root size="md" w="full" minW={0}>
          <NativeSelect.Field
            name="countryCode"
            value={countryCode}
            onChange={(event: ChangeEvent<HTMLSelectElement>) =>
              onCountryCodeChange(event.target.value)
            }
            h="11"
            borderRadius="xl"
            bg="gray.50"
            borderColor="gray.200"
            focusRingColor="teal.200"
            borderWidth="1px"
            fontSize="md"
            color="gray.800"
            w="full"
          >
            {dialCodeOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </NativeSelect.Field>
          <NativeSelect.Indicator color="gray.500" />
        </NativeSelect.Root>
        {countryCodeDisplayValue ? (
          <Text mt={1.5} fontSize="xs" color="gray.500">
            Selecionado: {countryCodeDisplayValue}
          </Text>
        ) : null}
      </Box>

      <Box minW={0}>
        <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
          DDD
        </Text>
        <Input
          id="areaCode"
          name="areaCode"
          value={areaCode}
          onChange={(event: ChangeEvent<HTMLInputElement>) => {
            const digits = event.target.value.replace(/\D/g, "").slice(0, 2);
            onAreaCodeChange(digits);
          }}
          h="11"
          borderRadius="xl"
          bg="gray.50"
          borderColor="gray.200"
          focusRingColor="teal.200"
          placeholder="11"
          inputMode="numeric"
          pattern="\d{2}"
          minLength={2}
          maxLength={2}
          required
          w="full"
        />
      </Box>

      <Box minW={0}>
        <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
          Telefone
        </Text>
        <Input
          id="phoneNumber"
          name="phoneNumber"
          value={phoneNumber}
          onChange={(event: ChangeEvent<HTMLInputElement>) =>
            onPhoneNumberChange(formatPhoneNumber(event.target.value))
          }
          onBlur={(event: ChangeEvent<HTMLInputElement>) =>
            onPhoneNumberChange(formatPhoneNumber(event.target.value))
          }
          h="11"
          borderRadius="xl"
          bg="gray.50"
          borderColor="gray.200"
          focusRingColor="teal.200"
          placeholder="7898-7898 ou 78987-7898"
          inputMode="numeric"
          pattern="^(\d{4}-\d{4}|\d{5}-\d{4})$"
          maxLength={10}
          required
          w="full"
        />
      </Box>
    </Grid>
  );
}
