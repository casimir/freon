# syntax=docker/dockerfile:1

ARG VERSION

#-------------------------------------------------------------------------------- 
# Build server
#-------------------------------------------------------------------------------- 

FROM golang:1.22-alpine AS server

RUN apk update && apk add build-base

COPY ./server/go.* /src/server/
WORKDIR /src/server
RUN go mod download

COPY . /src
WORKDIR /src
RUN make server-headless

#-------------------------------------------------------------------------------- 
# Build UI
#-------------------------------------------------------------------------------- 

FROM debian:stable-slim AS ui

RUN apt-get update && apt-get install -y curl git make unzip zip
RUN git clone -b stable --depth 1 https://github.com/flutter/flutter.git /opt/flutter
RUN git config --global --add safe.directory /usr/local/flutter
ENV PATH="/opt/flutter/bin:/opt/flutter/bin/cache/dart-sdk/bin:$PATH"
RUN flutter config --no-analytics --enable-web \
        --no-enable-linux-desktop --no-enable-macos-desktop --no-enable-windows-desktop \
        --no-enable-android --no-enable-ios --no-enable-fuchsia
RUN flutter precache web

COPY ./ui/pubspec.* ./
WORKDIR /src/ui
RUN flutter pub get

COPY . /src
WORKDIR /src
RUN make ui

#-------------------------------------------------------------------------------- 
# Final image
#-------------------------------------------------------------------------------- 

FROM alpine:latest

COPY --from=server /src/server/build/freon-headless /usr/bin/freon
RUN mkdir -p /var/lib/freon
RUN mkdir -p /var/lib/freon/data
COPY --from=ui /src/ui/build/web /var/lib/freon/ui

ENV GIN_MODE=release
ENV FREON_DB_PATH=/var/lib/freon/data/freon.db
ENV FREON_UI_PATH=/var/lib/freon/ui

CMD ["/usr/bin/freon", "server"]