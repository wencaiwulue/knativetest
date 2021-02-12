#! /bin/bash

# shellcheck disable=SC2164
cd knative

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v

# shellcheck disable=SC2103
cd ..

docker build -t test:latest .

kubectl delete -f test.yaml

kubectl apply -f test.yaml
