from __future__ import annotations

from internal.demo.rebalance import moved_keys
from internal.routing.modulo import ModuloRouter


def main() -> None:
    keys = [f"key{i}" for i in range(30)]
    router_3 = ModuloRouter(["redis1", "redis2", "redis3"])
    router_4 = ModuloRouter(["redis1", "redis2", "redis3", "redis4"])

    print("=== MODULO HASHING DEMO ===")
    print("\nDistribution with 3 servers:")
    print(router_3.distribution(keys))

    print("\nDistribution with 4 servers:")
    print(router_4.distribution(keys))

    moved = moved_keys(router_3, router_4, keys)
    print(f"\nKeys moved after adding redis4: {len(moved)} / {len(keys)}")
    for item in moved[:10]:
        print(item)


if __name__ == "__main__":
    main()
