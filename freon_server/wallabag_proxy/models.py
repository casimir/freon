from datetime import timedelta

from accounts.models import User
from django.db import models, transaction
from django.utils import timezone


class WallabagToken(models.Model):
    access_token = models.CharField(max_length=255)
    expires_at = models.DateTimeField()
    refresh_token = models.CharField(max_length=255)


# Username and password should not be stored in the database but the current implementation
# of OAuth2 in wallabag has issues with the refresh token.


class WallabagCredentials(models.Model):
    user = models.OneToOneField(
        User, on_delete=models.CASCADE, related_name="wallabag_credentials"
    )
    server_url = models.URLField()
    client_id = models.CharField(max_length=255)
    client_secret = models.CharField(max_length=255)
    username = models.CharField(max_length=255)
    password = models.CharField(max_length=255)
    token = models.OneToOneField(
        WallabagToken, on_delete=models.CASCADE, blank=True, null=True
    )

    class Meta:
        verbose_name_plural = "wallabag credentials"

    @property
    def has_valid_session(self) -> bool:
        if self.token is None:
            return False
        return self.token.expires_at > timezone.now()

    def update_token(self, payload: dict):
        """
        Update the token for the credentials based on the payload from the wallabag API.

        The old token is deleted and the new token is created in a single transaction.

        Args:
            payload: The payload from the wallabag API /oauth/v2/token.
        """
        with transaction.atomic():
            if self.token is not None:
                self.token.delete()
            self.token = WallabagToken.objects.create(
                access_token=payload["access_token"],
                expires_at=timezone.now() + timedelta(seconds=payload["expires_in"]),
                refresh_token=payload["refresh_token"],
            )
            self.save()
