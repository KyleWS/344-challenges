#!/usr/bin/env bash
npm run build
docker build -t kylews/summary .
docker push kylews/summary
