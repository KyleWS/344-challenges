#!/usr/bin/env bash
./build.sh
ssh root@67.205.178.73 '
docker rm -f summary
docker pull kylews/messages
docker rm -f messages
docker run -d  \
-p 80:80 \
-p 443:443 \
--name messages \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
kylews/messages
'
