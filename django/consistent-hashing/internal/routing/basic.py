from __future__ import annotations

import bisect

from internal.routing.common import RouterBase, hash_value


class BasicConsistentHashRouter(RouterBase):
    def __init__(self, server_names: list[str]) -> None:
        super().__init__(server_names)
        self.ring_positions: list[int] = []
        self.ring_map: dict[int, str] = {}
        self._build_ring()

    def _build_ring(self) -> None:
        for server in self.server_names:
            position = hash_value(server)
            self.ring_positions.append(position)
            self.ring_map[position] = server
        self.ring_positions.sort()

    def locate_server(self, key: str) -> str:
        key_position = hash_value(key)
        index = bisect.bisect(self.ring_positions, key_position)
        if index == len(self.ring_positions):
            index = 0
        server_position = self.ring_positions[index]
        return self.ring_map[server_position]

    def print_ring(self) -> None:
        print("\n=== BASIC RING ===")
        for position in self.ring_positions:
            print(f"{str(position)[:12]} -> {self.ring_map[position]}")
