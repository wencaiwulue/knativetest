#! /bin/bash

kubectl label namespace test inject-empty-container=enabled
openssl req -x509 -nodes -new -sha256 -days 3650 -newkey rsa:2048 -subj "/CN=diyadmissionwebhook.test.svc" \
  -keyout ca.key \
  -out ca.crt

openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -config csr.conf

openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out server.crt -days 3650 \
  -extensions v3_ext -extfile csr.conf

kubectl create secret generic webhook-secret -n test \
  --from-file=./ca.crt \
  --from-file=./server.key \
  --from-file=./server.crt

rm ca.key ca.crt server.key server.crt
