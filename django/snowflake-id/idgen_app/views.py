import os

from django.http import HttpResponse, JsonResponse

from internal.snowflake.snowflake import ClockWentBack, EpochInFuture, Generator, InvalidNodeID


def _must_get_int_env(name: str) -> int:
    raw = os.getenv(name)
    if raw is None or raw == "":
        raise RuntimeError(f"{name} is required")
    try:
        return int(raw)
    except ValueError as e:
        raise RuntimeError(f"{name} must be an integer") from e


# Create one generator for this whole process. It keeps state (sequence, last_ms).
MACHINE_ID = _must_get_int_env("MACHINE_ID")
_GENERATOR = Generator(MACHINE_ID)


def health(_request):
    return HttpResponse("ok", content_type="text/plain")


def new_id(_request):
    try:
        new_id_val = _GENERATOR.next_id()
    except (InvalidNodeID, EpochInFuture, ClockWentBack) as e:
        return JsonResponse({"error": str(e)}, status=500)
    return JsonResponse({"id": new_id_val, "machine_id": MACHINE_ID})

