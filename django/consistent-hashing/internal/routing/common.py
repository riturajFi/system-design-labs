from __future__ import annotations

import hashlib
from collections import Counter

from internal.storage.redis_store import build_clients


def hash_value(value: str) -> int:
    return int(hashlib.sha1(value.encode()).hexdigest(), 16)


class RouterBase:
    def __init__(self, server_names: list[str]) -> None:
        self.server_names = server_names
        self.clients = build_clients()

    def locate_server(self, key: str) -> str:
        raise NotImplementedError

    def add_node(self, server_name: str) -> None:
        if server_name not in self.server_names:
            self.server_names.append(server_name)

    def remove_node(self, server_name: str) -> None:
        if server_name in self.server_names:
            self.server_names.remove(server_name)

    def set(self, key: str, value: str) -> str:
        server = self.locate_server(key)
        self.clients[server].set(key, value)
        return server

    def get(self, key: str) -> tuple[str, str | None]:
        server = self.locate_server(key)
        value = self.clients[server].get(key)
        return server, value

    def distribution(self, keys: list[str]) -> dict[str, int]:
        counts: Counter[str] = Counter()
        for key in keys:
            counts[self.locate_server(key)] += 1
        return dict(counts)

    def flush_all(self) -> None:
        for server_name in self.clients:
            self.clients[server_name].flushdb()
