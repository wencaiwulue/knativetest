#!/bin/bash

go env -w GOPROXY=https://goproxy.cn,direct
go mod download
CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o "$1" knativetest/cmd/"$1"

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
IMAGE=$1:latest
docker build -t "$IMAGE" "$DIR"/../ -f "$DIR"/../deploy/Dockerfile_"$1"
rm "$1"

# restart the pods
for POD in $(kubectl get pods -n test | awk -F ' ' '{print$1}' | grep -v NAME); do
  # shellcheck disable=SC2046
  if [ "$IMAGE" == $(kubectl get pods "$POD" -n test -o jsonpath='{.spec.containers[0].image}') ]; then
    kubectl delete pods "$POD" -n test
    break
  fi
done
