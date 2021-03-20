#!/bin/bash

go env -w GOPROXY=https://goproxy.cn,direct
go mod download
CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server knativetest/cmd/"$1"

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
IMAGE=$1:latest
docker build -t "$IMAGE" "$DIR"/../ -f "$DIR"/../deploy/Dockerfile_"$1"
