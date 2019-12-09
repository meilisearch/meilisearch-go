#!/usr/bin/env bash

docker kill testmeili \
; docker run -d -p 7700:7700 --name testmeili --rm getmeili/meilisearch \
; sleep 1 \
; go test -v ./...