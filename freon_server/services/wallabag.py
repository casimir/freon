from accounts.models import User
from wallabag_proxy.models import WallabagCredentials
from wallabag_proxy.wallabag import request_wallabag


class WallabagService:
    def __init__(self, user: User):
        self.user = user

    async def list_all_entry_ids(self):
        credentials = await WallabagCredentials.objects.select_related("token").aget(
            user=self.user
        )
        entry_ids = []

        pages = 1
        page = 0
        while page < pages:
            resp = await request_wallabag(
                credentials,
                "GET",
                "/api/entries",
                {"perPage": 500, "detail": "metadata", "page": page + 1},
            )
            data = resp.json()
            entry_ids += [x["id"] for x in data["_embedded"]["items"]]
            pages = data["pages"]
            page += 1

        return entry_ids
