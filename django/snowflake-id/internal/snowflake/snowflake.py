from __future__ import annotations

from dataclasses import dataclass
from datetime import datetime, timezone
import threading
import time


class InvalidNodeID(ValueError):
    pass


class ClockWentBack(RuntimeError):
    pass


class EpochInFuture(RuntimeError):
    pass


_TIMESTAMP_BITS = 41
_NODE_BITS = 10
_SEQ_BITS = 12

_MAX_NODE = (1 << _NODE_BITS) - 1
_MAX_SEQ = (1 << _SEQ_BITS) - 1

_NODE_SHIFT = _SEQ_BITS
_TIME_SHIFT = _NODE_BITS + _SEQ_BITS

_DEFAULT_EPOCH = datetime(2024, 1, 1, 0, 0, 0, 0, tzinfo=timezone.utc)


@dataclass
class Generator:
    node_id: int
    epoch: datetime = _DEFAULT_EPOCH

    def __post_init__(self) -> None:
        if self.node_id < 0 or self.node_id > _MAX_NODE:
            raise InvalidNodeID("invalid node id")
        if self.epoch.tzinfo is None:
            # Keep it simple: force timezone-aware datetimes (UTC).
            self.epoch = self.epoch.replace(tzinfo=timezone.utc)
        self.epoch = self.epoch.astimezone(timezone.utc)

        now = datetime.now(tz=timezone.utc)
        if self.epoch > now:
            raise EpochInFuture("epoch is in the future")

        self._lock = threading.Lock()
        self._last_ms = -1
        self._seq = 0

    def next_id(self) -> int:
        with self._lock:
            now_ms = int((datetime.now(tz=timezone.utc) - self.epoch).total_seconds() * 1000)
            if now_ms < 0:
                raise EpochInFuture("epoch is in the future")

            if now_ms < self._last_ms:
                raise ClockWentBack("clock moved backwards")

            if now_ms == self._last_ms:
                self._seq = (self._seq + 1) & _MAX_SEQ
                if self._seq == 0:
                    # We used up 0..4095 within the same millisecond; wait for the next ms.
                    while now_ms <= self._last_ms:
                        time.sleep(0.0002)  # 200 microseconds
                        now_ms = int((datetime.now(tz=timezone.utc) - self.epoch).total_seconds() * 1000)
            else:
                self._seq = 0

            self._last_ms = now_ms
            ts = now_ms & ((1 << _TIMESTAMP_BITS) - 1)

            return (ts << _TIME_SHIFT) | (self.node_id << _NODE_SHIFT) | self._seq

