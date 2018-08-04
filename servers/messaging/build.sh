#!/usr/bin/env bash
set -e
docker build -t kylews/messaging .
docker push kylews/messaging
go clean
