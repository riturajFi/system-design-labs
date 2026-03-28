import hashlib
from collections import Counter

import redis


class ModuloRouter:
    """
    Very simple router:
        server_index = hash(key) % number_of_servers

    This is here only for comparison.
    This is the approach that causes large remapping when server count changes.
    """

    def __init__(self, server_names):
        self.server_names = server_names
        self.clients = {
            "redis1": redis.Redis(host="redis1", port=6379, decode_responses=True),
            "redis2": redis.Redis(host="redis2", port=6379, decode_responses=True),
            "redis3": redis.Redis(host="redis3", port=6379, decode_responses=True),
            "redis4": redis.Redis(host="redis4", port=6379, decode_responses=True),
        }

    def _hash(self, key: str) -> int:
        return int(hashlib.sha1(key.encode()).hexdigest(), 16)

    def locate_server(self, key: str) -> str:
        index = self._hash(key) % len(self.server_names)
        return self.server_names[index]

    def set(self, key: str, value: str):
        server = self.locate_server(key)
        self.clients[server].set(key, value)
        return server

    def get(self, key: str):
        server = self.locate_server(key)
        value = self.clients[server].get(key)
        return server, value

    def distribution(self, keys):
        count = Counter()
        for key in keys:
            count[self.locate_server(key)] += 1
        return dict(count)
