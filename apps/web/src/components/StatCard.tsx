type StatCardProps = {
  label: string;
  value: string | number;
  detail?: string;
};

export function StatCard({ label, value, detail }: StatCardProps) {
  return (
    <article className="rounded-lg border border-slate-800 bg-slate-900/80 p-4">
      <p className="text-xs font-semibold uppercase text-slate-400">{label}</p>
      <strong className="mt-2 block text-3xl font-semibold text-white">{value}</strong>
      {detail ? <span className="mt-2 block text-sm text-slate-400">{detail}</span> : null}
    </article>
  );
}
