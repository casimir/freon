from datetime import datetime, timedelta, timezone

from django.test import TestCase

from accounts.models import User

from .models import ReadProgress
from .read_progress import ReadProgressService, ReadProgressUpdate, ReadProgressUpdates


class ReadProgressTest(TestCase):
    @classmethod
    def setUpTestData(cls):
        cls.user = User.objects.create_user(username="testuser", password="testpass123")

    def setUp(self):
        self.ref_moment = datetime(2025, 1, 1, 12, 0, 0, tzinfo=timezone.utc)

        ReadProgress.objects.create(
            user=self.user,
            article_id=1,
            progress=0.5,
            updated_at=self.ref_moment + timedelta(hours=1),
        )
        ReadProgress.objects.create(
            user=self.user,
            article_id=3,
            progress=0.1,
            updated_at=self.ref_moment + timedelta(hours=2),
        )

    def _make_incoming_progresses(self) -> ReadProgressUpdates:
        return [
            ReadProgressUpdate(
                updated_at=self.ref_moment + timedelta(hours=3),
                article_id=2,
                progress=0.3,
            ),
            ReadProgressUpdate(
                updated_at=self.ref_moment,
                article_id=3,
                progress=0.7,
            ),
        ]

    async def test_compute_db_update(self):
        incoming = self._make_incoming_progresses()
        update = await ReadProgressService(self.user).compute_db_update(incoming)

        self.assertEqual(len(update), 2)
        self.assertEqual(update[0].article_id, 2)
        self.assertEqual(update[0].progress, 0.3)
        self.assertEqual(update[1].article_id, 3)
        self.assertEqual(update[1].progress, 0.7)

    async def test_apply_db_update(self):
        incoming = self._make_incoming_progresses()
        service = ReadProgressService(self.user)
        update = await service.compute_db_update(incoming)
        await service.apply_db_update(update)

        progresses = [
            it async for it in ReadProgress.objects.of(self.user).order_by("article_id")
        ]

        self.assertEqual(len(progresses), 3)
        self.assertEqual(progresses[0].article_id, 1)
        self.assertEqual(progresses[0].progress, 0.5)
        self.assertEqual(progresses[1].article_id, 2)
        self.assertEqual(progresses[1].progress, 0.3)
        self.assertEqual(progresses[2].article_id, 3)
        self.assertEqual(progresses[2].progress, 0.7)

    async def test_clear_stalled_progresses(self):
        service = ReadProgressService(self.user)
        deleted_count = await service.clear_stalled_progresses([1, 100])

        self.assertEqual(deleted_count, 1)
        self.assertEqual(await ReadProgress.objects.of(self.user).acount(), 1)

    async def test_clear_stalled_progresses_noop(self):
        service = ReadProgressService(self.user)
        deleted_count = await service.clear_stalled_progresses([1, 3])

        self.assertEqual(deleted_count, 0)
        self.assertEqual(await ReadProgress.objects.of(self.user).acount(), 2)
