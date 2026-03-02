from dataclasses import dataclass
from datetime import datetime

from accounts.models import User
from services.models import ReadProgress


@dataclass
class ReadProgressRecord:
    updated_at: datetime
    article_id: int
    progress: float


class ReadProgressService:
    def __init__(self, user: User):
        self.user = user

    async def list_since(self, dt: datetime) -> list[ReadProgressRecord]:
        return [
            ReadProgressRecord(
                updated_at=it.updated_at,
                article_id=it.article_id,
                progress=it.progress,
            )
            async for it in ReadProgress.objects.of(self.user).since(dt)
        ]

    async def compute_db_update(
        self, incoming: list[ReadProgressRecord]
    ) -> list[ReadProgressRecord]:
        oldest_update = min(incoming, key=lambda x: x.updated_at)
        progress_index = await (
            ReadProgress.objects.of(self.user)
            .since(oldest_update.updated_at)
            .ain_bulk(field_name="article_id")
        )

        update_set = []
        for it in incoming:
            update = None
            if progress := progress_index.get(it.article_id):
                if it.updated_at > progress.updated_at:
                    progress.updated_at = it.updated_at
                    update = progress
            if update is None:
                update_set.append(
                    ReadProgress(
                        user=self.user,
                        article_id=it.article_id,
                        progress=it.progress,
                    )
                )
        return update_set

    async def apply_db_update(self, records: list[ReadProgressRecord]):
        update_set = [
            ReadProgress(
                user=self.user,
                article_id=it.article_id,
                progress=it.progress,
                updated_at=it.updated_at,
            )
            for it in records
        ]
        await ReadProgress.objects.abulk_create(
            update_set,
            update_conflicts=True,
            update_fields=["updated_at", "progress"],
            unique_fields=["article_id"],
        )

    async def clear_stalled_progresses(self, server_entry_ids: list[int]) -> int:
        user_progresses = ReadProgress.objects.of(self.user)

        server_ids = set(server_entry_ids)
        local_ids = {
            it async for it in user_progresses.values_list("article_id", flat=True)
        }
        stalled_ids = local_ids - server_ids

        result = await user_progresses.filter(article_id__in=stalled_ids).adelete()
        return result[1].get(ReadProgress._meta.label, 0)
