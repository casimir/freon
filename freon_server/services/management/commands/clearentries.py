import asyncio

from django.core.management.base import BaseCommand

from accounts.models import User
from services.read_progress import ReadProgressService
from services.wallabag import WallabagService


class Command(BaseCommand):
    help = (
        "Can be run as a cronjob or directly to clean out server deleted entries from "
        "the database."
    )

    def handle(self, **options):
        count = asyncio.run(self._ahandle(**options))
        self.stdout.write(f"Cleared {count} stalled entries")

    async def _ahandle(self, **options):
        total_deleted = 0
        wallabag_users = User.objects.filter(wallabag_credentials__isnull=False)

        async for user in wallabag_users.aiterator():
            entry_ids = await WallabagService(user).list_all_entry_ids()
            total_deleted += await ReadProgressService(user).clear_stalled_progresses(
                entry_ids
            )

        return total_deleted
