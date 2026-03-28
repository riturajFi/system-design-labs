from __future__ import annotations

import bisect

from internal.routing.common import RouterBase, hash_value


class VNodeConsistentHashRouter(RouterBase):
    def __init__(self, server_names: list[str], virtual_nodes: int = 100) -> None:
        super().__init__(server_names)
        self.virtual_nodes = virtual_nodes
        self.ring_positions: list[int] = []
        self.ring_map: dict[int, str] = {}
        self._build_ring()

    def _build_ring(self) -> None:
        for server in self.server_names:
            for index in range(self.virtual_nodes):
                vnode_name = f"{server}#{index}"
                position = hash_value(vnode_name)
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

    def print_ring_sample(self, sample_size: int = 20) -> None:
        print(f"\n=== VNODE RING SAMPLE (first {sample_size}) ===")
        for position in self.ring_positions[:sample_size]:
            print(f"{str(position)[:12]} -> {self.ring_map[position]}")
