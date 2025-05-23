import logging

import httpx
from asgiref.sync import sync_to_async
from django.conf import settings

from .models import WallabagCredentials

USER_AGENT = f"freon/{settings.VERSION}"

logger = logging.getLogger(f"freon_server.{__name__}")


class WallabagApiError(Exception):
    def __init__(self, error: httpx.HTTPStatusError):
        self._inner = error
        self.status = error.response.status_code
        self.reason = error.response.reason_phrase
        self.body = error.response.text

    def __str__(self):
        return f"WAPI error: {self.status} {self.reason}: {self.body}"

    def __repr__(self):
        return f"WallabagApiError(status={self.status}, reason={self.reason}, body={self.body})"


async def request_wallabag(
    credentials: WallabagCredentials,
    method: str,
    path: str,
    query: dict | None = None,
    body: bytes | None = None,
    payload: dict | None = None,
) -> httpx.Response:
    assert body is None or payload is None, "body and payload cannot be used together"

    headers = {
        "User-Agent": USER_AGENT,
        "Content-Type": "application/json",
    }

    if not path.endswith("/info") and not path.startswith("/oauth/"):
        if not credentials.has_valid_session:
            try:
                await refresh_session(credentials)
            except WallabagApiError as e:
                logger.error(e)
                raise
        headers["Authorization"] = f"Bearer {credentials.token.access_token}"

    async with httpx.AsyncClient() as client:
        resp = await client.request(
            method,
            f"{credentials.server_url.rstrip('/')}{path}",
            content=body,
            json=payload,
            params=query,
            headers=headers,
        )
    try:
        resp.raise_for_status()
    except httpx.HTTPStatusError as e:
        raise WallabagApiError(e)
    return resp


async def refresh_session(credentials: WallabagCredentials):
    # wallabag's implementation of OAuth2 has issues with the refresh token reliability.
    # The only way to have a consistent behavior is ignore the refresh token and to
    # systematically use the password grant type.
    logger.info("attempting to refresh wallabag session")

    payload = {
        "grant_type": "password",
        "client_id": credentials.client_id,
        "client_secret": credentials.client_secret,
        "username": credentials.username,
        "password": credentials.password,
    }

    resp = await request_wallabag(
        credentials,
        "POST",
        "/oauth/v2/token",
        payload=payload,
    )
    data = resp.json()
    await sync_to_async(credentials.update_token)(data)
    logger.info("wallabag session refreshed")
