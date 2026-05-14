import type { Draw } from "../types";
import { formatDrawDate } from "../utils/format";
import { NumberBall } from "./NumberBall";
import { StarBall } from "./StarBall";

type DrawTableProps = {
  draws: Draw[];
  onEdit: (draw: Draw) => void;
  onDelete: (draw: Draw) => void;
};

export function DrawTable({ draws, onEdit, onDelete }: DrawTableProps) {
  return (
    <div className="overflow-hidden rounded-lg border border-slate-800 bg-slate-900/80">
      <table className="w-full min-w-[760px] text-left text-sm">
        <thead className="bg-slate-950/80 text-xs uppercase text-slate-400">
          <tr>
            <th className="px-4 py-3">Fecha</th>
            <th className="px-4 py-3">Números</th>
            <th className="px-4 py-3">Estrellas</th>
            <th className="px-4 py-3 text-right">Acciones</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-800">
          {draws.map((draw) => (
            <tr key={draw.id}>
              <td className="px-4 py-3 font-medium text-white">{formatDrawDate(draw.drawDate)}</td>
              <td className="px-4 py-3">
                <div className="flex gap-2">
                  {[draw.n1, draw.n2, draw.n3, draw.n4, draw.n5].map((number) => (
                    <NumberBall key={number} value={number} />
                  ))}
                </div>
              </td>
              <td className="px-4 py-3">
                <div className="flex gap-2">
                  {[draw.star1, draw.star2].map((star) => (
                    <StarBall key={star} value={star} />
                  ))}
                </div>
              </td>
              <td className="px-4 py-3">
                <div className="flex justify-end gap-2">
                  <button className="btn-secondary" type="button" onClick={() => onEdit(draw)}>
                    Editar
                  </button>
                  <button className="btn-danger" type="button" onClick={() => onDelete(draw)}>
                    Eliminar
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
