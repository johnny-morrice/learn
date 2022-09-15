#!/bin/sh
set -e

/app/learnblog --database "$DATABASE" --migrations "$MIGRATIONS" $@