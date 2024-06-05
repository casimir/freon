#!/bin/sh

awk -F'[ +]' '/version:/{ print $2 }' ui/pubspec.yaml 