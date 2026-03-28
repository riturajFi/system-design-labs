# Scaling Concepts Lab

Minimal Python + Postgres demo.

What it has:

- `server.py` for the HTTP server, database logic, and cache sync loop
- `cmd/server/main.py` as the thin entrypoint
- `internal/` modules for HTTP, storage, cache, config, and background sync
- one Postgres primary
- two Postgres replicas
- one Redis cache
- HAProxy in front of two Python app containers
- HAProxy also balances read-only DB traffic across replicas

## Endpoints

- `GET /health`
- `POST /shorten` with `long_url`
- `GET /resolve?code=...`
- `GET /state`

## Run

```bash
docker compose up --build
```

Open `http://127.0.0.1:8000`.

## Layout

```text
scaling-concepts/
├── cmd/server/main.py
├── internal/cache/
├── internal/config/
├── internal/http/
├── internal/shortener/
└── internal/worker/
```

## Behavior

- writes go to the primary only
- reads try Redis first, then the HAProxy DB read pool, then the primary
- cache updates come from an append-only DB log table via a background thread started by the server itself
