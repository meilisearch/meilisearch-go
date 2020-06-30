#!/usr/bin/env bash

docker kill testmeili \
; docker run -d -p 7700:7700 getmeili/meilisearch:latest ./meilisearch --master-key=masterKey --no-analytics=true\
; sleep 1 \
; go test -v -count=1 ./...
