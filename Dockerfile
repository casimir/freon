# syntax=docker/dockerfile:1

ARG VERSION=unknown

#-------------------------------------------------------------------------------- 
# Build virtualenv
#-------------------------------------------------------------------------------- 

FROM python:3.13 AS venv

ENV UV_PROJECT_ENVIRONMENT=/opt/venv

COPY --from=ghcr.io/astral-sh/uv:0.10.6 /uv /usr/local/bin/uv

WORKDIR /src
COPY freon_server/uv.lock freon_server/pyproject.toml ./
RUN uv sync --frozen --no-dev --group webserver --no-install-project --link-mode=copy

#-------------------------------------------------------------------------------- 
# Final image
#-------------------------------------------------------------------------------- 

FROM python:3.13-slim

ARG VERSION

COPY --from=venv /opt/venv /opt/venv
ENV VIRTUAL_ENV=/opt/venv
ENV PATH="/opt/venv/bin:${PATH}"

COPY freon_server /src
WORKDIR /src

ENV DEBUG=false
ENV FREON_DB_PATH=/var/lib/freon/data/freon.db
ENV GRANIAN_HOST=0.0.0.0
ENV GRANIAN_STATIC_PATH_MOUNT=/src/staticfiles/
ENV LOAD_DOTENV=false
ENV PYTHONUNBUFFERED=1
ENV VERSION=${VERSION}

RUN mkdir -p $(dirname ${FREON_DB_PATH})
RUN ONESHOT_SECRET_KEY=true python manage.py collectstatic --noinput

CMD ["/src/run.sh", "--help"]