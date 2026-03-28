from __future__ import annotations

from internal.routing.common import RouterBase, hash_value


class BasicConsistentHashRouter(RouterBase):
    def __init__(self, server_names: list[str]) -> None:
        super().__init__(server_names)
        self.ring_positions: list[int] = []
        self.ring_map: dict[int, str] = {}
        self._build_ring()

    def _build_ring(self) -> None:
        self.ring_positions.clear()
        self.ring_map.clear()
        for server in self.server_names:
            position = hash_value(server)
            self.ring_positions.append(position)
            self.ring_map[position] = server
        self.ring_positions.sort()

    def add_node(self, server_name: str) -> None:
        if server_name in self.server_names:
            return
        super().add_node(server_name)
        position = hash_value(server_name)
        self.ring_map[position] = server_name
        self.ring_positions.append(position)
        self.ring_positions.sort()

    def remove_node(self, server_name: str) -> None:
        if server_name not in self.server_names:
            return
        super().remove_node(server_name)
        position = hash_value(server_name)
        self.ring_map.pop(position, None)
        self.ring_positions = [p for p in self.ring_positions if p != position]

    def locate_server(self, key: str) -> str:
        key_position = hash_value(key)
        for server_position in self.ring_positions:
            if server_position >= key_position:
                return self.ring_map[server_position]

        return self.ring_map[self.ring_positions[0]]

    def print_ring(self) -> None:
        print("\n=== BASIC RING ===")
        for position in self.ring_positions:
            print(f"{str(position)[:12]} -> {self.ring_map[position]}")
