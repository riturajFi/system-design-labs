from modulo_router import ModuloRouter


def moved_keys(router_a, router_b, keys):
    moved = []
    for key in keys:
        old_server = router_a.locate_server(key)
        new_server = router_b.locate_server(key)
        if old_server != new_server:
            moved.append((key, old_server, new_server))
    return moved


if __name__ == "__main__":
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
