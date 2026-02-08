const baseUrl = process.env.NEXT_PUBLIC_WEB_URL ?? "http://localhost:3000";
const address = process.env.TEST_ADDRESS ?? "49048-390";

const run = async () => {
  const url = `${baseUrl}/api/geocode?address=${encodeURIComponent(address)}`;
  const response = await fetch(url);
  const data = await response.json();

  if (!response.ok) {
    console.error("Falha na geocodificação:", {
      status: response.status,
      data,
    });
    process.exit(1);
  }

  console.log("Geocode ok:", data);
};

run().catch((error) => {
  console.error("Erro ao testar geocode:", error);
  process.exit(1);
});
