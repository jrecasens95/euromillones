import type { GenerationStrategy } from "../types";

const strategies: { value: GenerationStrategy; label: string }[] = [
  { value: "balanced", label: "Balanced" },
  { value: "hot", label: "Hot" },
  { value: "cold", label: "Cold" },
  { value: "delayed", label: "Delayed" },
  { value: "random", label: "Random" },
  { value: "anti_human", label: "Anti-human" }
];

type StrategySelectorProps = {
  value: GenerationStrategy;
  onChange: (value: GenerationStrategy) => void;
};

export function StrategySelector({ value, onChange }: StrategySelectorProps) {
  return (
    <div className="grid grid-cols-2 gap-2 md:grid-cols-3">
      {strategies.map((strategy) => (
        <button
          key={strategy.value}
          type="button"
          onClick={() => onChange(strategy.value)}
          className={[
            "rounded-lg border px-3 py-3 text-sm font-semibold transition",
            value === strategy.value
              ? "border-cyan-300 bg-cyan-300 text-slate-950"
              : "border-slate-800 bg-slate-900 text-slate-300 hover:border-slate-600"
          ].join(" ")}
        >
          {strategy.label}
        </button>
      ))}
    </div>
  );
}
