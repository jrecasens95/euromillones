import { NavLink } from "react-router-dom";
import { AppRouter } from "./router";

export function App() {
  const navItems = [
    { href: "/", label: "Dashboard" },
    { href: "/draws", label: "Sorteos" },
    { href: "/stats", label: "Estadísticas" },
    { href: "/generator", label: "Generador" },
    { href: "/settings", label: "Configuración" }
  ];

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100">
      <div className="mx-auto flex min-h-screen max-w-7xl flex-col px-4 py-5 lg:px-6">
        <header className="mb-6 flex flex-col gap-4 border-b border-slate-800 pb-5 xl:flex-row xl:items-end xl:justify-between">
          <div>
            <p className="text-xs font-semibold uppercase text-cyan-300">Euromillones interno</p>
            <h1 className="mt-2 text-3xl font-semibold text-white md:text-4xl">Analizador estadístico</h1>
          </div>
          <nav className="flex gap-2 overflow-x-auto">
            {navItems.map((item) => (
              <NavLink
                key={item.href}
                to={item.href}
                className={({ isActive }) =>
                  [
                    "rounded-lg border px-3 py-2 text-sm font-semibold transition",
                    isActive
                      ? "border-cyan-300 bg-cyan-300 text-slate-950"
                      : "border-slate-800 bg-slate-900 text-slate-300 hover:border-slate-600"
                  ].join(" ")
                }
              >
                {item.label}
              </NavLink>
            ))}
          </nav>
        </header>
        <main className="pb-8">
          <AppRouter />
        </main>
      </div>
    </div>
  );
}
