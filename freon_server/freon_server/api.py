import orjson
from ninja import NinjaAPI
from ninja.parser import Parser
from ninja.renderers import BaseRenderer
from wallabag_proxy.api import router as wallabag_proxy_router

from api.api import router as api_router


class ORJSONParser(Parser):
    def parse_body(self, request):
        return orjson.loads(request.body)


class ORJSONRenderer(BaseRenderer):
    media_type = "application/json"

    def render(self, request, data, *, response_status):
        return orjson.dumps(data)


api = NinjaAPI(
    title="Freon API",
    openapi_url="/api/openapi.json",
    docs_url="/api/docs",
    parser=ORJSONParser(),
    renderer=ORJSONRenderer(),
)

api.add_router("/api/", api_router)
api.add_router("/wallabag/", wallabag_proxy_router)
