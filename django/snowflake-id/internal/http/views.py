from __future__ import annotations

from django.http import HttpResponse, JsonResponse

from internal.service.id_service import next_id
from internal.snowflake.snowflake import ClockWentBack, EpochInFuture, InvalidNodeID


def health(_request):
    return HttpResponse("ok", content_type="text/plain")


def new_id(_request):
    try:
        payload = next_id()
    except (InvalidNodeID, EpochInFuture, ClockWentBack) as exc:
        return JsonResponse({"error": str(exc)}, status=500)
    return JsonResponse(payload)
