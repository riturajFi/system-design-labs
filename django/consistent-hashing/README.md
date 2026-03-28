# Consistent Hashing Demo

Small Docker-based learning project for:

- modulo routing with `% N`
- basic consistent hashing
- consistent hashing with virtual nodes
- manual key migration after topology changes

## Project Layout

```text
consistent-hashing/
├── docker-compose.yml
└── app/
    ├── Dockerfile
    ├── requirements.txt
    ├── modulo_router.py
    ├── basic_consistent_router.py
    ├── vnode_consistent_router.py
    ├── rebalance.py
    ├── demo_modulo.py
    ├── demo_basic.py
    ├── demo_vnodes.py
    └── demo_basic_remove.py
```

## Run

From [django/consistent-hashing](/home/riturajtripathy/Documents/_Code/personal_projects/Learning Project/system-design-labs/django/consistent-hashing):

```bash
docker compose up --build -d
docker compose exec router bash
```

Then inside the router container:

```bash
python demo_modulo.py
python demo_basic.py
python demo_vnodes.py
python demo_basic_remove.py
```

## Manual Redis Inspection

```bash
docker compose exec redis1 redis-cli KEYS '*'
docker compose exec redis2 redis-cli KEYS '*'
docker compose exec redis3 redis-cli KEYS '*'
docker compose exec redis4 redis-cli KEYS '*'
```

## Notes

- These demos focus on routing and simple migration mechanics.
- Migration currently handles string values only.
- TTLs, non-string Redis types, concurrent writes, and cutover safety are intentionally out of scope.
