import * as Form from "@radix-ui/react-form";

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
    <div className="space-y-4 rounded-4xl border border-slate-100 bg-slate-50/60 p-4 sm:p-6">
      <div className="grid gap-4 sm:grid-cols-2">
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
              step="0.01"
              value={latitude}
              onChange={(event) => onLatitudeChange(event.target.value)}
              className="h-11 w-full rounded-2xl border border-slate-200 bg-white px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200 disabled:bg-slate-100 disabled:text-slate-500"
              placeholder="-23.550520"
              disabled={latitudeDisabled}
              required
            />
          </Form.Control>
        </Form.Field>

        <Form.Field className="space-y-2" name="longitude">
          <div className="flex items-baseline justify-between">
            <Form.Label className="text-sm font-medium">Longitude</Form.Label>
            <Form.Message
              className="text-xs text-rose-500"
              match="valueMissing"
            >
              Obrigatório
            </Form.Message>
          </div>
          <Form.Control asChild>
            <input
              id="longitude"
              type="number"
              step="0.01"
              value={longitude}
              onChange={(event) => onLongitudeChange(event.target.value)}
              className="h-11 w-full rounded-2xl border border-slate-200 bg-white px-4 text-sm text-slate-700 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200 disabled:bg-slate-100 disabled:text-slate-500"
              placeholder="-46.633308"
              disabled={longitudeDisabled}
              required
            />
          </Form.Control>
        </Form.Field>
      </div>

      <div className="flex flex-wrap items-center gap-3">
        <button
          type="button"
          onClick={onGeocode}
          disabled={!canGeocode || isGeocoding}
          className="inline-flex items-center justify-center rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-4 py-2 text-xs font-semibold text-white shadow-lg shadow-cyan-200 transition disabled:cursor-not-allowed disabled:opacity-60"
        >
          {isGeocoding ? "Buscando coordenadas..." : "Buscar coordenadas"}
        </button>
        {geocodeError ? (
          <span className="text-xs text-rose-500">{geocodeError}</span>
        ) : null}
      </div>
    </div>
  );
}
