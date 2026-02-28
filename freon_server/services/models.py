from datetime import datetime
from typing import ClassVar

from django.db import models

from accounts.models import User


class ScopedManager[M: models.Model, R: models.QuerySet[M]](models.Manager[M]):

    def of(self, user: User | int, user_field: str = "user") -> R:
        field = user_field
        if isinstance(user, int):
            field += "_id"
        return self.get_queryset().filter(**{field: user})


class ReadProgressQuerySet(models.QuerySet["ReadProgress"]):

    def since(self, dt: datetime):
        return self.filter(updated_at__gte=dt)


ReadProgressManager = ScopedManager.from_queryset(ReadProgressQuerySet)


class ReadProgress(models.Model):
    updated_at = models.DateTimeField(auto_now=True)
    user = models.ForeignKey(User, on_delete=models.CASCADE)
    article_id = models.IntegerField(unique=True)
    progress = models.FloatField()

    objects: ClassVar[ScopedManager["ReadProgress", ReadProgressQuerySet]] = (
        ReadProgressManager()
    )

    class Meta:
        verbose_name_plural = "read progresses"
