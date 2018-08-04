#!/usr/bin/env bash
./build.sh
ssh root@67.205.178.99 '
docker pull kylews/messaging
docker rm -f messagingSVC
docker run -d  \
--name messagingSVC \
--network apinetwork \
-e DBADDR=mongodb:27017 \
-e MSGADDR=messagingSVC:5000 \
kylews/messaging
'
