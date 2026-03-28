from __future__ import annotations

import time
from http.server import ThreadingHTTPServer

from handler import AppHandler
from storage import ensure_schema
from worker import start_background_workers


def wait_for_db() -> None:
    while True:
        try:
            ensure_schema()
            return
        except Exception:
            time.sleep(2)


def main() -> None:
    wait_for_db()
    start_background_workers()
    server = ThreadingHTTPServer(("0.0.0.0", 8000), AppHandler)
    server.serve_forever()


if __name__ == "__main__":
    main()
