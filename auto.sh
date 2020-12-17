#! /bin/bash

cd knative

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v

cd ..

docker build -t test:latest .

kubectl delete -f test.yaml

kubectl apply -f test.yaml
