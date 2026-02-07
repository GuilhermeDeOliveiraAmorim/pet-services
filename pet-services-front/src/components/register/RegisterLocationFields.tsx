type RegisterLocationFieldsProps = {
  latitude: string;
  onLatitudeChange: (value: string) => void;
  longitude: string;
  onLongitudeChange: (value: string) => void;
};

export default function RegisterLocationFields({
  latitude,
  onLatitudeChange,
  longitude,
  onLongitudeChange,
}: RegisterLocationFieldsProps) {
  return (
    <div className="grid gap-4 sm:grid-cols-2">
      <div className="space-y-2">
        <label className="text-sm font-medium" htmlFor="latitude">
          Latitude
        </label>
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
      </div>
      <div className="space-y-2">
        <label className="text-sm font-medium" htmlFor="longitude">
          Longitude
        </label>
        <input
          id="longitude"
          type="number"
          step="0.000001"
          value={longitude}
          onChange={(event) => onLongitudeChange(event.target.value)}
          className="h-11 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
          placeholder="-46.633308"
          required
        />
      </div>
    </div>
  );
}
