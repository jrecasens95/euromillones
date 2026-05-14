import { useQuery } from "@tanstack/react-query";
import { api } from "../api/client";
import { NumberBall } from "../components/NumberBall";

export function StatsPage() {
  const frequencies = useQuery({ queryKey: ["frequencies"], queryFn: api.frequencies });
  const positions = useQuery({ queryKey: ["positions"], queryFn: api.positions });
  const delays = useQuery({ queryKey: ["delays"], queryFn: api.delays });
  const pairs = useQuery({ queryKey: ["pairs"], queryFn: () => api.pairs() });

  const maxNumberCount = Math.max(1, ...(frequencies.data?.numbers.map((item) => item.count) ?? [1]));
  const maxStarCount = Math.max(1, ...(frequencies.data?.stars.map((item) => item.count) ?? [1]));

  return (
    <div className="space-y-6">
      <section className="rounded-lg border border-slate-800 bg-slate-900/80 p-4">
        <h2 className="text-lg font-semibold text-white">Mapa de números</h2>
        <div className="mt-4 grid grid-cols-5 gap-2 sm:grid-cols-10">
          {frequencies.data?.numbers.map((item) => (
            <div
              key={item.value}
              className="rounded-lg border border-slate-800 p-2 text-center"
              style={{ backgroundColor: `rgba(34, 211, 238, ${0.08 + item.count / maxNumberCount / 1.5})` }}
            >
              <NumberBall value={item.value} muted />
              <span className="mt-2 block text-xs text-slate-300">{item.count}</span>
            </div>
          ))}
        </div>
      </section>

      <section className="rounded-lg border border-slate-800 bg-slate-900/80 p-4">
        <h2 className="text-lg font-semibold text-white">Mapa de estrellas</h2>
        <div className="mt-4 grid grid-cols-4 gap-2 md:grid-cols-12">
          {frequencies.data?.stars.map((item) => (
            <div
              key={item.value}
              className="rounded-lg border border-slate-800 p-2 text-center text-sm font-semibold text-white"
              style={{ backgroundColor: `rgba(251, 191, 36, ${0.08 + item.count / maxStarCount / 1.5})` }}
            >
              {item.value}
              <span className="mt-2 block text-xs text-slate-300">{item.count}</span>
            </div>
          ))}
        </div>
      </section>

      <section className="grid gap-4 xl:grid-cols-2">
        <div className="table-card">
          <h2 className="table-title">Frecuencia por posición</h2>
          <table className="data-table">
            <thead>
              <tr>
                <th>Posición</th>
                <th>Top 5</th>
              </tr>
            </thead>
            <tbody>
              {[...(positions.data?.numbers ?? []), ...(positions.data?.stars ?? [])].map((position) => (
                <tr key={position.position}>
                  <td>{position.position}</td>
                  <td>
                    {position.values
                      .slice()
                      .sort((a, b) => b.count - a.count)
                      .slice(0, 5)
                      .map((item) => `${item.value} (${item.count})`)
                      .join(", ")}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div className="table-card">
          <h2 className="table-title">Retrasos</h2>
          <table className="data-table">
            <thead>
              <tr>
                <th>Valor</th>
                <th>Sorteos sin salir</th>
              </tr>
            </thead>
            <tbody>
              {delays.data?.numbers
                .slice()
                .sort((a, b) => b.delay - a.delay)
                .slice(0, 15)
                .map((item) => (
                  <tr key={item.value}>
                    <td>{item.value}</td>
                    <td>{item.delay}</td>
                  </tr>
                ))}
            </tbody>
          </table>
        </div>
      </section>

      <section className="table-card">
        <h2 className="table-title">Pares frecuentes</h2>
        <table className="data-table">
          <thead>
            <tr>
              <th>Par</th>
              <th>Apariciones</th>
            </tr>
          </thead>
          <tbody>
            {pairs.data?.map((pair) => (
              <tr key={`${pair.a}-${pair.b}`}>
                <td>
                  {pair.a} + {pair.b}
                </td>
                <td>{pair.count}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
    </div>
  );
}
