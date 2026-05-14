import type { GeneratedCombination } from "../types";
import { NumberBall } from "./NumberBall";
import { StarBall } from "./StarBall";

export function CombinationCard({ combination }: { combination: GeneratedCombination }) {
  return (
    <article className="rounded-lg border border-slate-800 bg-slate-900/80 p-4">
      <div className="flex flex-wrap items-center gap-2">
        {combination.numbers.map((number) => (
          <NumberBall key={number} value={number} />
        ))}
        <span className="mx-1 h-8 w-px bg-slate-700" />
        {combination.stars.map((star) => (
          <StarBall key={star} value={star} />
        ))}
      </div>
      <div className="mt-4 flex items-center justify-between gap-3">
        <span className="rounded bg-slate-800 px-2 py-1 text-xs font-semibold uppercase text-cyan-200">
          {combination.strategy}
        </span>
        <strong className="text-lg text-white">{combination.score}/100</strong>
      </div>
      <p className="mt-3 text-sm text-slate-400">{combination.explanation}</p>
    </article>
  );
}
