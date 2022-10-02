#!/usr/bin/env bash
set -e

docker run --network learnblog_learnblog-network --env-file env/dev/learnblog.env -it --rm learnblog /app/entrypoint.sh --command migrate-up