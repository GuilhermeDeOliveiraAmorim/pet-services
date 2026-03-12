export const BRAZIL_COUNTRY_CODE = "BR";

export const isBrazilCountry = (value: string) => {
  const normalized = value.trim().toUpperCase();
  return normalized === "BR" || normalized === "BRASIL";
};

export const fetchCoordinatesByZipCode = async (zipCodeValue: string) => {
  const response = await fetch(
    `/api/geocode?address=${encodeURIComponent(zipCodeValue)}`,
  );
  const data = await response.json();

  if (!response.ok) {
    const message =
      typeof data?.error === "string"
        ? data.error
        : "Não foi possível buscar coordenadas pelo CEP.";
    throw new Error(message);
  }

  const nextLatitude = Number(data?.latitude);
  const nextLongitude = Number(data?.longitude);

  const hasValidCoordinates =
    Number.isFinite(nextLatitude) &&
    Number.isFinite(nextLongitude) &&
    nextLatitude >= -90 &&
    nextLatitude <= 90 &&
    nextLongitude >= -180 &&
    nextLongitude <= 180;

  if (!hasValidCoordinates) {
    throw new Error("Coordenadas inválidas retornadas para o CEP informado.");
  }

  return { latitude: nextLatitude, longitude: nextLongitude };
};

export const fetchAddressByZipCode = async (zipCodeValue: string) => {
  const digits = zipCodeValue.replace(/\D/g, "");

  if (digits.length !== 8) {
    throw new Error("Informe um CEP válido com 8 dígitos.");
  }

  const response = await fetch(`https://viacep.com.br/ws/${digits}/json/`);
  const data = await response.json();

  if (!response.ok || data?.erro) {
    throw new Error("Não foi possível localizar o endereço pelo CEP.");
  }

  return {
    street: String(data?.logradouro ?? "").trim(),
    neighborhood: String(data?.bairro ?? "").trim(),
    city: String(data?.localidade ?? "").trim(),
    state: String(data?.uf ?? "").trim(),
    country: BRAZIL_COUNTRY_CODE,
  };
};
