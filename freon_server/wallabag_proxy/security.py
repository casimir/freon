from accounts.security import TokenAuth, TokenScopeConfig
from asgiref.sync import sync_to_async

from .models import WallabagCredentials


class WallabagProxyAuth(TokenAuth):
    scopes = [
        TokenScopeConfig(name="wallabag", description="Wallabag API access"),
    ]

    async def authenticate(self, request, token) -> WallabagCredentials | None:
        account_token = await super().authenticate(request, token)
        user = account_token.user
        try:
            return await sync_to_async(
                WallabagCredentials.objects.select_related("token").get
            )(user=user)
        except WallabagCredentials.DoesNotExist:
            return None
