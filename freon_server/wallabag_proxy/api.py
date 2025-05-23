import logging

from django.http import HttpRequest, HttpResponse
from ninja import Router

from .models import WallabagCredentials
from .security import WallabagProxyAuth
from .wallabag import request_wallabag

logger = logging.getLogger(f"freon_server.{__name__}")
router = Router(auth=WallabagProxyAuth(), tags=["wallabag"])


@router.api_operation(["GET", "POST", "PATCH", "PUT", "DELETE"], "/api/{path:target}")
async def forward_to_wallabag(request: HttpRequest, target: str):
    credentials: WallabagCredentials = request.auth
    path = f"/api/{target}"

    resp = await request_wallabag(
        credentials,
        request.method,
        path,
        query=request.GET,
        body=request.body,
    )
    elapsed = resp.elapsed.total_seconds() * 1000
    logger.info(
        f"wallabag -> {request.method} {path} {resp.status_code} ({elapsed:.2f} ms)"
    )

    return HttpResponse(
        resp.content,
        status=resp.status_code,
        headers={"X-Wallabag-Duration": f"{elapsed:.2f} ms"},
    )
