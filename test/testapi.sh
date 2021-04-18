#! /bin/bash

curl -X POST localhost:8080/CreateServiceAction -H 'Content-type: application/json' \
  -d '{"Name":"test", "Namespace":"test"}'

curl -X POST localhost:8080/CreateClusterTask -H 'Content-type: application/json' \
  -d '{"Name":"test", "Namespace":"test"}'

curl -X POST localhost:8080/CreateTaskRun -H 'Content-type: application/json' \
  -d '{"Name":"test", "Namespace":"test"}'

curl 127.0.0.1:32302/CreateServiceAction -d '{"Namespace":"test", "Name":"b"}'

curl -X POST 10.32.0.19:80/KvAction -H 'Content-type:applicatoon/json' \
  -d '{"Key":"testkey", "Value":"testvalue"}'
