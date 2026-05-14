import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { api } from "../api/client";
import { CombinationCard } from "../components/CombinationCard";
import { StrategySelector } from "../components/StrategySelector";
import type { GeneratedCombination, GenerationStrategy } from "../types";

export function GeneratorPage() {
  const [strategy, setStrategy] = useState<GenerationStrategy>("balanced");
  const [count, setCount] = useState(5);
  const [combinations, setCombinations] = useState<GeneratedCombination[]>([]);
  const [error, setError] = useState("");

  const generateMutation = useMutation({
    mutationFn: () => api.generate(strategy, count),
    onSuccess: (data) => {
      setCombinations(data.combinations);
      setError("");
    },
    onError: (err) => setError(err.message)
  });

  return (
    <div className="grid gap-6 xl:grid-cols-[360px_1fr]">
      <section className="rounded-lg border border-slate-800 bg-slate-900/80 p-4">
        <h2 className="text-lg font-semibold text-white">Generador</h2>
        <div className="mt-4 space-y-5">
          <StrategySelector value={strategy} onChange={setStrategy} />
          <label className="field-label">
            Cantidad
            <input
              className="field-input"
              type="number"
              min={1}
              max={50}
              value={count}
              onChange={(event) => setCount(Number(event.target.value))}
            />
          </label>
          {error ? <p className="text-sm text-rose-300">{error}</p> : null}
          <button
            className="btn-primary w-full"
            type="button"
            disabled={generateMutation.isPending}
            onClick={() => generateMutation.mutate()}
          >
            Generar combinaciones
          </button>
        </div>
      </section>

      <section className="grid content-start gap-4 lg:grid-cols-2">
        {combinations.map((combination, index) => (
          <CombinationCard key={`${combination.strategy}-${index}`} combination={combination} />
        ))}
      </section>
    </div>
  );
}
