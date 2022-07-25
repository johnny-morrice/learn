#!/usr/bin/env bash
set -e

docker build . -t learnblog
docker run -p 8080:8080/tcp learnblog