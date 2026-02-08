import { NextResponse } from "next/server";

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const address = searchParams.get("address");

  if (!address) {
    return NextResponse.json(
      { error: "Endereço obrigatório" },
      { status: 400 },
    );
  }

  const apiKey = process.env.GOOGLE_MAPS_API_KEY;

  if (!apiKey) {
    return NextResponse.json(
      { error: "Chave de API não configurada" },
      { status: 500 },
    );
  }

  const response = await fetch(
    `https://maps.googleapis.com/maps/api/geocode/json?address=${encodeURIComponent(address)}&region=br&components=country:BR&key=${apiKey}`,
  );

  const data = await response.json();

  if (data.status !== "OK" || !data.results?.length) {
    return NextResponse.json(
      { error: "Geocodificação falhou", status: data.status },
      { status: 400 },
    );
  }

  const location = data.results[0].geometry?.location;

  if (!location) {
    return NextResponse.json(
      { error: "Localização não encontrada" },
      { status: 400 },
    );
  }

  return NextResponse.json({
    latitude: location.lat,
    longitude: location.lng,
  });
}
