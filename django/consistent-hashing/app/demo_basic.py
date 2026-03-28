from basic_consistent_router import BasicConsistentHashRouter
from rebalance import print_server_keys, rebalance_keys


def moved_keys(router_a, router_b, keys):
    moved = []
    for key in keys:
        old_server = router_a.locate_server(key)
        new_server = router_b.locate_server(key)
        if old_server != new_server:
            moved.append((key, old_server, new_server))
    return moved


if __name__ == "__main__":
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
