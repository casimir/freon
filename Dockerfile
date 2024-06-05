# syntax=docker/dockerfile:1

ARG VERSION

#-------------------------------------------------------------------------------- 
# Build server
#-------------------------------------------------------------------------------- 

FROM golang:1.22-alpine as server

RUN apk update && apk add build-base

WORKDIR /src/server

COPY ./server/go.mod ./server/go.sum ./
RUN go mod download

COPY ./server/ .
RUN CGO_ENABLED=1 go build -o ./build/freon -ldflags="-X 'buildinfo.Version=${VERSION}'" .

#-------------------------------------------------------------------------------- 
# Build UI
#-------------------------------------------------------------------------------- 

FROM debian:stable-slim as ui

RUN apt-get update && apt-get install -y curl git unzip zip
RUN git clone -b stable --depth 1 https://github.com/flutter/flutter.git /opt/flutter
RUN git config --global --add safe.directory /usr/local/flutter
ENV PATH="/opt/flutter/bin:/opt/flutter/bin/cache/dart-sdk/bin:$PATH"
RUN flutter config --no-analytics --enable-web \
        --no-enable-linux-desktop --no-enable-macos-desktop --no-enable-windows-desktop \
        --no-enable-android --no-enable-ios --no-enable-fuchsia
RUN flutter precache web

WORKDIR /src/ui

COPY ./ui/pubspec.* ./
RUN flutter pub get

COPY ./ui/ .
RUN flutter build web --base-href /ui/

#-------------------------------------------------------------------------------- 
# Final image
#-------------------------------------------------------------------------------- 

FROM alpine:latest

COPY --from=server /src/server/build/freon /usr/bin/freon
COPY --from=ui /src/ui/build/web /var/lib/freon-ui

ENV GIN_MODE=release
ENV FREON_UI_PATH=/var/lib/freon-ui

CMD ["/usr/bin/freon", "server"]