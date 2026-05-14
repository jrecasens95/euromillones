import { useQuery } from "@tanstack/react-query";
import { api } from "../api/client";
import { FrequencyChart } from "../components/FrequencyChart";
import { StatCard } from "../components/StatCard";
import { formatDrawDate } from "../utils/format";

export function DashboardPage() {
  const dashboard = useQuery({ queryKey: ["dashboard"], queryFn: api.dashboard });
  const frequencies = useQuery({ queryKey: ["frequencies"], queryFn: api.frequencies });
  const hotCold = useQuery({ queryKey: ["hot-cold"], queryFn: api.hotCold });

  return (
    <div className="space-y-6">
      <section className="grid gap-4 md:grid-cols-5">
        <StatCard label="Sorteos" value={dashboard.data?.totalDraws ?? 0} detail="histórico guardado" />
        <StatCard
          label="Número top"
          value={dashboard.data?.mostFrequentNumber?.value ?? "-"}
          detail={`${dashboard.data?.mostFrequentNumber?.count ?? 0} apariciones`}
        />
        <StatCard
          label="Estrella top"
          value={dashboard.data?.mostFrequentStar?.value ?? "-"}
          detail={`${dashboard.data?.mostFrequentStar?.count ?? 0} apariciones`}
        />
        <StatCard
          label="Más retrasado"
          value={dashboard.data?.mostDelayedNumber?.value ?? "-"}
          detail={`${dashboard.data?.mostDelayedNumber?.delay ?? 0} sorteos`}
        />
        <StatCard label="Último sorteo" value={formatDrawDate(dashboard.data?.lastDrawDate ?? "")} />
      </section>

      <section className="grid gap-4 xl:grid-cols-2">
        <FrequencyChart title="Frecuencia de números" data={frequencies.data?.numbers ?? []} />
        <FrequencyChart title="Frecuencia de estrellas" data={frequencies.data?.stars ?? []} />
      </section>

      <section className="grid gap-4 xl:grid-cols-2">
        <FrequencyChart title="Hot numbers" data={hotCold.data?.hotNumbers ?? []} compact />
        <FrequencyChart title="Cold numbers" data={hotCold.data?.coldNumbers ?? []} compact />
      </section>
    </div>
  );
}
