#!/usr/bin/env bash
set -e

docker build . -t learnblog
docker run learnblog