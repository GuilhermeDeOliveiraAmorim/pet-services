import * as Form from "@radix-ui/react-form";

type RegisterAddressFieldsProps = {
  street: string;
  onStreetChange: (value: string) => void;
  addressNumber: string;
  onAddressNumberChange: (value: string) => void;
  neighborhood: string;
  onNeighborhoodChange: (value: string) => void;
  city: string;
  onCityChange: (value: string) => void;
  state: string;
  onStateChange: (value: string) => void;
  zipCode: string;
  onZipCodeChange: (value: string) => void;
  country: string;
  onCountryChange: (value: string) => void;
};

export default function RegisterAddressFields({
  street,
  onStreetChange,
  addressNumber,
  onAddressNumberChange,
  neighborhood,
  onNeighborhoodChange,
  city,
  onCityChange,
  state,
  onStateChange,
  zipCode,
  onZipCodeChange,
  country,
  onCountryChange,
}: RegisterAddressFieldsProps) {
  return (
    <>
      <div className="grid gap-4 sm:grid-cols-2">
        <Form.Field className="space-y-2" name="street">
          <div className="flex items-baseline justify-between">
            <Form.Label className="text-sm font-medium">Rua</Form.Label>
            <Form.Message
              className="text-xs text-rose-500"
              match="valueMissing"
            >
              Obrigatório
            </Form.Message>
          </div>
          <Form.Control asChild>
            <input
              id="street"
              value={street}
              onChange={(event) => onStreetChange(event.target.value)}
              className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
              placeholder="Rua Exemplo"
              required
            />
          </Form.Control>
        </Form.Field>

        <Form.Field className="space-y-2" name="addressNumber">
          <div className="flex items-baseline justify-between">
            <Form.Label className="text-sm font-medium">Número</Form.Label>
            <Form.Message
              className="text-xs text-rose-500"
              match="valueMissing"
            >
              Obrigatório
            </Form.Message>
          </div>
          <Form.Control asChild>
            <input
              id="addressNumber"
              value={addressNumber}
              onChange={(event) => onAddressNumberChange(event.target.value)}
              className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
              placeholder="123"
              required
            />
          </Form.Control>
        </Form.Field>
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <Form.Field className="space-y-2" name="neighborhood">
          <div className="flex items-baseline justify-between">
            <Form.Label className="text-sm font-medium">Bairro</Form.Label>
            <Form.Message
              className="text-xs text-rose-500"
              match="valueMissing"
            >
              Obrigatório
            </Form.Message>
          </div>
          <Form.Control asChild>
            <input
              id="neighborhood"
              value={neighborhood}
              onChange={(event) => onNeighborhoodChange(event.target.value)}
              className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
              placeholder="Centro"
              required
            />
          </Form.Control>
        </Form.Field>

        <Form.Field className="space-y-2" name="city">
          <div className="flex items-baseline justify-between">
            <Form.Label className="text-sm font-medium">Cidade</Form.Label>
            <Form.Message
              className="text-xs text-rose-500"
              match="valueMissing"
            >
              Obrigatório
            </Form.Message>
          </div>
          <Form.Control asChild>
            <input
              id="city"
              value={city}
              onChange={(event) => onCityChange(event.target.value)}
              className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
              placeholder="São Paulo"
              required
            />
          </Form.Control>
        </Form.Field>
      </div>

      <div className="grid gap-4 sm:grid-cols-3">
        <Form.Field className="space-y-2" name="state">
          <div className="flex items-baseline justify-between">
            <Form.Label className="text-sm font-medium">Estado</Form.Label>
            <Form.Message
              className="text-xs text-rose-500"
              match="valueMissing"
            >
              Obrigatório
            </Form.Message>
          </div>
          <Form.Control asChild>
            <input
              id="state"
              value={state}
              onChange={(event) => onStateChange(event.target.value)}
              className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
              placeholder="SP"
              required
            />
          </Form.Control>
        </Form.Field>

        <Form.Field className="space-y-2" name="zipCode">
          <div className="flex items-baseline justify-between">
            <Form.Label className="text-sm font-medium">CEP</Form.Label>
            <Form.Message
              className="text-xs text-rose-500"
              match="valueMissing"
            >
              Obrigatório
            </Form.Message>
          </div>
          <Form.Control asChild>
            <input
              id="zipCode"
              value={zipCode}
              onChange={(event) => onZipCodeChange(event.target.value)}
              className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
              placeholder="00000-000"
              required
            />
          </Form.Control>
        </Form.Field>

        <Form.Field className="space-y-2" name="country">
          <div className="flex items-baseline justify-between">
            <Form.Label className="text-sm font-medium">País</Form.Label>
            <Form.Message
              className="text-xs text-rose-500"
              match="valueMissing"
            >
              Obrigatório
            </Form.Message>
          </div>
          <Form.Control asChild>
            <input
              id="country"
              value={country}
              onChange={(event) => onCountryChange(event.target.value)}
              className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
              placeholder="Brasil"
              required
            />
          </Form.Control>
        </Form.Field>
      </div>

    </>
  );
}
