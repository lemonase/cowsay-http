#!/usr/bin/env bash
#
# Build and push docker image for amd64 and arm64

docker buildx build --platform linux/amd64 -t jamesdixon/cowsay-http:latest-amd64 . && docker push jamesdixon/cowsay-http:latest-amd64
docker buildx build --platform linux/arm64 -t jamesdixon/cowsay-http:latest-arm64 . && docker push jamesdixon/cowsay-http:latest-arm64
