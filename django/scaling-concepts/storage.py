from __future__ import annotations

import secrets
from urllib.parse import urlparse

import psycopg

from cache_store import get_cached_long_url
from config import PRIMARY_DSN, READ_DSN, REPLICA_DSNS


def normalize_long_url(raw: str) -> str:
    value = raw.strip()
    if not value:
        raise ValueError("long_url is required")

    parsed = urlparse(value)
    if not parsed.scheme:
        value = f"https://{value}"
        parsed = urlparse(value)

    if parsed.scheme not in {"http", "https"} or not parsed.netloc:
        raise ValueError("long_url must be an http or https URL")

    return value


def cache_key(code: str) -> str:
    return f"shortlink:{code}"


def primary_db() -> psycopg.Connection:
    return psycopg.connect(PRIMARY_DSN)


def read_db() -> psycopg.Connection:
    return psycopg.connect(READ_DSN)


def ensure_schema() -> None:
    statements = [
        """
        CREATE TABLE IF NOT EXISTS short_links (
            id BIGSERIAL PRIMARY KEY,
            code VARCHAR(16) UNIQUE NOT NULL,
            long_url TEXT NOT NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
        )
        """,
        """
        CREATE TABLE IF NOT EXISTS short_link_log (
            id BIGSERIAL PRIMARY KEY,
            code VARCHAR(16) NOT NULL,
            long_url TEXT NOT NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
            applied_at TIMESTAMPTZ NULL
        )
        """,
        "CREATE INDEX IF NOT EXISTS short_links_code_idx ON short_links (code)",
        "CREATE INDEX IF NOT EXISTS short_link_log_applied_idx ON short_link_log (applied_at)",
        "CREATE INDEX IF NOT EXISTS short_link_log_created_idx ON short_link_log (created_at)",
    ]
    with primary_db() as conn:
        with conn.cursor() as cur:
            for statement in statements:
                cur.execute(statement)
        conn.commit()


def unique_code() -> str:
    return secrets.token_hex(3)


def create_short_link(long_url: str) -> dict[str, str]:
    normalized = normalize_long_url(long_url)
    for _ in range(8):
        code = unique_code()
        try:
            with primary_db() as conn:
                with conn.cursor() as cur:
                    cur.execute(
                        "INSERT INTO short_links (code, long_url) VALUES (%s, %s)",
                        (code, normalized),
                    )
                    cur.execute(
                        "INSERT INTO short_link_log (code, long_url) VALUES (%s, %s)",
                        (code, normalized),
                    )
                conn.commit()
            return {"code": code, "long_url": normalized, "write_db": "primary"}
        except psycopg.IntegrityError:
            continue
    raise RuntimeError("could not allocate a unique short code")


def resolve_from_db(code: str) -> str | None:
    with read_db() as conn:
        with conn.cursor() as cur:
            cur.execute("SELECT long_url FROM short_links WHERE code = %s", (code,))
            row = cur.fetchone()
    return row[0] if row else None


def resolve_short_link(code: str) -> dict[str, str]:
    lookup_code = code.strip()
    if not lookup_code:
        raise ValueError("code is required")

    cached = get_cached_long_url(lookup_code)
    if cached is not None:
        return {"code": lookup_code, "long_url": cached, "served_by": "cache"}

    found = resolve_from_db(lookup_code)
    if found is None:
        raise LookupError(lookup_code)
    return {"code": lookup_code, "long_url": found, "served_by": "db-read"}


def pending_log_count() -> int:
    with primary_db() as conn:
        with conn.cursor() as cur:
            cur.execute("SELECT COUNT(*) FROM short_link_log WHERE applied_at IS NULL")
            return int(cur.fetchone()[0])


def cluster_state() -> dict[str, object]:
    with primary_db() as conn:
        with conn.cursor() as cur:
            cur.execute("SELECT COUNT(*) FROM short_links")
            primary_count = int(cur.fetchone()[0])

    replicas = []
    for alias, dsn in REPLICA_DSNS:
        try:
            with psycopg.connect(dsn) as conn:
                with conn.cursor() as cur:
                    cur.execute("SELECT COUNT(*) FROM short_links")
                    count = int(cur.fetchone()[0])
            healthy = True
        except Exception:
            count = None
            healthy = False
        replicas.append({"alias": alias, "healthy": healthy, "row_count": count})

    try:
        from cache_store import redis_client

        cache = redis_client()
        cache_state = {
            "healthy": True,
            "backend": "redis",
            "pending_log_rows": pending_log_count(),
            "keys": cache.dbsize(),
        }
    except Exception:
        cache_state = {
            "healthy": False,
            "backend": "redis",
            "pending_log_rows": pending_log_count(),
            "keys": None,
        }

    return {
        "primary": {"alias": "primary", "healthy": True, "row_count": primary_count},
        "replicas": replicas,
        "cache": cache_state,
    }
