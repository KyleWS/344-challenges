#!/usr/bin/env bash
./build.sh
ssh root@67.205.178.99 '
docker network rm apinetwork
docker network create apinetwork
docker rm -f redisServer
docker run -d --name redisServer --network apinetwork redis
docker rm -f mongodb
docker run -d --name mongodb --network apinetwork mongo
docker pull kylews/gateway
docker rm -f gateway
docker run -d  \
-p 443:443 \
--name gateway \
--network apinetwork \
-v /etc/letsencrypt:/letsencrypt:ro \
-e TLSCERT=/letsencrypt/live/api.kylewilliamscreates.com/fullchain.pem \
-e TLSKEY=/letsencrypt/live/api.kylewilliamscreates.com/privkey.pem \
-e SESSIONKEY=8b3f95a3bb29d578eb4544607856e4de \
-e REDISADDR=redisServer:6379 \
-e DBADDR=mongodb:27017 \
-e MESSAGESVCADDR=messagingSVC:5000 \
-e SUMMARYSVCADDR=summarySVC:4001 \
kylews/gateway
'
