# Featherweight Media Pipeline

Upload an image, and it's compressed, resized, and converted into modern
web-ready formats (JPEG, PNG, WebP, AVIF) in the background — with a live
progress view and side-by-side results when it's done.

- **Backend:** Go (Gin), clean architecture, PostgreSQL, Cloudinary
- **Frontend:** React + Vite + Tailwind CSS

---

## Table of contents

- [How it works](#how-it-works)
- [Architecture](#architecture)
- [Project structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Getting started](#getting-started)
  - [Option A — Docker Compose](#option-a--docker-compose-recommended)
  - [Option B — Run manually](#option-b--run-manually)
- [Configuration](#configuration)
- [API reference](#api-reference)
- [Frontend notes](#frontend-notes)
- [Known issues / housekeeping](#known-issues--housekeeping)

---

## How it works

1. The user drops an image into the frontend and picks target width, quality,
   and output formats.
2. The frontend uploads it to the Go backend (`POST /api/v1/media/upload`),
   which stores the original, creates a `Job` record, and hands it to a
   background worker pool.
3. A worker compresses/resizes/converts the image into each requested format
   and uploads the results to Cloudinary.
4. The frontend polls `GET /api/v1/jobs/:id` every ~900ms until the job is
   `completed` or `failed`, then renders a card per output format with the
   size saved and a download link.
5. Originals and processed files are periodically cleaned up after a
   configurable retention window.

## Architecture

The backend follows a clean-architecture layout so business logic stays
independent of frameworks and delivery mechanisms:

```
domain      → core entities (Job, Output) and repository interfaces.
              No knowledge of HTTP, JSON, or Postgres.
usecase     → application logic (create a job, fetch a job) that
              orchestrates the domain via interfaces.
repository  → in-memory / Postgres implementations of the domain's
              repository interfaces (internal/infrastructure/postgres).
delivery    → HTTP layer: Gin handlers, routing, and DTOs that translate
              domain entities to/from the JSON the frontend expects
              (internal/delivery/http, /handler, /dto).
processor   → image encoding/decoding (JPEG/PNG/WebP/AVIF).
worker      → background job pool + retention cleanup.
infrastructure → external services: Postgres connection, Cloudinary uploader.
```

The `delivery/dto` package is what maps domain structs to the wire format
(`job_id`, `file_name`, `status`, `outputs[].url`, `outputs[].size_bytes`,
etc.) — the domain layer itself carries no JSON tags or knowledge of the API
contract.

## Project structure

```
Featherweight_Media_Pipeline/
├── Dockerfile
├── docker-compose.yml
├── .env                          # shared env vars for docker-compose
├── backend/
│   ├── cmd/api/main.go            # entrypoint
│   ├── internal/
│   │   ├── domain/                # Job, Output, repository interfaces
│   │   ├── usecase/                # MediaUseCase, JobUseCase
│   │   ├── delivery/
│   │   │   ├── http/router.go
│   │   │   ├── handler/            # upload_handler.go, job_handler.go, health_handler.go
│   │   │   └── dto/                 # job.go, output.go, mapper.go, error_response.go
│   │   ├── processor/              # jpeg/png/webp/avif encoders
│   │   ├── worker/                 # job pool + cleanup
│   │   ├── infrastructure/
│   │   │   ├── postgres/            # connection, job_repository, migrations
│   │   │   └── cloudinary/          # uploader
│   │   └── config/config.go
│   └── uploads/                    # local scratch space (originals/processed)
└── frontend/
    ├── src/
    │   ├── pages/                  # Home, Upload, About
    │   ├── components/             # Dropzone, JobStatus, ResultsPanel, ResultCard, ...
    │   ├── hooks/                  # useJobPolling, useUploadHistory
    │   ├── api/client.js           # upload/getJob/getHealth
    │   └── context/ThemeContext.jsx
    └── vite.config.js
```

## Prerequisites

- Go 1.25+
- Node.js 18+ and npm
- PostgreSQL 16 (or Docker, to run it in a container)
- A [Cloudinary](https://cloudinary.com/console) account (free tier is fine)
  for storing processed images
- Docker + Docker Compose (optional, for the containerized setup)

## Getting started

### Option A — Docker Compose (recommended)

> ⚠️ At the moment `docker-compose.yml` in this repo is empty and the compose
> definition lives in `.gitignore` instead — the two files got swapped at
> some point. See [Known issues](#known-issues--housekeeping) below for the
> one-line fix before using this option.

1. Copy the compose config into `docker-compose.yml` (see note above).
2. Create a `.env` file at the project root with your Cloudinary credentials:

   ```env
   CLOUDINARY_CLOUD_NAME=your-cloud-name
   CLOUDINARY_API_KEY=your-api-key
   CLOUDINARY_API_SECRET=your-api-secret
   CLOUDINARY_FOLDER=featherweight
   ```

3. Start everything:

   ```bash
   docker compose up --build
   ```

   This brings up Postgres and the API (`http://localhost:8080`). Migrations
   in `backend/internal/infrastructure/postgres/migrations` run against the
   `postgres` service automatically on connect.

4. Run the frontend separately (Vite isn't containerized here):

   ```bash
   cd frontend
   npm install
   npm run dev
   ```

   The dev server runs at `http://localhost:5173` and talks to the API via
   `VITE_API_BASE_URL` (see `frontend/.env`).

### Option B — Run manually

**1. Start PostgreSQL** (locally installed, or via Docker):

```bash
docker run -d --name featherweight-postgres \
  -e POSTGRES_USER=featherweight \
  -e POSTGRES_PASSWORD=featherweight \
  -e POSTGRES_DB=featherweight \
  -p 5432:5432 postgres:16-alpine
```

**2. Configure the backend**

```bash
cd backend
cp .env.example .env   # if present — otherwise create one, see Configuration below
```

Fill in your `DB_*` and `CLOUDINARY_*` values.

**3. Run the backend**

```bash
go mod download
go run ./cmd/api
```

The API starts on `SERVER_ADDRESS` (default `:8080`) and logs which job store
it's using.

**4. Run the frontend**

```bash
cd frontend
npm install
cp .env.example .env   # adjust VITE_API_BASE_URL if the backend isn't on :8080
npm run dev
```

Open `http://localhost:5173`.

## Configuration

### Backend (`backend/.env`, loaded via `godotenv`)

| Variable | Default | Description |
|---|---|---|
| `SERVER_ADDRESS` | `:8080` | HTTP listen address |
| `JOB_STORE` | `postgres` | Job storage backend |
| `DATABASE_URL` | *(built from `DB_*` below if unset)* | Full Postgres connection string |
| `DB_HOST` | `localhost` | Postgres host |
| `DB_PORT` | `5432` | Postgres port |
| `DB_USER` | `featherweight` | Postgres user |
| `DB_PASSWORD` | `featherweight` | Postgres password |
| `DB_NAME` | `featherweight` | Postgres database name |
| `DB_SSLMODE` | `disable` | Postgres SSL mode |
| `DATABASE_MAX_CONNS` | `10` | Max pool connections |
| `DATABASE_CONNECT_TIMEOUT_SECONDS` | `10` | Connect timeout |
| `CLOUDINARY_CLOUD_NAME` | — | Cloudinary cloud name |
| `CLOUDINARY_API_KEY` | — | Cloudinary API key |
| `CLOUDINARY_API_SECRET` | — | Cloudinary API secret |
| `CLOUDINARY_FOLDER` | `featherweight` | Upload folder in Cloudinary |
| `UPLOAD_DIRECTORY` | `uploads/originals` | Local scratch dir for originals |
| `OUTPUT_DIRECTORY` | `uploads/processed` | Local scratch dir for processed files |
| `WORKER_COUNT` | `3` | Concurrent background workers |
| `MAX_UPLOAD_SIZE_MB` | `32` | Max accepted upload size |
| `MAX_IMAGE_WIDTH` | `5000` | Max resize width allowed |
| `FILE_RETENTION_HOURS` | `24` | How long processed files are kept |
| `CLEANUP_INTERVAL_MINUTES` | `60` | How often the retention sweep runs |
| `COMMAND_TIMEOUT_SECONDS` | `60` | Timeout for external encoder commands |

### Frontend (`frontend/.env`)

| Variable | Default | Description |
|---|---|---|
| `VITE_API_BASE_URL` | `http://localhost:8080` | Backend base URL |
| `VITE_MAX_UPLOAD_MB` | `10` | Client-side size shown in validation messages — keep ≤ `MAX_UPLOAD_SIZE_MB` |

## API reference

Base path: `/api/v1`

### `GET /health`

```json
{ "status": "ok", "webp_enabled": true, "avif_enabled": true }
```

### `POST /media/upload`

`multipart/form-data`:

| Field | Type | Notes |
|---|---|---|
| `file` | file | required |
| `width` | int | `0` keeps original size |
| `quality` | int | default `80` |
| `formats` | string (repeated) | e.g. `formats=webp&formats=avif`; defaults to all four |

**Response** — `202 Accepted`, a `JobResponse`:

```json
{
  "job_id": "b6f1...",
  "file_name": "photo.jpg",
  "original_size_bytes": 4213556,
  "status": "pending",
  "outputs": [],
  "width": 1200,
  "quality": 80,
  "formats": ["jpeg", "webp", "avif", "png"],
  "created_at": "2026-08-17T08:00:00Z",
  "updated_at": "2026-08-17T08:00:00Z"
}
```

### `GET /jobs/:id`

Same `JobResponse` shape, reflecting current status. Once `status` is
`completed`, `outputs` is populated:

```json
{
  "job_id": "b6f1...",
  "status": "completed",
  "outputs": [
    { "format": "webp", "url": "https://res.cloudinary.com/...", "size_bytes": 812044, "width": 1200, "height": 800 }
  ],
  "...": "..."
}
```

`status` is one of `pending`, `processing`, `completed`, `failed`. On
`failed`, `error` holds a human-readable message.

Errors follow `{ "error": "message" }` (`dto.ErrorResponse`).

## Frontend notes

- `src/hooks/useJobPolling.js` polls `GET /jobs/:id` every 900ms until a
  terminal status, with exponential-free retry on transient errors.
- `src/pages/Upload.jsx` drives the flow: dropzone → options → submit →
  progress → results, plus a local "today's uploads" history
  (`useUploadHistory`, persisted to `localStorage`).
- The four result cards (`ResultsPanel` → `ResultCard`) render once
  `job.status === 'completed'`, one per entry in `job.outputs`.
- `src/api/client.js` centralizes all backend calls and resolves
  Cloudinary/relative download URLs via `resolveDownloadUrl`.

## Known issues / housekeeping

- **`.gitignore` and `docker-compose.yml` appear to have been swapped** —
  `.gitignore` currently contains the Compose service definitions, and
  `docker-compose.yml` is empty. Swap their contents back before relying on
  `docker compose up`.
- **Committed credentials**: the root `.env` and `backend/.env` in this repo
  currently contain live Cloudinary API keys/secrets. Rotate them in the
  Cloudinary console and make sure `.env` files are actually excluded from
  version control (once `.gitignore` is fixed) before pushing further.
- `internal/delivery/dto/` is the source of truth for the API's JSON shape —
  if you add a field the frontend needs, add it to `dto.JobResponse` /
  `dto.OutputResponse` and the corresponding mapper in `dto/mapper.go`,
  rather than adding JSON tags back onto the domain structs.