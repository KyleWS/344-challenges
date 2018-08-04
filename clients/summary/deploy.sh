#!/usr/bin/env bash
./build.sh
ssh root@67.205.178.73 '
docker pull kylews/summary
docker rm -f summary
docker run -d  \
-p 80:80 \
-p 443:443 \
--name summary \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
kylews/summary
'
