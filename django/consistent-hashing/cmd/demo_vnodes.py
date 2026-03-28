from __future__ import annotations

from internal.demo.rebalance import moved_keys, print_server_keys, rebalance_keys
from internal.routing.vnode import VNodeConsistentHashRouter


def main() -> None:
    keys = [f"key{i}" for i in range(100)]

    print("=== VNODE CONSISTENT HASHING WITH MIGRATION ===")

    old_router = VNodeConsistentHashRouter(["redis1", "redis2", "redis3"], virtual_nodes=50)
    old_router.flush_all()
    old_router.print_ring_sample()

    print("\nWriting keys using OLD vnode topology (3 servers):")
    for key in keys:
        old_router.set(key, f"value-for-{key}")

    print("\nDistribution with OLD topology:")
    print(old_router.distribution(keys))
    print_server_keys(old_router, "REDIS STATE BEFORE ADDING redis4")

    new_router = VNodeConsistentHashRouter(["redis1", "redis2", "redis3", "redis4"], virtual_nodes=50)

    print("\nDistribution with NEW topology:")
    print(new_router.distribution(keys))

    predicted_moves = moved_keys(old_router, new_router, keys)
    print(f"\nPredicted moved keys after adding redis4: {len(predicted_moves)} / {len(keys)}")
    print("First 20 predicted moves:")
    for item in predicted_moves[:20]:
        print(item)

    result = rebalance_keys(old_router, new_router)

    print("\nActual migration complete.")
    print(f"Total keys seen: {result['total_seen']}")
    print(f"Moved keys: {result['moved_count']}")
    print("First 20 actual moves:")
    for item in result["moved"][:20]:
        print(item)

    print_server_keys(new_router, "REDIS STATE AFTER MIGRATION TO NEW TOPOLOGY")

    print("\nReading first 20 keys using NEW topology:")
    for key in keys[:20]:
        server, value = new_router.get(key)
        print(f"{key} -> {server} -> {value}")


if __name__ == "__main__":
    main()
