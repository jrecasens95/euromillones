import { Route, Routes } from "react-router-dom";
import { DashboardPage } from "../pages/DashboardPage";
import { DrawsPage } from "../pages/DrawsPage";
import { GeneratorPage } from "../pages/GeneratorPage";
import { SettingsPage } from "../pages/SettingsPage";
import { StatsPage } from "../pages/StatsPage";

export function AppRouter() {
  return (
    <Routes>
      <Route path="/" element={<DashboardPage />} />
      <Route path="/draws" element={<DrawsPage />} />
      <Route path="/stats" element={<StatsPage />} />
      <Route path="/generator" element={<GeneratorPage />} />
      <Route path="/settings" element={<SettingsPage />} />
    </Routes>
  );
}
