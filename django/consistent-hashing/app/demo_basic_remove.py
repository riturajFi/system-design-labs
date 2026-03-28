from basic_consistent_router import BasicConsistentHashRouter
from rebalance import print_server_keys, rebalance_keys


if __name__ == "__main__":
    keys = [f"key{i}" for i in range(20)]

    print("=== BASIC CONSISTENT HASHING: REMOVE SERVER WITH MIGRATION ===")

    old_router = BasicConsistentHashRouter(["redis1", "redis2", "redis3", "redis4"])
    old_router.flush_all()

    print("\nWriting keys using OLD topology (4 servers):")
    for key in keys:
        server = old_router.set(key, f"value-for-{key}")
        print(f"{key} -> {server}")

    print_server_keys(old_router, "REDIS STATE BEFORE REMOVING redis4")

    new_router = BasicConsistentHashRouter(["redis1", "redis2", "redis3"])
    result = rebalance_keys(old_router, new_router)

    print("\nActual migration complete.")
    print(f"Total keys seen: {result['total_seen']}")
    print(f"Moved keys: {result['moved_count']}")
    for item in result["moved"]:
        print(item)

    print_server_keys(new_router, "REDIS STATE AFTER MIGRATION TO 3-SERVER TOPOLOGY")

    print("\nReading all keys using NEW topology:")
    for key in keys:
        server, value = new_router.get(key)
        print(f"{key} -> {server} -> {value}")
