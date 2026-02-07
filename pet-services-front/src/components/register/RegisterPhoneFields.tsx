import * as Form from "@radix-ui/react-form";

type RegisterPhoneFieldsProps = {
  countryCode: string;
  onCountryCodeChange: (value: string) => void;
  areaCode: string;
  onAreaCodeChange: (value: string) => void;
  phoneNumber: string;
  onPhoneNumberChange: (value: string) => void;
};

export default function RegisterPhoneFields({
  countryCode,
  onCountryCodeChange,
  areaCode,
  onAreaCodeChange,
  phoneNumber,
  onPhoneNumberChange,
}: RegisterPhoneFieldsProps) {
  return (
    <div className="grid gap-4 sm:grid-cols-3">
      <Form.Field className="space-y-2" name="countryCode">
        <div className="flex items-baseline justify-between">
          <Form.Label className="text-sm font-medium">DDI</Form.Label>
          <Form.Message
            className="text-xs text-rose-500"
            match="valueMissing"
          >
            Obrigatório
          </Form.Message>
        </div>
        <Form.Control asChild>
          <input
            id="countryCode"
            value={countryCode}
            onChange={(event) => onCountryCodeChange(event.target.value)}
            className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
            placeholder="55"
            required
          />
        </Form.Control>
      </Form.Field>

      <Form.Field className="space-y-2" name="areaCode">
        <div className="flex items-baseline justify-between">
          <Form.Label className="text-sm font-medium">DDD</Form.Label>
          <Form.Message
            className="text-xs text-rose-500"
            match="valueMissing"
          >
            Obrigatório
          </Form.Message>
        </div>
        <Form.Control asChild>
          <input
            id="areaCode"
            value={areaCode}
            onChange={(event) => onAreaCodeChange(event.target.value)}
            className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
            placeholder="11"
            required
          />
        </Form.Control>
      </Form.Field>

      <Form.Field className="space-y-2" name="phoneNumber">
        <div className="flex items-baseline justify-between">
          <Form.Label className="text-sm font-medium">Telefone</Form.Label>
          <Form.Message
            className="text-xs text-rose-500"
            match="valueMissing"
          >
            Obrigatório
          </Form.Message>
        </div>
        <Form.Control asChild>
          <input
            id="phoneNumber"
            value={phoneNumber}
            onChange={(event) => onPhoneNumberChange(event.target.value)}
            className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
            placeholder="99999-9999"
            required
          />
        </Form.Control>
      </Form.Field>
    </div>
  );
}
