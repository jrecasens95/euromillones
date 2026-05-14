export function SettingsPage() {
  return (
    <div className="grid gap-4 lg:grid-cols-2">
      <section className="rounded-lg border border-slate-800 bg-slate-900/80 p-5">
        <h2 className="text-lg font-semibold text-white">Proyecto</h2>
        <dl className="mt-4 space-y-3 text-sm">
          <div className="flex justify-between gap-4 border-b border-slate-800 pb-3">
            <dt className="text-slate-400">Backend</dt>
            <dd className="text-white">Go + Fiber + PostgreSQL</dd>
          </div>
          <div className="flex justify-between gap-4 border-b border-slate-800 pb-3">
            <dt className="text-slate-400">Frontend</dt>
            <dd className="text-white">React + TailwindCSS</dd>
          </div>
          <div className="flex justify-between gap-4">
            <dt className="text-slate-400">Estado</dt>
            <dd className="text-white">MVP interno</dd>
          </div>
        </dl>
      </section>
      <section className="rounded-lg border border-amber-300/30 bg-amber-300/10 p-5">
        <h2 className="text-lg font-semibold text-amber-100">Aviso</h2>
        <p className="mt-3 text-sm leading-6 text-amber-50">
          Esta aplicación analiza datos históricos, pero no puede predecir resultados futuros. El Euromillones es un
          juego aleatorio.
        </p>
      </section>
    </div>
  );
}
