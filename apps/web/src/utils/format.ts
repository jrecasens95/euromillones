export function formatDrawDate(value: string) {
  if (!value) {
    return "Sin datos";
  }
  return new Intl.DateTimeFormat("es-ES", {
    day: "2-digit",
    month: "short",
    year: "numeric"
  }).format(new Date(value));
}

export function toDateInputValue(value: string) {
  if (!value) {
    return "";
  }
  return value.slice(0, 10);
}
