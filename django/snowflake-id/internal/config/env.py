from __future__ import annotations

import os


def must_get_int_env(name: str) -> int:
    raw = os.getenv(name)
    if raw is None or raw == "":
        raise RuntimeError(f"{name} is required")
    try:
        return int(raw)
    except ValueError as exc:
        raise RuntimeError(f"{name} must be an integer") from exc
