# LILA BLACK — Player Journey Visualization Tool

A web tool for Level Designers to explore production telemetry from **LILA BLACK**: player paths, combat events, loot pickups, storm deaths, and heatmaps on top of official minimaps.

## Live Demo

| Service  | URL |
|----------|-----|
| Frontend | `https://lila-player-viz.vercel.app` *(replace after Vercel deploy)* |
| Backend  | `https://lila-player-viz-api.onrender.com` |

> Frontend: deploy to Vercel (steps below), then update this URL.

## Tech Stack

| Layer | Choice | Why |
|-------|--------|-----|
| Backend | Go + Gin + DuckDB | Fast REST API, queries parquet in-place (no ETL) |
| Frontend | React + TypeScript + Vite | Component model, type safety, fast dev/build |
| Map | Leaflet `CRS.Simple` + leaflet.heat | Correct pixel-space minimap rendering + heatmaps |
| Deploy | Render (API) + Vercel (UI) | Docker backend on Render free tier + Vercel CDN for frontend |

## Features

- Load & parse parquet telemetry via DuckDB
- World → minimap coordinate mapping (per-map strategy)
- Player journey polylines (humans solid blue, bots dashed gray)
- Distinct markers: kill, death, loot, storm death
- Filters: map, date folder, match
- Timeline playback (play / pause / scrub)
- Heatmap overlays: traffic, kills, deaths
- Architecture doc + data insights included

## Project Structure

```
lila_assignment/
├── backend/          # Go API (domain → repository → service → HTTP)
├── frontend/         # React UI
├── data/player_data/ # Parquet files + minimaps (from player_data.zip)
├── scripts/ingest.py # Schema validation script
├── ARCHITECTURE.md
├── INSIGHTS.md
├── docker-compose.yml
└── render.yaml
```

## Prerequisites

- Go 1.22+ (CGO enabled — required for DuckDB)
- Node.js 20+
- `player_data.zip` extracted to `data/player_data/`

```bash
unzip player_data.zip -d data
# → data/player_data/February_*/ + minimaps/
```

## Local Development

### 1. Backend

```bash
cd backend
go mod download
CGO_ENABLED=1 DATA_PATH=../data/player_data go run ./cmd/server
# → http://localhost:8080
```

### 2. Frontend

```bash
cd frontend
cp .env.example .env.local   # VITE_API_URL=http://localhost:8080
npm install
npm run dev
# → http://localhost:5173
```

### 3. Docker Compose (optional)

```bash
docker compose up --build
# Backend :8080, Frontend :5173
```

### 4. Validate data

```bash
pip install pyarrow duckdb pandas
DATA_PATH=data/player_data python scripts/ingest.py
```

## Environment Variables

### Backend

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP port |
| `DATA_PATH` | `../data/player_data` | Parquet root directory |
| `MINIMAP_PATH` | `$DATA_PATH/minimaps` | Static minimap images |
| `ALLOWED_ORIGINS` | `*` | CORS origins (comma-separated) |

### Frontend

