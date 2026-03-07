import { type ChangeEvent } from "react";
import { Box, Button, Grid, Input, Text } from "@chakra-ui/react";

type RegisterAdditionalFieldsProps = {
  latitude: string;
  onLatitudeChange: (value: string) => void;
  longitude: string;
  onLongitudeChange: (value: string) => void;
  onGeocode: () => void;
  isGeocoding: boolean;
  geocodeError?: string;
  canGeocode: boolean;
  latitudeDisabled?: boolean;
  longitudeDisabled?: boolean;
};

export default function RegisterAdditionalFields({
  latitude,
  onLatitudeChange,
  longitude,
  onLongitudeChange,
  onGeocode,
  isGeocoding,
  geocodeError,
  canGeocode,
  latitudeDisabled,
  longitudeDisabled,
}: RegisterAdditionalFieldsProps) {
  return (
    <Box
      borderRadius="3xl"
      borderWidth="1px"
      borderColor="gray.200"
      bg="gray.50"
      p={{ base: 4, sm: 6 }}
    >
      <Grid gap={4} templateColumns={{ base: "1fr", sm: "1fr 1fr" }}>
        <Box>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Latitude
          </Text>
          <Input
            id="latitude"
            name="latitude"
            type="number"
            step="0.01"
            value={latitude}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              onLatitudeChange(event.target.value)
            }
            h="11"
            borderRadius="xl"
            bg="white"
            borderColor="gray.200"
            focusRingColor="teal.200"
            placeholder="-23.550520"
            disabled={latitudeDisabled}
            required
          />
        </Box>

        <Box>
          <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={2}>
            Longitude
          </Text>
          <Input
            id="longitude"
            name="longitude"
            type="number"
            step="0.01"
            value={longitude}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              onLongitudeChange(event.target.value)
            }
            h="11"
            borderRadius="xl"
            bg="white"
            borderColor="gray.200"
            focusRingColor="teal.200"
            placeholder="-46.633308"
            disabled={longitudeDisabled}
            required
          />
        </Box>
      </Grid>

      <Box mt={4}>
        <Button
          type="button"
          onClick={onGeocode}
          disabled={!canGeocode || isGeocoding}
          size="sm"
          borderRadius="full"
          bg="white"
          color="gray.700"
          borderWidth="1px"
          borderColor="gray.300"
          _hover={{ bg: "gray.100", borderColor: "green.300" }}
          _disabled={{ opacity: 0.6, cursor: "not-allowed" }}
        >
          {isGeocoding ? "Buscando coordenadas..." : "Buscar coordenadas"}
        </Button>
        {geocodeError ? (
          <Text mt={2} fontSize="xs" color="red.500">
            {geocodeError}
          </Text>
        ) : null}
      </Box>
    </Box>
  );
}
