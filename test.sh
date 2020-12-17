curl -X POST localhost:8080/CreateServiceAction -H 'Content-type: application/json' -d '{"Name":"test", "Namespace":"test"}'

curl -X POST localhost:8080/createClusterTask -H 'Content-type: application/json' -d '{"Name":"test", "Namespace":"test"}'

curl -X POST localhost:8080/createTaskRun -H 'Content-type: application/json' -d '{"Name":"test", "Namespace":"test"}'
