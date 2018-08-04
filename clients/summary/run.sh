#!/usr/bin/env bash
docker rm -f summary
docker run -d --name summary -p 80:80 -p 443:443 kylews/summary