| Variable | Default | Description |
|----------|---------|-------------|
| `VITE_API_URL` | `http://localhost:8080` | Backend base URL |

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/api/maps` | Map configs + minimap URLs |
| GET | `/api/dates` | Available date folders |
| GET | `/api/matches?map=&date=` | Match list (filtered) |
| GET | `/api/matches/:matchId/journey` | Paths + events for playback |
| GET | `/api/heatmap?map=&date=&match=&type=` | Heatmap points (`traffic`/`kills`/`deaths`) |
| GET | `/minimaps/*` | Minimap image assets |

## Deployment

### Backend → Render

1. Push this repo to **GitHub** (include `data/player_data/` — required for Docker build).
2. Go to [render.com](https://render.com) → **New +** → **Blueprint**.
3. Connect your GitHub repo — Render reads `render.yaml` automatically.
4. Deploy the `lila-player-viz-api` web service (Docker build, ~10–15 min first time due to DuckDB/CGO).
5. Copy your Render URL, e.g. `https://lila-player-viz-api.onrender.com`.
6. Test: `curl https://YOUR-RENDER-URL.onrender.com/health`

**Render env vars** (set in dashboard if not using blueprint defaults):

| Variable | Value |
|----------|-------|
| `DATA_PATH` | `/data/player_data` |
| `MINIMAP_PATH` | `/data/player_data/minimaps` |
| `ALLOWED_ORIGINS` | `https://your-app.vercel.app` *(after Vercel deploy)* |

> **Note:** Render free tier spins down after ~15 min idle. First request after sleep takes ~30–60s to wake up.

**Manual Render setup** (without Blueprint):
- New **Web Service** → connect repo
- **Environment:** Docker
- **Dockerfile Path:** `backend/Dockerfile`
- **Docker Context:** `.` (repo root)
- **Health Check Path:** `/health`

### Frontend → Vercel

1. Go to [vercel.com](https://vercel.com) → **Add New Project** → import your GitHub repo.
2. Set **Root Directory** to `frontend`.
3. Framework preset: **Vite** (auto-detected).
4. Add environment variable:
   ```
   VITE_API_URL=https://YOUR-RENDER-URL.onrender.com
   ```
5. Deploy → copy your Vercel URL (e.g. `https://lila-player-viz.vercel.app`).
6. Go back to **Render dashboard** → update `ALLOWED_ORIGINS` to your Vercel URL → redeploy backend.

### Post-deploy verification

```bash
API_URL=https://YOUR-RENDER-URL.onrender.com ./scripts/verify.sh
```

Then open your Vercel URL and walk through the [reviewer walkthrough](#walkthrough-for-reviewers) below.

### 5. Verify all features

```bash
./scripts/verify.sh
# Runs checklist: health, maps, dates, matches, journey, heatmaps, coordinate mapping
```

## Requirements Checklist

| Requirement | Status | How |
|-------------|--------|-----|
| Load & parse parquet | ✅ | DuckDB `read_parquet()` glob |
| World → minimap coordinates | ✅ | Per-map Strategy transformer (unit tested) |
| Human vs bot paths | ✅ | Solid blue vs dashed gray polylines |
| Kill / death / loot / storm markers | ✅ | Distinct shapes: ▲ ✕ ◆ ● |
| Filter by map, date, match | ✅ | Toolbar dropdowns + live API refresh |
| Timeline playback | ✅ | Play / pause / scrub (30s normalized) |
| Heatmap overlays | ✅ | Traffic, kills, deaths toggles |
| Hosted URL | ⏳ | Deploy to Render + Vercel (see below) |
| ARCHITECTURE.md | ✅ | SOLID, patterns, coordinate mapping |
| INSIGHTS.md (3 insights) | ✅ | Stats-backed level design insights |

## Design Patterns (see ARCHITECTURE.md)

- **Strategy** — coordinate transformer per map
- **Repository** — DuckDB behind `domain.Repository` interface
- **Dependency Inversion** — services depend on interfaces, not concrete DB
- **Container/Presenter** — React hooks vs presentational components

## Walkthrough (for reviewers)

1. Open the frontend URL.
2. Select **Ambrose Valley** → pick a **date** → choose a **match**.
3. Observe human (blue) vs bot (gray dashed) paths on the minimap.
4. Click **Play** on the timeline to scrub match progression.
5. Toggle **Heatmap → Traffic / Kills / Deaths** at map scope (works without a match too).
6. Use **Show Humans / Bots** toggles to isolate player types.

## Docs

- [ARCHITECTURE.md](./ARCHITECTURE.md) — design decisions, data flow, coordinate mapping
- [INSIGHTS.md](./INSIGHTS.md) — three data-driven level design insights

## License

Built as a take-home assignment for Lila Games.
