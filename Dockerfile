# syntax=docker/dockerfile:1

ARG VERSION=unknown

#-------------------------------------------------------------------------------- 
# Build virtualenv
#-------------------------------------------------------------------------------- 

FROM python:3.13 AS venv

ENV POETRY_VERSION=2.1.3
ENV VIRTUAL_ENV=/opt/venv

RUN pip install --upgrade pip
RUN pip install poetry==${POETRY_VERSION}
RUN poetry config virtualenvs.in-project true

WORKDIR /src
COPY freon_server/poetry.lock freon_server/pyproject.toml ./
RUN poetry install --no-interaction --only main,webserver

#-------------------------------------------------------------------------------- 
# Final image
#-------------------------------------------------------------------------------- 

FROM python:3.13-slim

ARG VERSION

COPY --from=venv /src/.venv /opt/venv
ENV VIRTUAL_ENV=/opt/venv
ENV PATH="/opt/venv/bin:${PATH}"

COPY freon_server /src
WORKDIR /src

ENV DEBUG=false
ENV FREON_DB_PATH=/var/lib/freon/data/freon.db
ENV GRANIAN_STATIC_PATH_MOUNT=/src/staticfiles
ENV GRANIAN_STATIC_PATH_EXPIRES=3600
ENV LOAD_DOTENV=false
ENV PYTHONUNBUFFERED=1
ENV VERSION=${VERSION}

RUN mkdir -p $(dirname ${FREON_DB_PATH})
RUN ONESHOT_SECRET_KEY=true python manage.py collectstatic --noinput

CMD ["/src/run.sh", "--help"]