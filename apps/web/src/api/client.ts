import type {
  DashboardStats,
  DelaysResponse,
  Draw,
  DrawPayload,
  FrequenciesResponse,
  GeneratedCombination,
  GenerationStrategy,
  HotColdResponse,
  PaginatedDraws,
  PairStat,
  PositionsResponse
} from "../types";

const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:4000";

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_URL}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...options?.headers
    },
    ...options
  });

  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: "Error inesperado" }));
    throw new Error(body.error ?? "Error inesperado");
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
}

export const api = {
  dashboard: () => request<DashboardStats>("/api/stats/dashboard"),
  draws: (page = 1, limit = 10) => request<PaginatedDraws>(`/api/draws?page=${page}&limit=${limit}`),
  draw: (id: number) => request<Draw>(`/api/draws/${id}`),
  createDraw: (payload: DrawPayload) =>
    request<Draw>("/api/draws", { method: "POST", body: JSON.stringify(payload) }),
  updateDraw: (id: number, payload: DrawPayload) =>
    request<Draw>(`/api/draws/${id}`, { method: "PUT", body: JSON.stringify(payload) }),
  deleteDraw: (id: number) => request<void>(`/api/draws/${id}`, { method: "DELETE" }),
  frequencies: () => request<FrequenciesResponse>("/api/stats/frequencies"),
  positions: () => request<PositionsResponse>("/api/stats/positions"),
  hotCold: () => request<HotColdResponse>("/api/stats/hot-cold"),
  delays: () => request<DelaysResponse>("/api/stats/delays"),
  pairs: () => request<PairStat[]>("/api/stats/pairs"),
  generate: (strategy: GenerationStrategy, count: number) =>
    request<{ combinations: GeneratedCombination[] }>("/api/generate", {
      method: "POST",
      body: JSON.stringify({ strategy, count })
    })
};
