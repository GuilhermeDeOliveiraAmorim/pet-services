import * as Form from "@radix-ui/react-form";

import RadixSelectField from "@/components/common/RadixSelectField";

type RegisterAddressFieldsProps = {
  street: string;
  onStreetChange: (value: string) => void;
  addressNumber: string;
  onAddressNumberChange: (value: string) => void;
  neighborhood: string;
  onNeighborhoodChange: (value: string) => void;
  zipCode: string;
  onZipCodeChange: (value: string) => void;
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

        <RadixSelectField
          name="cityId"
          label="Cidade"
          value={cityId ? String(cityId) : ""}
          onValueChange={(value) =>
            onCityIdChange(value ? Number(value) : undefined)
          }
          options={cities.map((city) => ({
            value: String(city.id),
            label: city.name,
          }))}
          placeholder="Selecione"
          searchable
          searchPlaceholder="Buscar cidade"
          required
          disabled={cityDisabled}
        />
      </div>

      <div className="grid gap-4 sm:grid-cols-3">
        <RadixSelectField
          name="stateId"
          label="Estado"
          value={stateId ? String(stateId) : ""}
          onValueChange={(value) =>
            onStateIdChange(value ? Number(value) : undefined)
          }
          options={states.map((state) => ({
            value: String(state.id),
            label: state.name,
          }))}
          placeholder="Selecione"
          searchable
          searchPlaceholder="Buscar estado"
          required
        />

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

        <RadixSelectField
          name="countryCode"
          label="País"
          value={countryCode}
          onValueChange={onCountryCodeChange}
          options={countryOptions}
          searchable
          searchPlaceholder="Buscar país"
          required
        />
      </div>
    </>
  );
}
