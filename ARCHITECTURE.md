# Architecture

## Stack Choices

| Component | Choice | Rationale |
|-----------|--------|-----------|
| Query engine | DuckDB | Reads parquet directly via glob — zero ingestion pipeline |
| API | Go + Gin | Small binary, strong typing, easy Docker deploy |
| UI | React + Leaflet | Leaflet `CRS.Simple` maps 1024×1024 pixel space natively |
| Layering | domain → repository → service → HTTP | SOLID dependency inversion; swap DuckDB without touching handlers |

## SOLID Principles & Design Patterns

| Principle / Pattern | Where applied |
|---------------------|---------------|
| **S** — Single Responsibility | Each layer has one job: `repository/` reads data, `service/` transforms/buckets, `handlers/` encode HTTP |
| **O** — Open/Closed | New maps added via `TransformerRegistry.Register()` without modifying existing transformers |
| **L** — Liskov Substitution | `DuckDBRepository` satisfies `domain.Repository`; any repo impl works in services |
| **I** — Interface Segregation | `Repository` and `CoordinateTransformer` are narrow interfaces in `domain/` (no fat interfaces) |
| **D** — Dependency Inversion | Services depend on `domain.Repository` interface, not DuckDB directly; wired in `main.go` |
| **Strategy** | One `CoordinateTransformer` per map (`linearTransformer`) — swap algorithm per map |
| **Repository** | `DuckDBRepository` encapsulates all SQL/parquet access behind `domain.Repository` |
| **Registry** | `TransformerRegistry` holds map strategies; lookup by `map_id` at runtime |
| **Container/Presenter** (frontend) | Hooks fetch data (`useMatchData`, `useHeatmapData`); components only render props |

### Backend package layout

```
domain/          → entities + interfaces ONLY (no imports from outer layers)
repository/      → DuckDBRepository implements Repository
service/         → MatchService, HeatmapService, TransformerRegistry
transport/http/  → thin Gin handlers (decode → service → encode)
```

### Frontend layout

```
hooks/           → data fetching & playback state (smart)
components/      → pure rendering (dumb): MapView, EventMarker, Filters
context/         → FilterContext shares filter state without prop drilling
```

## Data Flow

```
parquet files (February_*/{user}_{match}.nakama-0)
        │
        ▼
DuckDB read_parquet(glob, filename=true)     ← repository layer
        │
        ▼
service layer (coordinate transform, heatmap bucketing)
        │
        ▼
Gin REST handlers  ──JSON──▶  React hooks  ──▶  Leaflet layers
```

1. **Repository** runs SQL against the parquet glob. The `filename=true` flag exposes folder names for date filtering (`February_10`, etc.).
2. **Service** converts world `(x, z)` → pixel `(px, py)` via a per-map `CoordinateTransformer` strategy, and buckets heatmap cells (32×32 grid).
3. **HTTP** returns JSON; frontend hooks fetch and pass data to presentational map components.
4. **Leaflet** renders minimap as `ImageOverlay` with `CRS.Simple` bounds `[0,0]–[1024,1024]`. Positions use `[y, x]` (Leaflet lat/lng = pixel row/col).

## Coordinate Mapping

Each map defines `scale`, `originX`, `originZ` (from README). The `y` column is elevation — **ignored for 2D plotting**.

```
u = (x - originX) / scale
v = (z - originZ) / scale
pixelX = u × 1024
pixelY = (1 - v) × 1024        ← Y flipped (image origin = top-left)
```

**Strategy pattern** — one transformer per map (`AmbroseValley`, `GrandRift`, `Lockdown`). Adding a 4th map = register a new strategy; no changes to repository or handlers.

**Example** (Ambrose Valley, x=-301.45, z=-355.55):
- u = 0.0762, v = 0.1305 → pixel (78, 890) ✓ matches README
- Verified by unit test: `coordinate_transformer_test.go`

## Assumptions

| Ambiguity | Assumption |
|-----------|------------|
| Date filter | Uses **ingestion folder** (`February_10`), not wall-clock date. Same match can appear in multiple folders; journey endpoint loads **all files** for a `match_id`. |
| Bot detection | UUID regex = human; numeric `user_id` = bot (per README). |
| Event bytes | DuckDB `CAST(event AS VARCHAR)` decodes binary event names. |
| Timestamps | **Timestamp unit mismatch**: The `ts` column contains Unix timestamps in **seconds** representing actual 2026 play dates (e.g., `1770754537` = Feb 10, 2026 20:15:37 UTC). However, the parquet metadata types it as `timestamp[ms]`, causing engines (DuckDB/PyArrow) to interpret them as milliseconds since 1970 (yielding `1970-01-21`). This scales all time spans down by 1,000. By recognizing this unit mismatch, we corrected the duration calculations so matches represent real duration (e.g., 523 seconds = 8.7 minutes) instead of sub-second timelines. Playback is mapped to a 30s wall-clock timeline for smooth scrubbing while showing the actual in-game match timer. |
| Heatmap scope | Works at map/date level without selecting a match; optional `match` param narrows to one game. |
| Event categories | `BotKill`→kill, `BotKilled`→death, `KilledByStorm`→storm (merged for marker display). |

## Tradeoffs

| Decision | Considered | Chose | Why |
|----------|-----------|-------|-----|
| Query engine | Postgres import, Polars batch | DuckDB in-process | No ETL step; assignment data is read-only parquet |
| Coord transform location | Frontend | Backend | Single source of truth; frontend stays dumb renderer |
| Heatmap bucketing | Client-side kernel | Server 32×32 grid | Smaller payloads; consistent across clients |
| Playback | Real-time playback (~8–12 min) | 30s wall-clock + In-game timer | Real-time is too slow for designer review; a 30s playback budget with a ticking in-game clock provides the best of both: rapid review and exact context of when events occur. |
| Data in Docker | S3 mount, download on boot | Bundle in image | ~34 MB fits; simpler Render deploy |
| Match list date filter | Match-level date metadata | Folder-based | Only date signal available in filenames |
| Event markers | Same circle shape | Distinct shapes per type | Level designers scan kill vs loot vs storm faster |

## Layer Diagram

```
transport/http/handlers  →  service/  →  domain interfaces
                                ↑
                         repository/duckdb (implements Repository)
                                ↑
                         service/coordinate_transformer (implements CoordinateTransformer)
```

Frontend mirrors this: **hooks** (data) → **presentational components** (MapView, EventMarker, etc.) with shared **FilterContext**.
