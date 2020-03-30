#!/usr/bin/env bash

docker kill testmeili \
; docker run -d -p 7700:7700 --name testmeili -e MEILI_NO_ANALYTICS=1 --rm getmeili/meilisearch:v0.9.0\
; sleep 1 \
; go test -v ./...
