def scan_keys(client, pattern="*"):
    """
    Small helper around Redis SCAN.
    """

    keys = []
    cursor = 0

    while True:
        cursor, batch = client.scan(cursor=cursor, match=pattern, count=100)
        keys.extend(batch)
        if cursor == 0:
            break

    return keys


def collect_all_keys(router, pattern="*"):
    """
    Collect all keys from all active Redis servers in this router.
    """

    result = {}
    for server_name in router.server_names:
        result[server_name] = scan_keys(router.clients[server_name], pattern=pattern)
    return result


def rebalance_keys(old_router, new_router, pattern="*"):
    """
    Rebalance keys from old topology to new topology.

    Assumptions:
    - string values only
    - keys currently live on servers chosen by old_router
    """

    moved = []
    total_seen = 0

    keys_by_server = collect_all_keys(old_router, pattern=pattern)

    for scanned_server, keys in keys_by_server.items():
        for key in keys:
            total_seen += 1

            old_server = old_router.locate_server(key)
            new_server = new_router.locate_server(key)

            if old_server == new_server:
                continue

            old_client = old_router.clients[old_server]
            new_client = new_router.clients[new_server]
            value = old_client.get(key)

            if value is None:
                continue

            new_client.set(key, value)
            old_client.delete(key)

            moved.append((key, old_server, new_server))

    return {
        "total_seen": total_seen,
        "moved_count": len(moved),
        "moved": moved,
    }


def print_server_keys(router, title):
    print(f"\n=== {title} ===")
    for server_name in router.server_names:
        keys = sorted(scan_keys(router.clients[server_name]))
        print(f"{server_name}: {keys}")
