import type { NumberFrequency } from "../types";

type FrequencyChartProps = {
  title: string;
  data: NumberFrequency[];
  compact?: boolean;
};

export function FrequencyChart({ title, data, compact }: FrequencyChartProps) {
  const max = Math.max(1, ...data.map((item) => item.count));

  return (
    <section className="rounded-lg border border-slate-800 bg-slate-900/80 p-4">
      <h2 className="text-base font-semibold text-white">{title}</h2>
      <div className="mt-4 flex h-52 items-end gap-1 overflow-x-auto pb-2">
        {data.map((item) => (
          <div key={item.value} className="flex min-w-4 flex-1 flex-col items-center gap-2">
            <div
              className="w-full rounded-t bg-cyan-400"
              style={{ height: `${Math.max(4, (item.count / max) * 180)}px` }}
              title={`${item.value}: ${item.count}`}
            />
            {!compact ? <span className="text-[10px] text-slate-500">{item.value}</span> : null}
          </div>
        ))}
      </div>
    </section>
  );
}
