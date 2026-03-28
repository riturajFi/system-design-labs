from __future__ import annotations

from internal.routing.common import RouterBase, hash_value


class ModuloRouter(RouterBase):
    def locate_server(self, key: str) -> str:
        index = hash_value(key) % len(self.server_names)
        return self.server_names[index]
