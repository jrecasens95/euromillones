# Euromillones

Aplicación web para consultar, mantener y analizar históricos de sorteos de Euromillones. El proyecto combina una API en Go con una interfaz en React para visualizar estadísticas, explorar sorteos guardados y generar combinaciones con distintas estrategias.

> Este proyecto tiene fines educativos y estadísticos. No garantiza resultados ni debe interpretarse como recomendación de juego.

## Características

- Dashboard con resumen del histórico, último sorteo, números frecuentes y números retrasados.
- CRUD de sorteos sobre ficheros JSON anuales.
- Estadísticas de frecuencias, posiciones, retrasos y pares frecuentes.
- Generador de combinaciones con estrategias configurables.
- Frontend SPA con React, Vite, React Query, React Router, Tailwind CSS y Radix UI.
- Backend REST con Go y Fiber.
- Persistencia simple basada en archivos JSON, sin base de datos externa.

## Stack

| Capa | Tecnología |
| --- | --- |
| Frontend | React 18, Vite, TypeScript, Tailwind CSS, Radix UI |
| Estado remoto | TanStack React Query |
| Backend | Go, Fiber |
| Persistencia | JSON por año |
| Monorepo | npm workspaces/scripts por aplicación |

## Estructura

```text
.
├── apps
│   ├── api
│   │   ├── cmd/server
│   │   ├── data
│   │   └── internal
│   └── web
│       └── src
├── packages
├── package.json
└── README.md
```

## Requisitos

- Go 1.23 o superior.
- Node.js 18 o superior.
- npm.

## Puesta En Marcha

Instala las dependencias del frontend:

```bash
cd apps/web
npm install
```

Desde la raíz del repositorio puedes ejecutar cada parte con los scripts principales.

API:

```bash
npm run dev:api
```

Frontend:

```bash
npm run dev:web
```

Por defecto, la API escucha en `http://localhost:4000` y el frontend de Vite en `http://localhost:5173`.

## Scripts

| Comando | Descripción |
| --- | --- |
| `npm run dev:web` | Arranca el frontend en modo desarrollo. |
| `npm run build:web` | Compila el frontend para producción. |
| `npm run dev:api` | Arranca la API con `go run`. |
| `npm run build:api` | Compila los paquetes Go. |
| `npm run deps:api` | Ejecuta `go mod tidy` en la API. |

## Variables De Entorno

### API

| Variable | Valor por defecto | Descripción |
| --- | --- | --- |
| `PORT` | `4000` | Puerto del servidor HTTP. |
| `DRAWS_DATA_DIR` | `data` | Directorio con los archivos JSON de sorteos. |
| `JWT_SECRET` | vacío | Reservada para configuración de autenticación. |

Ejemplo:

```bash
cd apps/api
DRAWS_DATA_DIR=/absolute/path/data PORT=4000 go run ./cmd/server
```

### Frontend

| Variable | Valor por defecto | Descripción |
| --- | --- | --- |
| `VITE_API_URL` | `http://localhost:4000` | URL base de la API. |

Ejemplo:

```bash
cd apps/web
VITE_API_URL=http://localhost:4000 npm run dev
```

## Datos

Los sorteos se cargan desde archivos JSON anuales en `apps/api/data`. Cada archivo se nombra con el año correspondiente, por ejemplo `2024.json`.

Formato base:

```json
[
  {
    "sorteo": 1,
    "fecha": "2-ene",
    "numeros": [10, 18, 21, 33, 45],
    "estrellas": [3, 8]
  }
]
```

El cargador acepta fechas con mes numérico o abreviaturas en español, con o sin año explícito. Algunos históricos incluyen sorteos de frontera de año, por ejemplo `"sorteo": "2016/001"`.

## API REST

Todas las rutas cuelgan de `/api`.

| Método | Ruta | Descripción |
| --- | --- | --- |
| `GET` | `/draws` | Lista sorteos paginados. |
| `GET` | `/draws/:id` | Obtiene un sorteo por ID. |
| `POST` | `/draws` | Crea un sorteo. |
| `PUT` | `/draws/:id` | Actualiza un sorteo. |
| `DELETE` | `/draws/:id` | Elimina un sorteo. |
| `GET` | `/stats/dashboard` | Resumen general del histórico. |
| `GET` | `/stats/frequencies` | Frecuencias de números y estrellas. |
| `GET` | `/stats/positions` | Frecuencia por posición. |
| `GET` | `/stats/hot-cold` | Números calientes y fríos. |
| `GET` | `/stats/delays` | Retrasos por número. |
| `GET` | `/stats/pairs` | Pares frecuentes. |
| `POST` | `/generate` | Genera combinaciones. |

## Desarrollo

Ejecuta los tests de la API desde `apps/api`:

```bash
go test ./...
```

Compila el frontend:

```bash
npm run build:web
```

## Licencia

Este proyecto se distribuye bajo la licencia MIT. Consulta el archivo [LICENSE](LICENSE) para más detalles.
