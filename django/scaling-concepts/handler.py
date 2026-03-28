from __future__ import annotations

import json
from http.server import BaseHTTPRequestHandler
from urllib.parse import parse_qs, urlparse

from storage import cluster_state, create_short_link, resolve_short_link


def json_response(handler: BaseHTTPRequestHandler, status: int, payload: dict[str, object]) -> None:
    body = json.dumps(payload).encode("utf-8")
    handler.send_response(status)
    handler.send_header("Content-Type", "application/json")
    handler.send_header("Content-Length", str(len(body)))
    handler.end_headers()
    handler.wfile.write(body)


class AppHandler(BaseHTTPRequestHandler):
    def do_GET(self) -> None:
        parsed = urlparse(self.path)
        if parsed.path == "/":
            json_response(self, 200, {"name": "scaling-concepts", "endpoints": ["/health", "/shorten", "/resolve?code=...", "/state"]})
            return
        if parsed.path == "/health":
            body = b"ok"
            self.send_response(200)
            self.send_header("Content-Type", "text/plain")
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)
            return
        if parsed.path == "/resolve":
            code = parse_qs(parsed.query).get("code", [""])[0]
            try:
                json_response(self, 200, resolve_short_link(code))
            except ValueError as exc:
                json_response(self, 400, {"error": str(exc)})
            except LookupError as exc:
                json_response(self, 404, {"error": str(exc)})
            return
        if parsed.path == "/state":
            json_response(self, 200, cluster_state())
            return

        json_response(self, 404, {"error": "not found"})

    def do_POST(self) -> None:
        parsed = urlparse(self.path)
        if parsed.path != "/shorten":
            json_response(self, 404, {"error": "not found"})
            return

        length = int(self.headers.get("Content-Length", "0"))
        raw = self.rfile.read(length).decode("utf-8") if length else "{}"
        try:
            payload = json.loads(raw)
        except json.JSONDecodeError:
            json_response(self, 400, {"error": "invalid JSON body"})
            return

        try:
            result = create_short_link(payload.get("long_url", ""))
            json_response(self, 201, result)
        except ValueError as exc:
            json_response(self, 400, {"error": str(exc)})
        except Exception as exc:
            json_response(self, 503, {"error": str(exc)})

    def log_message(self, format: str, *args) -> None:
        return
