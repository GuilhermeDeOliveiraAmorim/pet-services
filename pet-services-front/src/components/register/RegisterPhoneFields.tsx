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
      <div className="space-y-2">
        <label className="text-sm font-medium" htmlFor="countryCode">
          DDI
        </label>
        <input
          id="countryCode"
          value={countryCode}
          onChange={(event) => onCountryCodeChange(event.target.value)}
          className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
          placeholder="55"
          required
        />
      </div>
      <div className="space-y-2">
        <label className="text-sm font-medium" htmlFor="areaCode">
          DDD
        </label>
        <input
          id="areaCode"
          value={areaCode}
          onChange={(event) => onAreaCodeChange(event.target.value)}
          className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
          placeholder="11"
          required
        />
      </div>
      <div className="space-y-2">
        <label className="text-sm font-medium" htmlFor="phoneNumber">
          Telefone
        </label>
        <input
          id="phoneNumber"
          value={phoneNumber}
          onChange={(event) => onPhoneNumberChange(event.target.value)}
          className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
          placeholder="99999-9999"
          required
        />
      </div>
    </div>
  );
}
