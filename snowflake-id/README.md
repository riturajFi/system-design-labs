# Snowflake ID (Go + Docker)

This repo contains a simple **Snowflake-style ID generator** implemented as:

1) a **Go library** (`internal/snowflake`) that generates 64-bit IDs
2) a **Go HTTP server** (`cmd/idgen`) that returns a new ID on demand
3) a **Docker image** (via `Dockerfile`) that runs the server in a container
4) a **docker-compose** setup that runs multiple containers (each container = one “machine”)

The key idea: **each container behaves like a separate machine** because we give it a different `MACHINE_ID`.

---

## Project Structure

```
internal/snowflake/   # Library (ID generator). No Docker/HTTP knowledge.
cmd/idgen/            # Runnable program (HTTP server) that uses the library.
Dockerfile            # Builds a small runtime image that runs the server.
docker-compose.yml    # Runs multiple containers with different MACHINE_IDs.
```

---

## How the Snowflake ID Is Built (Bit Layout)

Each generated ID is a `uint64` built from 3 parts:

```
41 bits timestamp(ms) | 10 bits machine_id | 12 bits sequence
```

- **timestamp(ms)**: milliseconds since a custom epoch (default: `2024-01-01T00:00:00Z`)
- **machine_id**: identifies which machine/container generated the ID (`0..1023`)
- **sequence**: counter for IDs generated in the *same* millisecond (`0..4095`)

Why this works:
- IDs are time-ordered (mostly increasing) because timestamp is the high bits.
- Multiple machines can generate IDs without collisions as long as they use different `machine_id`s.
- Within one machine, sequence prevents collisions when generating many IDs in the same millisecond.

---

## Library: `internal/snowflake`

Files:
- `internal/snowflake/snowflake.go`

Main API:
- `snowflake.New(nodeID uint16) (*Generator, error)`
- `(*Generator).NextID() (uint64, error)`

Important behavior:
- `nodeID` must be `0..1023` or `New` returns an error.
- `NextID()` is safe to call concurrently; it uses a mutex to protect internal state.
- If the system clock moves backwards (time goes earlier than the last generated timestamp), `NextID()` returns an error.
- If the generator produces more than 4096 IDs in one millisecond, it waits for the next millisecond.

---

## Server: `cmd/idgen`

This is the runnable program. Instead of printing IDs in a loop, it runs an HTTP server and generates IDs **when requested**.

### Environment variables

- `MACHINE_ID` (required): `0..1023`
- `PORT` (optional): server port (default `8080`)

### Endpoints

- `GET /health`
  - Returns: `ok`
- `GET /id`
  - Returns JSON:
    ```json
    {"id": 275238510129057792, "machine_id": 1}
    ```

---

## Run Locally (No Docker)

From the repo root:

```bash
MACHINE_ID=1 PORT=8080 go run ./cmd/idgen
```

In another terminal:

```bash
curl -s localhost:8080/health
curl -s localhost:8080/id
curl -s localhost:8080/id
```

### If `go test ./...` fails with a Go build cache permission error

Some environments have a restricted default Go cache directory. If you see a permission error like:

```
permission denied ... ~/.cache/go-build/...
```

Run commands with:

```bash
mkdir -p /tmp/go-build
GOCACHE=/tmp/go-build go test ./...
GOCACHE=/tmp/go-build go run ./cmd/idgen
```

---

## Run With Docker (Single Container)

Build the image:

```bash
docker build -t snowflake-id .
```

Run the container (map host port 8080 to container port 8080):

```bash
docker run --rm -p 8080:8080 -e MACHINE_ID=1 -e PORT=8080 snowflake-id
```

Call the server:

```bash
curl -s localhost:8080/id
```

---

## Run Multiple “Machines” With Docker Compose

`docker-compose.yml` starts 3 containers, each with a different `MACHINE_ID`, and exposes them on different host ports:

- machine 1: host `8081` → container `8080`
- machine 2: host `8082` → container `8080`
- machine 3: host `8083` → container `8080`

Start them:

```bash
docker compose up --build
```

Query each machine:

```bash
curl -s localhost:8081/id
curl -s localhost:8082/id
curl -s localhost:8083/id
```

You should see different `machine_id` values in responses, and no collisions between machines.

---

## (Optional) Decode an ID (See the Parts)

Given an ID:

- `sequence = id & ((1<<12)-1)`
- `machine  = (id >> 12) & ((1<<10)-1)`
- `ts_ms    = id >> 22` (milliseconds since the custom epoch)

Quick Python helper:

```bash
ID=275238510129057792 python3 - <<'PY'
import os
id = int(os.environ["ID"])
seq = id & ((1<<12)-1)
machine = (id >> 12) & ((1<<10)-1)
ts_ms = id >> 22
print("id:", id)
print("machine_id:", machine)
print("sequence:", seq)
print("timestamp_ms_since_epoch:", ts_ms)
PY
```

---

## Notes

- The generator’s custom epoch is hardcoded in `internal/snowflake/snowflake.go`.
- This design is intentionally minimal: no persistence, no clustering, no leader election—just a clean demonstration of Snowflake-style IDs locally and via containers.
