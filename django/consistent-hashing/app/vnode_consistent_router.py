import bisect
import hashlib
from collections import Counter

import redis


class VNodeConsistentHashRouter:
    """
    Advanced consistent hashing with virtual nodes.

    Why virtual nodes:
    - One server point on the ring can create uneven partitions.
    - Many virtual points per server make distribution more balanced.
    """

    def __init__(self, server_names, virtual_nodes=100):
        self.server_names = server_names
        self.virtual_nodes = virtual_nodes

        self.clients = {
            "redis1": redis.Redis(host="redis1", port=6379, decode_responses=True),
            "redis2": redis.Redis(host="redis2", port=6379, decode_responses=True),
            "redis3": redis.Redis(host="redis3", port=6379, decode_responses=True),
            "redis4": redis.Redis(host="redis4", port=6379, decode_responses=True),
        }

        self.ring_positions = []
        self.ring_map = {}

        self._build_ring()

    def _hash(self, value: str) -> int:
        return int(hashlib.sha1(value.encode()).hexdigest(), 16)

    def _build_ring(self):
        for server in self.server_names:
            for i in range(self.virtual_nodes):
                vnode_name = f"{server}#{i}"
                position = self._hash(vnode_name)
                self.ring_positions.append(position)
                self.ring_map[position] = server

        self.ring_positions.sort()

    def locate_server(self, key: str) -> str:
        key_position = self._hash(key)
        idx = bisect.bisect(self.ring_positions, key_position)

        if idx == len(self.ring_positions):
            idx = 0

        server_position = self.ring_positions[idx]
        return self.ring_map[server_position]

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

    def print_ring_sample(self, sample_size=20):
        print(f"\n=== VNODE RING SAMPLE (first {sample_size}) ===")
        for pos in self.ring_positions[:sample_size]:
            print(f"{str(pos)[:12]} -> {self.ring_map[pos]}")

    def flush_all(self):
        for server_name in self.clients:
            self.clients[server_name].flushdb()
