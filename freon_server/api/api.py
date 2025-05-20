from django.conf import settings
from ninja import Router, Schema

from accounts.security import TokenAuth

router = Router(auth=TokenAuth())


class InfoOut(Schema):
    appname: str
    version: str


@router.get("/info", response=InfoOut, auth=None)
async def info(request):
    return {"appname": "freon", "version": settings.VERSION}


class UserOut(Schema):
    username: str


@router.get("/me", response=UserOut)
async def index(request):
    return request.auth.user
