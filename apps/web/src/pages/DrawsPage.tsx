import { FormEvent, useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "../api/client";
import { DrawTable } from "../components/DrawTable";
import type { Draw, DrawPayload } from "../types";
import { toDateInputValue } from "../utils/format";

const emptyForm: DrawPayload = {
  drawDate: "",
  n1: 0,
  n2: 0,
  n3: 0,
  n4: 0,
  n5: 0,
  star1: 0,
  star2: 0
};

export function DrawsPage() {
  const [page, setPage] = useState(1);
  const [editing, setEditing] = useState<Draw | null>(null);
  const [form, setForm] = useState<DrawPayload>(emptyForm);
  const [error, setError] = useState("");
  const queryClient = useQueryClient();
  const draws = useQuery({ queryKey: ["draws", page], queryFn: () => api.draws(page, 10) });

  const saveMutation = useMutation({
    mutationFn: (payload: DrawPayload) => (editing ? api.updateDraw(editing.id, payload) : api.createDraw(payload)),
    onSuccess: () => {
      setForm(emptyForm);
      setEditing(null);
      setError("");
      queryClient.invalidateQueries();
    },
    onError: (err) => setError(err.message)
  });

  const deleteMutation = useMutation({
    mutationFn: (draw: Draw) => api.deleteDraw(draw.id),
    onSuccess: () => {
      queryClient.invalidateQueries();
    }
  });

  const totalPages = useMemo(() => {
    if (!draws.data) {
      return 1;
    }
    return Math.max(1, Math.ceil(draws.data.total / draws.data.limit));
  }, [draws.data]);

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    saveMutation.mutate(form);
  }

  function editDraw(draw: Draw) {
    setEditing(draw);
    setForm({
      drawDate: toDateInputValue(draw.drawDate),
      n1: draw.n1,
      n2: draw.n2,
      n3: draw.n3,
      n4: draw.n4,
      n5: draw.n5,
      star1: draw.star1,
      star2: draw.star2
    });
  }

  return (
    <div className="grid gap-6 xl:grid-cols-[380px_1fr]">
      <section className="rounded-lg border border-slate-800 bg-slate-900/80 p-4">
        <h2 className="text-lg font-semibold text-white">{editing ? "Editar sorteo" : "Añadir sorteo"}</h2>
        <form className="mt-4 space-y-4" onSubmit={handleSubmit}>
          <label className="field-label">
            Fecha
            <input
              className="field-input"
              type="date"
              value={form.drawDate}
              onChange={(event) => setForm({ ...form, drawDate: event.target.value })}
              required
            />
          </label>
          <div className="grid grid-cols-5 gap-2">
            {(["n1", "n2", "n3", "n4", "n5"] as const).map((key) => (
              <label key={key} className="field-label">
                {key.toUpperCase()}
                <input
                  className="field-input"
                  min={1}
                  max={50}
                  type="number"
                  value={form[key] || ""}
                  onChange={(event) => setForm({ ...form, [key]: Number(event.target.value) })}
                  required
                />
              </label>
            ))}
          </div>
          <div className="grid grid-cols-2 gap-2">
            {(["star1", "star2"] as const).map((key) => (
              <label key={key} className="field-label">
                {key === "star1" ? "Estrella 1" : "Estrella 2"}
                <input
                  className="field-input"
                  min={1}
                  max={12}
                  type="number"
                  value={form[key] || ""}
                  onChange={(event) => setForm({ ...form, [key]: Number(event.target.value) })}
                  required
                />
              </label>
            ))}
          </div>
          {error ? <p className="text-sm text-rose-300">{error}</p> : null}
          <div className="flex gap-2">
            <button className="btn-primary" type="submit" disabled={saveMutation.isPending}>
              {editing ? "Guardar" : "Crear"}
            </button>
            {editing ? (
              <button
                className="btn-secondary"
                type="button"
                onClick={() => {
                  setEditing(null);
                  setForm(emptyForm);
                }}
              >
                Cancelar
              </button>
            ) : null}
          </div>
        </form>
      </section>

      <section className="space-y-4 overflow-x-auto">
        <DrawTable
          draws={draws.data?.draws ?? []}
          onEdit={editDraw}
          onDelete={(draw) => deleteMutation.mutate(draw)}
        />
        <div className="flex items-center justify-between">
          <button className="btn-secondary" type="button" disabled={page <= 1} onClick={() => setPage(page - 1)}>
            Anterior
          </button>
          <span className="text-sm text-slate-400">
            Página {page} de {totalPages}
          </span>
          <button
            className="btn-secondary"
            type="button"
            disabled={page >= totalPages}
            onClick={() => setPage(page + 1)}
          >
            Siguiente
          </button>
        </div>
      </section>
    </div>
  );
}
