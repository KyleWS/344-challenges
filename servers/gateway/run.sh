#!/usr/bin/env bash
docker rm -f gateway
docker run -d  \
-p 443:443 \
--name gateway \
-v $(pwd)/tls:/tls:ro \
-e TLSCERT=/tls/fullchain.pem \
-e TLSKEY=/tls/privkey.pem \
kylews/gateway
