#!/usr/bin/env bash
./build.sh
ssh root@67.205.178.99 '
docker pull kylews/summary
docker rm -f summarySVC
docker run -d  \
--name summarySVC \
--network apinetwork \
-e DBADDR=mongodb:27017 \
-e SUMMARYADDR=summarySVC:4001 \
kylews/summary
'
