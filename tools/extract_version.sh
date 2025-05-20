#!/bin/sh

awk -F'[ +]' '/version = "/{ print $3 }' freon_server/pyproject.toml | tr -d '"'