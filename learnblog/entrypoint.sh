#!/usr/bin/env bash
set -e

/app/learnblog --database $DATABASE --migrations $MIGRATIONS $@