#!/usr/bin/env bash
set -e
CGO_ENABLED=0 go build -a
docker build -t kylews/summary .
docker push kylews/summary
go clean
