import uuid
from dataclasses import dataclass

from ninja.errors import AuthorizationError
from ninja.security import HttpBearer

from . import logger
from .models import Token, TokenScope


@dataclass
class TokenScopeConfig:
    name: str
    description: str


class TokenAuth(HttpBearer):
    scopes: list[TokenScopeConfig] | None = None

    def __init__(self, scopes: list[TokenScopeConfig] | None = None):
        super().__init__()
        self._scopes = scopes or self.scopes
        self._scopes_checked = False

    async def _ensure_scopes(self):
        if self._scopes:
            for scope in self._scopes:
                _, created = await TokenScope.objects.aget_or_create(
                    name=scope.name, defaults={"description": scope.description}
                )
                if created:
                    logger.info(f"created token scope: {scope.name}")
        self._scopes_checked = True

    def _parse_token(self, token: str) -> uuid.UUID | None:
        try:
            return uuid.UUID(token)
        except (TypeError, ValueError):
            return None

    async def matches_scopes(self, account_token: Token | None) -> bool:
        if account_token is None:
            return False

        if not self._scopes:
            return True

        scopes_count = await account_token.scopes.filter(
            name__in=[scope.name for scope in self._scopes]
        ).acount()
        return scopes_count == len(self._scopes)

    async def authenticate(self, request, token):
        scopes_fut = None
        if not self._scopes_checked:
            scopes_fut = self._ensure_scopes()

        account_token = await (
            Token.objects.filter(id=self._parse_token(token))
            .select_related("user")
            .afirst()
        )

        if scopes_fut is not None:
            await scopes_fut

        if account_token is not None and account_token.is_expired():
            raise AuthorizationError(
                message=f"token expired since {account_token.expires_at}"
            )

        if not await self.matches_scopes(account_token):
            raise AuthorizationError()

        return account_token
