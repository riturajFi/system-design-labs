from __future__ import annotations

from internal.demo.rebalance import moved_keys, print_server_keys, rebalance_keys
from internal.routing.basic import BasicConsistentHashRouter


def main() -> None:
    keys = [f"key{i}" for i in range(20)]

    print("=== BASIC CONSISTENT HASHING WITH MIGRATION ===")

    old_router = BasicConsistentHashRouter(["redis1", "redis2", "redis3"])
    old_router.flush_all()
    old_router.print_ring()

    print("\nWriting keys using OLD topology (3 servers):")
    for key in keys:
        server = old_router.set(key, f"value-for-{key}")
        print(f"{key} -> {server}")

    print_server_keys(old_router, "REDIS STATE BEFORE ADDING redis4")

    new_router = BasicConsistentHashRouter(["redis1", "redis2", "redis3", "redis4"])
    new_router.print_ring()

    predicted_moves = moved_keys(old_router, new_router, keys)
    print(f"\nPredicted moved keys after adding redis4: {len(predicted_moves)} / {len(keys)}")
    for item in predicted_moves:
        print(item)

    result = rebalance_keys(old_router, new_router)

    print("\nActual migration complete.")
    print(f"Total keys seen: {result['total_seen']}")
    print(f"Moved keys: {result['moved_count']}")
    for item in result["moved"]:
        print(item)

    print_server_keys(new_router, "REDIS STATE AFTER MIGRATION TO NEW TOPOLOGY")

    print("\nReading all keys using NEW topology:")
    for key in keys:
        server, value = new_router.get(key)
        print(f"{key} -> {server} -> {value}")


if __name__ == "__main__":
    main()
