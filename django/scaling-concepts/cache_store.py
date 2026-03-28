from __future__ import annotations

import redis

from config import REDIS_URL


def cache_key(code: str) -> str:
    return f"shortlink:{code}"


def redis_client() -> redis.Redis:
    return redis.from_url(REDIS_URL, decode_responses=True)


def get_cached_long_url(code: str) -> str | None:
    value = redis_client().get(cache_key(code))
    return value if value else None
