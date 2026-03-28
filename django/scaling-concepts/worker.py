from __future__ import annotations

import threading
import time

from cache_store import cache_key, redis_client
from config import SYNC_INTERVAL_SECONDS
from storage import primary_db


_worker_started = False
_worker_lock = threading.Lock()


def sync_cache_once(limit: int = 500) -> int:
    client = redis_client()
    with primary_db() as conn:
        with conn.cursor() as cur:
            cur.execute(
                """
                SELECT id, code, long_url
                FROM short_link_log
                WHERE applied_at IS NULL
                ORDER BY id
                LIMIT %s
                """,
                (limit,),
            )
            rows = cur.fetchall()

        processed = 0
        for log_id, code, long_url in rows:
            client.set(cache_key(code), long_url)
            with conn.cursor() as cur:
                cur.execute(
                    "UPDATE short_link_log SET applied_at = NOW() WHERE id = %s AND applied_at IS NULL",
                    (log_id,),
                )
            processed += 1

        conn.commit()
    return processed


def sync_cache_forever() -> None:
    lock_key = "shortener:sync-lock"
    client = redis_client()
    while True:
        try:
            if client.set(lock_key, "1", nx=True, ex=SYNC_INTERVAL_SECONDS):
                try:
                    sync_cache_once()
                finally:
                    client.delete(lock_key)
        except Exception:
            pass
        time.sleep(SYNC_INTERVAL_SECONDS)


def start_background_workers() -> None:
    global _worker_started

    with _worker_lock:
        if _worker_started:
            return
        threading.Thread(target=sync_cache_forever, daemon=True).start()
        _worker_started = True
