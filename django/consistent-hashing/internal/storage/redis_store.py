from __future__ import annotations

import redis


def build_clients() -> dict[str, redis.Redis]:
    return {
        "redis1": redis.Redis(host="redis1", port=6379, decode_responses=True),
        "redis2": redis.Redis(host="redis2", port=6379, decode_responses=True),
        "redis3": redis.Redis(host="redis3", port=6379, decode_responses=True),
        "redis4": redis.Redis(host="redis4", port=6379, decode_responses=True),
    }
