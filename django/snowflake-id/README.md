# Snowflake ID (Django + Docker)

This folder recreates the Go `snowflake-id` project using **Django**.

You get:

1) a small **Snowflake-style ID generator** (pure Python code)
2) a **Django HTTP server** that returns a new ID on demand
3) Docker + docker-compose so you can run multiple containers, where each container acts like a separate machine

The key idea: each container is given a different `MACHINE_ID`, so IDs do not collide across containers.

---

## How To Study This Repo (Recommended Reading Flow)

If you are a beginner, follow this order.

### Step 1: Understand the ID generator first (core logic)

Read: `internal/snowflake/snowflake.py`

Look for:
- what state the generator keeps (`last_ms`, `seq`, `node_id`, `epoch`)
- how it builds the final number: timestamp + machine_id + sequence
- why there is a lock (many requests can happen at the same time)

### Step 2: Read the Django server (how requests become IDs)

Read:
- `idgen_app/views.py` (endpoints)
- `idgen_project/urls.py` (URL routing)

Look for:
- how `MACHINE_ID` is read from the environment
- how `/id` calls the generator and returns JSON

### Step 3: Read Docker last (packaging + "machines")

Read:
- `Dockerfile.dev` (runs `python manage.py runserver`, like local dev)
- `Dockerfile.simple` (runs gunicorn in a single image)
- `Dockerfile.prod` / `Dockerfile` (multi-stage; small-ish runtime)
- `docker-compose.yml` (starts 3 containers with different `MACHINE_ID`s)

---

## Endpoints

- `GET /health` -> `ok`
- `GET /id` -> JSON:
  ```json
  {"id": 275238510129057792, "machine_id": 1}
  ```

---

## Run Locally (No Docker)

From `django/snowflake-id/`:

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

MACHINE_ID=1 PORT=8000 python3 manage.py runserver 0.0.0.0:8000
```

In another terminal:

```bash
curl -s localhost:8000/health
curl -s localhost:8000/id
```

---

## Run With Docker Compose (3 "Machines")

Pick which Dockerfile to use by setting `DOCKERFILE`:

- `Dockerfile.dev`    (slow start, easiest to understand)
- `Dockerfile.simple` (single-stage, runs gunicorn)
- `Dockerfile.prod`   (multi-stage, smaller runtime image)

Example:

```bash
DOCKERFILE=Dockerfile.dev docker compose up --build
```

Then query each machine:

```bash
curl -s localhost:8091/id
curl -s localhost:8092/id
curl -s localhost:8093/id
```

Stop:

```bash
docker compose down
```
