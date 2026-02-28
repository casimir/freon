#!/bin/sh

uv run python manage.py collectstatic --noinput
uv run granian --interface asginl freon_server/asgi.py:application \
    --reload --reload-ignore-patterns '.+\.sqlite3?(-journal)?$' \
    --reload-ignore-worker-failure \
    --access-log \
    --log-level debug