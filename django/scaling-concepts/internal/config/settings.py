PRIMARY_DSN = "postgresql://app:app_password@db-primary:5432/scaling_concepts"
READ_DSN = "postgresql://app:app_password@lb:5433/scaling_concepts"
REDIS_URL = "redis://redis:6379/0"
SYNC_INTERVAL_SECONDS = 30

REPLICA_DSNS = [
    ("replica_1", "postgresql://app:app_password@db-replica-1:5432/scaling_concepts"),
    ("replica_2", "postgresql://app:app_password@db-replica-2:5432/scaling_concepts"),
]
