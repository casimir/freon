from datetime import datetime
from http import HTTPStatus

from django.conf import settings
from django.http import HttpRequest, HttpResponse
from ninja import Router, Schema

from accounts.security import TokenAuth, TokenAuthHttpRequest
from services.read_progress import ReadProgressService, ReadProgressUpdates

router = Router(auth=TokenAuth())


class InfoOut(Schema):
    appname: str
    version: str


@router.get("/info", response=InfoOut, auth=None)
async def info(request: HttpRequest):
    return {"appname": "freon", "version": settings.VERSION}


class UserOut(Schema):
    username: str


@router.get("/me", response=UserOut)
async def index(request: TokenAuthHttpRequest):
    return request.auth.user


class GetReadProgressOut(Schema):
    updates: ReadProgressUpdates


@router.get("/read-progress", response=GetReadProgressOut)
async def get_read_progress(request: TokenAuthHttpRequest, since: datetime):
    service = ReadProgressService(request.auth.user)
    return await service.list_since(since)


class PutReadProgressIn(Schema):
    updates: ReadProgressUpdates


@router.put("/read-progress")
async def put_read_progress(request: TokenAuthHttpRequest, data: PutReadProgressIn):
    service = ReadProgressService(request.auth.user)
    update_set = await service.compute_db_update(data.updates)
    await service.apply_db_update(update_set)
    return HttpResponse(status=HTTPStatus.NO_CONTENT)
