type RegisterSubmitRowProps = {
  longitude: string;
  onLongitudeChange: (value: string) => void;
  isPending: boolean;
};

export default function RegisterSubmitRow({
  longitude,
  onLongitudeChange,
  isPending,
}: RegisterSubmitRowProps) {
  return (
    <div className="grid gap-4 sm:grid-cols-2">
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
      <div className="flex items-end">
        <button
          type="submit"
          disabled={isPending}
          className="inline-flex h-11 w-full items-center justify-center rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-4 text-sm font-semibold text-white shadow-lg shadow-cyan-200 transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
        >
          {isPending ? "Criando..." : "Criar conta"}
        </button>
      </div>
    </div>
  );
}
