from __future__ import annotations

from internal.config.env import must_get_int_env
from internal.snowflake.snowflake import Generator


MACHINE_ID = must_get_int_env("MACHINE_ID")
GENERATOR = Generator(MACHINE_ID)


def next_id() -> dict[str, int]:
    return {"id": GENERATOR.next_id(), "machine_id": MACHINE_ID}
