#!/bin/sh

awk '/^version = "/{ print $3 }' freon_server/pyproject.toml | tr -d '"'