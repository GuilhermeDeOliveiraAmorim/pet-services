import * as Form from "@radix-ui/react-form";

type RegisterAdditionalFieldsProps = {
  complement: string;
  onComplementChange: (value: string) => void;
  latitude: string;
  onLatitudeChange: (value: string) => void;
};

export default function RegisterAdditionalFields({
  complement,
  onComplementChange,
  latitude,
  onLatitudeChange,
}: RegisterAdditionalFieldsProps) {
  return (
    <div className="grid gap-4 sm:grid-cols-2">
      <Form.Field className="space-y-2" name="complement">
        <Form.Label className="text-sm font-medium">Complemento</Form.Label>
        <Form.Control asChild>
          <input
            id="complement"
            value={complement}
            onChange={(event) => onComplementChange(event.target.value)}
            className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
            placeholder="Apartamento, bloco, etc"
          />
        </Form.Control>
      </Form.Field>

      <Form.Field className="space-y-2" name="latitude">
        <div className="flex items-baseline justify-between">
          <Form.Label className="text-sm font-medium">Latitude</Form.Label>
          <Form.Message
            className="text-xs text-rose-500"
            match="valueMissing"
          >
            Obrigatório
          </Form.Message>
        </div>
        <Form.Control asChild>
          <input
            id="latitude"
            type="number"
            step="0.000001"
            value={latitude}
            onChange={(event) => onLatitudeChange(event.target.value)}
            className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
            placeholder="-23.550520"
            required
          />
        </Form.Control>
      </Form.Field>
    </div>
  );
}
