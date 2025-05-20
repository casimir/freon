from datetime import timedelta

from django.test import TestCase
from django.utils import timezone
from ninja.errors import AuthorizationError

from .models import Token, TokenScope, User
from .security import TokenAuth, TokenScopeConfig


class TokenValidationTest(TestCase):
    def setUp(self):
        self.user = User.objects.create_user(
            username="testuser", password="testpass123"
        )

    async def test_valid_token_authentication(self):
        auth = TokenAuth()
        valid_token = await Token.objects.acreate(user=self.user)

        token = await auth.authenticate(None, str(valid_token.id))
        self.assertEqual(token, valid_token)

    async def test_expired_token_authentication(self):
        auth = TokenAuth()
        expired_token = await Token.objects.acreate(
            user=self.user, expires_at=timezone.now() - timedelta(days=1)
        )

        with self.assertRaises(AuthorizationError) as context:
            await auth.authenticate(None, str(expired_token.id))
        self.assertIn("token expired since ", context.exception.message)

    async def test_invalid_token(self):
        auth = TokenAuth()

        with self.assertRaises(AuthorizationError):
            await auth.authenticate(None, "00000000-0000-0000-0000-000000000000")
        with self.assertRaises(AuthorizationError):
            await auth.authenticate(None, "invalid-token")
        with self.assertRaises(AuthorizationError):
            await auth.authenticate(None, "")


class TokenScopeConfigTest(TestCase):
    def _create_token(self, user, scopes=None) -> str:
        token = Token.objects.create(user=user)
        if scopes:
            token.scopes.set(scopes)
        return str(token.id)

    def setUp(self):
        self.user = User.objects.create_user(
            username="testuser", password="testpass123"
        )

        self.read_scope_config = TokenScopeConfig(
            name="read", description="Read access"
        )
        read_scope = TokenScope.objects.create(
            name=self.read_scope_config.name,
            description=self.read_scope_config.description,
        )
        self.write_scope_config = TokenScopeConfig(
            name="write", description="Write access"
        )
        write_scope = TokenScope.objects.create(
            name=self.write_scope_config.name,
            description=self.write_scope_config.description,
        )

        self.token_with_read = self._create_token(self.user, [read_scope])
        self.token_with_write = self._create_token(self.user, [write_scope])
        self.token_with_both = self._create_token(self.user, [read_scope, write_scope])
        self.token_with_none = self._create_token(self.user)

    async def test_auth_without_scopes(self):
        auth = TokenAuth()

        self.assertIsNotNone(await auth.authenticate(None, self.token_with_read))
        self.assertIsNotNone(await auth.authenticate(None, self.token_with_write))
        self.assertIsNotNone(await auth.authenticate(None, self.token_with_both))
        self.assertIsNotNone(await auth.authenticate(None, self.token_with_none))

    async def test_auth_single_scope(self):
        auth = TokenAuth(scopes=[self.read_scope_config])

        self.assertIsNotNone(await auth.authenticate(None, self.token_with_read))
        with self.assertRaises(AuthorizationError):
            await auth.authenticate(None, self.token_with_write)
        self.assertIsNotNone(await auth.authenticate(None, self.token_with_both))
        with self.assertRaises(AuthorizationError):
            await auth.authenticate(None, self.token_with_none)

    async def test_auth_multiple_scopes(self):
        auth = TokenAuth(scopes=[self.read_scope_config, self.write_scope_config])

        with self.assertRaises(AuthorizationError):
            await auth.authenticate(None, self.token_with_read)
        with self.assertRaises(AuthorizationError):
            await auth.authenticate(None, self.token_with_write)
        self.assertIsNotNone(await auth.authenticate(None, self.token_with_both))
        with self.assertRaises(AuthorizationError):
            await auth.authenticate(None, self.token_with_none)
